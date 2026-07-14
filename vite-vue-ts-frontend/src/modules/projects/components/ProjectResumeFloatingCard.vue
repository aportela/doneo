<script setup lang="ts">
    import { ref, reactive, onMounted, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NDrawer, NDrawerContent, NCollapse, NCollapseItem, NButton, NButtonGroup, NTag, NSpin, NIcon, NTimeline, NTimelineItem, NDivider, NFlex, NDropdown } from 'naive-ui';
    import { IconAlertTriangle, IconCalendarBolt, IconCalendarCheck, IconCalendarDue, IconCalendarTime, IconFilter2, IconLink, IconMessage2, IconPaperclip, IconReport, IconSortDescending, IconStatusChange, IconUser } from '@tabler/icons-vue';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectService } from '../services/project.ts';
    import { handleAPIError } from '../../../api/client/errorHandler.ts';
    import { appBus } from '../../../shared/composables/bus.ts';

    import { Project } from '../models/project';
    import { type ProjectResponse } from '../types/dto';

    import { getNaiveUITagColorProperty } from '../../../shared/composables/color.ts';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import ToggleMarkDownEditor from '../../../shared/components/form-blocks/ToggleMarkDownEditor.vue';
    import NoteItem from '../../notes/components/NoteItem.vue';
    import { IDate } from '../../../shared/types/idate.ts';
    import { Note } from '../../notes/models/note.ts';

    interface IProps {
        projectId: string;
    }

    const props = defineProps<IProps>();

    const { t } = useI18n();

    const project = ref<Project>(new Project());

    const show = defineModel<boolean>("show", { default: false });

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
            if (response.id === projectId) {
                project.value = new Project(response);
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
        onGet(props.projectId);
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

    const hasDetails = computed(() => project.value.attachmentsCount > 0);
</script>

<template>
    <n-drawer v-model:show="show" :width="768" placement="right">
        <n-drawer-content :native-scrollbar="false">
            <template #header>
                Project {{ project.slug }}
            </template>
            <n-spin :show="state.ajaxRunning">
                <div v-show="!state.ajaxRunning">
                    <h3>{{ project.summary }}</h3>

                    <p>Description: {{ project.description }}</p>

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
                    <n-divider v-if="hasDetails" />
                    <n-collapse class="doneo-disable-user-select" v-if="hasDetails">
                        <n-collapse-item title="Attachments" key="attachments" v-if="project.attachmentsCount > 0">
                            <template #header>
                                <n-icon :component="IconPaperclip" /> Attachments
                            </template>
                            <template #header-extra>
                                ({{ project.attachmentsCount }})
                            </template>
                        </n-collapse-item>
                        <n-collapse-item title="Relations" key="relations" v-if="false">
                            <template #header>
                                <n-icon :component="IconLink" /> Relations
                            </template>
                            <template #header-extra>
                                (0)
                            </template>
                        </n-collapse-item>
                    </n-collapse>
                    <n-divider />
                    <h4>Properties</h4>
                    <p><n-icon :component="IconStatusChange" /> State:
                        <n-tag size="tiny" :color="getNaiveUITagColorProperty(project.status.hexColor ?? '#888888')">{{
                            project.status.name }}</n-tag>
                    </p>
                    <p><n-icon :component="IconAlertTriangle" />Priority:
                        <n-tag size="tiny"
                            :color="getNaiveUITagColorProperty(project.priority.hexColor ?? '#888888')">{{
                                project.priority.name }}</n-tag>
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

                    <ToggleMarkDownEditor v-model:value="noteBody" hide-preview placeholder="Add comment" />
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
                                <NoteItem :note="defaultNote" />
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