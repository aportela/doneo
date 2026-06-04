<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard } from 'naive-ui';

    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import type { FormMode } from '../../../shared/types/form-mode';
    import type { SearchRequest, TaskPriorityResponse } from '../types/dto';
    import type { TaskPrioritiesTableFilters } from '../types/task-priorities-table-filters.ts';

    import { Sort } from '../../../shared/types/models/sort';
    import { TaskPriority } from '../models/task-priority';

    import { taskPriorityService } from '../services/task-priority';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import TaskPriorityForm from '../components/TaskPriorityForm.vue';
    import TaskPrioritiesTable from '../components/TaskPrioritiesTable.vue';

    const { t } = useI18n();
    const { notify } = useNotify();
    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<TaskPriority[]>([]);

    const sort = ref<Sort>(new Sort("name", "ASC"));

    const filters = reactive<ProjectTypesTableFilters>({
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

    const selectedItem = ref<TaskPriority>(new TaskPriority());

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

    const onShowUpdateForm = (taskPriority: TaskPriority, _index: number) => {
        selectedItem.value = taskPriority;
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
                    field: sort.value.field,
                    sort: sort.value.order,
                },
                filter: {
                    //name: filters.name.length > 0 ? filters.name : undefined,
                }
            };
            const response = await taskPriorityService.search(payload);
            items.value = response.taskPriorities.map((taskPriority: TaskPriorityResponse) => new TaskPriority(taskPriority))
        } catch (error: unknown) {
            items.value = [];
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageTaskPrioritiesPage.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.taskPriority.components.ManageTaskPrioritiesPage.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.taskPriority.components.ManageTaskPrioritiesPage.errors.refreshError");
                    console.error("Unhandled API error", { file: "ManageTaskPrioritiesPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onDelete = async (taskPriority: TaskPriority, _index?: number) => {
        if (taskPriority.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await taskPriorityService.delete(taskPriority.id);
                notify('success', t("modules.taskPriority.components.ManageTaskPrioritiesPage.notifications.taskPriorityDeleted", { name: taskPriority.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                selectedItem.value = taskPriority;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageTaskPrioritiesPage.onDelete" } });
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.taskPriority.components.ManageTaskPrioritiesPage.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.taskPriority.components.ManageTaskPrioritiesPage.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.taskPriority.components.ManageTaskPrioritiesPage.errors.deleteError");
                        console.error("Unhandled API error", { file: "ManageTaskPrioritiesPage.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("task priority id not set", { file: "ManageTaskPrioritiesPage.vue", method: "onDelete" });
        }
    };

    const onAdd = (taskPriority: TaskPriority) => {
        showModal.value = false;
        notify('success', t("modules.taskPriority.components.ManageTaskPrioritiesPage.notifications.taskPriorityAdded", { name: taskPriority.name }));
        onRefresh();
    };

    const onUpdate = (taskPriority: TaskPriority) => {
        showModal.value = false;
        notify('success', t("modules.taskPriority.components.ManageTaskPrioritiesPage.notifications.taskPriorityUpdated", { name: taskPriority.name }));
        onRefresh();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageTaskPrioritiesPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageTaskPrioritiesPage.onDelete")) {
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
        <TaskPriorityForm :mode="modalFormMode == 'add' ? 'add' : 'update'" :taskPriorityId="selectedItem.id"
            style="width: 40%;" @add="onAdd" @update="onUpdate" @cancel="onCancelForm" />
    </n-modal>

    <n-card :title="t('modules.taskPriority.components.ManageTaskPrioritiesPage.header.title')">
        <TaskPrioritiesTable :items="filteredItems" :disabled="state.ajaxRunning" @refresh="onRefresh"
            @add="onShowAddForm" @update="onShowUpdateForm" @delete="onDelete" :sort="sort" @sort="onSort"
            v-model:filters="filters" />
    </n-card>
</template>

<style lang="css" scoped></style>