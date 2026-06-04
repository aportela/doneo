<script setup lang="ts">
    import { h, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NEmpty, NTag } from 'naive-ui';
    import { IconTrash } from '@tabler/icons-vue';

    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import type { Sort } from '../../../shared/types/models/sort.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { ProjectStatusesTableFilters } from '../types/project-statuses-table-filters.ts';
    import { ProjectStatus } from '../models/project-status';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import ClearFiltersTableButton from '../../../shared/components/tables/ClearFiltersTableButton.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import { getNaiveUITagColorProperty } from '../../../shared/composables/color';

    interface Props {
        disabled: boolean;
        items: ProjectStatus[];
        sort?: Sort;
    }

    const { t } = useI18n();
    const dialog = useDialog();

    const emit = defineEmits(['refresh', 'add', 'update', 'delete', 'sort']);

    const props = defineProps<Props>();

    const filters = defineModel<ProjectStatusesTableFilters>("filters", {
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
            label: t("modules.projectStatus.components.ProjectStatusesTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByName.value,
        }
    ]);

    const onSort = (sort: Sort) => {
        emit("sort", sort);
    };

    const onRefresh = () => {
        emit("refresh");
    };

    const onAdd = () => {
        emit("add");
    };

    const onUpdate = (projectStatus: ProjectStatus, index: number) => {
        emit("update", projectStatus, index);
    };

    const onConfirmDelete = (projectStatus: ProjectStatus, index: number) => {
        dialog.warning({
            title: t("modules.projectStatus.components.ProjectStatusesTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    t("modules.projectStatus.components.ProjectStatusesTable.dialogs.deleteConfirmation.message", { name: projectStatus.name }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                emit("delete", projectStatus, index)
            },
        });
    };

    const onClearFilters = () => {
        filters.value.name = "";
    };
</script>

<template>
    <ManageTable size="small" :columns="columns" :current-sort="sort" @sort="onSort" @refresh="onRefresh" @add="onAdd">
        <template #thead>
            <tr>
                <th>
                    <TextFilterInput clearable size="small"
                        :placeholder="t('modules.projectStatus.components.ProjectStatusesTable.filters.name.placeholder')"
                        v-model:value="filters.name" />
                </th>
                <th class="doneo-text-center">
                    <ClearFiltersTableButton @clear="onClearFilters" :disabled="props.disabled || !hasFilters" />
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="projectStatus, index in items" :key="projectStatus.id ?? index">
                <td>
                    <n-tag :color="getNaiveUITagColorProperty(projectStatus.hexColor ?? '#888888')">{{
                        projectStatus.name }}</n-tag>
                </td>
                <td class="doneo-text-center">
                    <ManageTableActionButtons show-update show-delete :update-disabled="props.disabled"
                        :delete-disabled="props.disabled" @update="onUpdate(projectStatus, index)"
                        @delete="onConfirmDelete(projectStatus, index)" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="items.length < 1 && !props.disabled">
                    <n-empty
                        :description="t('modules.projectStatus.components.ProjectStatusesTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>