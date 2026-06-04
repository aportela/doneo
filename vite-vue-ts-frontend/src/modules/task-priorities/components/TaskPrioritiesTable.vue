<script setup lang="ts">
    import { h, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NEmpty, NTag } from 'naive-ui';
    import { IconTrash } from '@tabler/icons-vue';

    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { TaskPrioritiesTableFilters } from '../types/task-priorities-table-filters.ts';
    import { TaskPriority } from '../models/task-priority';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import ClearFiltersTableButton from '../../../shared/components/tables/ClearFiltersTableButton.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import { getNaiveUITagColorProperty } from '../../../shared/composables/color';

    interface Props {
        disabled: boolean;
        items: TaskPriority[];
    }

    const { t } = useI18n();
    const dialog = useDialog();

    const emit = defineEmits(['refresh', 'add', 'update', 'delete']);

    const props = defineProps<Props>();

    const filters = defineModel<TaskPrioritiesTableFilters>("filters", {
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
            label: t("modules.taskPriority.components.TaskPrioritiesTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByName.value,
        },
    ]);

    const onRefresh = () => {
        emit("refresh");
    };

    const onAdd = () => {
        emit("add");
    };

    const onUpdate = (taskPriority: TaskPriority, index: number) => {
        emit("update", taskPriority, index);
    };

    const onConfirmDelete = (taskPriority: TaskPriority, index: number) => {
        dialog.warning({
            title: t("modules.taskPriority.components.TaskPrioritiesTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash)(24),
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
                emit("delete", taskPriority, index)
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
                        :placeholder="t('modules.taskPriority.components.TaskPrioritiesTable.filters.name.placeholder')"
                        v-model:value="filters.name" />
                </th>
                <th class="doneo-text-center">
                    <ClearFiltersTableButton @clear="onClearFilters" :disabled="props.disabled || !hasFilters" />
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="taskPriority, index in items" :key="taskPriority.id ?? index">
                <td>
                    <n-tag :color="getNaiveUITagColorProperty(taskPriority.hexColor ?? '#888888')">{{ taskPriority.name
                    }}</n-tag>
                </td>
                <td class="doneo-text-center">
                    <ManageTableActionButtons show-update show-delete :update-disabled="props.disabled"
                        :delete-disabled="props.disabled" :disabled="props.disabled"
                        @update="onUpdate(taskPriority, index)" @delete="onConfirmDelete(taskPriority, index)" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="items.length < 1 && !props.disabled">
                    <n-empty
                        :description="t('modules.taskPriority.components.TaskPrioritiesTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>