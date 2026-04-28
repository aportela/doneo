<script setup lang="ts">
  import "@tabler/core/dist/css/tabler.min.css";
  import "@tabler/core/dist/js/tabler.min.js";

  const theme = localStorage.getItem("theme");
  if (theme === "dark") {
    document.documentElement.setAttribute(
      'data-bs-theme', 'dark'
    )
  }
  const fluidLayout = localStorage.getItem("fluidLayout")

  if (fluidLayout == "true") {
    document.body.classList.add("layout-fluid");
  }

</script>

<template>
  <router-view />
</template>