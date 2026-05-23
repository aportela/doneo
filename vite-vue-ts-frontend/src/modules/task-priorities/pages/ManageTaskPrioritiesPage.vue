<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard } from 'naive-ui';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';
    import { taskPriorityService } from '../services/task-priority';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { TaskPriorityResponse, SearchRequest } from '../types/dto';
    import { TaskPriority } from '../models/task-priority';
    import TaskPrioritiesTable from '../components/TaskPrioritiesTable.vue';
    import TaskPriorityForm from '../components/TaskPriorityForm.vue';
    import { Sort } from '../../../shared/types/models/sort';
    import type { FormMode } from '../../../shared/types/form-mode';

    const { t } = useI18n();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<TaskPriority[]>([]);

    const sort = ref<Sort>(new Sort("name", "ASC"));

    const nameFilter = ref<string>("");

    const showForm = ref<boolean>(false);
    const formMode = ref<FormMode>("add");

    const selectedItem = ref<TaskPriority>(new TaskPriority({
        id: "",
        name: "",
        hexColor: "",
    }));

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

    const onShowUpdateForm = (taskPriority: TaskPriority, _index: number) => {
        selectedItem.value = taskPriority;
        formMode.value = "update";
        showForm.value = true;
    };

    const onAdd = (taskPriority: TaskPriority) => {
        showForm.value = false;
        notify('success', t("modules.taskPriority.components.ManageTaskPrioritiesPage.notifications.taskPriorityAdded", { name: taskPriority.name }));
        onRefresh();
    };

    const onUpdate = (taskPriority: TaskPriority) => {
        showForm.value = false;
        notify('success', t("modules.taskPriority.components.ManageTaskPrioritiesPage.notifications.taskPriorityUpdated", { name: taskPriority.name }));
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
            const response = await taskPriorityService.search(payload);
            items.value = response.taskPriorities.map((taskPriority: TaskPriorityResponse) => new TaskPriority(taskPriority))
        } catch (error: unknown) {
            items.value.length = 0;
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
        }
    };

    const onDelete = async (taskPriority: TaskPriority, _index?: number) => {
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
        }
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
    <n-modal v-model:show="showForm">
        <TaskPriorityForm :mode="formMode == 'add' ? 'add' : 'update'" :taskPriorityId="selectedItem.id"
            style="width: 40%;" @add="onAdd" @update="onUpdate" @cancel="onCancel" />
    </n-modal>

    <n-card :title="t('modules.taskPriority.components.ManageTaskPrioritiesPage.header.title')">
        <TaskPrioritiesTable :task-priorities="items" :loading="state.ajaxRunning" @refresh="onRefresh"
            @add="onShowAddForm" @update="onShowUpdateForm" @delete="onDelete" @textfilter-keydown-enter="onRefresh"
            :sort-field="sort.field" :sort-order="sort.order" @toggle-sort="onToggleSort"
            v-model:task-priority-name-filter="nameFilter" />
    </n-card>
</template>

<style lang="css" scoped></style>