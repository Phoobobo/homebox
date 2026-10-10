import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { IntlMessageFormat } from "intl-messageformat";

function flatten(object, prefix = "") {
  return Object.fromEntries(
    Object.entries(object).flatMap(([key, value]) =>
      typeof value === "object"
        ? Object.entries(flatten(value, `${prefix}${key}.`))
        : [[`${prefix}${key}`, value]]
    )
  );
}

function argumentsIn(nodes, names = new Set()) {
  for (const node of nodes) {
    if (node.type >= 1 && node.type <= 6) names.add(node.value);
    for (const option of Object.values(node.options ?? {})) argumentsIn(option.value, names);
    if (node.children) argumentsIn(node.children, names);
  }
  return [...names].sort();
}

const english = flatten(JSON.parse(readFileSync(new URL("../locales/en.json", import.meta.url), "utf8")));
const chinese = flatten(JSON.parse(readFileSync(new URL("../locales/zh-CN.json", import.meta.url), "utf8")));
for (const [key, message] of Object.entries(english)) {
  assert.ok(chinese[key], `Missing Chinese translation: ${key}`);
  const source = new IntlMessageFormat(message, "en");
  const translated = new IntlMessageFormat(chinese[key], "zh-CN");
  assert.deepEqual(argumentsIn(translated.getAst()), argumentsIn(source.getAst()), `Placeholder mismatch: ${key}`);
}
console.log(`Validated ${Object.keys(english).length} Chinese translations and ICU placeholders.`);
