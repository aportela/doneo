<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useLoadingStore } from '../../../stores/loading';


    import { appBus } from '../../../shared/composables/bus';


    import type { Order } from '../../../shared/types/order.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';

    import { HistoryOperation } from '../models/history-operation.ts';

    import { useUserSettingsStore } from '../../../stores/userSettings.ts';

    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { historyOperationsService } from '../services/history-operations.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import HistoryOperationSelect from '../../../shared/components/selectors/HistoryOperationSelect.vue';
    import type { TimestampRange } from '../../../shared/composables/timestamps.ts';
    import { renderLabel } from '../../../shared/composables/naive-ui-helpers.ts';
    import type { SearchResponse } from '../types/dto.ts';

    interface Props {
        id?: string;
        projectId: string;
        taskId?: string;
    }

    const props = withDefaults(defineProps<Props>(), { id: "HistoryOperationsTable" });

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    const { t } = useI18n();
    const userSettingsStore = useUserSettingsStore();
    const loadingStore = useLoadingStore();
    const tableSettingsStore = useTableSettingsStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    watch(
        () => state.ajaxRunning,
        (ajaxRunning) => {
            loadingStore.set(ajaxRunning);
        }
    );


    const items = shallowRef<HistoryOperation[]>([]);

    const showNoItemsWarningMessage = ref<boolean>(false);


    const currentOrder = reactive<Order>({ field: "name", direction: "ASC" });

    const onSort = (newOrder: Order) => {
        currentOrder.field = newOrder.field;
        currentOrder.direction = newOrder.direction;
        // we have all results, use local sorting for avoiding server load
        if (currentOrder.direction === "ASC") {
            items.value = [...items.value].sort((a, b) =>
                a.id.localeCompare(b.id)
            );
        } else {
            items.value = [...items.value].sort((a, b) =>
                b.id.localeCompare(a.id)
            );
        }
    };

    const createdAtFilterRef = ref<InstanceType<typeof DateFilterSelect>[] | null>(null);

    interface HistoryOperationsTableFilters {
        userId: string | null;
        createdAt: TimestampRange;
        operationType: number | null;
    }
    const filters = reactive<HistoryOperationsTableFilters>(
        {
            userId: null,
            createdAt: {
                from: null,
                to: null,
            },
            operationType: null,
        }
    );

    const isFilteredByUser = computed<boolean>(() => filters.userId !== null);
    const isFilteredByCreatedAt = computed<boolean>(() => filters.createdAt.from != null || filters.createdAt.to != null);
    const isFilteredByOperationType = computed<boolean>(() => filters.operationType !== null);

    const onClearFilters = () => {
        filters.userId = null;
        filters.operationType = null;
        if (createdAtFilterRef.value) {
            createdAtFilterRef.value[0]?.reset();
        }
    };


    const localFilteredItems = computed(() => {
        return items.value.filter((historyOperation: HistoryOperation) => {
            return (
                (filters.userId === null || filters.userId === historyOperation.createdBy.id) &&
                (filters.operationType === null || filters.operationType === historyOperation.operationType) &&
                ((filters.createdAt.from === null && filters.createdAt.to === null) || (historyOperation.createdAt.msTimestamp != null && filters.createdAt.from != null && filters.createdAt.from <= historyOperation.createdAt.msTimestamp && filters.createdAt.to != null && filters.createdAt.to >= historyOperation.createdAt.msTimestamp))
            );
        });
    });

    const columnDefinitions = reactive<TableHeaderColumn<HistoryOperation>[]>(
        [
            {
                label: t("modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.header.columns.operationDate"),
                field: "createdAt",
                visible: true,
                sortable: true,
                isFiltered: () => isFilteredByCreatedAt.value,
                render: (row: HistoryOperation) => renderLabel(row.createdAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) ?? ""),
            },
            {
                label: t("modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.header.columns.operationType"),
                field: "operationType",
                visible: true,
                sortable: true,
                isFiltered: () => isFilteredByOperationType.value,
                render: (row: HistoryOperation) => renderLabel(row.getOperationTypeLabel()),
            },
            {
                label: t("modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.header.columns.user"),
                field: "createdBy",
                visible: true,
                sortable: true,
                isFiltered: () => isFilteredByUser.value,
                render: (row: HistoryOperation) => {
                    return h(AvatarUserName, { userId: row.createdBy.id, userName: row.createdBy.name });
                }
            },
        ]
    );

    // create (if not found) default settings for this table (column order & visibility)
    tableSettingsStore.register(props.id, { columns: columnDefinitions.map((column) => { return { field: column.field, visible: column.visible } }) ?? [] });

    // restore previous settings
    const tableSettings = tableSettingsStore.get(props.id);

    // build columns based on saved order visibility settings
    const columns = computed<TableHeaderColumn<Attachment>[]>(() =>
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

    const onRefresh = async () => {
        if (props.projectId) {
            if (props.taskId) {
                onRefreshTaskHistoryOperations();
            } else {
                onRefreshProjectHistoryOperations();
            }
        }
    };

    const onRefreshProjectHistoryOperations = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const results: SearchResponse = await historyOperationsService.getProjectHistoryOperations(props.projectId);
            items.value = results.historyOperations.map((operation) => new HistoryOperation(operation));
            itemCount.value = items.value?.length ?? 0;
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectAttachmentsTab.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                    console.error("Unhandled API error", { file: "ProjectAttachmentsTab.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onRefreshTaskHistoryOperations = async () => {
        if (props.taskId) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                const results: SearchResponse = await historyOperationsService.getTaskHistoryOperations(props.projectId, props.taskId);
                items.value = results.historyOperations.map((operation) => new HistoryOperation(operation));
                itemCount.value = items.value?.length ?? 0;
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectAttachmentsTab.onRefresh" } });
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                        console.error("Unhandled API error", { file: "ProjectAttachmentsTab.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("HistoryOperationsTable.onRefresh")) {
                onRefresh();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });

</script>

<template>
    <ManageTable :id="props.id" size="small" :disabled="state.ajaxRunning" :rows="localFilteredItems"
        :row-key="row => row.id" :columns="columns" :order="currentOrder"
        :show-no-items-warning-message="showNoItemsWarningMessage || (items.length > 0 && localFilteredItems.length === 0)"
        :no-items-warning-message="t('modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.warnings.noItemsFound')"
        @sort="onSort" @refresh="onRefresh" @clear-filters="onClearFilters" :buttons="['refresh', 'settings']">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">

                <DateFilterSelect v-if="column.field === 'createdAt'" clearable v-model:range="filters.createdAt"
                    ref="createdAtFilterRef" :disabled="state.ajaxRunning" />
                <HistoryOperationSelect v-else-if="column.field === 'operationType'"
                    v-model:history-operation-type="filters.operationType" :disabled="state.ajaxRunning" size="small"
                    clearable :show-only-task-history-operations="!!props.taskId"
                    :placeholder="t('modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.filters.operationType.placeholder')" />
                <UserSelector v-else-if="column.field === 'createdBy'" v-model:id="filters.userId"
                    :disabled="state.ajaxRunning" size="small" hide-avatar clearable
                    :placeholder="t('modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.filters.user.placeholder')" />
            </th>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>