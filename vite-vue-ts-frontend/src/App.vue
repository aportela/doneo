<script setup lang="ts">
  import { NConfigProvider, NNotificationProvider } from 'naive-ui';
  import { darkTheme } from 'naive-ui'

  import { useThemeStore } from './stores/theme';

  const themeStore = useThemeStore();

</script>

<template>
  <n-config-provider :theme="themeStore.dark ? darkTheme : undefined"
    :theme-overrides="themeStore.dark ? undefined : undefined">
    <n-notification-provider>
      <router-view />
    </n-notification-provider>
  </n-config-provider>
</template>

<style lang="css" scoped></style>