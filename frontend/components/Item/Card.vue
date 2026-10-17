<template>
  <Card class="relative overflow-hidden transition-colors hover:border-primary/40">
    <div v-if="tableRow" class="absolute left-1 top-1 z-10">
      <Checkbox
        class="size-5 bg-accent hover:bg-background-accent"
        :model-value="tableRow.getIsSelected()"
        :aria-label="$t('components.item.view.selectable.select_card')"
        @update:model-value="tableRow.toggleSelected()"
      />
    </div>
    <NuxtLink :to="`/item/${item.id}`">
      <div class="relative h-[200px]">
        <img
          v-if="item.imageId"
          class="absolute h-[200px] w-full bg-muted/30"
          :class="objectContain ? 'object-contain' : 'object-cover'"
          loading="lazy"
          :src="imageUrl"
          :alt="item.name"
        />
        <div v-else class="flex h-full items-center justify-center bg-secondary/40 text-primary/40">
          <MdiPackageVariant class="size-14" />
        </div>
        <div class="absolute inset-x-1 bottom-1">
          <Badge class="text-wrap bg-secondary text-secondary-foreground hover:bg-secondary/70 hover:underline">
            <NuxtLink v-if="item.parent" :to="`/location/${item.parent.id}`">
              {{ locationString }}
            </NuxtLink>
          </Badge>
        </div>
      </div>
      <div class="col-span-4 flex grow flex-col gap-y-2 p-4">
        <h2 class="line-clamp-2 text-ellipsis text-wrap text-base font-semibold">
          {{ item.name }}
        </h2>

        <Markdown class="mb-2 line-clamp-2 text-sm text-muted-foreground" :source="item.description" />
        <div class="-mr-1 mt-auto flex flex-wrap justify-start gap-2">
          <TagChip v-for="tag in itemTags" :key="tag.id" :tag="tag" size="sm" :ancestors="tag.ancestors" />
        </div>
      </div>
    </NuxtLink>
  </Card>
</template>

<script setup lang="ts">
  import MdiPackageVariant from "~icons/mdi/package-variant";
  import type { EntityOut, EntitySummary } from "~~/lib/api/types/data-contracts";
  import { Badge } from "@/components/ui/badge";
  import { Card } from "@/components/ui/card";
  import Markdown from "@/components/global/Markdown.vue";
  import TagChip from "@/components/Tag/Chip.vue";
  import type { Row } from "@tanstack/vue-table";
  import { Checkbox } from "@/components/ui/checkbox";

  const api = useUserApi();
  const preferences = useViewPreferences();

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

  const objectContain = computed(() => imageUrl.value !== "/no-image.jpg" && !preferences.value.legacyImageFit);

  const locationString = computed(
    () => props.locationFlatTree.find(l => l.id === props.item.parent?.id)?.treeString || props.item.parent?.name
  );
</script>

<style lang="css"></style>
