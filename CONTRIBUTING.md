# Contributing

This file is for whoever changes this template. The
[README](README.md) is for whoever stamps a project out of it.

## Getting set up

- **Go 1.25 or newer.**
- **[bun](https://bun.sh)** for the tooling that isn't Go — commitlint,
  Prettier, markdownlint, [Redocly](https://redocly.com/docs/cli), and the
  [lefthook](https://lefthook.dev) that runs the git hooks. There's a
  `package.json`, but nothing here is JavaScript; it exists only so those
  tools resolve and stay pinned.
- **[golangci-lint](https://golangci-lint.run) v2.12.2**, which the
  pre-commit hook runs from your `PATH` while CI runs it pinned. Install
  that version rather than whichever is current: when the two disagree, the
  hook passes and the pipeline fails, and the reason isn't obvious from the
  failure.
- **[Vale](https://vale.sh)** on your `PATH`, for the style tier of the
  prose lint:

  ```sh
  go install github.com/errata-ai/vale/v3/cmd/vale@latest
  ```

  `ltex-cli-plus` needs nothing installed: the hook fetches and caches it
  on first use.

One command installs the linters and the git hooks:

```sh
bun install
```

An uninstalled hook silently does nothing, which is worse than not having
one, so the `prepare` script runs `lefthook install` for you. You find out
at the pipeline otherwise, not at the commit.

## Everyday commands

Every one of these is what a hook or CI runs — see `lefthook.yml` and
`.github/workflows/*.yml` for exactly which.

```sh
go build ./...
go vet ./...
go test ./...
golangci-lint run
golangci-lint fmt          # the fixer; `run` stays the check
go run golang.org/x/vuln/cmd/govulncheck@v1.7.0 ./...

bun run format:check       # prettier --check, add --write to fix
bun run lint:md
bun run lint:api           # redocly lint, bare — no path argument
bun run lint:prose         # vale
bun run lint:mechanics     # ltex-cli-plus
```

## How it fits together

`internal/api` holds the handler, and `cmd/scaffold-go-api/main.go` is the
composition root that wires it up and starts the server — per `go.md`'s "a
server keeps everything in `internal/` and its commands in `cmd/`", since
there's nothing here worth exporting. No finer `internal/domain`/
`internal/adapter` split on top of that: that shape earns its keep the day a
second resource needs it, not on day one of a template nobody's used yet.

## The contract

`api/openapi.yaml` describes the API and is handwritten, not generated
from the handlers — see the header comment for why. `redocly lint` checks
the document is valid OpenAPI; nothing yet checks the handler still
matches it, the way `form-handler`'s `openapi_test.go` does for a real
service. Add that test the day this scaffold's example resource stops
being an example.

## Commit messages

[Conventional Commits](https://www.conventionalcommits.org/):
`type(scope): description`, types `feat`/`fix`/`docs`/`style`/`refactor`/
`perf`/`test`/`build`/`ci`/`chore`/`revert`. Subject under 50 characters,
lowercase, no trailing full stop. commitlint enforces the shape at
commit-msg and again in CI; the length and case rules are tighter than
what it checks, so hold to them anyway.

## Branching, review, and release

Every change goes through a pull request — nothing is pushed straight to
`main`, including the bootstrapping that built this repo. GitHub's branch
protection needs a paid plan this account doesn't have, so nothing enforces
that mechanically here; it's discipline, not a gate.

Once a pull request's checks are green, squash-merge it and delete the
branch. [release-please](https://github.com/googleapis/release-please)
reads the Conventional Commits on `main` and keeps a release pull request
open with the next version and changelog entry; merging that one tags the
release, and [goreleaser](https://goreleaser.com) builds the binaries onto
it. Nobody picks a version by hand.
