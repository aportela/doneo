<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NCard } from 'naive-ui';

    import ProjectTypesTable from '../components/ProjectTypesTable.vue';

    const { t } = useI18n();
</script>

<template>
    <n-card :title="t('modules.projectType.components.ManageProjectTypesPage.header.title')">
        <ProjectTypesTable />
    </n-card>
</template>

<style lang="css" scoped></style>