# scaffold-go-api

A GitHub template repo, not a running service. It's built from
`~/.config/claude/CLAUDE.md` and `~/.config/claude/rules/*.md` — read those
for the "why" behind everything below. This file only says what's specific
to this repo.

## What this is

Bootstrapped in the PR sequence documented in
`~/.config/claude/plans/adaptive-conjuring-karp.md`: chassis, prose tooling,
docs, spec, handler, CI, prose/secret CI, release automation, Dependabot —
each in its own PR. Keep new work in that shape: one concern per PR.

## Commands

```sh
go build ./... && go vet ./... && go test ./...
golangci-lint run                  # golangci-lint fmt is the fixer
bun run format:check               # bun run lint:md, lint:api, lint:prose, lint:mechanics too
```

Full list and what each one does: [CONTRIBUTING.md](CONTRIBUTING.md).

## Gotchas

- **No branch protection.** GitHub requires a paid plan for it on a private
  repo, which this account doesn't have. The PR-only discipline is
  enforced by nobody but whoever's committing — never push straight to
  `main`.
- **`internal/api` holds the handler, `cmd/scaffold-go-api` starts it.**
  Per `go.md`'s "a server keeps everything in `internal/` and its commands
  in `cmd/`" — there's nothing here worth exporting, since the API this
  service offers is its endpoints, not its Go packages. Don't pre-build
  `internal/domain`/`internal/adapter` inside that on top of it — that finer
  split earns its place the day a second resource needs it, not on day one.
- **The example resource is a placeholder.** `GET /widgets/{id}` exists so
  the spec-first chain (spec, handler, tests, CI) has something real to
  run end to end. A project stamped from this template replaces it with
  its first real endpoint — spec first, per `rules/api.md`.
- **`LICENSE` is deliberately unpicked.** Don't default it to GPL-3.0 or
  anything else; that's a decision the project stamped from this template
  makes, not this template.
- **Renovate can't reach this repo.** It's GitHub-primary; Dependabot
  (`.github/dependabot.yml`) is what raises dependency pull requests here.
