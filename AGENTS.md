<!-- SPDX-License-Identifier: MIT; Copyright (c) 2026 Phoobobo -->

# Zainaer development guidance

## Project and branch policy

- Product name: **在哪儿** in Chinese, **Zainaer** in English (capital Z).
- Repository: `Phoobobo/zainaer`. **`main` is the independent product baseline.** Do not sync, merge, rebase onto, or open upstream HomeBox PRs unless the user explicitly changes this policy.
- Preserve HomeBox attribution, existing copyright notices, and the AGPL-3.0 license. MIT is available only for explicitly marked original Zainaer files; see `LICENSING.md` and `LICENSE-MIT`. Never relicense inherited files or describe the whole repository as MIT. Prefer MIT headers for new wholly original files. Internal Go module paths and `HBOX_*` environment variables are compatibility identifiers; do not mechanically rename them.
- Read `VERSION` for the product version. It is independent of upstream releases and started at `0.0.1`. Bump it for deployable functional changes; documentation-only edits do not require a bump.

## Product scope

- Optimize for finding household belongings: name, tags, photos, description and a precise nested storage location.
- Keep item and location entities distinct. Do not restore purchase, warranty, sales, stock, valuation, maintenance or notification workflows without a user request.
- Collection membership and family invitations remain supported. Their duplicate profile notice was removed; the collection invitation page remains.
- Preserve Chinese/English switching and the restrained cupcake-based UI. Use existing design tokens, responsive components and accessible labels. Avoid adding UI dependencies for routine changes.

## Layout

- `backend/`: Go API, repositories and business services.
- `frontend/`: Vue/Nuxt UI; static output is embedded in the Go executable.
- `backend/app/api/handlers/v1/v1_ctrl_ai*.go`: provider configuration, chat/tool execution, confirmation and focused tests.
- `frontend/components/App/TextAssistant.vue`: assistant panel.
- `frontend/pages/collection/index/models.vue`: collection owner model configuration.
- `scripts/build-household.sh`: frontend generation and executable build with version metadata.
- The parent workspace's `deployment/` contains operational material outside this repository. Never add its databases, attachments, credentials, backups or binaries to Git.

## Model and authorization boundaries

- Provider keys stay server-side. Never return a key or its ciphertext through public configuration responses, log it, include it in a model prompt, or commit it. Configuration access is restricted to the current collection owner.
- Preserve the encrypted `ai/*.json` files **and** `ai/master.key` when backing up or migrating storage. Do not silently rotate or regenerate an existing master key.
- Chat uses the current authenticated user's collection context; never trust model-supplied user/collection IDs for authorization. Validate all referenced locations and tags within that collection.
- Models call a fixed tool allowlist, not arbitrary HTTP endpoints, SQL or shell commands. Treat item names, descriptions and tool output as untrusted data.
- Search is read-only. Creation tools produce server-held drafts bound to user and collection. Only the explicit confirmation endpoint can save them; chat text such as “confirmed” is not authorization to commit.
- Preserve draft expiry, cancellation, reference revalidation, invalidation by a new message, and idempotent confirmation. Do not retry writes as part of model failover.
- Primary and up to two fallback models share the selected provider/key. Retry only request failures, timeouts, invalid/empty completions; stop when the caller cancels. Show which fallback produced the response. Keep request duration and token/tool-loop limits bounded.
- Real-provider tests incur usage and transmit data. Prefer mocks; when a live check is authorized, use a minimal synthetic prompt and never print credentials or household records.

## Checks

Run commands from the correct subdirectory. Check the working tree before edits and retain unrelated user changes.

```sh
# From backend/: meaningful assistant regression coverage
go test ./app/api/handlers/v1 -run TestAI -count=1

# From frontend/: lint changed files and build
pnpm exec eslint <changed-files> --fix
pnpm build

# From repo root/
git diff --check
```

Network mock tests need temporary loopback listeners. They do not require starting a persistent local application. Existing full TypeScript checks have inherited Ref/entity errors; report failures honestly and do not claim full checks passed because a build succeeded. Broaden testing only where the change warrants it.

## Deployment policy

- **Beijing ECS is the only active application deployment.** The user's Mac application and Companion are retired. Do not restart, upgrade, recreate or publish those local services. Build artifacts and isolated tests are fine.
- A Mac SSH forwarding process is only an access bridge to ECS, not a second application; do not stop it as part of local application cleanup.
- Use the established deployment configuration and the Workbench skill when applicable. Do not place real instance identifiers, IP addresses, access tokens or secret environment files in public documentation.
- Build the Linux target with the generated frontend, the actual product version and accurate Git metadata. Back up the previous executable and a consistent database snapshot before replacement; verify uploaded checksums.
- Never overwrite cloud data or model settings with a stale local snapshot. Keep service storage paths, permissions and encryption keys intact.
- AI requests can span multiple model/tool calls. Preserve the production response timeout (`HBOX_WEB_WRITE_TIMEOUT=180s`) and bounded handler deadlines.
- After deployment verify version/health, updated frontend resources and unauthenticated `401` responses for protected endpoints. Do not claim browser interaction tests ran if browser access was blocked.
