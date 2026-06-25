<script setup lang="ts">
    import { ref, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NEmpty } from 'naive-ui';

    import { useUserSettingsStore } from '../../../stores/userSettings.ts';
    import { HistoryOperation } from '../models/history-operation.ts';

    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { HistoryOperationsTableFilters } from '../types/history-operations-table-filters.ts';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import ClearFiltersTableButton from '../../../shared/components/tables/ClearFiltersTableButton.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import type { DateFilterSelectComponent } from '../../users/components/date-filter-select-component.ts';
    import HistoryOperationSelect from '../../../shared/components/selectors/HistoryOperationSelect.vue';

    interface HistoryOperationsTableProps {
        disabled: boolean;
        items: HistoryOperation[];
        projectId: string;
        taskId?: string;
        errorMessage?: string | null;
    }

    const { t } = useI18n();
    const userSettingsStore = useUserSettingsStore();

    const emit = defineEmits(['refresh']);

    const props = defineProps<HistoryOperationsTableProps>();

    const createdAtFilterRef = ref<DateFilterSelectComponent | undefined>();

    const filters = defineModel<HistoryOperationsTableFilters>("filters", {
        default: () => ({
            userId: "",
            createdAt: {
                from: null,
                to: null,
            },
            operationType: null,
        })
    });

    const isFilteredByUser = computed<boolean>(() => filters.value.userId !== null);
    const isFilteredByCreatedAt = computed<boolean>(() => filters.value.createdAt.from != null || filters.value.createdAt.to != null);
    const isFilteredByOperationType = computed<boolean>(() => filters.value.operationType !== null);

    const hasFilters = computed<boolean>(() =>
        isFilteredByUser.value || isFilteredByOperationType.value || isFilteredByCreatedAt.value
    );

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.header.columns.operationDate"),
            field: "createdAt",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByCreatedAt.value,
        },
        {
            label: t("modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.header.columns.operationType"),
            field: "operationType",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByOperationType.value,
        },
        {
            label: t("modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.header.columns.user"),
            field: "createdBy",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByUser.value,
        },
    ]);

    const onRefresh = () => {
        emit("refresh");
    };

    const onClearFilters = () => {
        filters.value.userId = null;
        filters.value.operationType = null;
        createdAtFilterRef.value?.reset();
    };
</script>

<template>
    <ManageTable size="small" :columns="columns" hide-add @refresh="onRefresh">
        <template #thead>
            <tr>
                <th>
                    <DateFilterSelect clearable v-model:range="filters.createdAt" ref="createdAtFilterRef"
                        :disabled="props.disabled" />
                </th>
                <th>
                    <HistoryOperationSelect v-model:history-operation-type="filters.operationType"
                        :disabled="props.disabled" clearable
                        :placeholder="t('modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.filters.operationType.placeholder')" />
                </th>
                <th>
                    <UserSelector v-model:id="filters.userId" :disabled="props.disabled" hide-avatar clearable
                        :placeholder="t('modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.filters.user.placeholder')" />
                </th>
                <th class="doneo-text-center">
                    <ClearFiltersTableButton @clear="onClearFilters" :disabled="props.disabled || !hasFilters" />
                </th>
            </tr>
        </template>
        <template #tbody v-if="!props.errorMessage">
            <tr v-for="projectHistoryOperation, index in items" :key="projectHistoryOperation.id ?? index">
                <td>{{ projectHistoryOperation.createdAt.toCustomMaskString(userSettingsStore.currentDatetimeMask) }}
                </td>
                <td>{{ projectHistoryOperation.getOperationTypeLabel() }}</td>
                <td>
                    <AvatarUserName :user-id="projectHistoryOperation.createdBy?.id ?? ''"
                        :user-name="projectHistoryOperation.createdBy?.name ?? ''" />
                </td>

                <td class="doneo-text-center">
                    <ManageTableActionButtons show-preview :disabled="true || props.disabled" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="items.length < 1 && !props.disabled">
                    <n-empty
                        :description="t('modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>