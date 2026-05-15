#!/usr/bin/env python3
"""Build and smoke-test the Docker example matrix locally without publishing.

Examples:
    python3 scripts/check-wolfi-examples.py
    python3 scripts/check-wolfi-examples.py --example wolfi-java-runtime --example wolfi-kafka
    python3 scripts/check-wolfi-examples.py --platform linux/amd64 --platform linux/arm64

The default platform matches CI (linux/amd64). Every build pulls base images and
reruns the builder stage, including APK index fetching and slice extraction.
Go compilation layers may be reused. Requires Python 3, Make, Docker Buildx,
network access, and emulation support when testing a different architecture.
"""

import argparse
from concurrent.futures import ThreadPoolExecutor
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import tempfile
import uuid


REPO_ROOT = Path(__file__).resolve().parents[1]


def matrix_examples():
    workflow = (REPO_ROOT / ".github/workflows/wolfi-ci.yaml").read_text()
    job = re.search(r"(?ms)^  docker-examples:\n(.*?)(?=^  [\w-]+:\n|\Z)", workflow)
    if not job:
        raise ValueError("cannot find the Docker example job in wolfi-ci.yaml")
    examples = re.findall(r"^\s+- (wolfi-[\w-]+)\s*$", job[1], re.MULTILINE)
    if not examples or len(examples) != len(set(examples)):
        raise ValueError("missing or duplicate Docker example matrix entries")
    for name in examples:
        if not (REPO_ROOT / "examples" / name / "Makefile").is_file():
            raise ValueError(f"missing Makefile for {name}")
    return examples


def positive_int(value):
    number = int(value)
    if number < 1:
        raise argparse.ArgumentTypeError("must be a positive integer")
    return number


def main():
    examples = matrix_examples()
    parser = argparse.ArgumentParser(
        description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument("--example", action="append", choices=examples, metavar="EXAMPLE", help="repeat to select examples; default: all CI examples")
    parser.add_argument("--platform", action="append", choices=("linux/amd64", "linux/arm64"), help="repeat to test both; default: linux/amd64")
    parser.add_argument("--jobs", type=positive_int, default=2, help="parallel builds (default: 2)")
    parser.add_argument("--output-dir", type=Path, help="new directory for logs and summary.json")
    args = parser.parse_args()
    selected = list(dict.fromkeys(args.example or examples))
    platforms = list(dict.fromkeys(args.platform or ["linux/amd64"]))

    if args.output_dir:
        output = args.output_dir.resolve()
        output.mkdir(parents=True, exist_ok=False)
    else:
        output = Path(tempfile.mkdtemp(prefix="chisel-wolfi-examples-"))
    print(f"Fresh builder stages; reports: {output}", flush=True)
    cases = [(name, platform) for name in selected for platform in platforms]
    run_id = uuid.uuid4().hex[:12]

    def check(case):
        name, platform = case
        arch = platform.split("/")[1]
        log_path = output / f"{name}-{arch}.log"
        # Separate tags prevent tests for different platforms from racing.
        image = f"chisel-local-checks:{run_id}-{name}-{arch}"
        command = [
            "make", "-C", str(REPO_ROOT / "examples" / name),
            f"IMAGE={image}", "BUILD_FLAGS=--pull --no-cache-filter=builder",
            "build", "test", "size",
        ]
        env = dict(os.environ, DOCKER_DEFAULT_PLATFORM=platform, DOCKER_BUILDKIT="1", BUILDKIT_PROGRESS="plain")
        with log_path.open("w") as log:
            try:
                result = subprocess.run(
                    command, cwd=REPO_ROOT, env=env,
                    stdout=log, stderr=subprocess.STDOUT,
                )
                code = result.returncode
            except OSError as error:
                print(error, file=log)
                code = 1
        status = "PASS" if code == 0 else "FAIL"
        print(f"{status} {name} ({platform}): {log_path}", flush=True)
        return {"example": name, "platform": platform, "exit_code": code, "log": str(log_path), "image": image}

    with ThreadPoolExecutor(max_workers=args.jobs) as pool:
        results = list(pool.map(check, cases))
    (output / "summary.json").write_text(json.dumps(results, indent=2) + "\n")
    passed = sum(result["exit_code"] == 0 for result in results)
    print(f"{passed}/{len(results)} cases passed. Reports: {output}", flush=True)
    return 0 if passed == len(results) else 1


if __name__ == "__main__":
    try:
        sys.exit(main())
    except (OSError, ValueError) as error:
        print(f"error: {error}", file=sys.stderr)
        sys.exit(1)
