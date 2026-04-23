<script setup lang="ts">

  const theme = localStorage.getItem("theme");
  if (theme === "dark") {
    document.documentElement.setAttribute(
      'data-bs-theme', 'dark'
    )
  }
</script>

<template>
  <router-view />
</template>