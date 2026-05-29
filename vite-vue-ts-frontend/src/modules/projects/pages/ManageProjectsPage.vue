<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NCard, NModal } from 'naive-ui';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';
    import { projectService } from '../services/project';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { ProjectResponse, SearchRequest } from '../types/dto';
    import { Project } from '../models/project';
    import ProjectsTable from '../components/ProjectsTable.vue';
    import Pager from '../../../shared/components/tables/Pager.vue';
    import NewProjectForm from '../components/NewProjectForm.vue';
    import { Sort } from '../../../shared/types/models/sort';
    import type { FormMode } from '../../../shared/types/form-mode';

    const { t } = useI18n();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<Project[]>([]);

    const sort = ref<Sort>(new Sort("name", "ASC"));

    const keyFilter = ref<string>("");

    const showForm = ref<boolean>(false);
    const formMode = ref<FormMode>("add");

    const selectedItem = ref<Project>(new Project());

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    const onToggleSort = (field: string) => {
        sort.value.toggleSort(field);
        onRefresh();
    };

    const onShowAddForm = () => {
        formMode.value = "add";
        showForm.value = true;
    };

    const currentPage = ref(1);
    const pageSize = ref(10);
    const totalResults = ref(0);
    const totalPages = ref(0);

    watch([keyFilter], () => {
        currentPage.value = 1;
    });

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

    const onShowUpdateForm = (project: Project, _index: number) => {
        selectedItem.value = project;
        formMode.value = "update";
        showForm.value = true;
    };


    const onAdd = (project: Project) => {
        showForm.value = false;
        notify('success', t("modules.project.components.ManageProjectsPage.notifications.projectStatusAdded", { summary: project.summary }));
        onRefresh();
    };

    const onCancel = () => {
        showForm.value = false;
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
                    field: sort.value.field,
                    sort: sort.value.order,
                },
                filter: {
                    key: keyFilter.value,
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

    <n-modal v-model:show="showForm">
        <NewProjectForm style="width: 40%;" @add="onAdd" @cancel="onCancel" />
    </n-modal>
    <n-card :title="t('modules.project.components.ManageProjectsPage.header.title')">
        <Pager v-model:current-page="currentPage" v-model:page-size="pageSize" :total-pages="totalPages"
            :total-results="totalResults" class="doneo-pager-container">
            <template #total-results-label="{ totalResults }">
                {{ t("modules.project.components.ManageProjectsPage.pager.totalItemsLabel", { total: totalResults }) }}
            </template>
        </Pager>
        <ProjectsTable :projects="items" :loading="state.ajaxRunning" @refresh="onRefresh" @add="onShowAddForm"
            @update="onShowUpdateForm" @delete="onDelete" @textfilter-keydown-enter="onRefresh" :sort-field="sort.field"
            :sort-order="sort.order" @toggle-sort="onToggleSort" />
    </n-card>
</template>

<style lang="css" scoped></style>