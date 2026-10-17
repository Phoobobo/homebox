<template>
  <Card>
    <NuxtLink :to="`/location/${location.id}`" class="group/location-card transition duration-300">
      <div
        :class="{
          'p-4': !dense,
          'px-3 py-2': dense,
        }"
      >
        <h2 class="flex items-center justify-between gap-2">
          <MdiMapMarkerOutline class="size-5 shrink-0 text-primary" />
          <span class="mr-auto font-medium">{{ location.name }}</span>
          <Badge :class="{ 'opacity-0': !hasCount }">
            {{ count }}
          </Badge>
        </h2>
      </div>
    </NuxtLink>
  </Card>
</template>

<script lang="ts" setup>
  import type { EntityOut, EntitySummary } from "~~/lib/api/types/data-contracts";
  import MdiMapMarkerOutline from "~icons/mdi/map-marker-outline";
  import { Card } from "@/components/ui/card";
  import { Badge } from "@/components/ui/badge";

  const props = defineProps({
    location: {
      type: Object as () => EntitySummary | EntityOut,
      required: true,
    },
    dense: {
      type: Boolean,
      default: false,
    },
  });

  const hasCount = computed(() => {
    return !!(props.location as EntitySummary).itemCount;
  });

  const count = computed(() => {
    return hasCount.value ? (props.location as EntitySummary).itemCount : undefined;
  });
</script>
