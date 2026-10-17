<script lang="ts" setup>
  import { cn } from "~/lib/utils";

  const props = defineProps({
    cmp: {
      type: String,
      default: "div",
    },
    class: {
      type: String,
      default: "",
    },
  });

  const cmpClass = computed(() => cn("container mx-auto my-6 max-w-7xl px-4 sm:px-6 lg:px-8", props.class));
</script>

<template>
  <component :is="cmp" :class="cmpClass">
    <slot />
  </component>
</template>
