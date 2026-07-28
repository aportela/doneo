<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, NTooltip, type ButtonSize } from 'naive-ui'

    import Github from "../icons/Github.vue";

    import { GITHUB_PROJECT_URL, DEFAULT_SINGLE_ACTION_BUTTON_SIZE, DEFAULT_BUTTON_ICON_SIZE } from '../../../constants';

    interface Props {
        size?: ButtonSize;
        iconSize?: number,
        disabled?: boolean;
    };

    const props = withDefaults(defineProps<Props>(), {
        size: DEFAULT_SINGLE_ACTION_BUTTON_SIZE,
        iconSize: DEFAULT_BUTTON_ICON_SIZE,
        disabled: false,
    });

    const { t } = useI18n();
</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button :size="props.size" quaternary tag="a" :href="GITHUB_PROJECT_URL" target="_blank"
                :disabled="props.disabled">
                <n-icon :size="props.iconSize" :component="Github" />
            </n-button>
        </template>
        {{ t("shared.components.buttons.github.toolTip") }}
    </n-tooltip>
</template>

<style lang="css" scoped></style>