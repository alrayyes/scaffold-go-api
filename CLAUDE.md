# scaffold-go-api

A Forgejo template repo, not a running service. It's built from
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

- **Branch protection isn't turned on**, though Forgejo supports it for
  free on a private repo (unlike GitHub). If it gets enabled, remember the
  "Enable Push" toggle is a separate step from the push whitelist under it
  — missing either one silently breaks the release bot's ability to push.
  Until then, the PR-only discipline is enforced by nobody but whoever's
  committing — never push straight to `main`.
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
- **Renovate raises the dependency pull requests here**, not Dependabot —
  it's a GitHub-only feature, and this repo is hosted on Forgejo instead.
  `renovatebot` is a repo collaborator; `renovate.json` only carries the
  overrides on top of whatever shared config the bot already runs with.
