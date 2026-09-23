# Contributing

Thank you for improving notify-jev-router.

## Development setup

- Go 1.22+
- From the repository root: `go test ./...` (no API keys required)

Optional live Jev tests:

```bash
export TYPESAFE_API_KEY=...
export JEV_MODEL=jev-x.y.z   # pinned, not jev-latest
go test -tags=integration -count=1 ./...
```

## Pull requests

1. Fork and create a feature branch from `main`.
2. Keep changes focused; update docs when behavior or public API changes.
3. Ensure `go test ./...` passes.
4. Describe shadow or integration impact in the PR body when relevant.

## PII and reviews

- Do not add fields that carry contact information or free-text user content to `RoutingContext`.
- Reject or redact PII in examples, fixtures, and logs.
- Reviewers should block changes that weaken validation or log sensitive state.

## Dependency licenses

Runtime dependency: [`github.com/havlan/jev-go`](https://github.com/havlan/jev-go). Confirm its license remains compatible before upgrading. See [docs/jev-fixtures.md](docs/jev-fixtures.md) for fixture refresh.

## Code of conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
