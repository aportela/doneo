<template>
  <iframe id="doneo-pdf-wrapper-container" :src="url" :class="props.innerContentClass"></iframe>
</template>

<script setup lang="ts">
  interface IProps {
    url: string
    innerContentClass?: string;
  };

  const props = defineProps<IProps>();
</script>

<style lang="css" scoped>
  .doneo-pdf-wrapper-container {
    display: block;
    width: 100%;
    margin: 0px auto;
  }
</style>