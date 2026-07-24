<script setup lang="ts">
    import { ref } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NTabs, NTabPane, NIcon, NCard } from 'naive-ui';
    import { IconList, IconLayoutKanban, IconCalendarWeek } from '@tabler/icons-vue';

    import TasksTable from '../components/TasksTable.vue';
    import TasksKanban from '../components/TasksKanban.vue';
    import TasksCalendar from '../components/TasksCalendar.vue';

    const { t } = useI18n();

    const tab = ref<string>("Kanban");
</script>

<template>
    <n-tabs v-model="tab">
        <n-tab-pane name="List" tab="List" type="segment" animated v-model="tab" display-directive="show:lazy">
            <template #tab>
                <n-icon size="16" class="tab-icon">
                    <IconList />
                </n-icon>
                List
            </template>
            <template #default>
                <n-card :title="t('modules.task.components.ManageTasksPage.header.title')">
                    <TasksTable />
                </n-card>
            </template>
        </n-tab-pane>
        <n-tab-pane name="Kanban" tab="Kanban" display-directive="show:lazy">
            <template #tab>
                <n-icon size="16" class="tab-icon">
                    <IconLayoutKanban />
                </n-icon>
                Kanban
            </template>
            <template #default>
                <TasksKanban />
            </template>
        </n-tab-pane>
        <n-tab-pane name="Calendar" tab="Calendar" display-directive="show:lazy">
            <template #tab>
                <n-icon size="16" class="doneo-mr-4px" :component="IconCalendarWeek" />
                Calendar
            </template>
            <template #default>
                <TasksCalendar />
            </template>
        </n-tab-pane>
    </n-tabs>
</template>

<style lang="css" scoped></style>