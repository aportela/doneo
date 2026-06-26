<script setup lang="ts">
    import { ref, reactive, computed, watch, nextTick, onMounted } from 'vue';
    import { useI18n } from 'vue-i18n';
    import { useRoute, useRouter } from 'vue-router';

    import { NTabs, NTabPane, type TabsInst } from 'naive-ui';

    import { useLoadingStore } from '../../../stores/loading.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { taskService } from '../services/task.ts';
    import { handleAPIError } from '../../../api/client/errorHandler.ts';
    import { appBus } from '../../../shared/composables/bus.ts';

    import TaskMetadataTab from '../components/tabs/Metadata.vue';
    import TaskAttachmentsTab from '../components/tabs/Attachments.vue';
    import TaskNotesTab from '../components/tabs/Notes.vue';
    import TaskTimeTrackingsTab from '../components/tabs/TimeTrackings.vue';
    import TaskHistoryTab from '../components/tabs/History.vue';
    import { Task } from '../models/tasks.ts';
    import type { TaskResponse } from '../types/dto.ts';

    const { t } = useI18n();
    const loadingStore = useLoadingStore();
    const route = useRoute();
    const router = useRouter();

    const projectId = route.params.projectId as string ?? ""
    const taskId = route.params.taskId as string ?? ""
    const task = ref<Task>(new Task());

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const serverErrors = ref<Record<string, string>>({});

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

    const attachmentsTabLabel = computed(() => t("modules.task.components.TaskPage.tabs.attachments.label", task.value.attachmentsCount));
    const notesTabLabel = computed(() => t("modules.task.components.TaskPage.tabs.notes.label", task.value.notesCount));
    const timeTrackingsTabLabel = computed(() => t("modules.task.components.TaskPage.tabs.timeTrackings.label", task.value.timeTrackingsCount));
    const historyTabLabel = computed(() => t("modules.task.components.TaskPage.tabs.history.label", task.value.historyOperationsCount));

    const tabsRef = ref<TabsInst>();

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    // recalc bar position on dynamic tab labels changes
    watch(
        () => [attachmentsTabLabel.value, notesTabLabel.value, timeTrackingsTabLabel.value, historyTabLabel.value],
        async () => {
            await nextTick();
            tabsRef.value?.syncBarPosition();
        }
    );

    const onGet = async (projectId: string, taskId: string) => {
        serverErrors.value = {};
        let notFoundError = false;
        let deletedError = false;
        let accessDeniedError = false;
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: TaskResponse = await taskService.get(projectId, taskId);
            if (response.id === taskId) {
                task.value = new Task(response);
            } else {
                state.ajaxErrorMessage = t("modules.task.components.TaskPage.errors.loadError");
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectPage.onGet" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            accessDeniedError = true;
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.task.components.TaskPage.errors.notFoundError");
                            notFoundError = true;
                            break;
                        case 410:
                            state.ajaxErrorMessage = t("modules.task.components.TaskPage.errors.deletedError");
                            deletedError = true;
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.task.components.TaskPage.errors.loadError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.task.components.TaskPage.errors.loadError");
                    console.error("Unhandled API error", { file: "TaskMetadataTab.vue", method: "onGet" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage, denyCloseDialog: notFoundError || deletedError || accessDeniedError } });
            }
        }
    };

    onMounted(() => {
        if (projectId && taskId) {
            onGet(projectId, taskId);
        }
    });

</script>

<template>
    <n-tabs placement="top" type="line" ref="tabsRef" animated v-model:value="tab">
        <n-tab-pane name="metadata" display-directive="show" key="metadata" :disabled="!projectId || !taskId">
            <template #tab>
                Metadata
            </template>
            <TaskMetadataTab v-if="task.id" v-model:task="task" :read-only="!task.allowedOperations.updateTask" />
        </n-tab-pane>
        <n-tab-pane name="notes" :tab="notesTabLabel" display-directive="show:lazy" key="notes"
            :disabled="!projectId || !taskId || (!task.allowedOperations.updateTask)">
            <TaskNotesTab v-if="projectId && taskId" :project-id="projectId" :task-id="taskId"
                v-model:item-count="task.notesCount" />
        </n-tab-pane>
        <n-tab-pane name="attachments" :tab="attachmentsTabLabel" display-directive="show:lazy" key="attachments"
            :disabled="!projectId || !taskId || (!task.allowedOperations.updateTask)">
            <TaskAttachmentsTab v-if="projectId && taskId" :project-id="projectId" :task-id="taskId"
                v-model:item-count="task.attachmentsCount" />
        </n-tab-pane>
        <n-tab-pane name="timetrackings" :tab="timeTrackingsTabLabel" display-directive="show:lazy" key="timetrackings"
            :disabled="!projectId || !taskId || (!task.allowedOperations.updateTask)">
            <TaskTimeTrackingsTab v-if="projectId && taskId" :project-id="projectId" :task-id="taskId"
                v-model:item-count="task.timeTrackingsCount" />
        </n-tab-pane>
        <n-tab-pane name="history" :tab="historyTabLabel" display-directive="show:lazy" key="history"
            :disabled="!projectId || !taskId || (!task.allowedOperations.updateTask)">
            <TaskHistoryTab v-if="projectId && taskId" :project-id="projectId" :task-id="taskId"
                :key="task.historyOperationsCount" />
        </n-tab-pane>
    </n-tabs>
</template>

<style lang="css" scoped>
    .avatar {
        margin-right: 4px;
    }
</style>