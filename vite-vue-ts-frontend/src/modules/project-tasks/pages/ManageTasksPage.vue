<script setup lang="ts">
    import { NTabs, NTabPane, NIcon } from 'naive-ui';
    import { IconList, IconLayoutKanban, IconCalendarWeek } from '@tabler/icons-vue';
</script>

<template>


    <n-tabs>
        <n-tab-pane name="List" tab="List" type="segment" animated>
            <template #tab>
                <n-icon size="16" class="tab-icon">
                    <IconList />
                </n-icon>
                List
            </template>
            <template #default>
                LIST
            </template>
        </n-tab-pane>
        <n-tab-pane name="Kanban" tab="Kanban">
            <template #tab>
                <n-icon size="16" class="tab-icon">
                    <IconLayoutKanban />
                </n-icon>
                Kanban
            </template>
            <template #default>
                KANBAN
            </template>
        </n-tab-pane>
        <n-tab-pane name="Calendar" tab="Calendar">
            <template #tab>
                <n-icon size="16" class="tab-icon">
                    <IconCalendarWeek />
                </n-icon>
                Calendar
            </template>
            <template #default>
                CALENDAR
            </template>
        </n-tab-pane>
    </n-tabs>
</template>

<style lang="css" scoped>

    .tab-icon {
        margin-right: 4px;
    }
</style>