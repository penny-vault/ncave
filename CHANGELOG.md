# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.2] - 2026-05-26

### Changed
- Upgrade pvbt dependency to v0.10.2
- Update other dependencies to latest

## [0.2.1] - 2026-05-05

### Changed
- Upgrade pvbt dependency to v0.9.3

## [0.2.0] - 2026-05-04

### Changed
- Default benchmark from VFINX to SPY
- Upgrade pvbt dependency to v0.9.2

### Removed
- Snapshot-based unit tests; pvbt 0.9.x captures the universe-wide metrics view, producing snapshots that exceed GitHub's per-file size limit for ncave's universe-screening strategy

## [0.1.4] - 2026-05-01

### Changed
- Upgrade pvbt dependency to v0.8.1

## [0.1.3] - 2026-04-25

### Changed
- Upgrade pvbt dependency to v0.8.0
- Regenerate testdata snapshot for pvbt's v5 snapshot schema

## [0.1.2] - 2026-04-23

### Changed
- Upgrade pvbt dependency to v0.7.7

## [0.1.1] - 2026-04-21

### Fixed
- Remove local pvbt replace directive so the module resolves correctly outside the monorepo



## [0.1.0] - 2026-04-21

### Added
- Initial release of Net Current Asset Value Effect strategy
- Implements NCAV (Net Current Asset Value) fundamental screen with configurable universe parameter and regime-change asset overlay
- Uses FetchFundamentalsByDateKey for accurate March 31 formation date fundamentals

[0.1.0]: https://github.com/penny-vault/ncave/releases/tag/v0.1.0

[0.1.1]: https://github.com/penny-vault/ncave/compare/v0.1.0...v0.1.1
[0.1.2]: https://github.com/penny-vault/ncave/compare/v0.1.1...v0.1.2
[0.1.3]: https://github.com/penny-vault/ncave/compare/v0.1.2...v0.1.3
[0.1.4]: https://github.com/penny-vault/ncave/compare/v0.1.3...v0.1.4
[0.2.0]: https://github.com/penny-vault/ncave/compare/v0.1.4...v0.2.0
[0.2.1]: https://github.com/penny-vault/ncave/compare/v0.2.0...v0.2.1
