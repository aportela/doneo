<script setup lang="ts">
    import { h, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NFlex, NEmpty, NTag } from 'naive-ui';
    import { IconTrash } from '@tabler/icons-vue';

    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import type { Sort } from '../../../shared/types/models/sort.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { ProjectTypesTableFilters } from '../types/project-types-table-filters.ts';
    import { ProjectType } from '../models/project-type';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import ClearFiltersTableButton from '../../../shared/components/tables/ClearFiltersTableButton.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import { getNaiveUITagColorProperty } from '../../../shared/composables/color';

    interface Props {
        disabled: boolean;
        items: ProjectType[];
        sort?: Sort;
    }

    const { t } = useI18n();
    const dialog = useDialog();

    const emit = defineEmits(['refresh', 'add', 'update', 'delete', 'sort']);

    const props = defineProps<Props>();


    const filters = defineModel<ProjectTypesTableFilters>("filters", {
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
            label: t("modules.projectType.components.ProjectTypesTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByName.value,
        },
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

    const onUpdate = (projectType: ProjectType, index: number) => {
        emit("update", projectType, index);
    };

    const onConfirmDelete = (projectType: ProjectType, index: number) => {
        dialog.warning({
            title: t("modules.projectType.components.ProjectTypesTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    t("modules.projectType.components.ProjectTypesTable.dialogs.deleteConfirmation.message", { name: projectType.name }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                emit("delete", projectType, index)
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
                        :placeholder="t('modules.projectType.components.ProjectTypesTable.filters.name.placeholder')"
                        v-model:value="filters.name" />
                </th>
                <th class="doneo-text-center">
                    <ClearFiltersTableButton @clear="onClearFilters" :disabled="props.disabled || !hasFilters" />
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="projectType, index in items" :key="projectType.id ?? index">
                <td>
                    <n-tag :color="getNaiveUITagColorProperty(projectType.hexColor ?? '#888888')">{{ projectType.name
                        }}</n-tag>
                </td>
                <td class="doneo-text-center">
                    <ManageTableActionButtons show-update show-delete :update-disabled="props.disabled"
                        :delete-disabled="props.disabled" @update="onUpdate(projectType, index)"
                        @delete="onConfirmDelete(projectType, index)" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="items.length < 1 && !props.disabled">
                    <n-empty
                        :description="t('modules.projectType.components.ProjectTypesTable.warnings.noItemsFound')" />
                </td>
            </tr>

        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>