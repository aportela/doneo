<script setup lang="ts">
  import { NConfigProvider, NNotificationProvider } from 'naive-ui';
  import { darkTheme } from 'naive-ui'

  import { useColorSchemeStore } from './stores/colorScheme';

  const colorSchemeStore = useColorSchemeStore();

  import { type GlobalThemeOverrides } from "naive-ui";

  const draculaTheme: GlobalThemeOverrides = {
    common: {
      // colores globales
      primaryColor: "#BD93F9",
      primaryColorHover: "#CDA7FF",
      primaryColorPressed: "#A77BF3",

      bodyColor: "#282A36",
      cardColor: "#44475A",
      modalColor: "#44475A",

      textColorBase: "#F8F8F2",
      textColor1: "#F8F8F2",
      textColor2: "#F8F8F2",
      textColor3: "#6272A4",

      borderColor: "#6272A4",

      inputColor: "#44475A",
      tableColor: "#282A36",
    },
  };
</script>

<template>
  <n-config-provider :theme="colorSchemeStore.dark ? darkTheme : undefined"
    :theme-overrides="colorSchemeStore.dark ? undefined : undefined">
    <n-notification-provider>
      <router-view />
    </n-notification-provider>
  </n-config-provider>
</template>

<style lang="css" scoped></style>