<script setup lang="ts">
  import { NConfigProvider, NNotificationProvider } from 'naive-ui';
  import { darkTheme } from 'naive-ui'

  import { useUserSettingsStore } from './stores/userSettings';

  const userSettingsStore = useUserSettingsStore();

</script>

<template>
  <n-config-provider :theme="userSettingsStore.darkTheme ? darkTheme : undefined"
    :theme-overrides="userSettingsStore.darkTheme ? undefined : undefined">
    <n-notification-provider>
      <router-view />
    </n-notification-provider>
  </n-config-provider>
</template>

<style lang="css" scoped></style>