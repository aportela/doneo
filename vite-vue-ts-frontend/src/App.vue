<script setup lang="ts">
  import { NConfigProvider, NNotificationProvider } from 'naive-ui';
  import { darkTheme } from 'naive-ui'

  import { useColorSchemeStore } from './stores/colorScheme';

  const colorSchemeStore = useColorSchemeStore();

</script>

<template>
  <n-config-provider :theme="colorSchemeStore.dark ? darkTheme : null">
    <n-notification-provider>
      <router-view />
    </n-notification-provider>
  </n-config-provider>
</template>

<style lang="css" scoped></style>