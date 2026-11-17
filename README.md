# scaffold-go-api

[![CI](https://github.com/alrayyes/scaffold-go-api/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/alrayyes/scaffold-go-api/actions/workflows/ci.yml)
[![Codecov](https://codecov.io/gh/alrayyes/scaffold-go-api/graph/badge.svg)](https://codecov.io/gh/alrayyes/scaffold-go-api)
[![release](https://img.shields.io/github/v/release/alrayyes/scaffold-go-api?sort=semver)](https://github.com/alrayyes/scaffold-go-api/releases/latest)
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
first if you're replacing the example resource with a real one. A
Scalar docs page renders that same spec and publishes to
GitHub Pages on every merge to `main`:
<https://alrayyes.github.io/scaffold-go-api/>.

### Reports

Every merge to `main` publishes the test and coverage reports beside the
docs, in the same Pages deployment:

- Tests:
  <https://apis.ryankes.eu/scaffold-go-api/reports/tests/unit.xml> (JUnit XML)
- Coverage: <https://apis.ryankes.eu/scaffold-go-api/reports/coverage/> (HTML),
  with `coverage.xml` (Cobertura) and `coverage.out` (Go profile) beside it

A project generated from this template gets the same layout under its own
Pages address. A `<owner>.github.io` repo serves from the root, so its
reports land at `/reports/` with no repo-name prefix, and the workflow's
smoke step uses the deployed `page_url`, so it holds for either. Pull
requests run the assembly but deploy nothing.

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

## Actions event policy

The workflows here use `pull_request_target`, which GitHub's default Actions
policy blocks on public repos from 2026-12-02. A new public repo made from this
template needs a per-repo policy allowing the events its workflows use, or
those workflows end in `startup_failure`:

```sh
gh api -X POST repos/<owner>/<repo>/actions/policies --input - <<'JSON'
{"name":"allow-pull_request_target","enforcement":"active",
 "rules":[{"type":"restrict_action_events",
           "parameters":{"allowed_events":["pull_request_target","pull_request","push"]}}]}
JSON
```

`allowed_events` is an allow-list: add every other event a workflow uses
(`schedule`, `workflow_dispatch`, `workflow_run`, `release` and so on), or it
stops firing. A later change to an existing policy is made in the UI under
Settings > Actions > Policies. See the
[GitHub docs](https://docs.github.com/en/actions/reference/security/securely-using-pull_request_target).

## Licence

No licence has been chosen yet — see [`LICENSE`](LICENSE). Pick one before
a project stamped from this template goes anywhere public.
