# scaffold-go-api

[![CI](https://git.higherlearning.eu/alrayyes/scaffold-go-api/actions/workflows/ci.yml/badge.svg)](https://git.higherlearning.eu/alrayyes/scaffold-go-api/actions)
[![release](https://img.shields.io/gitea/v/release/alrayyes/scaffold-go-api?gitea_url=https%3A%2F%2Fgit.higherlearning.eu&sort=semver)](https://git.higherlearning.eu/alrayyes/scaffold-go-api/releases)
[![licence](https://img.shields.io/badge/licence-unlicensed-lightgrey)](LICENSE)

A Forgejo template for a Go backend API. Generate a project from it
(Settings → repository → "Generate a new repo" on `scaffold-go-api`, or
`POST /repos/alrayyes/scaffold-go-api/generate` against the API) and you get
a project with the conventions already wired in — pinned tooling, a
spec-first OpenAPI layout, prose linting, secret scanning, and release
automation — rather than a blank directory and a checklist to work
through by hand.

It isn't a product on its own. The one endpoint it ships,
`GET /widgets/{id}`, exists so the whole chain — spec, handler, tests,
hooks, CI — has something real to run against. Replace it with your
first real resource and delete this paragraph.

Hosted here on Forgejo: semantic-release, Renovate, `.forgejo/workflows/`.
See [FORGEJO.md](FORGEJO.md) for the reasoning and the GitHub-hosted
alternative for each, if a project stamped from this template ever needs
to move there instead.

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
git clone https://git.higherlearning.eu/alrayyes/scaffold-go-api.git
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

### Docker

The [Dockerfile](Dockerfile) builds a static binary into a distroless,
non-root image. This repo ships the image but doesn't run it, so the
runtime hardening `rules/containers.md` asks for lives here rather than in
a compose file — verified locally, and the binary needs neither a
capability nor a path to write to:

```sh
docker build -t scaffold-go-api .
docker run --rm -p 8080:8080 \
  --cap-drop=ALL --security-opt=no-new-privileges --read-only \
  --memory=64m --cpus=0.5 \
  scaffold-go-api
curl localhost:8080/healthz
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for the toolchain, the hooks, and
how a change gets reviewed and released.

## Licence

No licence has been chosen yet — see [`LICENSE`](LICENSE). Pick one before
a project stamped from this template goes anywhere public.
