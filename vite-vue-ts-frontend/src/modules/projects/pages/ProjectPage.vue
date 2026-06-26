<script setup lang="ts">
    import { ref, reactive, computed, watch, nextTick, onMounted } from 'vue';
    import { useI18n } from 'vue-i18n';
    import { useRoute, useRouter } from 'vue-router';

    import { NTabs, NTabPane, type TabsInst, NIcon } from 'naive-ui';

    import { IconAlertTriangle } from '@tabler/icons-vue';

    import { useLoadingStore } from '../../../stores/loading.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectService } from '../services/project.ts';
    import { handleAPIError } from '../../../api/client/errorHandler.ts';
    import { appBus } from '../../../shared/composables/bus.ts';

    import ProjectMetadataTab from '../components/tabs/Metadata.vue';
    import ProjectTasksTab from '../components/tabs/Tasks.vue';
    import ProjectPermissionsTab from '../components/tabs/Permissions.vue';
    import ProjectAttachmentsTab from '../components/tabs/Attachments.vue';
    import ProjectNotesTab from '../components/tabs/Notes.vue';
    import ProjectHistoryTab from '../components/tabs/History.vue';
    import type { ProjectResponse } from '../types/dto.ts';
    import { Project } from '../models/project.ts';

    const { t } = useI18n();
    const loadingStore = useLoadingStore();
    const route = useRoute();
    const router = useRouter();

    const projectId = route.params.projectId as string

    const project = ref<Project>(new Project());

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const serverErrors = ref<Record<string, string>>({});

    // TODO: set tab with type (type tab = "metadata" | "permissions"....)
    const tab = computed({
        // TODO: invalid tab route ???
        get: () => route.params.tab as string,
        set: (value: string) => {
            router.push({
                name: 'projectTab',
                params: {
                    projectId: route.params.projectId,
                    tab: value
                }
            });
        }
    });

    const permissionsTabLabel = computed(() => t("modules.project.components.ProjectPage.tabs.permissions.label", project.value.permissionsCount));
    const attachmentsTabLabel = computed(() => t("modules.project.components.ProjectPage.tabs.attachments.label", project.value.attachmentsCount));
    const notesTabLabel = computed(() => t("modules.project.components.ProjectPage.tabs.notes.label", project.value.notesCount));
    const historyTabLabel = computed(() => t("modules.project.components.ProjectPage.tabs.history.label", project.value.historyOperationsCount));
    const tasksTabLabel = computed(() => t("modules.project.components.ProjectPage.tabs.tasks.label", project.value.tasksCount));

    const tabsRef = ref<TabsInst>();

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    // recalc bar position on dynamic tab labels changes
    watch(
        () => [permissionsTabLabel.value, attachmentsTabLabel.value, notesTabLabel.value, historyTabLabel.value, tasksTabLabel.value],
        async () => {
            await nextTick();
            tabsRef.value?.syncBarPosition();
        }
    );

    const tabPlacement = ref<"top" | "left">("top");

    const onToggleTabPlacement = () => {
        if (tabPlacement.value == "top") {
            tabPlacement.value = "left";
        } else {
            tabPlacement.value = "top";
        }
    };

    const onGet = async (id: string) => {
        serverErrors.value = {};
        let notFoundError = false;
        let deletedError = false;
        let accessDeniedError = false;
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
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            accessDeniedError = true;
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.notFoundError");
                            notFoundError = true;
                            break;
                        case 410:
                            state.ajaxErrorMessage = t("modules.project.components.ProjectPage.errors.deletedError");
                            deletedError = true;
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
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage, denyCloseDialog: notFoundError || deletedError || accessDeniedError } });
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
    <h1 class="doneo-cursor-pointer" @click="onToggleTabPlacement">PROJECT</h1>
    <n-tabs :placement="tabPlacement" type="line" animated ref="tabsRef" v-model:value="tab">
        <n-tab-pane name="metadata" display-directive="show" key="metadata" :disabled="!projectId">
            <template #tab>
                Metadata
            </template>
            <ProjectMetadataTab v-if="projectId" v-model:project="project"
                :read-only="!project.allowedOperations.updateProject" />
        </n-tab-pane>
        <n-tab-pane name="permissions" display-directive="show:lazy" key="permissions"
            :disabled="!projectId || (!project.allowedOperations.updateProject && project.permissionsCount === 0)">
            <template #tab>
                {{ permissionsTabLabel }}
                <n-icon :component="IconAlertTriangle" color="red" style="margin-left: 8px;"
                    v-if="project.permissionsCount < 1" />
            </template>
            <ProjectPermissionsTab v-if="projectId" :project-id="projectId"
                v-model:item-count="project.permissionsCount" :read-only="!project.allowedOperations.updateProject" />
        </n-tab-pane>
        <n-tab-pane name="notes" :tab="notesTabLabel" display-directive="show:lazy" key="notes"
            :disabled="!projectId || (!project.allowedOperations.updateProject && project.notesCount === 0)">
            <ProjectNotesTab v-if="projectId" :project-id="projectId" v-model:item-count="project.notesCount"
                :read-only="!project.allowedOperations.updateProject" />
        </n-tab-pane>
        <n-tab-pane name="attachments" :tab="attachmentsTabLabel" display-directive="show:lazy" key="attachments"
            :disabled="!projectId || (!project.allowedOperations.updateProject && project.attachmentsCount === 0)">
            <ProjectAttachmentsTab v-if="projectId" :project-id="projectId"
                v-model:item-count="project.attachmentsCount" :read-only="!project.allowedOperations.updateProject" />
        </n-tab-pane>
        <n-tab-pane name="history" :tab="historyTabLabel" display-directive="show:lazy" key="history"
            :disabled="!projectId || (!project.allowedOperations.updateProject && project.historyOperationsCount === 0)">
            <ProjectHistoryTab v-if="projectId" :project-id="projectId"
                v-model:item-count="project.historyOperationsCount" :key="project.historyOperationsCount" />
        </n-tab-pane>
        <n-tab-pane name="tasks" display-directive="show:lazy" key="tasks"
            :disabled="!projectId || (!project.allowedOperations.updateProject && project.tasksCount === 0)">
            <template #tab>
                {{ tasksTabLabel }}
                <n-icon :component="IconAlertTriangle" color="red" style="margin-left: 8px;"
                    v-if="project.tasksCount < 1" />
            </template>
            <ProjectTasksTab v-if="projectId" :project-id="projectId" v-model:item-count="project.tasksCount"
                :read-only="!project.allowedOperations.updateProject" />
        </n-tab-pane>
    </n-tabs>
</template>

<style lang="css" scoped>
    .avatar {
        margin-right: 4px;
    }
</style>