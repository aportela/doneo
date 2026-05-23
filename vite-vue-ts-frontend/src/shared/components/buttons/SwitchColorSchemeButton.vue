<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, NTooltip } from 'naive-ui';
    import { IconMoon, IconSun } from '@tabler/icons-vue';

    import { useColorSchemeStore } from "../../../stores/colorScheme";

    interface SwitchColorSchemeButtonProps {
        iconSize?: number,
    };

    withDefaults(defineProps<SwitchColorSchemeButtonProps>(), {
        iconSize: 20
    });

    const { t } = useI18n();

    const colorSchemeStore = useColorSchemeStore();

</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button quaternary @click.prevent="colorSchemeStore.toggle" @mousedown.prevent>
                <n-icon :size="iconSize" :component="colorSchemeStore.light ? IconMoon : IconSun" />
            </n-button>
        </template>
        {{
            t(colorSchemeStore.light ?
                "shared.components.buttons.colorScheme.darkMode.toolTip" :
                "shared.components.buttons.colorScheme.lightMode.toolTip")
        }}
    </n-tooltip>
</template>

<style lang="css" scoped></style>