<script setup lang="ts">
    import { NButton } from 'naive-ui'
    import { IconBrandGithub } from '@tabler/icons-vue';
    import { GITHUB_PROJECT_URL } from '../../constants';

    interface GithubButtonProps {
        size?: number,
    };

    withDefaults(defineProps<GithubButtonProps>(), {
        size: 20
    });
</script>

<template>
    <n-button quaternary>
        <a :href="GITHUB_PROJECT_URL" target="_blank">
            <IconBrandGithub :size="size" />
        </a>
    </n-button>
</template>

<style lang="css" scoped></style>