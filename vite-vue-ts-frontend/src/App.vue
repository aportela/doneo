<script setup lang="ts">
  import { NConfigProvider } from 'naive-ui';
  import { darkTheme } from 'naive-ui'

  import { useColorSchemeStore } from './stores/colorScheme';

  const colorSchemeStore = useColorSchemeStore();

</script>

<template>
  <n-config-provider :theme="colorSchemeStore.dark ? darkTheme : null">
    <router-view />
  </n-config-provider>
</template>

<style lang="css" scoped></style>