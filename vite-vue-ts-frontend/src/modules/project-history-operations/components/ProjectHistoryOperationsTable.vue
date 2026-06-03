<script setup lang="ts">
    import { computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NEmpty } from 'naive-ui';

    import { ProjectHistoryOperation } from '../models/project-history-operation.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    //import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import RefreshAddActionsColumn from '../../../shared/components/tables/RefreshAddActionsColumn.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import UserSelector from '../../users/components/UserSelector.vue';

    interface Props {
        loading: boolean;
        projectHistoryOperations: ProjectHistoryOperation[];
        errorMessage?: string | null;
    }

    const { t } = useI18n();

    const emit = defineEmits(['refresh']);

    const props = defineProps<Props>();

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.header.columns.createdAt"),
            field: "createdAt",
            sortable: false,
        },
        {
            label: t("modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.header.columns.operationType"),
            field: "operationType",
            sortable: false,
        },
        {
            label: t("modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.header.columns.user"),
            field: "createdBy",
            sortable: false,
        },
    ]);

    const userFilter = defineModel<string>("userFilter", {
        default: "",
    });

    const onRefresh = () => {
        emit("refresh");
    };

</script>

<template>
    <ManageTable size="small">
        <template #thead>
            <tr>
                <th v-for="column in columns" :key="column.field"
                    :class="{ 'doneo-text-center': column.align === 'center' }">
                    {{ column.label }}
                </th>
                <th class="doneo-table-actions-column">{{ t("shared.components.table.header.columns.actions") }}</th>
            </tr>
            <tr>
                <th></th>
                <th></th>
                <th>
                    <UserSelector v-model:id="userFilter" hide-avatar clearable
                        :placeholder="t('modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.filters.user.placeholder')" />
                </th>
                <th class="doneo-text-center">
                    <RefreshAddActionsColumn @refresh="onRefresh" hide-add />
                </th>
            </tr>
        </template>
        <template #tbody v-if="!props.errorMessage">
            <tr v-for="projectHistoryOperation, index in projectHistoryOperations"
                :key="projectHistoryOperation.createdAt.msTimestamp ?? index">
                <td>{{ projectHistoryOperation.createdAt.toLocaleString() }}</td>
                <td>{{ projectHistoryOperation.getOperationTypeLabel() }}</td>
                <td>
                    <AvatarUserName :user-id="projectHistoryOperation.createdBy?.id ?? ''"
                        :user-name="projectHistoryOperation.createdBy?.name ?? ''" />
                </td>

                <td class="doneo-text-center">
                    <ManageTableActionButtons show-preview />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="projectHistoryOperations.length < 1 && !props.loading">
                    <n-empty
                        :description="t('modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>