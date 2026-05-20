<script setup lang="ts">
    import { h, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NButtonGroup, NButton, NFlex, NEmpty, NIcon, NTag } from 'naive-ui';
    import { IconEdit, IconPlus, IconRefresh, IconTrash } from '@tabler/icons-vue';

    import { ProjectPriority } from '../models/project-priority';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import { type SortOrder } from '../../../shared/types/common';
    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import TableCellHeaderSortIcon from '../../../shared/components/tables/TableCellHeaderSortIcon.vue';
    import { getNaiveUITagColorProperty } from '../../../shared/composables/color';

    interface Props {
        loading: boolean;
        projectPriorities: ProjectPriority[];
        sortField: string;
        sortOrder: SortOrder;
    }

    const { t } = useI18n();

    const emit = defineEmits(['refresh', 'add', 'update', 'delete', 'toggleSort', 'textfilterKeydownEnter']);

    const props = defineProps<Props>();

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("ProjectPriorityNameTableHeader"),
            field: "name",
            sortable: true,
        },
    ]);


    const projectPriorityNameFilter = defineModel<string>("projectPriorityNameFilter", {
        default: "",
    });

    const dialog = useDialog();

    const onToggleSort = (field: string) => {
        emit("toggleSort", field);
    };

    const onRefresh = () => {
        emit("refresh");
    };

    const onAdd = () => {
        emit("add");
    };

    const onUpdate = (projectPriority: ProjectPriority, index: number) => {
        emit("update", projectPriority, index);
    };

    const onConfirmDelete = (projectPriority: ProjectPriority, index: number) => {
        dialog.warning({
            title: t("Delete project priority"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    t("deleteProjectPriorityConfirmation", { name: projectPriority.name }),
                    h('br'),
                    h('br'),
                    t("Do you want to continue ?"),
                ]),
            positiveText: t("Delete"),
            negativeText: t("Cancel"),
            onPositiveClick: () => {
                emit("delete", projectPriority, index)
            },
        });
    };

    const onTextFilterKeyDownEnter = () => {
        emit("textfilterKeydownEnter");
    };

</script>

<template>
    <ManageTable size="small">
        <template #thead>
            <tr class="table-header-click-action">
                <th v-for="column in columns" :key="column.field" @click="column.sortable && onToggleSort(column.field)"
                    :class="{ 'doneo-cursor-pointer': column.sortable, 'doneo-text-center': column.align === 'center' }">
                    <n-flex justify="space-between" v-if="column.sortable">
                        <span>{{ column.label }}</span>
                        <TableCellHeaderSortIcon v-if="props.sortField === column.field" :order="props.sortOrder" />
                    </n-flex>
                    <span v-else>{{ column.label }}</span>
                </th>
                <th class="doneo-table-actions-column">{{ t("Actions") }}</th>
            </tr>
            <tr class="hide-mobile">
                <th>
                    <TextFilterInput clearable size="small" :placeholder="t('searchByNameDefaultPlaceholder')"
                        v-model:value="projectPriorityNameFilter" @keydown-enter="onTextFilterKeyDownEnter" />
                </th>
                <th class="doneo-text-center">
                    <n-button-group size="small">
                        <n-button @click="onRefresh">
                            <template #icon>
                                <n-icon :size="22">
                                    <IconRefresh />
                                </n-icon>
                            </template>
                            {{ t("Refresh") }}
                        </n-button>
                        <n-button @click="onAdd">
                            <template #icon>
                                <n-icon :size="22">
                                    <IconPlus />
                                </n-icon>
                            </template>
                            {{ t("Add") }}
                        </n-button>
                    </n-button-group>
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="projectPriority, index in projectPriorities" :key="projectPriority.id">
                <td>
                    <n-tag :color="getNaiveUITagColorProperty(projectPriority.hexColor)">{{ projectPriority.name
                    }}</n-tag>
                </td>
                <td class="doneo-text-center">
                    <n-button-group size="small">
                        <n-button @click="onUpdate(projectPriority, index)" :disabled="props.loading">
                            {{ t("Update") }}
                            <template #icon>
                                <n-icon :size="22">
                                    <IconEdit />
                                </n-icon>
                            </template>
                        </n-button>
                        <n-button @click="onConfirmDelete(projectPriority, index)" :disabled="props.loading">
                            {{ t("Delete") }}
                            <template #icon>
                                <n-icon :size="22">
                                    <IconTrash />
                                </n-icon>
                            </template>
                        </n-button>
                    </n-button-group>
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="projectPriorities.length < 1 && !props.loading">
                    <n-empty :description="t('No project priorities found')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped>

    .table-header-click-action th:not(:last-of-type) .n-icon {
        margin-top: 4px;
    }

    @media (max-width: 768px) {
        .hide-mobile {
            display: none;
        }
    }
</style>