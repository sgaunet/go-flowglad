# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial SDK release
- `flowglad` package: full Flowglad REST API client with 17 resource groups
- `webhook` package: HMAC-SHA256 webhook signature verification
- `flowgladtest` package: fake HTTP server for consumer testing
- Cursor-based pagination via Go 1.23 `iter.Seq2` iterators
- Configurable retry with exponential back-off and jitter
- Typed `*Error` with HTTP status code, error code, message, and validation issues
- Pluggable `http.RoundTripper` via `WithHTTPClient` option
- Structured logging via `log/slog` with `WithLogger` option
