<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard } from 'naive-ui';

    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import type { FormMode } from '../../../shared/types/form-mode';
    import type { SearchRequest, TaskStatusResponse } from '../types/dto';
    import type { TaskStatusesTableFilters } from '../types/task-statuses-table-filters.ts';

    import { Sort } from '../../../shared/types/models/sort';
    import { TaskStatus } from '../models/task-status';

    import { taskStatusService } from '../services/task-status';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import TaskStatusForm from '../components/TaskStatusForm.vue';
    import TaskStatusesTable from '../components/TaskStatusesTable.vue';

    const { t } = useI18n();
    const { notify } = useNotify();
    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<TaskStatus[]>([]);

    const sort = reactive<Sort>(new Sort("index", "ASC"));

    const filters = reactive<TaskStatusesTableFilters>({
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

    const selectedItem = ref<TaskStatus>(new TaskStatus());

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    const onShowAddForm = () => {
        modalFormMode.value = "add";
        showModal.value = true;
    };

    const onShowUpdateForm = (taskStatus: TaskStatus, _index: number) => {
        selectedItem.value = taskStatus;
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
            const response = await taskStatusService.search(payload);
            items.value = response.taskStatuses.map((taskStatus: TaskStatusResponse) => new TaskStatus(taskStatus))
        } catch (error: unknown) {
            items.value = [];
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageTaskStatusesPage.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.taskStatus.components.ManageTaskStatusesPage.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.taskStatus.components.ManageTaskStatusesPage.errors.refreshError");
                    console.error("Unhandled API error", { file: "ManageTaskStatusesPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onDelete = async (taskStatus: TaskStatus, _index?: number) => {
        if (taskStatus.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await taskStatusService.delete(taskStatus.id);
                notify('success', t("modules.taskStatus.components.ManageTaskStatusesPage.notifications.taskStatusUpdated", { name: taskStatus.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                selectedItem.value = taskStatus;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageTaskStatusesPage.onDelete" } });
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.taskStatus.components.ManageTaskStatusesPage.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.taskStatus.components.ManageTaskStatusesPage.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.taskStatus.components.ManageTaskStatusesPage.errors.deleteError");
                        console.error("Unhandled API error", { file: "ManageTaskStatusesPage.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("task status id not set", { file: "ManageTaskStatusesPage.vue", method: "onDelete" });
        }
    };

    const onAdded = (taskStatus: TaskStatus) => {
        showModal.value = false;
        notify('success', t("modules.taskStatus.components.ManageTaskStatusesPage.notifications.taskStatusAdded", { name: taskStatus.name }));
        onRefresh();
    };

    const onUpdated = (taskStatus: TaskStatus) => {
        showModal.value = false;
        notify('success', t("modules.taskStatus.components.ManageTaskStatusesPage.notifications.taskStatusUpdated", { name: taskStatus.name }));
        onRefresh();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageTaskStatusesPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageTaskStatusesPage.onDelete")) {
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
        <TaskStatusForm :mode="modalFormMode == 'add' ? 'add' : 'update'" :taskStatusId="selectedItem.id"
            class="modal-form" @add="onAdded" @update="onUpdated" @cancel="onCancelForm" />
    </n-modal>

    <n-card :title="t('modules.taskStatus.components.ManageTaskStatusesPage.header.title')">
        <TaskStatusesTable :items="filteredItems" :disabled="state.ajaxRunning" @refresh="onRefresh"
            @add="onShowAddForm" @update="onShowUpdateForm" @delete="onDelete" v-model:filters="filters" />
    </n-card>
</template>

<style lang="css" scoped>
    .modal-form {
        width: 40%;
    }
</style>