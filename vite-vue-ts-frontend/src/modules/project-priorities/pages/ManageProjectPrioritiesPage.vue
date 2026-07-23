<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NCard } from "naive-ui";

    import ProjectPrioritiesTable from '../components/ProjectPrioritiesTable.vue';

    const { t } = useI18n();
</script>

<template>
    <n-card :title="t('modules.projectPriority.components.ManageProjectPrioritiesPage.header.title')">
        <ProjectPrioritiesTable />
    </n-card>
</template>

<style lang="css" scoped></style>