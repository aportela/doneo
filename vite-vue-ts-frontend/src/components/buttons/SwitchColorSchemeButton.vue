<script setup lang="ts">
    import { NButton } from 'naive-ui';
    import { IconMoon, IconSun } from '@tabler/icons-vue';
    import { useColorSchemeStore } from '../../stores/colorScheme';

    interface SwitchColorSchemeButtonProps {
        size?: number,
    };

    withDefaults(defineProps<SwitchColorSchemeButtonProps>(), {
        size: 20
    });

    const colorSchemeStore = useColorSchemeStore();
</script>

<template>
    <n-button quaternary>
        <IconMoon :size="size" v-if="colorSchemeStore.light" @click.prevent="colorSchemeStore.toggle"
            @mousedown.prevent />
        <IconSun :size="size" v-else @click="colorSchemeStore.toggle" @mousedown.prevent />
    </n-button>
</template>

<style lang="css" scoped></style>