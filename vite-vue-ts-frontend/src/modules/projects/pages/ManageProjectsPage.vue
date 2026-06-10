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
    import type { SearchRequest, ProjectResponse } from '../types/dto';
    import type { ProjectsTableFilters } from '../types/projects-table-filters.ts';

    import { Sort } from '../../../shared/types/models/sort';
    import { Project } from '../models/project';

    import { projectService } from '../services/project';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import NewProjectForm from '../components/NewProjectForm.vue';
    import ProjectsTable from '../components/ProjectsTable.vue';
    import Pager from '../../../shared/components/tables/Pager.vue';

    const router = useRouter();
    const { t } = useI18n();
    const { notify } = useNotify();
    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<Project[]>([]);

    const sort = reactive<Sort>(new Sort("createdAt", "DESC"));

    const resetPager = ref<boolean>(false);
    const currentPage = ref(1);
    const pageSize = ref(10);
    const totalResults = ref(0);
    const totalPages = ref(0);

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

    watch(pageSize, () => {
        if (currentPage.value != 1) {
            currentPage.value = 1;
        } else {
            onRefresh();
        }
    });

    watch(currentPage, () => {
        onRefresh();
    });

    const onSort = (newSort: Sort) => {
        sort.field = newSort.field;
        sort.order = newSort.order;
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
                    currentPage: currentPage.value,
                    resultsPage: pageSize.value,
                },
                order: {
                    field: sort.field,
                    sort: sort.order,
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
            totalPages.value = response.pager.totalPages;
            totalResults.value = response.pager.totalResults;
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

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageProjectsPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageProjectsPage.onDelete")) {
                onDelete(selectedItem.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>

    <n-modal v-model:show="showModal">
        <NewProjectForm class="modal-form" @add="onAdded" @cancel="onCancelForm" />
    </n-modal>
    <n-card :title="t('modules.project.components.ManageProjectsPage.header.title')">
        <Pager v-model:current-page="currentPage" v-model:page-size="pageSize" :total-pages="totalPages"
            :total-results="totalResults" class="doneo-pager-container">
            <template #total-results-label="{ totalResults }">
                {{ t("modules.project.components.ManageProjectsPage.pager.totalItemsLabel", { total: totalResults }) }}
            </template>
        </Pager>
        <ProjectsTable :items="items" :disabled="state.ajaxRunning" @refresh="onRefresh" @add="onShowAddForm"
            :sort="sort" @sort="onSort" v-model:filters="filters" />
    </n-card>
</template>

<style lang="css" scoped>
    .modal-form {
        width: 40%;
    }
</style>