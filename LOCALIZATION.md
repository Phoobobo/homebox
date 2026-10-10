# Chinese / English trial build

This fork's `feat/zh-en-switch` branch is based on HomeBox v0.26.2, matching the locally deployed stable version.

- The login page and application header provide direct **中文 / EN** buttons.
- The selected language uses HomeBox's existing view preferences, survives a reload, and participates in the existing preference synchronization.
- The full language selector remains available on the login and profile pages.
- Simplified Chinese translations cover all 816 English message keys, including item creation and entity type management.
- Item creation uses a translated heading instead of a hardcoded English title.
- Existing location names, tags, and user content remain as entered; switching the interface language does not rename database records.

## Validation

From `frontend`, after installing the locked dependencies:

```sh
node scripts/check-zh-cn.mjs
pnpm build
```

The translation check validates coverage, ICU syntax, and interpolation argument parity against English.
Browser checks cover Chinese and English switching, reload persistence, and the item creation dialog.
The upstream v0.26.2 frontend has pre-existing TypeScript diagnostics; comparison against the unmodified source showed no new diagnostics for the language controls.

## Native build

```sh
cd frontend
pnpm build
mkdir -p ../backend/app/api/static/public
cp -R .output/public/. ../backend/app/api/static/public/
cd ../backend
go build -o ../../deployment/homebox/homebox-bilingual ./app/api
```

The binary embeds the generated frontend. Preserve the existing data directory and API-key pepper when replacing a running installation.
