<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, NTooltip } from 'naive-ui';
    import { IconMoon, IconSun } from '@tabler/icons-vue';

    import { useColorSchemeStore } from "../../../stores/colorScheme";

    import { BUTTON_DEFAULT_ICON_SIZE } from '../../../constants';

    interface IProps {
        iconSize?: number,
        disabled?: boolean;
    };

    const props = withDefaults(defineProps<IProps>(), {
        iconSize: BUTTON_DEFAULT_ICON_SIZE,
        disabled: false,
    });

    const { t } = useI18n();

    const colorSchemeStore = useColorSchemeStore();
</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button quaternary @click.prevent="colorSchemeStore.toggle" @mousedown.prevent :disabled="props.disabled">
                <n-icon :size="props.iconSize" :component="colorSchemeStore.light ? IconMoon : IconSun" />
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