<script setup lang="ts">
    import { ref, reactive, watch, computed, type CSSProperties, nextTick, onMounted } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NCard, NForm, NFormItem, NInput, NButton, NButtonGroup, NIcon, type InputInst, NFlex, NEllipsis, NGrid, NGridItem } from 'naive-ui';

    import { useLoadingStore } from '../../../../stores/loading';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../../shared/types/ajaxState';
    import { taskService } from '../../services/task.ts';
    import { handleAPIError } from '../../../../api/client/errorHandler';
    import { appBus } from '../../../../shared/composables/bus';
    import type { TaskResponse, UpdateRequest } from '../../types/dto';

    import type { FormMode } from '../../../../shared/types/form-mode';
    import { Task, MAX_SUMMARY_LENGTH } from "../../models/tasks.ts";
    import TaskPrioritySelector from '../../../task-priorities/components/TaskPrioritySelector.vue';
    import TaskStatusSelector from '../../../task-statuses/components/TaskStatusSelector.vue';
    import AvatarUserName from '../../../../shared/components/AvatarUserName.vue';
    import { IconX, IconCheck, IconDeviceFloppy } from '@tabler/icons-vue';
    import { useMarkdown } from "../../../../shared/composables/useMarkdown.ts";
    import ToggleInput from '../../../../shared/components/ToggleInput.vue';
    import ToggleDateTimePicker from '../../../../shared/components/ToggleDateTimePicker.vue';
    import { IDate } from '../../../../shared/types/idate.ts';
    import ToggleTagSelector from '../../../../shared/components/ToggleTagSelector.vue';
    import EstimatedTimeInput from '../../../../shared/components/forms/EstimatedTimeInput.vue';
    import TaskSpentEstimatedPercent from '../../../../shared/components/progress/TaskSpentEstimatedPercent.vue';

    interface TaskMetadataTabProps {
        mode: FormMode;
        style?: string | CSSProperties;
        disabled?: boolean;
        projectId: string;
        taskId: string;
    }

    const task = ref<Task>(new Task());

    const props = defineProps<TaskMetadataTabProps>();

    const emit = defineEmits(["save"]);

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const serverErrors = ref<Record<string, string>>({});

    const { t } = useI18n();
    const loadingStore = useLoadingStore();
    const { render, toMarkdown } = useMarkdown();

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    const noteCount = defineModel<number>("noteCount", { default: 0 });
    const attachmentCount = defineModel<number>("attachmentCount", { default: 0 });
    const historyOperationCount = defineModel<number>("historyOperationCount", { default: 0 });
    const timeTrackingCount = defineModel<number>("timeTrackingCount", { default: 0 });

    const descriptionEditMode = ref<boolean>(false);

    const descriptionExpanded = ref<boolean>(true);

    const htmlMarkDownDescriptionPreview = computed(() => render(task.value.description ?? ""));

    const onGet = async (projectId: string, taskId: string) => {
        serverErrors.value = {};
        let notFoundError = false;
        let deletedError = false;
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: TaskResponse = await taskService.get(projectId, taskId);
            if (response.id === taskId) {
                task.value = new Task(response);
                noteCount.value = task.value.notesCount;
                attachmentCount.value = task.value.attachmentsCount;
                historyOperationCount.value = task.value.historyOperationsCount;
                timeTrackingCount.value = task.value.timeTrackingsCount;
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
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage, denyCloseDialog: notFoundError || deletedError } });
            }
        }
    };

    const onUpdate = async () => {
        serverErrors.value = {};
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                id: task.value.id ?? "",
                summary: task.value.summary ?? "",
                description: task.value.description,
                priority: {
                    id: task.value.priority.id ?? ""
                },
                status: {
                    id: task.value.status.id ?? ""
                },
                startedAt: task.value.startedAt?.msTimestamp ?? null,
                finishedAt: task.value.finishedAt?.msTimestamp ?? null,
                dueAt: task.value.dueAt?.msTimestamp ?? null,
                estimatedTime: task.value.estimatedTime ?? 0,
                tags: task.value.tags,
            };
            const response: TaskResponse = await taskService.update(props.projectId, payload);
            if (response.id === task.value.id) {
                task.value = new Task(response);
                noteCount.value = task.value.notesCount;
                attachmentCount.value = task.value.attachmentsCount;
                historyOperationCount.value = task.value.historyOperationsCount;
            } else {
                state.ajaxErrorMessage = t("modules.task.components.TaskPage.errors.updateError");
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectPage.onUpdate" } });
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.task.components.TaskPage.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.task.components.TaskPage.errors.updateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.task.components.TaskPage.errors.updateError");
                    console.error("Unhandled API error", { file: "TaskMetadataTab.vue", method: "onUpdate" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }

        }
    };

    const descriptionRef = ref<InputInst | null>(null);

    const onToggleDescriptionMode = () => {
        descriptionEditMode.value = !descriptionEditMode.value;
        if (descriptionEditMode.value) {
            nextTick(() => {
                descriptionRef.value?.focus();
            });
        }
    };

    const onFillEmptyStartDate = () => {
        if (!task.value.startedAt.hasValue()) {
            task.value.startedAt = new IDate(new Date().getTime())
        }
    };

    const onSetStartDate = () => {
        task.value.startedAt = new IDate(new Date().getTime())
    };

    const onFillEmptyFinishDate = () => {
        if (!task.value.finishedAt.hasValue()) {
            task.value.finishedAt = new IDate(new Date().getTime())
        }
    };

    const onSetFinishDate = () => {
        if (!task.value.finishedAt) {
            task.value.finishedAt = new IDate(new Date().getTime())
        }
    };

    const insertAtCursor = (value: string) => {
        const el = document.activeElement as HTMLTextAreaElement
        if (!el) {
            task.value.description += value
            return
        }

        const start = el.selectionStart ?? task.value.description?.length
        const end = el.selectionEnd ?? task.value.description?.length

        task.value.description =
            task.value.description?.slice(0, start) +
            value +
            task.value.description?.slice(end)

        // restore cursor
        requestAnimationFrame(() => {
            el.selectionStart = el.selectionEnd = start + value.length
        })
    }

    const onPaste = (e: ClipboardEvent) => {
        const clipboard = e.clipboardData
        if (!clipboard) return

        const html = clipboard.getData('text/html')
        const plain = clipboard.getData('text/plain')

        let markdown = plain

        if (html) {
            markdown = toMarkdown(html)
        }

        e.preventDefault()

        insertAtCursor(markdown)
    };

    onMounted(() => {
        if (props.projectId && props.taskId) {
            onGet(props.projectId, props.taskId);
        }
    });

</script>

<template>
    <!-- TODO: add missing i18n labels -->
    <n-card bordered :style="props.style">
        <n-form-item label="Created by">
            <div class="note-user">
                <AvatarUserName :user-id="task.createdBy.id" :user-name="task.createdBy.name" />
            </div>
        </n-form-item>
        <n-flex>
            <n-form-item label="Created at">
                <span class="doneo-datetime-label-readonly">
                    {{ task.createdAt.toLocaleString() }}
                </span>
            </n-form-item>
            <n-form-item label="Updated at">
                <span class="doneo-datetime-label-readonly">
                    {{ task.updatedAt?.toLocaleString() }}
                </span>
            </n-form-item>
            <n-form-item label="Started at">
                <ToggleDateTimePicker clearable v-model:value="task.startedAt.msTimestamp" :disabled="props.disabled" />
            </n-form-item>
            <n-form-item label="Finished at">
                <ToggleDateTimePicker clearable v-model:value="task.finishedAt.msTimestamp"
                    :disabled="props.disabled" />
            </n-form-item>
            <n-form-item label="Due at">
                <ToggleDateTimePicker clearable v-model:value="task.dueAt.msTimestamp" :disabled="props.disabled" />
            </n-form-item>
        </n-flex>
        <n-form>
            <n-flex>
                <n-form-item label="Priority">
                    <TaskPrioritySelector v-model:id="task.priority.id" :disabled="props.disabled" />
                </n-form-item>
                <n-form-item label="Status">
                    <TaskStatusSelector v-model:id="task.status.id" :disabled="props.disabled"
                        @fill-empty-start-date="onFillEmptyStartDate" @set-start-date="onSetStartDate"
                        @fill-empty-finish-date="onFillEmptyFinishDate" @set-finish-date="onSetFinishDate" />
                </n-form-item>
            </n-flex>
            <n-form-item label="Summary">
                <ToggleInput v-model:value="task.summary" show-count :max-length="MAX_SUMMARY_LENGTH"
                    :disabled="props.disabled" />
            </n-form-item>
            <n-form-item label="description">
                <template #label>
                    <n-flex align="center">
                        <span>Description</span>
                    </n-flex>
                </template>
                <div v-if="descriptionEditMode" style="width: 100%;">
                    <n-input v-model:value="task.description" type="textarea" clearable :disabled="props.disabled"
                        @paste="onPaste" ref="descriptionRef" :rows="8" />
                    <n-flex justify="end">
                        <n-button-group>
                            <n-button @click="onToggleDescriptionMode" :disabled="props.disabled">
                                <template #icon>
                                    <n-icon :component="IconCheck" />
                                </template>
                            </n-button>
                            <n-button @click="onToggleDescriptionMode" :disabled="props.disabled">
                                <template #icon>
                                    <n-icon :component="IconX" />
                                </template>
                            </n-button>
                        </n-button-group>
                    </n-flex>
                </div>
                <div v-else v-html="htmlMarkDownDescriptionPreview"
                    class="doneo-project-description-markdown-preview doneo-cursor-pointer"
                    :class="{ 'doneo-project-description-markdown-preview-expanded': descriptionExpanded }"
                    @click="onToggleDescriptionMode" />
                <!-- TODO: test alternatives -->
                <n-ellipsis v-if="false" expand-trigger="click" line-clamp="4" :tooltip="false" class="ellipsis"
                    v-html="htmlMarkDownDescriptionPreview">
                </n-ellipsis>
            </n-form-item>
            <n-form-item label="Tags">
                <ToggleTagSelector v-model:value="task.tags" />
            </n-form-item>
            <n-grid :cols="2" :x-gap="8">
                <n-grid-item>
                    <TaskSpentEstimatedPercent type="line" :estimated="task.estimatedTime" :spent="task.totalSpentTime"
                        :height="24" />
                </n-grid-item>
                <n-grid-item>
                    <EstimatedTimeInput v-model:seconds="task.estimatedTime" />
                </n-grid-item>

            </n-grid>
        </n-form>
        <n-button @click="onUpdate" :disabled="props.disabled">
            <template #icon>
                <n-icon :component="IconDeviceFloppy" color="red"></n-icon>
            </template>
            {{ t("shared.buttons.Save.label") }}
        </n-button>
    </n-card>
</template>

<style lang="css" scoped>
    .doneo-project-description-markdown-preview {
        width: 100%;
        border: 1px solid #e0e0e6;
        border-radius: var(--n-border-radius);
        padding: 4px 12px;
        color: var(--n-text-color);
        min-height: 1.5em;
        overflow: hidden;
        max-height: 12em;
        transition: max-height 0.3s ease;
    }

    .doneo-project-description-markdown-preview-expanded {
        max-height: unset;
    }
</style>