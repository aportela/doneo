<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NCard } from "naive-ui";

    import TaskPrioritiesTable from "../components/TaskPrioritiesTable.vue";

    const { t } = useI18n();
</script>

<template>
    <n-card :title="t('modules.taskPriority.components.ManageTaskPrioritiesPage.header.title')">
        <TaskPrioritiesTable />
    </n-card>
</template>

<style lang="css" scoped></style>