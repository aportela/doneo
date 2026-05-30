<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch, computed } from 'vue';
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
    import { Sort } from '../../../shared/types/models/sort';
    import type { FormMode } from '../../../shared/types/form-mode';

    const { t } = useI18n();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectStatus[]>([]);

    const sort = ref<Sort>(new Sort("name", "ASC"));

    const nameFilter = ref<string>("");

    const nameFilterLowerCase = computed(() =>
        nameFilter.value?.toLowerCase() ?? ''
    );

    const filteredItems = computed(() => {
        return items.value.filter((projectPriority) => {
            const name = projectPriority.name?.toLowerCase();
            return ((!name || name?.includes(nameFilterLowerCase.value))
            );
        });
    });


    const showForm = ref<boolean>(false);
    const formMode = ref<FormMode>("add");

    const selectedItem = ref<ProjectStatus>(new ProjectStatus());

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

    const onShowUpdateForm = (projectStatus: ProjectStatus, _index: number) => {
        selectedItem.value = projectStatus;
        formMode.value = "update";
        showForm.value = true;
    };

    const onAdd = (projectStatus: ProjectStatus) => {
        showForm.value = false;
        notify('success', t("modules.projectStatus.components.ManageProjectStatusesPage.notifications.projectStatusAdded", { name: projectStatus.name }));
        onRefresh();
    };

    const onUpdate = (projectStatus: ProjectStatus) => {
        showForm.value = false;
        notify('success', t("modules.projectStatus.components.ManageProjectStatusesPage.notifications.projectStatusUpdated", { name: projectStatus.name }));
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
                    currentPage: 1,
                    resultsPage: 0,
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
    <n-modal v-model:show="showForm">
        <ProjectStatusForm :mode="formMode == 'add' ? 'add' : 'update'" :project-status-id="selectedItem.id"
            style="width: 40%;" @add="onAdd" @update="onUpdate" @cancel="onCancel" />
    </n-modal>

    <n-card :title="t('modules.projectStatus.components.ManageProjectStatusesPage.header.title')">
        <ProjectStatusesTable :projectStatuses="filteredItems" :loading="state.ajaxRunning" @refresh="onRefresh"
            @add="onShowAddForm" @update="onShowUpdateForm" @delete="onDelete" :sort-field="sort.field"
            :sort-order="sort.order" @toggle-sort="onToggleSort" v-model:project-status-name-filter="nameFilter" />
    </n-card>
</template>

<style lang="css" scoped></style>