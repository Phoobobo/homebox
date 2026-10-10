<script setup lang="ts">
  import { useI18n } from "vue-i18n";
  import { Button } from "~/components/ui/button";
  import { useViewPreferences } from "~~/composables/use-preferences";

  const { locale, t } = useI18n();
  const preferences = useViewPreferences();

  function setLanguage(language: "zh-CN" | "en") {
    locale.value = language;
    preferences.value.language = language;
  }
</script>

<template>
  <div class="flex shrink-0 items-center gap-1" role="group" :aria-label="t('profile.language')">
    <Button
      type="button"
      size="sm"
      :variant="locale === 'zh-CN' ? 'default' : 'outline'"
      :aria-pressed="locale === 'zh-CN'"
      aria-label="切换到简体中文"
      lang="zh-CN"
      @click="setLanguage('zh-CN')"
    >
      中文
    </Button>
    <Button
      type="button"
      size="sm"
      :variant="locale === 'en' ? 'default' : 'outline'"
      :aria-pressed="locale === 'en'"
      aria-label="Switch to English"
      lang="en"
      @click="setLanguage('en')"
    >
      EN
    </Button>
  </div>
</template>
