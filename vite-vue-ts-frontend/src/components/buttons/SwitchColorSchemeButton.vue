<script setup lang="ts">
    import { NButton } from 'naive-ui';
    import { IconMoon, IconSun } from '@tabler/icons-vue';
    import { useColorSchemeStore } from '../../stores/colorScheme';

    const colorSchemeStore = useColorSchemeStore();
</script>

<template>
    <n-button quaternary>
        <IconMoon v-if="colorSchemeStore.light" @click.prevent="colorSchemeStore.toggle" @mousedown.prevent />
        <IconSun v-else @click="colorSchemeStore.toggle" @mousedown.prevent />
    </n-button>
</template>

<style lang="css" scoped></style>