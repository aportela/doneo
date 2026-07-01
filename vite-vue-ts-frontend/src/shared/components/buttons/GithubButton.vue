<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, NTooltip } from 'naive-ui'
    import { IconBrandGithub } from '@tabler/icons-vue';

    import { GITHUB_PROJECT_URL, BUTTON_DEFAULT_ICON_SIZE } from '../../../constants';

    interface IProps {
        iconSize?: number,
    };

    const props = withDefaults(defineProps<IProps>(), {
        iconSize: BUTTON_DEFAULT_ICON_SIZE
    });

    const { t } = useI18n();
</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button quaternary tag="a" :href="GITHUB_PROJECT_URL" target="_blank">
                <n-icon :size="props.iconSize" :component="IconBrandGithub" />
            </n-button>
        </template>
        {{ t("shared.components.buttons.github.toolTip") }}
    </n-tooltip>
</template>

<style lang="css" scoped></style>