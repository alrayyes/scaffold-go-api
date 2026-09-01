# Forgejo alternatives

This template defaults to GitHub-primary tooling: `release-please`,
Dependabot, `.github/workflows/`, and a couple of smaller pieces below. That's
the right default for a repo that lives on `github.com`. If a project stamped
from this template ends up hosted on a Forgejo instance instead, here is what
needs swapping, with working configuration rather than hand-waving — grounded
in two real, already-migrated siblings: `scaffold-go-cli` (same language,
same release shape, minus the OpenAPI spec) and `scaffold-php-api` (a
different language, but the same OpenAPI-spec-first API shape this repo has).

None of this applies if the repo stays on GitHub. Skip the whole file.

## What ships GitHub-specific, and what replaces it

| GitHub-specific piece                                                                  | Forgejo replacement                                                                                                                                                                                                        |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `release-please` + `release-please-config.json` + `.release-please-manifest.json`      | `semantic-release` + `release.config.mjs` (already committed — see §1)                                                                                                                                                     |
| `.github/dependabot.yml`                                                               | `renovate.json` (see §2)                                                                                                                                                                                                   |
| `.github/workflows/dependabot-auto-merge.yml`                                          | Nothing — Renovate has `automerge: true` built in (see §2)                                                                                                                                                                 |
| `.github/workflows/*.yml`                                                              | `.forgejo/workflows/*.yml` (see §3)                                                                                                                                                                                        |
| `.github/workflows/pr-title.yml` (`amannn/action-semantic-pull-request`)               | A commitlint-on-title job — no Forgejo-hosted mirror of that action exists (see §3)                                                                                                                                        |
| `.github/workflows/release.yml`'s `attest-build-provenance` step                       | Nothing — see §1's last paragraph                                                                                                                                                                                          |
| `.github/PULL_REQUEST_TEMPLATE.md`, `.github/ISSUE_TEMPLATE/*.yml`                     | `.forgejo/PULL_REQUEST_TEMPLATE.md`, `.forgejo/ISSUE_TEMPLATE/*.yml` (same shape, new directory — confirmed against `scaffold-php-api`'s own layout)                                                                       |
| `redocly lint --format=github-actions` (inside `.github/workflows/ci.yml`'s `api` job) | Bare `redocly lint`, no `--format` flag (see §4)                                                                                                                                                                           |
| No image registry wired up today                                                       | Forgejo's own container registry, if a project stamped from this template adds one (see §5)                                                                                                                                |
| README badges, `git clone https://github.com/...`, `gh repo create --template`         | Point at `git.higherlearning.eu` instead; `scaffold-go-cli`'s own README is the worked example — "Generate from it (Forgejo's web UI: this repo's page → 'Use this template')" replaces the `gh repo create` line entirely |

## 1. Release automation: `release-please` → `semantic-release`

**Why:** `release-please` is GitHub-only — there's no Forgejo port. Anything
on a Forgejo instance wants `semantic-release` in front of `goreleaser`
instead. The split doesn't change: one tool reads the Conventional Commits
and decides the version, the other cross-compiles and attaches the binaries.
Only the version-deciding tool is different.

**What stays exactly the same:** `.goreleaser.yml` — the Linux-only build
(`amd64`/`arm64`, this is a deployed service, not something that runs on
someone's own machine), `CGO_ENABLED=0`, the `ldflags` that stamp
`main.version` from the tag, the tarball archive, `changelog.disable: true`
(release-please/semantic-release already wrote the notes; a second shortlog
from goreleaser would either duplicate them or, with `release.mode: replace`,
throw the curated copy away). The only thing that needs adding to
`.goreleaser.yml` is a `release.gitea:` block — see the note at the end of
this section.

### `release.config.mjs` — already committed, this is the step

This file already exists at the repo root. Its own header comment explains
the swap in detail (`@ribbon-studios/semantic-release-forgejo` replacing
`release-please`, why there's no `@semantic-release/exec` step — a Go binary
has no manifest to bump; `goreleaser`'s `ldflags` read the version straight
off the tag `semantic-release` pushes). **Swapping it in is one step: delete
`release-please-config.json` and `.release-please-manifest.json`
(semantic-release keeps no manifest of its own — it derives the last release
from `git tag`), then point `.forgejo/workflows/release.yml` at
`release.config.mjs`** the way the workflow below does. Nothing in the file
itself needs editing to go live.

What it still needs, because none of this repo's current `package.json`
carries it — `release-please` needed no npm packages at all, only the GitHub
Action — is the `semantic-release` toolchain as `devDependencies`. Pinned
versions from `scaffold-go-cli`'s own working installation (exact versions,
per this house's rule of pinning everything):

```json
{
  "devDependencies": {
    "@ribbon-studios/semantic-release-forgejo": "0.1.3",
    "@semantic-release/changelog": "7.0.0",
    "@semantic-release/commit-analyzer": "13.0.1",
    "@semantic-release/git": "11.0.1",
    "@semantic-release/release-notes-generator": "14.1.1",
    "conventional-changelog-conventionalcommits": "10.4.0",
    "semantic-release": "25.0.9"
  },
  "overrides": {
    "conventional-changelog-writer": "9.2.1"
  }
}
```

The `overrides` entry pins a transitive dependency that otherwise floats —
carry it across too, not just the top-level packages.

### Picking `tagFormat`

`semantic-release` defaults to `tagFormat: "v${version}"`. Get this wrong and
it can't find the repo's own last release, concludes nothing has ever
shipped, and republishes `1.0.0` over whatever's already tagged.

This repo already has real tags: `git tag --list` on it shows `v0.2.0` and
`v0.3.0`, both cut by `release-please` before this migration
(`release-please-config.json` sets `"include-v-in-tag": true`), so the
default `tagFormat` (`v${version}`) matches and needs no explicit setting.
`release.config.mjs`'s own header comment predates those tags and still says
"no git tags exist yet" — stale now, not a reason to doubt the default;
`git tag --list` is the source of truth, not the comment. **If you're
stamping a fresh project from this template rather than migrating this repo
itself, check `git tag --list` on the actual project before assuming that
holds** — a project with no tags yet, or with bare (non-`v`-prefixed) tags
from being on `semantic-release` somewhere else already, needs
`tagFormat: "${version}"` instead. There's no
correct default independent of what's already there.

### `.forgejo/workflows/release.yml`

Adapted from `scaffold-go-cli`'s own working workflow, trimmed to this
repo's Linux-only build (no `windows`/`darwin` targets, no `.zip`
`format_overrides`) and its build-provenance gap (see the note after this
example).

```yaml
name: release

on:
  push:
    branches: [main]

concurrency:
  group: release
  cancel-in-progress: false

jobs:
  release:
    runs-on: docker
    container:
      # semantic-release has to run under real `node`; bun 1.3/1.4 report
      # themselves as node v24.3.0, which lands inside neither range
      # semantic-release accepts, and it checks at startup — a tool that
      # shells out and imports every plugin at runtime is not the one to
      # find a Node-compatibility gap on, on the single job that pushes to
      # main.
      image: oven/bun:1.4.0-alpine@sha256:07235578f79ef8c6f97d94aee7938e76f5cdba5f21ae5dbfdd3d3d38058437eb
    steps:
      # git and ca-certificates together: installing git flips checkout from
      # a node-fetched tarball to an HTTPS clone, and git reads the system CA
      # store, which this image hasn't got. nodejs is what actually runs
      # semantic-release.
      - run: apk add --no-cache ca-certificates git nodejs

      # fetch-depth: 0 is not optional — semantic-release walks back to the
      # last tag to work out what to release, and a shallow clone has none.
      - uses: https://code.forgejo.org/actions/checkout@v4
        with:
          fetch-depth: 0

      - run: bun install --frozen-lockfile

      - name: Release
        env:
          # How semantic-release itself pushes the tag and the changelog
          # commit — the generic "<username>:<password>" form, since it has
          # no Forgejo-specific auth of its own (GH_TOKEN and GL_TOKEN are
          # the ones it knows). Forgejo takes an access token as the
          # password half.
          GIT_CREDENTIALS: "${{ secrets.RELEASE_USER }}:${{ secrets.RELEASE_TOKEN }}"
          # release.config.mjs's own credential — forgejoToken, not left to
          # the environment. See its header comment for why: the runner
          # injects its own FORGEJO_TOKEN into every job automatically, and
          # leaving the plugin to fall back to the environment means the
          # release is authenticated by whichever of the two happens to be
          # in scope.
          RELEASE_TOKEN: ${{ secrets.RELEASE_TOKEN }}
        id: release
        run: |
          bunx semantic-release
          # No release_created-style output the way release-please-action
          # gives — cheapest reliable check: HEAD sits exactly on the tag
          # semantic-release just pushed if, and only if, it actually
          # released.
          tag=$(git describe --tags --exact-match 2>/dev/null || true)
          echo "tag=$tag" >> "$GITHUB_OUTPUT"
    outputs:
      tag: ${{ steps.release.outputs.tag }}

  artefacts:
    needs: release
    if: needs.release.outputs.tag != ''
    runs-on: docker
    container:
      image: golang:1.27-alpine@sha256:REPLACE_WITH_PINNED_DIGEST
    steps:
      - run: apk add --no-cache ca-certificates git nodejs
      - uses: https://code.forgejo.org/actions/checkout@v4
        with:
          ref: ${{ needs.release.outputs.tag }}
          fetch-depth: 0
      - name: goreleaser
        env:
          # goreleaser's Gitea/Forgejo release target reads GITEA_TOKEN and
          # GITEA_SERVER_URL — names goreleaser picked, not this house's, so
          # they don't collide with Forgejo's own reserved-prefix rule (§3).
          GITEA_TOKEN: ${{ secrets.RELEASE_TOKEN }}
          GITEA_SERVER_URL: ${{ vars.FORGEJO_URL }}
        run: |
          go install github.com/goreleaser/goreleaser/v2@v2.17.1
          goreleaser release --clean
```

`RELEASE_USER` and `RELEASE_TOKEN` are repository secrets, created under
Settings → Actions → Secrets: `RELEASE_TOKEN` a Forgejo access token with
read/write on the repository, `RELEASE_USER` the account it belongs to. Both
halves are needed because git wants a username and a password to push over
HTTPS, and an access token is only the password half.

`.goreleaser.yml` also needs a `release:` block added, since unlike the
GitHub target, `goreleaser` can't infer owner/name from a
`GITHUB_REPOSITORY`-style variable on Forgejo:

```yaml
release:
  gitea:
    owner: <account>
    name: <repo>
```

**The build-provenance step has no Forgejo equivalent, and this doc knows of
none.** `.github/workflows/release.yml`'s `attest-build-provenance` step
mints a Sigstore OIDC token through GitHub's own identity provider and
publishes the attestation to GitHub's attestations API — both GitHub-native,
with nothing on Forgejo standing in for either. The preceding workflow simply
drops it. Self-hosting `cosign`/Sigstore signing as a plain step is possible
if a project genuinely needs provenance, but there's no forge-native place to
publish the attestation to the way there is on GitHub.

## 2. Dependency updates: Dependabot → Renovate

**Why:** Dependabot is GitHub-native and doesn't run anywhere else. Renovate
is the Forgejo-side equivalent, but it needs to be told that it is talking to
a Forgejo (Gitea-compatible) API: `"platform": "gitea"` in its config, not
`"github"` and not a bare default.

### This repo already has a `renovate.json` — check before adding to it

A `renovate.json` sits at this repo's root already, left over from an earlier
attempt at this exact migration (see `git log -- renovate.json`). It carries
the commit-prefix package rules — `fix:` for `gomod`/`bun` (they ship in the
built binary; a bump is a change to what ships, same as any other fix),
`ci:` for `github-actions` — mirroring `.github/dependabot.yml`'s own
`commit-message.prefix` groups. Those rules are genuinely this repo's own
convention, not something a shared bot's autodiscovered config would know,
so they're worth keeping regardless of which of the two setups below applies.
The `matchManagers: ["github-actions"]` rule does not need renaming for
`.forgejo/workflows/` — Renovate's `github-actions` manager already matches
Forgejo/Gitea Actions workflow paths under that same manager name; there is
no separate `forgejo-actions` manager to switch to.

### If the Forgejo instance already runs a shared Renovate bot

This is the actual situation on `git.higherlearning.eu`: both `scaffold-
go-cli` and `scaffold-php-api`'s own `renovate.json` files carry next to
nothing — `scaffold-php-api`'s is just the `$schema` line, `scaffold-go-cli`'s
adds only a `customManagers` entry for pinning tool versions inside its own
`.forgejo/workflows/*.yml`. Neither sets `"platform"` nor `"autodiscover"`
itself — the instance's shared bot already carries that policy
(`config:best-practices`, semantic commits, vulnerability alerts, automerge
rules) across every repository it discovers, and a per-repo `renovate.json`
only needs to state what's genuinely repo-specific on top of that. **If this
holds for wherever this template lands, keep this repo's `renovate.json` down
to the commit-prefix rules given earlier — nothing else is needed.**

### Standalone `renovate.json`, if there's no shared bot

If the target instance has no shared bot, the file needs the platform and a
base policy stated explicitly:

```json
{
  "$schema": "https://docs.renovatebot.com/renovate-schema.json",
  "platform": "gitea",
  "extends": [
    "config:best-practices",
    ":semanticCommits",
    ":enableVulnerabilityAlertsWithLabel('security')"
  ],
  "packageRules": [
    {
      "matchManagers": ["gomod", "bun"],
      "commitMessagePrefix": "fix:"
    },
    {
      "matchManagers": ["github-actions"],
      "commitMessagePrefix": "ci:"
    },
    {
      "matchUpdateTypes": ["minor", "patch", "pin", "digest"],
      "automerge": true
    }
  ]
}
```

Either way: delete `.github/dependabot.yml` and
`.github/workflows/dependabot-auto-merge.yml` once Renovate is doing the job.
Renovate has `automerge: true` built in; Dependabot doesn't, which is the
entire reason `dependabot-auto-merge.yml` exists on the GitHub side — nothing
on Forgejo needs to replace it, it just goes away.

## 3. CI workflows: `.github/workflows/*.yml` → `.forgejo/workflows/*.yml`

Forgejo Actions is a GitHub-Actions-compatible runner: the YAML shape,
`jobs:`/`steps:`, and most context variables (`github.event_name`,
`github.ref`, `$GITHUB_STEP_SUMMARY`, `$GITHUB_OUTPUT`, `$GITHUB_ENV`) all
carry over unchanged. What actually needs rewriting, concretely:

- **`uses:` needs to name an action the runner can actually reach.** GitHub's
  shorthand (`actions/checkout@<sha>`) resolves against `github.com`. A
  Forgejo runner needs a full URL to a mirrored action —
  `https://code.forgejo.org/actions/checkout@v4` for checkout,
  `https://code.forgejo.org/actions/setup-go@<sha>` for Go,
  `https://code.forgejo.org/oven-sh/setup-bun@<sha>` for bun — or the step
  rewritten as plain shell where no mirror exists. Confirmed with no mirror:
  `golangci-lint-action` (404 on `code.forgejo.org/golangci/golangci-lint-
action`) and `goreleaser-action` — both fall back to `go install
<module>@<version>` in a `run:` step instead, the same way `gitleaks.yml`
  already installs `gitleaks` with `go install` rather than an action.
  `attest-build-provenance` has no fallback at all — see §1's last
  paragraph.

- **Secret and variable names can't start with `FORGEJO_`, `GITHUB_`, or
  `GITEA_`** — the same way GitHub reserves `GITHUB_`, Forgejo reserves all
  three. `secrets.FORGEJO_TOKEN` still _resolves_ — the run's own
  automatic token is published under that name — which is the trap: it looks
  like a secret is being read when it's actually the automatic one. This is
  why `release.config.mjs` and the preceding workflow name the release credential
  `RELEASE_TOKEN`, not anything `FORGEJO_`-prefixed.

- **Bare containers need `git`, `ca-certificates`, and `nodejs` installed
  before `actions/checkout` runs, not after.** `actions/checkout` is a
  JavaScript action — a container with no `node` can't execute it at all
  (`exec: "node": executable file not found`), and a container with no `git`
  doesn't fail: checkout silently falls back to an API-fetched tarball with
  no `.git` directory, which leaves `commits`, `gitleaks`, and anything else
  that reads history with nothing to read. `nodejs` reads the CA store too,
  which a bare Alpine image hasn't got — that's what `ca-certificates` is
  for. Order in every job that uses a non-default `container:` image: install
  those three first, then `uses: .../checkout`.

- **Pin container images by digest**, same as any other dependency:
  `image: hadolint/hadolint:v2.15.1-alpine@sha256:...`, not a bare tag.

- **`docker build` needs pointing at the runner's own dind daemon.** A job
  running under `runs-on: docker` has no `/var/run/docker.sock` to bind-mount
  — the runner is itself docker-in-docker, and a sibling `dind` daemon starts
  alongside every job container on a private network of its own. That
  container's own default gateway _is_ the daemon; read it at run time from
  `/proc/net/route`, don't hardcode it, and point `DOCKER_HOST` at
  `tcp://<gateway>:2375`. `scaffold-go-cli`'s `ci.yml` `docker-build` job has
  the worked shell for this, including a `curl .../_ping` sanity check before
  trusting it.

Job-by-job for this repo's own `.github/workflows/`:

- **`ci.yml`**'s `lint`, `test`, `build`, and `gomod` jobs translate as-is —
  `runs-on: docker` with no explicit `container:` (this instance's default
  runner image already carries Go), `golangci-lint` installed with `go
install` instead of the marketplace action.
- **`dockerfile`**: run `hadolint` as the job's own `container:` image
  (`hadolint/hadolint:v2.15.1-alpine@sha256:...`, the `-alpine` tag
  specifically — the default tag is `FROM scratch`, no shell, and can't run a
  checkout step) rather than shelling out to `docker run` from inside another
  container. `apk add --no-cache nodejs` first, then checkout, then `hadolint
Dockerfile` — no `-v .hadolint.yaml:/.hadolint.yaml` mount needed;
  `hadolint` picks up the config file from the working directory on its own.
- **`docker`**: needs the dind-gateway rewrite described earlier. `scaffold-go-cli`'s
  `docker-build` job is the direct template — `image:
docker:27-cli@sha256:...` as the container, `apk add --no-cache nodejs git
ca-certificates curl` first, then the gateway lookup, then `docker build -t
<name>:ci-check .`. Still no registry pushed to here — this is validation
  only, same as the GitHub version.
- **`prose`** (the `ci.yml` job — Prettier and markdownlint) and the
  standalone **`prose.yml`** (Vale and `ltex-cli-plus`, mechanics/style
  split) both translate as-is: `bun install --frozen-lockfile`, then the same
  `bun run format:check` / `bun run lint:md` / `./scripts/lint-*.sh` commands
  the hooks already run.
- **`api`**: see §4.
- **`commits`**: translates as-is, `fetch-depth: 0` and the same `bunx
commitlint --from ... --to ...` invocation.
- **`gitleaks.yml`**: translates as-is once `git`/`ca-certificates`/`nodejs`
  are installed ahead of checkout and `fetch-depth: 0` is set — a shallow
  clone finds nothing and goes green, which is worse than not running it at
  all.
- **`pr-title.yml`**: **no Forgejo-hosted mirror of
  `amannn/action-semantic-pull-request` exists.** The replacement both
  `scaffold-go-cli` and `scaffold-php-api` actually run is `commitlint`
  itself, reading the pull request title from `stdin` rather than from
  commit objects:

  ```yaml
  name: pr-title

  on:
    pull_request:
      types: [opened, reopened, edited, synchronize]

  jobs:
    check:
      runs-on: docker
      steps:
        - uses: https://code.forgejo.org/actions/checkout@v4
        - uses: https://code.forgejo.org/oven-sh/setup-bun@v2
          with:
            bun-version-file: package.json
        - run: bun install --frozen-lockfile
        - name: commitlint
          env:
            PR_TITLE: ${{ github.event.pull_request.title }}
          run: printf '%s\n' "$PR_TITLE" | bunx commitlint --verbose
  ```

  This checks the same `type(scope): description` shape the individual-commit
  `commitlint` job already checks — the two can't disagree about what's a
  valid type, since there's only the one config either reads.

- **`dependabot-auto-merge.yml`**: deleted, not translated — see §2.

## 4. Does the OpenAPI spec-lint step translate?

Yes, directly, with one flag dropped. `redocly lint` itself is not
GitHub-specific — `redocly.yaml` and `api/openapi.yaml` need no changes at
all. What's GitHub-specific is the flag `ci.yml`'s `api` job passes on top of
it: `bun run lint:api -- --format=github-actions` asks Redocly to emit
GitHub's `::error file=...,line=...::message` workflow-command syntax for
inline annotations. `scaffold-php-api` — this account's own Forgejo-native
OpenAPI-spec-first sibling — runs the identical `bun run lint:api` script
with no `--format` flag at all in its `.forgejo/workflows/ci.yml`, letting
Redocly fall back to its default `stylish` console output. That's the
confirmed, working answer here too:

```yaml
api:
  runs-on: docker
  steps:
    - uses: https://code.forgejo.org/actions/checkout@v4
    - uses: https://code.forgejo.org/oven-sh/setup-bun@v2
      with:
        bun-version-file: package.json
    - run: bun install --frozen-lockfile
    - run: bun run lint:api
```

`package.json`'s `lint:api` script (`"redocly lint"`, already bare — no
`--format` baked in there either) needs no change; the flag only ever lived
in the CI invocation.

## 5. GHCR → Forgejo's own container registry

**This template doesn't publish an image anywhere today, GHCR included.**
`ci.yml`'s `docker` job comment says so directly: "There's no registry to
push it to — this template doesn't wire one up; a project stamped from it
adds that the day it needs one." So there's nothing to swap in this repo as
it stands — but if a project stamped from this template adds image
publishing, here's the working shape, taken from `scaffold-go-cli`'s own
release job, which does push a real image to `git.higherlearning.eu`'s
registry (`rules/containers.md`'s expectation that a tool built here also
ships a Docker image, the same one the Docker section of this README already
describes for the build-only case):

- **`goreleaser`'s `dockers:` block** builds and tags the image; it expects a
  registry session to already exist, the same way a bare `docker push` does
  outside `goreleaser` entirely. Add to `.goreleaser.yml`:

  ```yaml
  dockers:
    - image_templates:
        - "git.higherlearning.eu/<owner>/<repo>:{{ .Tag }}"
        - "git.higherlearning.eu/<owner>/<repo>:latest"
      goos: linux
      goarch: amd64
  ```

- **Log in before running `goreleaser`**, inside the same `artefacts` job
  from §1 — reusing `RELEASE_TOKEN`/`RELEASE_USER`, no new credential needed:

  ```yaml
  - name: Log in to git.higherlearning.eu's registry
    run: echo "${{ secrets.RELEASE_TOKEN }}" | docker login git.higherlearning.eu -u "${{ secrets.RELEASE_USER }}" --password-stdin
  ```

  `RELEASE_TOKEN` needs `write:package` scope for the push to succeed, on top
  of whatever scope §1 already needs it to have — check that on the account
  the token belongs to if this job's push ever starts failing with a
  401/403.

- **The `artefacts` job needs the same dind-gateway rewrite** §3 describes
  for `docker build` — `goreleaser`'s `dockers:` step shells out to a real
  `docker` daemon the same way, so `docker-cli` and the `DOCKER_HOST` gateway
  lookup both need to be in place before `goreleaser release --clean` runs,
  not just before a plain `docker build`.

This isn't exercised by this repo itself, unlike §1 through §4, which are
either already committed here (`release.config.mjs`) or lifted directly from
a sibling's own working `.forgejo/workflows/`. Verify the `dockers:` block
against `goreleaser`'s current docs before relying on it.
