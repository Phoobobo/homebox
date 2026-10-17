<template>
  <Card class="flex flex-row items-center justify-between gap-4 bg-transparent p-5">
    <CardHeader class="p-0">
      <CardTitle class="text-sm font-medium">{{ title }}</CardTitle>
    </CardHeader>
    <CardContent class="p-0 text-3xl font-semibold tabular-nums text-primary">
      <Currency v-if="type === 'currency'" :amount="value" />
      <template v-if="type === 'number'">{{ value }}</template>
    </CardContent>
    <CardFooter v-if="subtitle">{{ subtitle }}</CardFooter>
  </Card>
</template>

<script setup lang="ts">
  import Currency from "../Currency.vue";
  import type { StatsFormat } from "./types";
  import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";

  type Props = {
    title: string;
    value: number;
    subtitle?: string;
    type?: StatsFormat;
  };

  withDefaults(defineProps<Props>(), {
    type: "number",
    subtitle: undefined,
  });
</script>

<style></style>
