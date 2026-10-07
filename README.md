# Project generator

---

## Creates a simple project structure with an example that can be run immediately

* First step is create a build (`make build`)
* Second step (`./project-gen {app_name}`)

Or in one step: `make gen NAME={app_name}`

#### Flags

* `-go {version}` — Go version written to `go.mod` and the Dockerfile base image.
  Defaults to the toolchain that built the generator, so pass it only when the
  generated project must build on an older one: `./project-gen -go 1.24 myapp`.
* `-module {path}` — module path for `go.mod` and every import.
  Defaults to `github.com/khusainnov/{app_name}`:
  `./project-gen -module example.com/team/billing billing`

Flags go before the project name.

#### For running generated project you need to use `deployment/local.env`

```sh
cd {app_name}
go mod tidy
set -a && . ./deployment/local.env && set +a
go run .
curl -s -X POST http://localhost:5050/jsonrpc/v2 \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","method":"Handler.Echo","params":[{"message":"hello"}],"id":1}'
```

`go mod tidy` is required on a freshly generated project: the generator writes
`go.mod` but no `go.sum`.
