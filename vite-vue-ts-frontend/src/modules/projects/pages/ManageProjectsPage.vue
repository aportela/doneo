<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch } from 'vue';
    import { useI18n } from "vue-i18n";
    import { useRouter } from 'vue-router';

    import { NCard, NModal } from 'naive-ui';

    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import type { FormMode } from '../../../shared/types/form-mode';
    import type { SearchRequest, ProjectResponse, PatchRequest } from '../types/dto';
    import type { ProjectsTableFilters } from '../types/projects-table-filters.ts';

    import type { Order } from '../../../shared/types/order.ts';
    import { Project } from '../models/project';

    import { projectService } from '../services/project';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import NewProjectForm from '../components/NewProjectForm.vue';
    import ProjectsTable from '../components/ProjectsTable.vue';
    import Pager from '../../../shared/components/tables/Pager.vue';
    import { type Pagination, PAGER_DEFAULT_RESULTS_PAGE } from '../../../shared/types/pager.ts';
    import type { ProjectStatus } from '../../project-statuses/models/project-status.ts';

    const router = useRouter();
    const { t } = useI18n();
    const { notify } = useNotify();
    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<Project[]>([]);

    const order = reactive<Order>({ field: "createdAt", direction: "DESC" });
    const pagination = reactive<Pagination>({ currentPage: 1, resultsPage: PAGER_DEFAULT_RESULTS_PAGE, totalPages: 1, totalResults: 0 });

    const resetPager = ref<boolean>(false);

    const filters = reactive<ProjectsTableFilters>({
        slug: "",
        typeId: null,
        priorityId: null,
        statusId: null,
        summary: "",
        createdAt: {
            from: null,
            to: null,
        },
        createdByUserId: null,
    });

    const showModal = ref<boolean>(false);
    const modalFormMode = ref<FormMode>("add");

    const selectedItem = ref<Project>(new Project());

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    watch(() => filters, () => {
        resetPager.value = true;
    }, { deep: true });

    watch([pagination.resultsPage], () => {
        if (pagination.currentPage != 1) {
            pagination.currentPage = 1;
        } else {
            onRefresh();
        }
    });

    watch([pagination.currentPage], () => {
        onRefresh();
    });

    const onSort = (newOrder: Order) => {
        order.field = newOrder.field;
        order.direction = newOrder.direction;
        onRefresh();
    };

    const onShowAddForm = () => {
        modalFormMode.value = "add";
        showModal.value = true;
    };

    const onCancelForm = () => {
        showModal.value = false;
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: SearchRequest = {
                pager: {
                    currentPage: pagination.currentPage,
                    resultsPage: pagination.resultsPage,
                },
                order: {
                    field: order.field,
                    direction: order.direction,
                },
                filter: {
                    slug: filters.slug.length > 0 ? filters.slug : undefined,
                    summary: filters.summary.length > 0 ? filters.summary : undefined,
                    typeId: filters.typeId !== null ? filters.typeId : undefined,
                    priorityId: filters.priorityId !== null ? filters.priorityId : undefined,
                    statusId: filters.statusId !== null ? filters.statusId : undefined,
                    createdAt: filters.createdAt,
                    createdByUserId: filters.createdByUserId !== null ? filters.createdByUserId : undefined,
                }
            };
            const response = await projectService.search(payload);
            pagination.totalPages = response.pager.totalPages;
            pagination.totalResults = response.pager.totalResults;
            items.value = response.projects.map((project: ProjectResponse) => new Project(project))
        } catch (error: unknown) {
            items.value = [];
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageProjectsPage.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.project.components.ManageProjectsPage.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.project.components.ManageProjectsPage.errors.refreshError");
                    console.error("Unhandled API error", { file: "ManageProjectsPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    // TODO:
    const onDelete = async (project: Project, _index?: number) => {
        if (project.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await projectService.delete(project.id);
                notify('success', t("modules.project.components.ManageProjectsPage.notifications.projectDeleted", { summary: project.summary }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                selectedItem.value = project;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageProjectsPage.onDelete" } });
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.project.components.ManageProjectsPage.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.project.components.ManageProjectsPage.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.project.components.ManageProjectsPage.errors.deleteError");
                        console.error("Unhandled API error", { file: "ManageProjectsPage.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("project id not set", { file: "ManageProjectsPage.vue", method: "onDelete" });
        }
    };

    const onAdded = (project: Project, openProjectAfterCreate: boolean) => {
        showModal.value = false;
        notify('success', t("modules.project.components.ManageProjectsPage.notifications.projectAdded", { summary: project.summary }));
        if (openProjectAfterCreate) {
            router.push(
                {
                    name: "projectTab",
                    params: {
                        projectId: project.id,
                        tab: "metadata",
                    }
                },
            ).catch((e) => {
                console.error(e);
            });
        } else {
            onRefresh();
        }
    };

    let updatedStatusProject: Project;
    let updatedStatus: ProjectStatus;

    const onStatusChanged = async (updatedProject: Project, status: ProjectStatus) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: PatchRequest = {
                id: updatedProject.id ?? "",
                status: {
                    id: status.id ?? ""
                },
            };
            const response: ProjectResponse = await projectService.patch(payload);
            if (response.id === updatedProject.id) {
                onRefresh();
                notify('success', t("modules.project.components.ManageProjectsPage.notifications.projectStatusUpdated", { summary: updatedProject.summary, status: status.name }));
            } else {
                state.ajaxErrorMessage = t("modules.project.components.ManageProjectsPage.errors.statusUpdateError", { summary: updatedProject.summary });
            }
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            updatedStatusProject = updatedProject;
                            updatedStatus = status;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageProjectsPage.onStatusChanged" } });
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.project.components.ManageProjectsPage.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.project.components.ManageProjectsPage.errors.statusUpdateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.project.components.ManageProjectsPage.errors.statusUpdateError");
                    console.error("Unhandled API error", { file: "ManageProjectsPage.vue", method: "onStatusChanged" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageProjectsPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageProjectsPage.onDelete")) {
                onDelete(selectedItem.value);
            } else if (payload.to.includes("ManageProjectsPage.onStatusChanged")) {
                onStatusChanged(updatedStatusProject, updatedStatus);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>

    <!-- TODO close icon-->
    <n-modal v-model:show="showModal">
        <NewProjectForm class="modal-form" @add="onAdded" @cancel="onCancelForm" />
    </n-modal>
    <n-card :title="t('modules.project.components.ManageProjectsPage.header.title')">
        <ProjectsTable :items="items" :disabled="state.ajaxRunning" @refresh="onRefresh" @add="onShowAddForm"
            :order="order" @sort="onSort" @status-changed="onStatusChanged" v-model:filters="filters" />
    </n-card>
</template>

<style lang="css" scoped>
    .modal-form {
        width: 40%;
    }
</style>