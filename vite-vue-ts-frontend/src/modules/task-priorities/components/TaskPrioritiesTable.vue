<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, useDialog, NIcon, NButton, NButtonGroup } from 'naive-ui';

    import type { Order } from '../../../shared/types/order.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';

    import { TaskPriority } from '../models/task-priority.ts';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';

    import { useLoadingStore } from '../../../stores/loading';
    import { useCacheStore } from '../../../stores/cache.ts';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { renderColoredTag, renderIcon, renderLabel } from '../../../shared/composables/naive-ui-helpers.ts';
    import { taskPriorityService } from '../services/task-priority.ts';
    import { handleAPIError } from '../../../api/client/errorHandler.ts';
    import { DONEO_ICON_ACTION_DELETE, DONEO_ICON_ACTION_EDIT } from '../../../shared/types/icons.ts';
    import type { TaskPriorityResponse, SearchRequest } from '../types/dto.ts';
    import TaskPriorityForm from './TaskPriorityForm.vue';

    interface Props {
        id?: string;
    }

    const props = withDefaults(defineProps<Props>(), { id: "TaskPrioritiesTable" });;

    const { t } = useI18n();
    const dialog = useDialog();

    const { notify } = useNotify();

    const loadingStore = useLoadingStore();
    const tableSettingsStore = useTableSettingsStore();
    const cacheStore = useCacheStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    watch(
        () => state.ajaxRunning,
        (ajaxRunning) => {
            loadingStore.set(ajaxRunning);
        }
    );

    const items = shallowRef<TaskPriority[]>([]);

    const tmpItem = ref<TaskPriority>(new TaskPriority());

    const showNoItemsWarningMessage = ref<boolean>(false);

    const currentOrder = reactive<Order>({ field: "name", direction: "ASC" });

    const onSort = (newOrder: Order) => {
        currentOrder.field = newOrder.field;
        currentOrder.direction = newOrder.direction;
        // we have all results, use local sorting for avoiding server load
        if (currentOrder.direction === "ASC") {
            items.value = [...items.value].sort((a, b) =>
                currentOrder.field === "name" ? a.name.localeCompare(b.name) : String(a.index).localeCompare(String(b.index))
            );
        } else {
            items.value = [...items.value].sort((a, b) =>
                currentOrder.field === "name" ? b.name.localeCompare(a.name) : String(b.index).localeCompare(String(a.index))
            );
        }
    };

    interface TaskPrioritiesTableFilters {
        name: string;
    };

    const filters = reactive<TaskPrioritiesTableFilters>(
        {
            name: "",
        }
    );

    const isFilteredByName = computed<boolean>(() => filters.name.length > 0);

    const onClearFilters = () => {
        filters.name = "";
    };

    const nameFilterLowerCase = computed(() => filters.name.toLowerCase());

    // we have all results, use local filtering for avoiding server load
    const localFilteredItems = computed<TaskPriority[]>(() => {
        return items.value.filter((taskType: TaskPriority) => {
            const name = taskType.name?.toLowerCase();
            return (
                (!name || name?.includes(nameFilterLowerCase.value))
            );
        });
    });

    const columnDefinitions = reactive<TableHeaderColumn<TaskPriority>[]>([
        {
            label: t("modules.taskPriority.components.TaskPrioritiesTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByName.value,
            render: (row: TaskPriority) => renderColoredTag(row.name, row.hexColor, true),
        },
        {
            label: t("modules.taskPriority.components.TaskPrioritiesTable.header.columns.index"),
            field: "index",
            visible: true,
            sortable: true,
            isFiltered: () => false,
            render: (row: TaskPriority) => renderLabel(row.index),
        },
    ]);

    // create (if not found) default settings for this table (column order & visibility)
    tableSettingsStore.register(props.id, { columns: columnDefinitions.map((column) => { return { field: column.field, visible: column.visible } }) ?? [] });

    // restore previous settings
    const tableSettings = tableSettingsStore.get(props.id);

    // build columns based on saved order visibility settings
    const columns = computed<TableHeaderColumn<TaskPriority>[]>(() =>
        tableSettings.columns.map((column) => { // get saved ordered columns
            const definition = columnDefinitions.find((c) => c.field === column.field);
            return {
                label: definition?.label ?? "",
                field: column.field,
                visible: column.visible,
                sortable: definition!.sortable,
                align: definition?.align,
                isFiltered: definition?.isFiltered ?? (() => false),
                render: definition?.render ?? (() => "")
            };
        })
    );

    const onDelete = async (taskPriority: TaskPriority) => {
        if (taskPriority.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await taskPriorityService.delete(taskPriority.id);
                cacheStore.clearTaskPrioritiesCache();
                notify('success', t("modules.taskPriority.components.TaskPrioritiesTable.notifications.taskPriorityDeleted", { name: taskPriority.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                tmpItem.value = taskPriority;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "TaskPrioritiesTable.onDelete" } });
                                break;
                            case 403:
                                state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPrioritiesTable.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPrioritiesTable.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPrioritiesTable.errors.deleteError");
                        console.error("Unhandled API error", { file: "TaskPrioritiesTable.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("task priority id not set", { file: "TaskPrioritiesTable.vue", method: "onDelete" });
        }
    };

    const onConfirmDelete = (taskPriority: TaskPriority) => {
        dialog.warning({
            title: t("modules.taskPriority.components.TaskPrioritiesTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(DONEO_ICON_ACTION_DELETE, { size: 24 }),
            content: () =>
                h('div', [
                    t("modules.taskPriority.components.TaskPrioritiesTable.dialogs.deleteConfirmation.message", { name: taskPriority.name }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                onDelete(taskPriority)
            },
        });
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        showNoItemsWarningMessage.value = false;
        try {
            const payload: SearchRequest = {
                pager: {
                    enabled: false,
                    currentPage: 1,
                    resultsPage: 0,
                },
                order: {
                    field: currentOrder.field,
                    direction: currentOrder.direction,
                },
                filter: {
                    //name: filters.name.length > 0 ? filters.name : undefined,
                }
            };
            const response = await taskPriorityService.search(payload);
            items.value = response.taskPriorities.map((taskPriority: TaskPriorityResponse) => new TaskPriority(taskPriority))
            showNoItemsWarningMessage.value = items.value.length === 0;
        } catch (error: unknown) {
            items.value = [];
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TaskPrioritiesTable.onRefresh" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPrioritiesTable.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.taskPriority.components.TaskPrioritiesTable.errors.refreshError");
                    console.error("Unhandled API error", { file: "TaskPrioritiesTable.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const showFormModal = ref<boolean>(false);

    const onAdd = () => {
        tmpItem.value = new TaskPriority();
        showFormModal.value = true;
    };

    const onUpdate = (taskPriority: TaskPriority) => {
        tmpItem.value = taskPriority;
        showFormModal.value = true;
    };

    const onTaskPriorityAdded = (taskPriority: TaskPriority) => {
        showFormModal.value = false;
        cacheStore.clearTaskPrioritiesCache();
        notify('success', t("modules.taskPriority.components.TaskPrioritiesTable.notifications.taskPriorityAdded", { name: taskPriority.name }));
        onRefresh();
    };

    const onTaskPriorityUpdated = (taskPriority: TaskPriority) => {
        showFormModal.value = false;
        cacheStore.clearTaskPrioritiesCache();
        notify('success', t("modules.taskPriority.components.TaskPrioritiesTable.notifications.taskPriorityUpdated", { name: taskPriority.name }));
        onRefresh();
    };

    const hideFormModal = () => {
        showFormModal.value = false;
        tmpItem.value = new TaskPriority();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("TaskPrioritiesTable.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("TaskPrioritiesTable.onDelete")) {
                onDelete(tmpItem.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-modal v-model:show="showFormModal">
        <TaskPriorityForm class="task-priority-form" :task-priority-id="tmpItem.id" @add="onTaskPriorityAdded"
            @update="onTaskPriorityUpdated" @cancel="hideFormModal" v-if="showFormModal" />
        <span v-else />
    </n-modal>
    <ManageTable :id="props.id" size="small" :disabled="state.ajaxRunning" :rows="localFilteredItems"
        :row-key="row => row.id" :columns="columns" :order="currentOrder"
        :show-no-items-warning-message="showNoItemsWarningMessage || (items.length > 0 && localFilteredItems.length === 0)"
        :no-items-warning-message="t('modules.taskPriority.components.TaskPrioritiesTable.warnings.noItemsFound')"
        @sort="onSort" @refresh="onRefresh" @add="onAdd" @clear-filters="onClearFilters">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <TextFilterInput v-if="column.field === 'name'" clearable :disabled="state.ajaxRunning" size="small"
                    :placeholder="t('modules.taskPriority.components.TaskPrioritiesTable.filters.name.placeholder')"
                    v-model:value="filters.name" />
            </th>
        </template>
        <template #rowactions="{ row }">
            <n-button-group class="doneo-table-actions-button-group" size="small">
                <n-button @click="onUpdate(row)" :disabled="state.ajaxRunning" class="doneo-table-actions-button">
                    {{ t("shared.buttons.Edit.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_EDIT" />
                    </template>
                </n-button>
                <n-button @click="onConfirmDelete(row)" :disabled="state.ajaxRunning"
                    class="doneo-table-actions-button">
                    {{ t("shared.buttons.Delete.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_DELETE" />
                    </template>
                </n-button>
            </n-button-group>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped>
    .task-priority-form {
        width: 95%;
        max-width: 640px;
    }
</style>