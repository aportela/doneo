<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NCard } from 'naive-ui';

    import ProjectStatusesTable from "../components/ProjectStatusesTable.vue";

    const { t } = useI18n();
</script>

<template>
    <n-card :title="t('modules.projectStatus.components.ManageProjectStatusesPage.header.title')">
        <ProjectStatusesTable />
    </n-card>
</template>

<style lang="css" scoped></style>