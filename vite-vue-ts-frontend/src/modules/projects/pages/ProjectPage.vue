<script setup lang="ts">
    import { ref, reactive, computed, nextTick, onMounted, onBeforeUnmount } from 'vue';
    import { useI18n } from "vue-i18n";
    import { useRoute, useRouter } from 'vue-router';

    import { NTabs, NTabPane } from 'naive-ui';

    import { Project } from "../models/project";

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
    import type { AddRequest, ProjectResponse, UpdateRequest } from '../types/dto';

    const { t } = useI18n();
    const route = useRoute();
    const router = useRouter();

    const projectId = route.params.id as string

    const project = ref<Project>(new Project());

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
                state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.loadError");
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
                            state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.loadError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.loadError");
                    console.error("Unhandled API error", { file: "ProjectPage.vue", method: "onGet" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                await nextTick();
            }
        }
    };

    const onSave = async () => {
        if (project.value.id) {
            onUpdate();
        } else {
            onAdd();
        }
    };

    const onAdd = async () => {
        serverErrors.value = {};
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                key: project.value.key ?? "",
                summary: project.value.summary ?? "",
                description: project.value.description
            };
            const response: ProjectResponse = await projectService.add(payload);
            if (response.id === project.value.id) {
                project.value = new Project(response);
            } else {
                state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.addError");
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectPage.onAdd" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.addError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.addError");
                    console.error("Unhandled API error", { file: "ProjectPage.vue", method: "onAdd" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    const onUpdate = async () => {
        serverErrors.value = {};
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                id: project.value.id ?? "",
                key: project.value.key ?? "",
                summary: project.value.summary ?? "",
                description: project.value.description
            };
            const response: ProjectResponse = await projectService.update(payload);
            if (response.id === project.value.id) {
                project.value = new Project(response);
            } else {
                state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.updateError");
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
                            state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.updateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.updateError");
                    console.error("Unhandled API error", { file: "ProjectPage.vue", method: "onUpdate" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        if (projectId) {
            onGet(projectId);
        }
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectPage.onGet")) {
                onGet(projectId);
            } else if (payload.to.includes("ProjectPage.onAdd")) {
                // TODO:
            } else if (payload.to.includes("ProjectPage.onUpdate")) {
                // TODO:
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });

    // TODO: skeleton while loading ???
    // TODO: loading error
</script>

<template>
    <n-tabs placement="top" type="line" animated v-model:value="tab">
        <n-tab-pane name="metadata" tab="Metadata" display-directive="show:lazy">
            <ProjectMetadataForm mode="add" :project-id="projectId" v-model:project="project" @save="onSave" />
        </n-tab-pane>
        <n-tab-pane name="permissions" :tab="permissionsTabLabel" display-directive="show:lazy">
            <ProjectPermissions :project-id="project.id" v-model:item-count="project.permissionsCount" />
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