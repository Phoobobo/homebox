<script setup lang="ts">
  import { computed, h } from "vue";
  import type { EntitySummary } from "~/lib/api/types/data-contracts";
  import DataTable from "./table/data-table.vue";
  import { makeColumns } from "./table/columns";
  import { useLocationStore } from "~/stores/locations";
  import { useI18n } from "vue-i18n";

  defineProps<{
    items: EntitySummary[];
  }>();

  const { t } = useI18n();
  const locations = useLocationStore();

  function locationPath(id: string, fallback: string) {
    function find(nodes: typeof locations.tree, path: string[] = []): string[] | undefined {
      for (const node of nodes ?? []) {
        const next = [...path, node.name];
        if (node.id === id) return next;
        const found = find(node.children, next);
        if (found) return found;
      }
    }
    return find(locations.tree)?.join(" → ") || fallback;
  }

  const columns = computed(() =>
    makeColumns({ t })
      .filter(c => ["name", "location"].includes(c.id ?? ""))
      .map(column =>
        column.id !== "location"
          ? column
          : {
              ...column,
              accessorFn: (item: EntitySummary) => (item.parent ? locationPath(item.parent.id, item.parent.name) : ""),
              cell: ({ row }: { row: { original: EntitySummary } }) => {
                const parent = row.original.parent;
                return parent
                  ? h(
                      "a",
                      { href: `/location/${parent.id}`, class: "text-sm hover:underline" },
                      locationPath(parent.id, parent.name)
                    )
                  : "—";
              },
            }
      )
  );
</script>

<template>
  <DataTable view="table" :data="items" :columns="columns" disable-controls />
</template>
