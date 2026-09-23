# Refreshing Jev JSON fixtures

CI uses recorded responses under `testdata/jev/` so unit tests do not call the live API.

To update a fixture after a pinned model change:

1. Set `TYPESAFE_API_KEY` and a **pinned** `JEV_MODEL` (not `jev-latest`).
2. Run `go test -tags=integration -count=1 ./...` and capture the response, or use your host’s logging of `answers_snapshot` from a staging shadow run.
3. Redact any fields that are not part of the public `jev.Response` schema.
4. Write the JSON to `testdata/jev/<scenario>.json` and run `go test ./...`.

Optional nightly or manual workflow: [`.github/workflows/integration.yml`](../.github/workflows/integration.yml).
