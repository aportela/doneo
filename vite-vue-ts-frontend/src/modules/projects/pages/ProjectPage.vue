<script setup lang="ts">
    import { ref, reactive, computed, nextTick, onMounted } from 'vue';
    import { useI18n } from "vue-i18n";
    import { useRoute, useRouter } from 'vue-router';

    import { NTabs, NTabPane } from 'naive-ui';

    import { Project } from "../models/project";
    import { ProjectType } from '../../project-types/models/project-type';
    import { ProjectStatus } from '../../project-statuses/models/project-status';
    import { ProjectPriority } from '../../project-priorities/models/project-priority';
    import { UserBase } from '../../users/models/user';

    import ProjectMetadataForm from '../components/ProjectMetadataForm.vue';
    import ProjectTasks from '../components/ProjectTasks.vue';
    import ProjectPermissions from '../components/ProjectPermissions.vue';
    import ProjectAttachments from '../components/ProjectAttachments.vue';
    import ProjectNotes from '../components/ProjectNotes.vue';
    import ProjectHistoryOperations from '../components/ProjectHistoryOperations.vue';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectService } from '../services/project';
    import { handleAPIError } from '../../../api/client/errorHandler';
    //import RemoteAPIAlert from '../../../shared/components/alerts/RemoteAPIAlert.vue';
    import { appBus } from '../../../shared/composables/bus';
    import type { ProjectResponse } from '../types/dto';

    const { t } = useI18n();
    const route = useRoute();
    const router = useRouter();

    const projectId = route.params.id as string

    const project = ref<Project>(new Project(
        {
            id: "",
            key: "",
            summary: "",
            description: "",
            type: new ProjectType({
                id: "",
                name: "",
                hexColor: "",
            }),
            priority: new ProjectPriority({
                id: "",
                name: "",
                hexColor: "",
            }),
            status: new ProjectStatus({
                id: "",
                name: "",
                hexColor: "",
            }),
            createdAt: 0,
            createdBy: new UserBase({
                id: "",
                name: "",
            }),
            tasksCount: 0,
            permissionsCount: 0,
            attachmentsCount: 0,
            notesCount: 0,
            historyOperationsCount: 0,
        }
    ));

    const tab = computed({
        get: () => route.params.tab as string,
        set: (value: string) => {
            router.push({
                name: 'project-tab',
                params: {
                    id: route.params.id,
                    tab: value
                }
            });
        }
    });

    const tasksTabLabel = computed(() => project.value.tasksCount > 0 ? `Tasks (${project.value.tasksCount})` : 'Tasks')
    const permissionsTabLabel = computed(() => project.value.permissionsCount > 0 ? `Permissions (${project.value.permissionsCount})` : 'Permissions')
    const attachmentsTabLabel = computed(() => project.value.attachmentsCount > 0 ? `Attachments (${project.value.attachmentsCount})` : 'Attachments')
    const notesTabLabel = computed(() => project.value.notesCount > 0 ? `Notes (${project.value.notesCount})` : 'Notes')
    const historyTabLabel = computed(() => project.value.historyOperationsCount > 0 ? `History (${project.value.historyOperationsCount})` : 'History')

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const serverErrors = ref<Record<string, string>>({});

    const onGet = async (id: string) => {
        serverErrors.value = {};
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: ProjectResponse = await projectService.get(id);
            if (response.id === id) {
                project.value = new Project(response);
            } else {
                state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.loadError");
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectPriorityForm.onGet" } });
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.loadError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectPriority.components.ProjectPriorityForm.errors.loadError");
                    console.error("Unhandled API error", { file: "ProjectPriorityForm.vue", method: "onGet" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                await nextTick();
            }
        }
    };

    onMounted(() => {
        if (projectId) {
            onGet(projectId);
        }
    });

</script>

<template>
    <n-tabs placement="top" type="line" animated v-model:value="tab">
        <n-tab-pane name="metadata" tab="Metadata" display-directive="show:lazy">
            <ProjectMetadataForm mode="add" :project-id="projectId" v-model:project="project" />
        </n-tab-pane>
        <n-tab-pane name="permissions" :tab="permissionsTabLabel" display-directive="show:lazy">
            <ProjectPermissions :project-id="project.id" />
        </n-tab-pane>
        <n-tab-pane name="notes" :tab="notesTabLabel" display-directive="show:lazy">
            <ProjectNotes :project-id="project.id" />
        </n-tab-pane>
        <n-tab-pane name="attachments" :tab="attachmentsTabLabel" display-directive="show:lazy">
            <ProjectAttachments :project-id="project.id" />
        </n-tab-pane>
        <n-tab-pane name="history" :tab="historyTabLabel" display-directive="show:lazy">
            <ProjectHistoryOperations :project-id="project.id" />
        </n-tab-pane>
        <n-tab-pane name="tasks" :tab="tasksTabLabel" display-directive="show:lazy">
            <ProjectTasks :project-id="project.id" />
        </n-tab-pane>
    </n-tabs>
</template>

<style lang="css" scoped>
    .avatar {
        margin-right: 4px;
    }
</style>