<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h, type Component } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, useDialog, NIcon, NButton, NButtonGroup, NTooltip } from 'naive-ui';

    import type { Order } from '../../../shared/types/order.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';

    import { TaskStatus } from '../models/task-status.ts';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';

    import { useLoadingStore } from '../../../stores/loading';
    import { useCacheStore } from '../../../stores/cache.ts';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { renderColoredTag, renderIcon, renderLabel } from '../../../shared/composables/naive-ui-helpers.ts';
    import { taskStatusService } from '../services/task-status.ts';
    import { handleAPIError } from '../../../api/client/errorHandler.ts';
    import { DONEO_ICON_ACTION_DELETE, DONEO_ICON_ACTION_EDIT, DONEO_ICON_CLEAR_DATE, DONEO_ICON_FILL_DATE, DONEO_ICON_FILL_EMTPY_DATE, DONEO_ICON_STAR } from '../../../shared/types/icons.ts';
    import type { TaskStatusResponse, SearchRequest } from '../types/dto.ts';
    import TaskStatusForm from './TaskStatusForm.vue';

    interface Props {
        id?: string;
    };

    const props = withDefaults(defineProps<Props>(), { id: "TaskStatusesTable" });;

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

    const items = shallowRef<TaskStatus[]>([]);

    const tmpItem = ref<TaskStatus>(new TaskStatus());

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

    interface TaskStatusTableFilters {
        name: string;
    };

    const filters = reactive<TaskStatusTableFilters>(
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
    const localFilteredItems = computed<TaskStatus[]>(() => {
        return items.value.filter((taskType: TaskStatus) => {
            const name = taskType.name?.toLowerCase();
            return (
                (!name || name?.includes(nameFilterLowerCase.value))
            );
        });
    });

    const FLAG_ICON_SIZE = 22;

    const createFlagTooltip = (
        enabled: boolean,
        icon: Component,
        enabledKey: string,
        disabledKey: string
    ) =>
        h(
            NTooltip,
            { trigger: "hover" },
            {
                trigger: () =>
                    h(NIcon, {
                        component: icon,
                        size: FLAG_ICON_SIZE,
                        class: [
                            "doneo-cursor-help",
                            { "doneo-disabled-icon": !enabled },
                        ],
                    }),
                default: () => t(enabled ? enabledKey : disabledKey),
            }
        );

    // get task status column from task status row
    const renderTaskStatusFlagsColumn = (row: TaskStatus) => {
        return h("div",
            {
                class: "task-status-flags",
            },
            [
                createFlagTooltip(
                    row.flags.defaultStatusOnCreation,
                    DONEO_ICON_STAR,
                    "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasDefaultStatusOnCreation",
                    "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasNotdefaultStatusOnCreation"
                ),
                createFlagTooltip(
                    row.flags.fillEmptyStartDate,
                    DONEO_ICON_FILL_EMTPY_DATE,
                    "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasFillEmptyStartDate",
                    "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasNotFillEmptyStartDate"
                ),
                createFlagTooltip(
                    row.flags.setStartDate,
                    DONEO_ICON_FILL_DATE,
                    "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasSetStartDate",
                    "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasNotSetStartDate"
                ),
                createFlagTooltip(
                    row.flags.fillEmptyFinishDate,
                    DONEO_ICON_FILL_EMTPY_DATE,
                    "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasFillEmptyFinishDate",
                    "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasNotFillEmptyFinishDate"
                ),
                createFlagTooltip(
                    row.flags.setFinishDate,
                    DONEO_ICON_FILL_DATE,
                    "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasSetFinishDate",
                    "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasNotSetFinishDate"
                ),
                createFlagTooltip(
                    row.flags.unsetFinishDateOnLeave,
                    DONEO_ICON_CLEAR_DATE,
                    "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasUnsetFinishDateOnLeave",
                    "modules.taskStatus.components.TaskStatusesTable.body.columns.permissionsHints.hasNotUnsetFinishDateOnLeave"
                ),
            ]
        );
    };

    const columnDefinitions = reactive<TableHeaderColumn<TaskStatus>[]>([
        {
            label: t("modules.taskStatus.components.TaskStatusesTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByName.value,
            render: (row: TaskStatus) => renderColoredTag(row.name, row.hexColor, true),
        },
        {
            label: t("modules.taskStatus.components.TaskStatusesTable.header.columns.index"),
            field: "index",
            visible: true,
            sortable: true,
            isFiltered: () => false,
            render: (row: TaskStatus) => renderLabel(row.index),
        },
        {
            label: t("modules.taskStatus.components.TaskStatusesTable.header.columns.flags"),
            field: "flags",
            visible: true,
            sortable: false,
            align: "center",
            isFiltered: () => false,
            render: (row: TaskStatus) => renderTaskStatusFlagsColumn(row),
        },
    ]);

    // create (if not found) default settings for this table (column order & visibility)
    tableSettingsStore.register(props.id, { columns: columnDefinitions.map((column) => { return { field: column.field, visible: column.visible } }) ?? [] });

    // restore previous settings
    const tableSettings = tableSettingsStore.get(props.id);

    // build columns based on saved order visibility settings
    const columns = computed<TableHeaderColumn<TaskStatus>[]>(() =>
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

    const onDelete = async (taskStatus: TaskStatus) => {
        if (taskStatus.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await taskStatusService.delete(taskStatus.id);
                cacheStore.clearTaskStatusesCache();
                notify('success', t("modules.taskStatus.components.TaskStatusesTable.notifications.taskStatusUpdated", { name: taskStatus.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                tmpItem.value = taskStatus;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "TaskStatusesTable.onDelete" } });
                                break;
                            case 403:
                                state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusesTable.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusesTable.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusesTable.errors.deleteError");
                        console.error("Unhandled API error", { file: "TaskStatusesTable.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("task status id not set", { file: "TaskStatusesTable.vue", method: "onDelete" });
        }
    };

    const onConfirmDelete = (taskStatus: TaskStatus) => {
        dialog.warning({
            title: t("modules.taskStatus.components.TaskStatusesTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(DONEO_ICON_ACTION_DELETE, { size: 24 }),
            content: () =>
                h('div', [
                    t("modules.taskStatus.components.TaskStatusesTable.dialogs.deleteConfirmation.message", { name: taskStatus.name }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                onDelete(taskStatus);
            },
        });
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
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
            const response = await taskStatusService.search(payload);
            items.value = response.taskStatuses.map((taskStatus: TaskStatusResponse) => new TaskStatus(taskStatus))
            showNoItemsWarningMessage.value = items.value.length === 0;
        } catch (error: unknown) {
            items.value = [];
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TaskStatusesTable.onRefresh" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusesTable.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.taskStatus.components.TaskStatusesTable.errors.refreshError");
                    console.error("Unhandled API error", { file: "TaskStatusesTable.vue", method: "onRefresh" }, { err: fatalError });
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
        tmpItem.value = new TaskStatus();
        showFormModal.value = true;
    };

    const onUpdate = (taskStatus: TaskStatus) => {
        tmpItem.value = taskStatus;
        showFormModal.value = true;
    };

    const onTaskStatusAdded = (taskStatus: TaskStatus) => {
        showFormModal.value = false;
        cacheStore.clearTaskStatusesCache();
        notify('success', t("modules.taskStatus.components.TaskStatusesTable.notifications.taskStatusAdded", { name: taskStatus.name }));
        onRefresh();
    };

    const onTaskStatusUpdated = (taskStatus: TaskStatus) => {
        showFormModal.value = false;
        cacheStore.clearTaskStatusesCache();
        notify('success', t("modules.taskStatus.components.TaskStatusesTable.notifications.taskStatusUpdated", { name: taskStatus.name }));
        onRefresh();
    };

    const hideFormModal = () => {
        showFormModal.value = false;
        tmpItem.value = new TaskStatus();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("TaskStatusesTable.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("TaskStatusesTable.onDelete")) {
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
        <TaskStatusForm class="task-status-form" :task-status-id="tmpItem.id" @add="onTaskStatusAdded"
            @update="onTaskStatusUpdated" @cancel="hideFormModal" v-if="showFormModal" />
        <span v-else />
    </n-modal>
    <ManageTable :id="props.id" size="small" :disabled="state.ajaxRunning" :rows="localFilteredItems"
        :row-key="row => row.id" :columns="columns" :order="currentOrder"
        :show-no-items-warning-message="showNoItemsWarningMessage || (items.length > 0 && localFilteredItems.length === 0)"
        :no-items-warning-message="t('modules.taskStatus.components.TaskStatusesTable.warnings.noItemsFound')"
        @sort="onSort" @refresh="onRefresh" @add="onAdd" @clear-filters="onClearFilters">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <TextFilterInput v-if="column.field === 'name'" clearable :disabled="state.ajaxRunning" size="small"
                    :placeholder="t('modules.taskStatus.components.TaskStatusesTable.filters.name.placeholder')"
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
    .task-status-form {
        width: 95%;
        max-width: 640px;
    }
</style>