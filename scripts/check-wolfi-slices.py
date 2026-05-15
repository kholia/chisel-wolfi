#!/usr/bin/env python3
"""Run the Wolfi slice matrix locally, retaining roots and logs for inspection.

Examples:
    python3 scripts/check-wolfi-slices.py
    python3 scripts/check-wolfi-slices.py --chisel ./chisel --arch x86_64
    python3 scripts/check-wolfi-slices.py --jobs 2 --output-dir /tmp/wolfi-checks

Requires Python 3, Make, network access, and Go unless --chisel is supplied.
This checks package extraction and the filesystem assertions from CI; it does
not execute the extracted programs or run the Docker example smoke tests.
"""

import argparse
from concurrent.futures import ThreadPoolExecutor
import os
from pathlib import Path
import re
import subprocess
import sys
import tempfile
import textwrap


REPO_ROOT = Path(__file__).resolve().parents[1]


def positive_int(value):
    result = int(value)
    if result < 1:
        raise argparse.ArgumentTypeError("must be a positive integer")
    return result


def matrix_cases():
    workflow = (REPO_ROOT / ".github/workflows/wolfi-ci.yaml").read_text()
    # Read the small, regular matrix directly so this script needs no PyYAML.
    cases = re.findall(
        r"- name: ([\w-]+)\n\s+arch: (\w+)\n\s+slice_set: (\w+)",
        workflow,
    )
    if not cases:
        raise ValueError("cannot find the Wolfi slice matrix in wolfi-ci.yaml")
    if len({name for name, _, _ in cases}) != len(cases):
        raise ValueError("duplicate names in the Wolfi slice matrix")
    return cases


def verification_script():
    lines = (REPO_ROOT / ".github/workflows/wolfi-ci.yaml").read_text().splitlines()
    start = next((
        i for i, line in enumerate(lines)
        if line.strip() == "- name: Verify ${{ matrix.name }}"
    ), None)
    if start is None:
        raise ValueError("cannot find the Wolfi filesystem verification step")
    for i in range(start + 1, len(lines)):
        if lines[i].strip() == "run: |":
            indent = len(lines[i]) - len(lines[i].lstrip())
            body = []
            for line in lines[i + 1:]:
                if line.strip() and len(line) - len(line.lstrip()) <= indent:
                    break
                body.append(line)
            if body:
                return textwrap.dedent("\n".join(body))
            break
    raise ValueError("cannot find the Wolfi filesystem verification script")


def main():
    parser = argparse.ArgumentParser(
        description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument("--chisel", type=Path, help="existing Chisel binary; otherwise build one")
    parser.add_argument("--arch", choices=("x86_64", "aarch64"), help="run only this architecture")
    parser.add_argument("--case", help="run only this matrix case, for example base-shell-x86_64")
    parser.add_argument("--jobs", type=positive_int, default=2, help="parallel cuts (default: 2)")
    parser.add_argument("--output-dir", type=Path, help="new directory for roots and logs")
    parser.add_argument("--cache-dir", type=Path, help="set XDG_CACHE_HOME for Chisel")
    args = parser.parse_args()

    cases = matrix_cases()
    if args.arch:
        cases = [case for case in cases if case[1] == args.arch]
    if args.case:
        cases = [case for case in cases if case[0] == args.case]
    if not cases:
        parser.error("no matrix cases match the selected filters")
    verify = verification_script()

    if args.output_dir:
        output = args.output_dir.resolve()
        output.mkdir(parents=True, exist_ok=False)
    else:
        output = Path(tempfile.mkdtemp(prefix="chisel-wolfi-checks-"))
    print(f"Reports and extracted roots: {output}", flush=True)

    if args.chisel:
        chisel = args.chisel.resolve()
        if not chisel.is_file() or not os.access(chisel, os.X_OK):
            parser.error(f"Chisel binary is not executable: {chisel}")
    else:
        chisel = output / "chisel"
        subprocess.run(
            ["go", "build", "-trimpath", "-o", str(chisel), "./cmd/chisel"],
            cwd=REPO_ROOT,
            check=True,
        )

    env = os.environ.copy()
    if args.cache_dir:
        env["XDG_CACHE_HOME"] = str(args.cache_dir.resolve())

    def check(case):
        name, arch, slice_set = case
        root = output / name
        root.mkdir()
        log_path = output / f"{name}.log"
        with log_path.open("w") as log:
            try:
                slices = subprocess.check_output(
                    ["make", "-s", "-f", "examples/wolfi-common-slices.mk", "print-slices", f"SLICE_SET={slice_set}"],
                    cwd=REPO_ROOT,
                    stderr=log,
                    text=True,
                ).split()
                subprocess.run(
                    [str(chisel), "cut", "--release", "examples/wolfi-release", "--root", str(root), "--arch", arch, *slices],
                    cwd=REPO_ROOT,
                    env=env,
                    stdout=log,
                    stderr=subprocess.STDOUT,
                    check=True,
                )
                subprocess.run(
                    ["sh", "-eu", "-c", verify.replace("${{ matrix.name }}", name)],
                    cwd=REPO_ROOT,
                    env=dict(env, CHISEL_ROOT=str(root)),
                    stdout=log,
                    stderr=subprocess.STDOUT,
                    check=True,
                )
            except (OSError, subprocess.CalledProcessError) as error:
                print(error, file=log)
                print(f"FAIL {name}: {log_path}", flush=True)
                return False
        print(f"PASS {name}", flush=True)
        return True

    with ThreadPoolExecutor(max_workers=args.jobs) as pool:
        results = list(pool.map(check, cases))
    passed = sum(results)
    print(f"{passed}/{len(results)} cases passed. Reports: {output}", flush=True)
    return 0 if all(results) else 1


if __name__ == "__main__":
    try:
        sys.exit(main())
    except (OSError, ValueError, subprocess.CalledProcessError) as error:
        print(f"error: {error}", file=sys.stderr)
        sys.exit(1)
