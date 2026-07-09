<script setup lang="ts">
    import { ref, reactive, onMounted } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NDrawer, NDrawerContent, NCollapse, NCollapseItem, NButton, NButtonGroup, NTag, NSpin, NIcon, NTimeline, NTimelineItem, NInput, NDivider, NFlex, NCard, NDropdown } from 'naive-ui';
    import { IconAlertTriangle, IconAlignCenter, IconAlignLeft, IconAlignRight, IconBold, IconCalendarBolt, IconCalendarCheck, IconCalendarDue, IconCalendarTime, IconDots, IconFilter2, IconItalic, IconLink, IconMessage2, IconPaperclip, IconReport, IconSortDescending, IconStatusChange, IconStrikethrough, IconUnderline, IconUser } from '@tabler/icons-vue';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { taskService } from '../services/task.ts';
    import { handleAPIError } from '../../../api/client/errorHandler.ts';
    import { appBus } from '../../../shared/composables/bus.ts';

    import { Task } from '../models/tasks';
    import { type TaskResponse } from '../types/dto';

    import TaskTimeProgress from '../../../shared/components/progress/TaskTimeProgress.vue';
    import { getNaiveUITagColorProperty } from '../../../shared/composables/color.ts';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';

    interface IProps {
        projectId: string;
        taskId: string;
    }

    const props = defineProps<IProps>();

    const { t } = useI18n();

    const task = ref<Task>(new Task());

    const show = defineModel<boolean>("show", { default: false });

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const serverErrors = ref<Record<string, string>>({});

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
        onGet(props.projectId, props.taskId);
    });

    const noteOpts = [
        {
            label: "Update",
            key: "update",
        },
        {
            label: "Delete",
            key: "delete"
        }
    ];

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
</script>

<template>
    <n-drawer v-model:show="show" :width="768" placement="right">
        <n-drawer-content :native-scrollbar="false">
            <template #header>
                Task {{ task.slug }}
            </template>
            <n-spin :show="state.ajaxRunning">
                <div v-show="!state.ajaxRunning">
                    <h3>{{ task.summary }}</h3>

                    <p>Description: {{ task.description }}</p>
                    <task-time-progress :estimated="task.estimatedTime" :spent="task.totalSpentTime" />

                    <n-button-group size="tiny" style="margin-top: 16px;">
                        <n-button round><template #icon><n-icon :component="IconPaperclip" /></template> Add
                            attachment</n-button>
                        <n-button round><template #icon><n-icon :component="IconLink" /></template> Add
                            relation</n-button>
                        <n-button round><template #icon><n-icon :component="IconReport" /></template> Add time
                            tracking</n-button>
                        <n-button round><template #icon><n-icon :component="IconMessage2" tag="a"
                                    href="#aa" /></template> Add
                            note</n-button>
                    </n-button-group>
                    <n-divider />
                    <n-collapse class="doneo-disable-user-select">
                        <n-collapse-item title="Attachments" key="attachments">
                            <template #header>
                                <n-icon :component="IconPaperclip" /> Attachments
                            </template>
                            <template #header-extra>
                                ({{ task.attachmentsCount }})
                            </template>
                        </n-collapse-item>
                        <n-collapse-item title="Relations" key="relations">
                            <template #header>
                                <n-icon :component="IconLink" /> Relations
                            </template>
                            <template #header-extra>
                                (0)
                            </template>
                        </n-collapse-item>
                        <n-collapse-item title="Time trackings" key="timetrackings">
                            <template #header>
                                <n-icon :component="IconReport" /> Time trackings
                            </template>
                            <template #header-extra>
                                ({{ task.timeTrackingsCount }})
                            </template>
                        </n-collapse-item>
                    </n-collapse>
                    <n-divider />
                    <h4>Properties</h4>
                    <p><n-icon :component="IconStatusChange" /> State:
                        <n-tag size="tiny" :color="getNaiveUITagColorProperty(task.status.hexColor ?? '#888888')">{{
                            task.status.name }}</n-tag>
                    </p>
                    <p><n-icon :component="IconAlertTriangle" />Priority:
                        <n-tag size="tiny" :color="getNaiveUITagColorProperty(task.priority.hexColor ?? '#888888')">{{
                            task.priority.name }}</n-tag>
                    </p>
                    <p><n-icon :component="IconUser" />Asignee: John doe</p>
                    <p><n-icon :component="IconCalendarBolt" /> Created at: {{ task.createdAt.toLocaleString() }} </p>
                    <p><n-icon :component="IconCalendarTime" /> Started at: {{ task.startedAt.toLocaleString() }} </p>
                    <p><n-icon :component="IconCalendarCheck" /> Finished at: {{ task.finishedAt.toLocaleString() }}
                    </p>
                    <p><n-icon :component="IconCalendarDue" /> Due at: {{ task.dueAt.toLocaleString() }} </p>
                    <p><n-flex align="center">
                            Created by:
                            <AvatarUserName :user-id="task.createdBy.id" :user-name="task.createdBy.name" />
                        </n-flex>
                    </p>
                    <p>Tags: <n-tag v-for="tag in task.tags" :key="tag" size="small" class="mr-tiny">{{ tag
                    }}</n-tag>
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
                    <div class="editor">
                        <n-input type="textarea" placeholder="Add comment" />
                        <div class="toolbar">
                            <n-icon :size="20" :component="IconBold" />
                            <n-icon :size="20" :component="IconItalic" />
                            <n-icon :size="20" :component="IconUnderline" />
                            <n-icon :size="20" :component="IconStrikethrough" />
                            <n-divider vertical />
                            <n-icon :size="20" :component="IconAlignLeft" />
                            <n-icon :size="20" :component="IconAlignRight" />
                            <n-icon :size="20" :component="IconAlignCenter" />
                        </div>
                    </div>
                    <n-timeline>
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
                                <n-card>
                                    <template #header>
                                        <n-flex align="center">
                                            <AvatarUserName :user-id="task.createdBy.id"
                                                :user-name="task.createdBy.name" />
                                            commented in less than a minute
                                        </n-flex>
                                    </template>
                                    <template #header-extra>
                                        <n-dropdown :options="noteOpts" trigger="click">
                                            <n-button size="tiny">
                                                <template #icon>
                                                    <n-icon :component="IconDots" />
                                                </template>
                                            </n-button>
                                        </n-dropdown>
                                    </template>
                                    <template #default>
                                        Summary: Optimize existing workflows with the goal of improving customer
                                        experience to improve time-to-market.

                                    </template>
                                </n-card>
                            </template>
                        </n-timeline-item>
                    </n-timeline>
                </div>
            </n-spin>
        </n-drawer-content>
    </n-drawer>
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