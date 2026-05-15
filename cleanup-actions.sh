gh repo set-default kholia/chisel-wolfi

gh run list --limit 1000 --json databaseId --jq '.[].databaseId' | xargs -I{} gh run delete {}
