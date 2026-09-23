# Release v0.1.0 (maintainer)

After merging the Issue #3 PR to `main`:

1. Confirm CI is green on `main`.
2. Create annotated tag on `main` HEAD:
   ```bash
   git tag -a v0.1.0 -m "v0.1.0 — Issue #1 MVP, Issue #3 OSS gates"
   git push origin v0.1.0
   ```
3. Set repository visibility to **Public** (Settings → General).
4. Create GitHub Release `v0.1.0` from the tag with notes: MVP routing core, shadow mode helpers, no breaking API changes from private pins.
