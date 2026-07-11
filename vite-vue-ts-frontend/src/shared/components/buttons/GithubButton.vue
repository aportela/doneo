<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, NTooltip, type ButtonSize } from 'naive-ui'

    import { GITHUB_PROJECT_URL, DEFAULT_SINGLE_ACTION_BUTTON_SIZE, DEFAULT_BUTTON_ICON_SIZE } from '../../../constants';

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
</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button :size="props.size" quaternary tag="a" :href="GITHUB_PROJECT_URL" target="_blank"
                :disabled="props.disabled">
                <n-icon :size="props.iconSize">
                    <!-- credits BrandGithub icon from tabler icons -->
                    <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
                        viewBox="0 0 24 24">
                        <path
                            d="M9 19c-4.3 1.4-4.3-2.5-6-3m12 5v-3.5c0-1 .1-1.4-.5-2c2.8-.3 5.5-1.4 5.5-6a4.6 4.6 0 0 0-1.3-3.2a4.2 4.2 0 0 0-.1-3.2s-1.1-.3-3.5 1.3a12.3 12.3 0 0 0-6.2 0C6.5 2.8 5.4 3.1 5.4 3.1a4.2 4.2 0 0 0-.1 3.2A4.6 4.6 0 0 0 4 9.5c0 4.6 2.7 5.7 5.5 6c-.6.6-.6 1.2-.5 2V21"
                            fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                            stroke-linejoin="round"></path>
                    </svg>
                </n-icon>
            </n-button>
        </template>
        {{ t("shared.components.buttons.github.toolTip") }}
    </n-tooltip>
</template>

<style lang="css" scoped></style>