# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Professional documentation suite: `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `SECURITY.md`, and `CHANGELOG.md`.
- Comprehensive Go Doc Comments for all exported items.
- CI/CD pipeline with Go 1.24, linting (golangci-lint), and security scanning (gosec).
- Standardized `Makefile` for project lifecycle management.
- `.cspell.json` to manage technical spelling across the repository.

### Changed
- Refined `internal/storage` tests to support Write-Ahead Log (WAL) dependency.
- Updated `Dockerfile` to pin Alpine version (3.21) for reproducibility.
- Translated `README.md` to English and added professional roadmap.

## [0.1.0] - 2026-03-27

### Added
- **Phase 1**: In-memory KV Store with thread-safety (`sync.RWMutex`).
- **Phase 2**: gRPC Server and Client implementation.
- **Phase 3**: Cluster synchronization and discovery logic.
- **Phase 4**: Leader election (Raft-like mechanism).
- **Phase 5**: Persistence via Write-Ahead Log (WAL).
- Basic project structure and Docker support.
