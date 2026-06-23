<script setup lang="ts">
    import { NBreadcrumb, NBreadcrumbItem, NIcon } from 'naive-ui';

    import { IconHome, IconUser } from '@tabler/icons-vue';
</script>

<template>
    <n-breadcrumb>
        <n-breadcrumb-item><n-icon :component="IconHome" /> Home</n-breadcrumb-item>
        <n-breadcrumb-item>
            <n-icon :component="IconUser" /> Account
        </n-breadcrumb-item>
        <n-breadcrumb-item>
            John Doe
        </n-breadcrumb-item>
    </n-breadcrumb>
</template>

<style lang="css" scoped></style>