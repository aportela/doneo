<script setup lang="ts">
    import { ref, reactive, onMounted, computed, shallowRef } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NCollapse, NCollapseItem, NButton, NButtonGroup, NTag, NSpin, NIcon, NTimeline, NTimelineItem, NDivider, NFlex, NDropdown } from 'naive-ui';
    import { IconAlertTriangle, IconBookmark, IconCalendarBolt, IconCalendarCheck, IconCalendarDue, IconCalendarTime, IconFilter2, IconMessage2, IconPaperclip, IconSortDescending, IconStatusChange, IconUser } from '@tabler/icons-vue';
    import { ListTodo } from '@lucide/vue';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectService } from '../services/project.ts';
    import { handleAPIError } from '../../../api/client/errorHandler.ts';
    import { appBus } from '../../../shared/composables/bus.ts';

    import { Project } from '../models/project';
    import { type ProjectResponse } from '../types/dto';

    import { getNaiveUITagColorProperty } from '../../../shared/composables/naive-ui-helpers.ts';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import ToggleMarkDownEditor from '../../../shared/components/form-blocks/ToggleMarkDownEditor.vue';
    import NoteItem from '../../notes/components/NoteItem.vue';
    import { IDate } from '../../../shared/types/idate.ts';
    import { Note } from '../../notes/models/note.ts';
    import type { SearchResponse } from '../../attachments/types/dto.ts';
    import { Attachment } from '../../attachments/models/attachment.ts';
    import { attachmentService } from '../../attachments/services/attachment.ts';
    import { formatBytes } from '../../../shared/composables/format.ts';
    import { DONEO_ICON_ACTION_DOWNLOAD, DONEO_ICON_ACTION_OPEN, DONEO_ICON_ACTION_PREVIEW } from '../../../shared/types/icons.ts';
    import { Task } from '../../tasks/models/tasks.ts';
    import type { SearchRequest } from '../../tasks/types/dto.ts';
    import { taskService } from '../../tasks/services/task.ts';
    import { HistoryOperation } from '../../history-operations/models/history-operation.ts';
    import { historyOperationsService } from '../../history-operations/services/history-operations.ts';
    import { noteService } from '../../notes/services/note.ts';

    interface Props {
        projectId: string;
    }

    const props = defineProps<Props>();

    const { t } = useI18n();

    const project = ref<Project>(new Project());

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const serverErrors = ref<Record<string, string>>({});

    const onGet = async (projectId: string) => {
        serverErrors.value = {};
        let notFoundError = false;
        let deletedError = false;
        let accessDeniedError = false;
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: ProjectResponse = await projectService.get(projectId);
            project.value = new Project(response);
            if (project.value.attachmentsCount > 0) {
                getProjectAttachments();
            }
            if (project.value.tasksCount > 0) {
                getProjectTasks();
            }
            if (project.value.historyOperationsCount > 0) {
                getProjectHistoryOperations();
            }
            if (project.value.notesCount > 0) {
                getProjectNotes();
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

    const attachments = shallowRef<Attachment[]>([]);

    const getProjectAttachments = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: SearchResponse = await attachmentService.getProjectAttachments(props.projectId);
            attachments.value = response.attachments.map((attachment) => new Attachment(attachment));
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "AttachmentsTable.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectAttachment.components.AttachmentsTable.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectAttachment.components.AttachmentsTable.errors.refreshError");
                    console.error("Unhandled API error", { file: "AttachmentsTable.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const tasks = shallowRef<Task[]>([]);

    const getProjectTasks = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: SearchRequest = {
                pager: { enabled: false, currentPage: 1, resultsPage: 0 },
                order: { field: "creationDate", direction: "DESC" },
                filter: {
                    projectId: props.projectId,
                }
            };
            const response = await taskService.search(null, payload);
            tasks.value = response.tasks.map((task) => new Task(task));
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageTasksPage.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.project.components.ManageTasksPage.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.project.components.ManageTasksPage.errors.refreshError");
                    console.error("Unhandled API error", { file: "ManageTasksPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const historyOperations = shallowRef<HistoryOperation[]>([]);

    const getProjectHistoryOperations = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const results = await historyOperationsService.getProjectHistoryOperations(props.projectId);
            historyOperations.value = results.historyOperations.map((operation) => new HistoryOperation(operation));
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "HistoryOperationsTable.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                    console.error("Unhandled API error", { file: "HistoryOperationsTable.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const notes = shallowRef<Note[]>([]);

    const getProjectNotes = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const results = await noteService.getProjectNotes(props.projectId);
            notes.value = results.notes.map((note) => new Note(note));
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectNotes.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                    console.error("Unhandled API error", { file: "ProjectNotes.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    onMounted(() => {
        if (props.projectId) {
            onGet(props.projectId);
        }
    });

    const activityFilterOpts = [
        {
            label: "Updates",
            key: "updates",
        },
        {
            label: "Notes",
            key: "notes",
        }
    ];

    const noteBody = ref<string>("");

    const defaultNote = new Note();
    defaultNote.id = "019f4908-2f06-7e36-8751-f7c2d9a3e7c2";
    defaultNote.createdBy.id = "019f4908-5229-740c-a4f1-284d512eb4a0";
    defaultNote.createdAt = new IDate(Date.now());
    defaultNote.body = "We recommend configuring it at the project entry point, such as in main.js for projects created with Vite. Avoid calling config within components!";

    const hasDetails = computed(() => project.value.attachmentsCount > 0 || project.value.tasksCount > 0);

    type TimelineItem =
        | {
            type: 'historyOperation'
            createdAt: IDate | null
            item: HistoryOperation
        }
        | {
            type: 'note'
            createdAt: IDate | null
            item: Note
        };


    const timeline = computed<TimelineItem[]>(() => [
        ...historyOperations.value.map(h => ({
            type: 'historyOperation' as const,
            createdAt: h.createdAt,
            item: h
        })),
        ...notes.value.map(n => ({
            type: 'note' as const,
            createdAt: n.createdAt,
            item: n
        }))
    ].sort((a, b) => (b.createdAt?.msTimestamp ?? 0) - (a.createdAt?.msTimestamp ?? 0)))
</script>

<template>
    <n-spin :show="state.ajaxRunning" style="height: 100vh;">
        <div v-show="!state.ajaxRunning">
            <h3>Summary: {{ project.summary }}</h3>

            <h3>Description:</h3>
            <ToggleMarkDownEditor read-only v-model:value="project.description"
                style="max-height: 32vh; overflow-y: scroll;" />

            <n-button-group size="tiny" style="margin-top: 16px;">
                <n-button round :disabled="!project.allowedOperations.updateProject"><template #icon><n-icon
                            :component="IconPaperclip" /></template>
                    Add
                    attachment</n-button>

                <n-button round tag="a" href="#new_note" :disabled="!project.allowedOperations.updateProject"><template
                        #icon><n-icon :component="IconMessage2" /></template>
                    Add
                    note</n-button>
                <n-button round :disabled="!project.allowedOperations.addTask"><template #icon><n-icon
                            :component="ListTodo" /></template> Add
                    task</n-button>
            </n-button-group>
            <n-divider v-if="hasDetails" />
            <n-collapse class="doneo-disable-user-select" v-if="hasDetails">
                <n-collapse-item title="Attachments" key="attachments" name="attachments"
                    v-if="project.attachmentsCount > 0">
                    <template #header>
                        <n-icon :component="IconPaperclip" /> Attachments
                    </template>
                    <template #header-extra>
                        ({{ project.attachmentsCount }})
                    </template>
                    <p v-for="attachment in attachments" :key="attachment.id"
                        style="display: flex; align-items: center;">
                        <n-button size="tiny" style="margin-right: 4px;">
                            <template #icon>
                                <n-icon :component="DONEO_ICON_ACTION_DOWNLOAD" />
                            </template>
                        </n-button>
                        <n-button size="tiny" style="margin-right: 4px;">
                            <template #icon>
                                <n-icon :component="DONEO_ICON_ACTION_PREVIEW" />
                            </template>
                        </n-button>
                        <span>
                            <strong>{{ attachment.name }}</strong> ({{ formatBytes(attachment.size) }})
                        </span>
                    </p>
                </n-collapse-item>
                <n-collapse-item title="Tasks" key="tasks" name="tasks" v-if="project.tasksCount > 0">
                    <template #header>
                        <n-icon :component="ListTodo" /> Tasks
                    </template>
                    <template #header-extra>
                        ({{ project.tasksCount }})
                    </template>
                    <p v-for="task in tasks" :key="task.id" style="display: flex; align-items: center;">
                        <n-button size="tiny" style="margin-right: 4px;">
                            <template #icon>
                                <n-icon :component="DONEO_ICON_ACTION_OPEN" />
                            </template>
                        </n-button>
                        <span>
                            <strong>{{ task.slug }}</strong> - {{ task.summary }}
                        </span>
                    </p>
                </n-collapse-item>
            </n-collapse>
            <n-divider />
            <h4>Properties</h4>
            <p>
                <n-icon :component="IconBookmark" />
                Type:
                <n-tag size="small" :color="getNaiveUITagColorProperty(project.type.hexColor ?? '#888888')">
                    {{
                        project.type.name
                    }}
                </n-tag>
            </p>
            <p>
                <n-icon :component="IconStatusChange" />
                State:
                <n-tag size="small" :color="getNaiveUITagColorProperty(project.status.hexColor ?? '#888888')">
                    {{
                        project.status.name
                    }}
                </n-tag>
            </p>
            <p>
                <n-icon :component="IconAlertTriangle" />
                Priority:
                <n-tag size="small" :color="getNaiveUITagColorProperty(project.priority.hexColor ?? '#888888')">
                    {{
                        project.priority.name
                    }}
                </n-tag>
            </p>
            <p><n-icon :component="IconUser" />Asignee: John doe</p>
            <p><n-icon :component="IconCalendarBolt" /> Created at: {{ project.createdAt.toLocaleString() }}
            </p>
            <p><n-icon :component="IconCalendarTime" /> Started at: {{ project.startedAt.toLocaleString() }}
            </p>
            <p><n-icon :component="IconCalendarCheck" /> Finished at: {{ project.finishedAt.toLocaleString() }}
            </p>
            <p><n-icon :component="IconCalendarDue" /> Due at: {{ project.dueAt.toLocaleString() }} </p>
            <p><n-flex align="center">
                    Created by:
                    <AvatarUserName :user-id="project.createdBy.id" :user-name="project.createdBy.name" />
                </n-flex>
            </p>
            <n-divider />
            <n-flex justify="space-between" align="center">
                <h4>Activity</h4>
                <n-button-group size="small">
                    <n-button><template #icon><n-icon :component="IconSortDescending" /></template></n-button>
                    <n-dropdown trigger="click" :options="activityFilterOpts">
                        <n-button><template #icon><n-icon :component="IconFilter2" /></template></n-button>
                    </n-dropdown>
                </n-button-group>
            </n-flex>

            <ToggleMarkDownEditor v-model:value="noteBody" hide-preview placeholder="Add note" id="new_note" />
            <n-divider />
            <n-timeline>
                <n-timeline-item v-for="tt in timeline" type="default" :key="tt.item.id ?? ''">
                    <template #icon v-if="tt.type === 'historyOperation'">
                        <n-icon :component="IconAlertTriangle" size="20" />
                    </template>
                    <template #icon v-else>
                        <n-icon :component="IconMessage2" size="22" />
                    </template>
                    <template #default v-if="tt.type === 'historyOperation'">
                        {{ t(tt.item.getOperationTypeLabel()) }} by {{ tt.item.createdBy.name }} on {{
                            tt.item.createdAt.toLocaleString() }}
                    </template>
                    <template #default v-else>
                        <NoteItem :note="tt.item" />
                    </template>
                </n-timeline-item>
                <!--
                <n-timeline-item type="default">
                    <template #icon>
                        <n-icon :component="IconAlertTriangle" size="20" />
                    </template>
                    <template #default>John doe created task about 1 hour ago</template>
                </n-timeline-item>
                <n-timeline-item type="default">
                    <template #icon>
                        <n-icon :component="IconMessage2" size="22" />
                    </template>
                    <template #default>
                        <NoteItem :note="defaultNote" />
                    </template>
                </n-timeline-item>
            -->
            </n-timeline>
        </div>
    </n-spin>
</template>

<style lang="css" scoped>

    div.editor {
        width: 100%;
        padding: 8px;
        margin-bottom: 16px;
        border: 1px solid #e4e4e4;
        border-radius: 4px;
        color: #888;
    }

    div.toolbar {
        width: 98%;
        margin: 0px auto;
        padding: 8px;
        background: #f4f4f4;
        border: 1px solid #e4e4e4;
    }

    .mr-tiny {
        margin-right: 4px;
    }
</style>