<script setup lang="ts">
    import { ref, computed, watch, nextTick } from 'vue';
    import { useI18n } from 'vue-i18n';
    import { useRoute, useRouter } from 'vue-router';

    import { NTabs, NTabPane, type TabsInst } from 'naive-ui';

    // TODO: move to tabs
    import TaskMetadataTab from '../components/TaskMetadataTab.vue';
    import TaskAttachmentsTab from '../components/tabs/Attachments.vue';
    import TaskNotesTab from '../components/tabs/Notes.vue';
    import TaskTimeTrackingsTab from '../components/tabs/TimeTrackings.vue';
    import TaskHistoryTab from '../components/tabs/History.vue';

    const { t } = useI18n();
    const route = useRoute();
    const router = useRouter();

    const projectId = route.params.projectId as string ?? ""
    const taskId = route.params.taskId as string ?? ""

    // TODO: set tab with type (type tab = "metadata" | "notes"....)
    const tab = computed({
        get: () => route.params.tab as string,
        set: (value: string) => {
            router.push({
                name: 'taskTab',
                params: {
                    projectId: route.params.projectId,
                    taskId: route.params.taskId,
                    tab: value
                }
            });
        }
    });

    const noteCount = ref<number>(0);
    const attachmentCount = ref<number>(0);
    const timeTrackingCount = ref<number>(0)
    const historyOperationCount = ref<number>(0);

    const attachmentsTabLabel = computed(() => t("modules.task.components.TaskPage.tabs.attachments.label", attachmentCount.value));
    const notesTabLabel = computed(() => t("modules.task.components.TaskPage.tabs.notes.label", noteCount.value));
    const timeTrackingsTabLabel = computed(() => t("modules.task.components.TaskPage.tabs.timeTrackings.label", timeTrackingCount.value));
    const historyTabLabel = computed(() => t("modules.task.components.TaskPage.tabs.history.label", historyOperationCount.value));

    const tabsRef = ref<TabsInst>();

    // recalc bar position on dynamic tab labels changes
    watch(
        () => [attachmentsTabLabel.value, notesTabLabel.value, historyTabLabel.value],
        async () => {
            await nextTick();
            tabsRef.value?.syncBarPosition();
        }
    );
</script>

<template>
    <n-tabs placement="top" type="line" ref="tabsRef" animated v-model:value="tab">
        <n-tab-pane name="metadata" display-directive="show" key="metadata" :disabled="!projectId || !taskId">
            <template #tab>
                Metadata
            </template>
            <!-- TODO: mode ???-->
            <TaskMetadataTab mode="add" :project-id="projectId" :task-id="taskId" v-model:note-count="noteCount"
                v-model:attachment-count="attachmentCount" v-model:history-operation-count="historyOperationCount"
                v-model:time-tracking-count="timeTrackingCount" />
        </n-tab-pane>
        <n-tab-pane name="notes" :tab="notesTabLabel" display-directive="show:lazy" key="notes"
            :disabled="!projectId || !taskId">
            <TaskNotesTab v-if="projectId && taskId" :project-id="projectId" :task-id="taskId"
                v-model:item-count="noteCount" />
        </n-tab-pane>
        <n-tab-pane name="attachments" :tab="attachmentsTabLabel" display-directive="show:lazy" key="attachments"
            :disabled="!projectId || !taskId">
            <TaskAttachmentsTab v-if="projectId && taskId" :project-id="projectId" :task-id="taskId"
                v-model:item-count="attachmentCount" />
        </n-tab-pane>
        <n-tab-pane name="timetrackings" :tab="timeTrackingsTabLabel" display-directive="show:lazy" key="timetrackings"
            :disabled="!projectId || !taskId">
            <TaskTimeTrackingsTab v-if="projectId && taskId" :project-id="projectId" :task-id="taskId"
                v-model:item-count="timeTrackingCount" />
        </n-tab-pane>
        <n-tab-pane name="history" :tab="historyTabLabel" display-directive="show:lazy" key="history"
            :disabled="!projectId || !taskId">
            <TaskHistoryTab v-if="projectId && taskId" :project-id="projectId" :task-id="taskId" />
        </n-tab-pane>
    </n-tabs>
</template>

<style lang="css" scoped>
    .avatar {
        margin-right: 4px;
    }
</style>