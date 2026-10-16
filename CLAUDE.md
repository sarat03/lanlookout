# Working rules

## Versioning (every release)
- One source of truth: `backend/internal/web/public/version` (`VERSION=X.Y.Z`). The About page reads it; nothing else holds the version.
- Any PR that changes `backend/`, `frontend/src/` or `Dockerfile` bumps it. Patch for fixes, minor for features, major for breaking changes. CI (`version-check.yml`) fails the PR otherwise.
- Same PR: add a `## [vX.Y.Z] - YYYY-MM-DD` entry at the top of `CHANGELOG.md`.
- After a frontend change, run `make all` in `frontend/` and commit the rebuilt `backend/internal/web/public/assets/`.
- Don't tag or bump Umbrel by hand: merging to main makes CI tag `vX.Y.Z`, create the GitHub release, tag the image `:X.Y.Z` and open the Umbrel PR with the same version.
