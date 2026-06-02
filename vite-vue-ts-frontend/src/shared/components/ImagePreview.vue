<script setup lang="ts">
    import { NImageGroup } from 'naive-ui';

    interface ImagePreviewProps {
        sources: string[];
    };

    const props = defineProps<ImagePreviewProps>();

    const show = defineModel<boolean>("show");
    const currentIndex = defineModel<number>("currentIndex");
</script>

<template>
    <n-image-group v-model:show="show" v-model:current="currentIndex" :src-list="sources" />
</template>

<style lang="css" scoped></style>