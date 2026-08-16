# scaffold-go-api

[![licence](https://img.shields.io/badge/licence-unlicensed-lightgrey)](LICENSE)

A GitHub template for a Go backend API. Run `gh repo create
my-real-project --template alrayyes/scaffold-go-api` and you get a
project with the conventions already wired in — pinned tooling, a
spec-first OpenAPI layout, prose linting, secret scanning, and release
automation — rather than a blank directory and a checklist to work
through by hand.

It isn't a product on its own. The one endpoint it ships,
`GET /widgets/{id}`, exists so the whole chain — spec, handler, tests,
hooks, CI — has something real to run against. Replace it with your
first real resource and delete this paragraph.

## Requirements

- **Go 1.25 or newer.**
- **[bun](https://bun.sh)**, for the tooling that isn't Go — commitlint,
  Prettier, markdownlint, [Redocly](https://redocly.com/docs/cli), and the
  [lefthook](https://lefthook.dev) that runs the git hooks. There's a
  `package.json`, but nothing here is JavaScript; it exists only so those
  tools resolve and stay pinned.
- **[golangci-lint](https://golangci-lint.run)**, pinned in
  [CONTRIBUTING.md](CONTRIBUTING.md#getting-set-up).
- No external services. The example resource is an in-memory map.

## Installation

```sh
git clone https://github.com/alrayyes/scaffold-go-api.git
cd scaffold-go-api
go build ./cmd/scaffold-go-api
```

## Usage

```sh
./scaffold-go-api
```

Listens on `:8080` by default; set `ADDR` to change it.

```sh
curl localhost:8080/healthz
curl localhost:8080/widgets/hammer
```

`api/openapi.yaml` is the contract both endpoints are held to — read it
first if you're replacing the example resource with a real one.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for the toolchain, the hooks, and
how a change gets reviewed and released.

## Licence

No licence has been chosen yet — see [`LICENSE`](LICENSE). Pick one before
a project stamped from this template goes anywhere public.
