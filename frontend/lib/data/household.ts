import type { EntityTypeSummary } from "~/lib/api/types/data-contracts";

// Match the server's GetDefault, even when type names have been customized.
export function householdEntityTypes(types: EntityTypeSummary[]) {
  const sorted = [...types].sort((a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime());
  return [false, true].flatMap(isLocation => {
    const type = sorted.find(type => type.isLocation === isLocation);
    return type ? [type] : [];
  });
}

export const householdItemMode = true;
