<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard } from 'naive-ui';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';
    import { projectStatusService } from '../services/project-status';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { ProjectStatusResponse, SearchRequest } from '../types/dto';
    import { ProjectStatus } from '../models/project-status';
    import ProjectStatusesTable from '../components/ProjectStatusesTable.vue';
    import ProjectStatusForm from '../components/ProjectStatusForm.vue';
    import Pager from '../../../shared/components/tables/Pager.vue';
    import { Sort } from '../../../shared/types/models/sort';
    import type { FormMode } from '../types/form-mode';

    const { notify } = useNotify();

    const { t } = useI18n();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectStatuses = shallowRef<ProjectStatus[]>([]);

    const sort = ref<Sort>(new Sort("index", "ASC"));

    const nameFilter = ref<string>("");

    const currentPage = ref(1);
    const pageSize = ref(0);
    const totalResults = ref(0);
    const totalPages = ref(0);


    const showProjectStatusDialogForm = ref<boolean>(false);
    const projectStatusDialogFormMode = ref<FormMode>("add");

    const selectedProjectStatus = ref<ProjectStatus>(new ProjectStatus({
        id: "",
        name: "",
        hexColor: "",
    }));

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
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

    const onToggleSort = (field: string) => {
        sort.value.toggleSort(field);
        onRefresh();
    };

    const onShowAddForm = () => {
        projectStatusDialogFormMode.value = "add";
        showProjectStatusDialogForm.value = true;
    };

    const onShowUpdateForm = (projectStatus: ProjectStatus, _index: number) => {
        selectedProjectStatus.value = projectStatus;
        projectStatusDialogFormMode.value = "update";
        showProjectStatusDialogForm.value = true;
    };

    const onAdd = (projectStatus: ProjectStatus) => {
        showProjectStatusDialogForm.value = false;
        notify('success', t("projectStatusAddedNotification", { name: projectStatus.name }));
        onRefresh();
    };

    const onUpdate = (projectStatus: ProjectStatus) => {
        showProjectStatusDialogForm.value = false;
        notify('success', t("projectStatusUpdatedNotification", { name: projectStatus.name }));
        onRefresh();
    };

    const onCancel = () => {
        showProjectStatusDialogForm.value = false;
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
                    name: nameFilter.value,
                }
            };
            const response = await projectStatusService.search(payload);
            totalPages.value = response.pager.totalPages;
            totalResults.value = response.pager.totalResults;
            projectStatuses.value = response.projectStatuses.map((projectStatus: ProjectStatusResponse) => new ProjectStatus(projectStatus))
        } catch (error: unknown) {
            projectStatuses.value.length = 0;
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageProjectStatusesPage.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while refreshing the project type list");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while refreshing the project type list");
                    console.error("Unhandled API error", { file: "ManageProjectStatusesPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
        }
    };

    const onDelete = async (projectStatus: ProjectStatus, _index?: number) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await projectStatusService.delete(projectStatus.id);
            notify('success', t("projectStatusIndexMovedNotification", { name: projectStatus.name }));
            onRefresh();
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            selectedProjectStatus.value = projectStatus;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageProjectStatusesPage.onDelete" } });
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("We couldn’t find the specified project type");
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while deleting the project type");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while deleting the project type");
                    console.error("Unhandled API error", { file: "ManageProjectStatusesPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageProjectStatusesPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageProjectStatusesPage.onDelete")) {
                onDelete(selectedProjectStatus.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-modal v-model:show="showProjectStatusDialogForm">
        <ProjectStatusForm :mode="projectStatusDialogFormMode == 'add' ? 'add' : 'update'"
            :project-status-id="selectedProjectStatus.id" style="width: 40%;" @add="onAdd" @update="onUpdate"
            @cancel="onCancel" />
    </n-modal>

    <n-card :title="t('Manage project statuses')">
        <Pager v-model:current-page="currentPage" v-model:page-size="pageSize" :total-pages="totalPages"
            :total-results="totalResults" class="doneo-pager-container">
            <template #total-results-label="{ totalResults }">
                {{ t("TotalProjectStatusesPagerLabel", { total: totalResults }) }}
            </template>
        </Pager>
        <ProjectStatusesTable :projectStatuses="projectStatuses" :loading="state.ajaxRunning" @refresh="onRefresh"
            @add="onShowAddForm" @update="onShowUpdateForm" @delete="onDelete" @textfilter-keydown-enter="onRefresh"
            :sort-field="sort.field" :sort-order="sort.order" @toggle-sort="onToggleSort"
            v-model:project-status-name-filter="nameFilter" />
    </n-card>
</template>

<style lang="css" scoped>
    .doneo-pager-container {
        margin-bottom: 4px;
    }
</style>