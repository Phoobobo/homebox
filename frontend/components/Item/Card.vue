<template>
  <Card class="relative isolate aspect-[4/5] min-h-[260px] overflow-hidden transition-colors hover:border-primary/40">
    <img
      v-if="item.imageId"
      class="absolute inset-0 size-full object-cover"
      loading="lazy"
      :src="imageUrl"
      :alt="item.name"
    />
    <div v-else class="absolute inset-0 flex items-center justify-center bg-secondary/40 pb-24 text-primary/40">
      <MdiPackageVariant class="size-14" />
    </div>

    <NuxtLink
      :to="`/item/${item.id}`"
      class="absolute inset-0 z-10 rounded-[inherit] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-primary"
      :aria-label="item.name"
    />

    <div v-if="tableRow" class="absolute left-3 top-3 z-30">
      <Checkbox
        class="size-5 bg-accent hover:bg-background-accent"
        :model-value="tableRow.getIsSelected()"
        :aria-label="$t('components.item.view.selectable.select_card')"
        @update:model-value="tableRow.toggleSelected()"
      />
    </div>

    <div
      class="pointer-events-none absolute inset-x-0 bottom-0 z-20 flex flex-col gap-2 border-t border-foreground/10 bg-background/85 p-4 text-foreground supports-[backdrop-filter]:bg-background/70 supports-[backdrop-filter]:backdrop-blur-xl"
    >
      <NuxtLink
        v-if="item.parent"
        :to="`/location/${item.parent.id}`"
        class="pointer-events-auto flex w-fit max-w-full items-start gap-1.5 text-xs font-medium leading-5 text-foreground/80 hover:underline"
      >
        <MdiMapMarkerOutline class="mt-0.5 size-4 shrink-0" />
        <span class="line-clamp-2 break-words">{{ locationString }}</span>
      </NuxtLink>
      <h2 class="line-clamp-2 break-words text-base font-semibold leading-snug">
        {{ item.name }}
      </h2>
      <Markdown v-if="item.description" class="line-clamp-2 text-sm text-foreground/75" :source="item.description" />
      <div v-if="itemTags.length" class="pointer-events-auto flex flex-wrap gap-1.5">
        <TagChip v-for="tag in itemTags" :key="tag.id" :tag="tag" size="sm" :ancestors="tag.ancestors" />
      </div>
    </div>
  </Card>
</template>

<script setup lang="ts">
  import MdiMapMarkerOutline from "~icons/mdi/map-marker-outline";
  import MdiPackageVariant from "~icons/mdi/package-variant";
  import type { EntityOut, EntitySummary } from "~~/lib/api/types/data-contracts";
  import { Card } from "@/components/ui/card";
  import Markdown from "@/components/global/Markdown.vue";
  import TagChip from "@/components/Tag/Chip.vue";
  import type { Row } from "@tanstack/vue-table";
  import { Checkbox } from "@/components/ui/checkbox";

  const api = useUserApi();

  const imageUrl = computed(() => {
    if (!props.item.imageId) {
      return "/no-image.jpg";
    }
    if (props.item.thumbnailId) {
      return api.authURL(`/entities/${props.item.id}/attachments/${props.item.thumbnailId}`);
    } else {
      return api.authURL(`/entities/${props.item.id}/attachments/${props.item.imageId}`);
    }
  });

  const itemTags = computed(() => {
    return useTagStore().withAncestors(props.item.tags);
  });

  const props = defineProps({
    item: {
      type: Object as () => EntityOut | EntitySummary,
      required: true,
    },
    locationFlatTree: {
      type: Array as () => FlatTreeItem[],
      required: false,
      default: () => [],
    },
    tableRow: {
      type: Object as () => Row<EntitySummary>,
      required: false,
      default: () => null,
    },
  });

  const locationString = computed(
    () => props.locationFlatTree.find(l => l.id === props.item.parent?.id)?.treeString || props.item.parent?.name
  );
</script>

<style lang="css"></style>
