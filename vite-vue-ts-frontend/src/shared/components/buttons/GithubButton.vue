<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, NTooltip } from 'naive-ui'
    import { IconBrandGithub } from '@tabler/icons-vue';

    import { GITHUB_PROJECT_URL } from '../../constants';

    interface GithubButtonProps {
        iconSize?: number,
    };

    withDefaults(defineProps<GithubButtonProps>(), {
        iconSize: 20
    });

    const { t } = useI18n();

</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button quaternary>
                <a :href="GITHUB_PROJECT_URL" class="link-color" target="_blank">
                    <n-icon :size="iconSize" :component="IconBrandGithub" />
                </a>
            </n-button>
        </template>
        {{ t("Open project page") }}
    </n-tooltip>
</template>

<style lang="css" scoped>
    .link-color {
        color: var(--n-text-color-base) !important;
    }
</style>