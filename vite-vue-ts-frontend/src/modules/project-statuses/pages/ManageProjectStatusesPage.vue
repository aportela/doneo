<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard } from 'naive-ui';

    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import type { FormMode } from '../../../shared/types/form-mode';
    import type { SearchRequest, ProjectStatusResponse } from '../types/dto';
    import type { ProjectStatusesTableFilters } from '../types/project-statuses-table-filters.ts';

    import { Sort } from '../../../shared/types/models/sort';
    import { ProjectStatus } from '../models/project-status';

    import { projectStatusService } from '../services/project-status';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import ProjectStatusForm from '../components/ProjectStatusForm.vue';
    import ProjectStatusesTable from '../components/ProjectStatusesTable.vue';

    const { t } = useI18n();
    const { notify } = useNotify();
    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectStatus[]>([]);

    const sort = reactive<Sort>(new Sort("index", "ASC"));

    const filters = reactive<ProjectStatusesTableFilters>({
        name: "",
    });

    const nameFilterLowerCase = computed(() =>
        filters.name.toLowerCase()
    );

    const filteredItems = computed(() => {
        return items.value.filter((projectPriority) => {
            const name = projectPriority.name?.toLowerCase();
            return ((!name || name?.includes(nameFilterLowerCase.value))
            );
        });
    });

    const showModal = ref<boolean>(false);
    const modalFormMode = ref<FormMode>("add");

    const selectedItem = ref<ProjectStatus>(new ProjectStatus());

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
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

    const onShowUpdateForm = (projectStatus: ProjectStatus, _index: number) => {
        selectedItem.value = projectStatus;
        modalFormMode.value = "update";
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
                    currentPage: 1,
                    resultsPage: 0,
                },
                order: {
                    field: sort.field,
                    sort: sort.order,
                },
                filter: {
                    //name: filters.name.length > 0 ? filters.name : undefined,
                }
            };
            const response = await projectStatusService.search(payload);
            items.value = response.projectStatuses.map((projectStatus: ProjectStatusResponse) => new ProjectStatus(projectStatus))
        } catch (error: unknown) {
            items.value = [];
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageProjectStatusesPage.onRefresh" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectStatus.components.ManageProjectStatusesPage.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectStatus.components.ManageProjectStatusesPage.errors.refreshError");
                    console.error("Unhandled API error", { file: "ManageProjectStatusesPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onDelete = async (projectStatus: ProjectStatus, _index?: number) => {
        if (projectStatus.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await projectStatusService.delete(projectStatus.id);
                notify('success', t("modules.projectStatus.components.ManageProjectStatusesPage.notifications.projectStatusUpdated", { name: projectStatus.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                selectedItem.value = projectStatus;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageProjectStatusesPage.onDelete" } });
                                break;
                            case 403:
                                state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.projectStatus.components.ManageProjectStatusesPage.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectStatus.components.ManageProjectStatusesPage.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectStatus.components.ManageProjectStatusesPage.errors.deleteError");
                        console.error("Unhandled API error", { file: "ManageProjectStatusesPage.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("project status id not set", { file: "ManageProjectStatusesPage.vue", method: "onDelete" });
        }
    };

    const onAdded = (projectStatus: ProjectStatus) => {
        showModal.value = false;
        notify('success', t("modules.projectStatus.components.ManageProjectStatusesPage.notifications.projectStatusAdded", { name: projectStatus.name }));
        onRefresh();
    };

    const onUpdated = (projectStatus: ProjectStatus) => {
        showModal.value = false;
        notify('success', t("modules.projectStatus.components.ManageProjectStatusesPage.notifications.projectStatusUpdated", { name: projectStatus.name }));
        onRefresh();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageProjectStatusesPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageProjectStatusesPage.onDelete")) {
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
        <ProjectStatusForm :mode="modalFormMode == 'add' ? 'add' : 'update'" :project-status-id="selectedItem.id"
            class="modal-form" @add="onAdded" @update="onUpdated" @cancel="onCancelForm" />
    </n-modal>

    <n-card :title="t('modules.projectStatus.components.ManageProjectStatusesPage.header.title')">
        <ProjectStatusesTable :items="filteredItems" :disabled="state.ajaxRunning" @refresh="onRefresh"
            @add="onShowAddForm" @update="onShowUpdateForm" @delete="onDelete" :sort="sort" @sort="onSort"
            v-model:filters="filters" />
    </n-card>
</template>

<style lang="css" scoped>
    .modal-form {
        width: 40%;
    }
</style>