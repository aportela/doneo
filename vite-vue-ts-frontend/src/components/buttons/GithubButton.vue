<script setup lang="ts">
    import { useI18n } from "vue-i18n";
    import { NButton, NTooltip } from 'naive-ui'
    import { IconBrandGithub } from '@tabler/icons-vue';
    import { GITHUB_PROJECT_URL } from '../../constants';

    const { t } = useI18n();

    interface GithubButtonProps {
        size?: number,
    };

    withDefaults(defineProps<GithubButtonProps>(), {
        size: 20
    });
</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button quaternary>
                <a :href="GITHUB_PROJECT_URL" class="link-color" target="_blank">
                    <IconBrandGithub :size="size" />
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