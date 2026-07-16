<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, NTooltip, type ButtonSize } from 'naive-ui';
    import { Sun, Moon } from "@lucide/vue";

    import { useThemeStore } from "../../../stores/theme";

    import { DEFAULT_SINGLE_ACTION_BUTTON_SIZE, DEFAULT_BUTTON_ICON_SIZE } from '../../../constants';

    interface IProps {
        size?: ButtonSize;
        iconSize?: number,
        disabled?: boolean;
    };

    const props = withDefaults(defineProps<IProps>(), {
        size: DEFAULT_SINGLE_ACTION_BUTTON_SIZE,
        iconSize: DEFAULT_BUTTON_ICON_SIZE,
        disabled: false,
    });

    const { t } = useI18n();

    const themeStore = useThemeStore();
</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button :size="props.size" quaternary @click.prevent="themeStore.toggle" @mousedown.prevent
                :disabled="props.disabled">
                <n-icon :size="props.iconSize" :component="themeStore.light ? Moon : Sun" />
            </n-button>
        </template>
        {{
            t(themeStore.light ?
                "shared.components.buttons.colorScheme.darkMode.toolTip" :
                "shared.components.buttons.colorScheme.lightMode.toolTip")
        }}
    </n-tooltip>
</template>

<style lang="css" scoped></style>