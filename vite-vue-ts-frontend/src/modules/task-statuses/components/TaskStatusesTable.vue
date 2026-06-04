<script setup lang="ts">
    import { h, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NEmpty, NTag } from 'naive-ui';
    import { IconTrash } from '@tabler/icons-vue';

    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import type { Sort } from '../../../shared/types/models/sort.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { TaskStatusesTableFilters } from '../types/task-statuses-table-filters.ts';
    import { TaskStatus } from '../models/task-status';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import ClearFiltersTableButton from '../../../shared/components/tables/ClearFiltersTableButton.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import { getNaiveUITagColorProperty } from '../../../shared/composables/color';

    interface Props {
        disabled: boolean;
        items: TaskStatus[];
    }

    const { t } = useI18n();
    const dialog = useDialog();

    const emit = defineEmits(['refresh', 'add', 'update', 'delete']);

    const props = defineProps<Props>();

    const filters = defineModel<TaskStatusesTableFilters>("filters", {
        default: () => ({
            name: "",
        })
    });

    const isFilteredByName = computed<boolean>(() => filters.value.name.length > 0);

    const hasFilters = computed<boolean>(() =>
        isFilteredByName.value
    );

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("modules.taskStatus.components.TaskStatusesTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByName.value,
        }
    ]);

    const onRefresh = () => {
        emit("refresh");
    };

    const onAdd = () => {
        emit("add");
    };

    const onUpdate = (taskStatus: TaskStatus, index: number) => {
        emit("update", taskStatus, index);
    };

    const onConfirmDelete = (taskStatus: TaskStatus, index: number) => {
        dialog.warning({
            title: t("modules.taskStatus.components.TaskStatusesTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash)(24),
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
                emit("delete", taskStatus, index)
            },
        });
    };

    const onClearFilters = () => {
        filters.value.name = "";
    };
</script>

<template>
    <ManageTable size="small" :columns="columns" @refresh="onRefresh" @add="onAdd">
        <template #thead>
            <tr>
                <th>
                    <TextFilterInput clearable :disabled="props.disabled" size="small"
                        :placeholder="t('modules.taskStatus.components.TaskStatusesTable.filters.name.placeholder')"
                        v-model:value="filters.name" />
                </th>
                <th class="doneo-text-center">
                    <ClearFiltersTableButton @clear="onClearFilters" :disabled="props.disabled || !hasFilters" />
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="taskStatus, index in items" :key="taskStatus.id ?? index">
                <td>
                    <n-tag :color="getNaiveUITagColorProperty(taskStatus.hexColor ?? '#888888')">{{ taskStatus.name
                        }}</n-tag>
                </td>
                <td class="doneo-text-center">
                    <ManageTableActionButtons show-update show-delete :update-disabled="props.disabled"
                        :delete-disabled="props.disabled" @update="onUpdate(taskStatus, index)"
                        @delete="onConfirmDelete(taskStatus, index)" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="items.length < 1 && !props.disabled">
                    <n-empty :description="t('modules.taskStatus.components.TaskStatusesTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>