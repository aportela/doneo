<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NCard } from 'naive-ui';

    import TaskStatusesTable from '../components/TaskStatusesTable.vue';

    const { t } = useI18n();
</script>

<template>
    <n-card :title="t('modules.taskStatus.components.ManageTaskStatusesPage.header.title')">
        <TaskStatusesTable />
    </n-card>
</template>

<style lang="css" scoped></style>