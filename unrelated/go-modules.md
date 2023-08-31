# Go Modules

## 1. Dependency Management

Make sure other developers can build your code using similar version

Go uses SIV (Semantic Import Versioning) which is based on Semantic Versioning

Highlights:

- Composed of 3 numbers: `v1.12.4`
- Major, minor, patch
- Major bumps == breaking change. Manual upgrade is necessary.
- Other bumps can be automated by tooling if needed.

Examples:

- Install `v1.12.4`, other devs may end up building with `v1.13.0` if it releases
- If releasing a new version from `v1.13.6` with breaking change, next version will prob be `v2.0.0`
- If `v2.0.0` releases our tooling will NOT upgrade to that automatically

Special case:

- v0 allows for any breaking changes. It is a special version
- `go get github.com/YourCompany/somelib` will NOT get the most recent version.
  - We need to run `go get github.com/YourCompany/somelib/v4` or similar.
  - Tooling around this isn't always perfect. Will hopefully improve over time.

## 2. Working outside of `GOPATH`

All Go code used to live inside a single directory on your computer - the `GOPATH`.

- See <semver.org>

Go Modules allow us to run code from anywhere, as long as we initialize a module.

## Setting up our module

```bash
go mod init github.com/YourCompany/somelib
#                        ^^
#                     use your own github handle!
```

This creates a `go.mod` file.
