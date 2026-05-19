<script setup lang="ts">
    import { h, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NButtonGroup, NButton, NFlex, NEmpty, NIcon, NTooltip } from 'naive-ui';
    import { IconEdit, IconPlus, IconRefresh, IconEyeCheck, IconSquarePlus, IconTrash } from '@tabler/icons-vue';

    import { Role } from '../models/role';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import { type SortOrder } from '../../../shared/types/common';
    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import TableCellHeaderSortIcon from '../../../shared/components/tables/TableCellHeaderSortIcon.vue';

    interface Props {
        loading: boolean;
        roles: Role[];
        sortField: string;
        sortOrder: SortOrder;
    }

    const { t } = useI18n();

    const emit = defineEmits(['refresh', 'add', 'update', 'delete', 'toggleSort', 'textfilterKeydownEnter']);

    const props = defineProps<Props>();

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("RoleNameTableHeader"),
            field: "name",
            sortable: true,
        },
        {
            label: t("RoleProjectPermissionsTableHeader"),
            field: "projectPermissions",
            sortable: false,
            align: "center",
        },
        {
            label: t("RoleTasksPermissionsTableHeader"),
            field: "taskPermissions",
            sortable: false,
            align: "center",
        },
    ]);

    const roleNameFilter = defineModel<string>("roleNameFilter", {
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

    const onUpdate = (role: Role, index: number) => {
        emit("update", role, index);
    };

    const onConfirmDelete = (role: Role, index: number) => {
        dialog.warning({
            title: t("Delete role"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    t("deleteRoleConfirmation", { name: role.name }),
                    h('br'),
                    h('br'),
                    t("Do you want to continue ?"),
                ]),
            positiveText: t("Delete"),
            negativeText: t("Cancel"),
            onPositiveClick: () => {
                emit("delete", role, index)
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
                        v-model:value="roleNameFilter" @keydown-enter="onTextFilterKeyDownEnter" />
                </th>
                <th></th>
                <th></th>
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
            <tr v-for="role, index in roles" :key="role.id">
                <td>
                    <div class="doneo-flex-center-align" style="gap: 8px;">
                        {{ role.name }}
                    </div>
                </td>
                <td class="doneo-text-center">
                    <n-tooltip trigger="hover" v-if="role.permissions.allowUpdateProject">
                        <template #trigger>
                            <n-icon :size="22">
                                <IconEdit class="doneo-cursor-help" />
                            </n-icon>
                        </template>
                        {{ t("Update project allowed") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover" v-if="role.permissions.allowDeleteProject">
                        <template #trigger>
                            <n-icon :size="22">
                                <IconTrash class="doneo-cursor-help" />
                            </n-icon>
                        </template>
                        {{ t("Delete project allowed") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover" v-if="role.permissions.allowViewProject">
                        <template #trigger>
                            <n-icon :size="22">
                                <IconEyeCheck class="doneo-cursor-help" />
                            </n-icon>
                        </template>
                        {{ t("View project allowed") }}
                    </n-tooltip>
                </td>
                <td class="doneo-text-center">
                    <n-tooltip trigger="hover" v-if="role.permissions.allowAddTask">
                        <template #trigger>
                            <n-icon :size="22">
                                <IconSquarePlus class="doneo-cursor-help" />
                            </n-icon>
                        </template>
                        {{ t("Add task allowed") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover" v-if="role.permissions.allowUpdateTask">
                        <template #trigger>
                            <n-icon :size="22">
                                <IconEdit class="doneo-cursor-help" />
                            </n-icon>
                        </template>
                        {{ t("Update task allowed") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover" v-if="role.permissions.allowDeleteTask">
                        <template #trigger>
                            <n-icon :size="22">
                                <IconTrash class="doneo-cursor-help" />
                            </n-icon>
                        </template>
                        {{ t("Delete task allowed") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover" v-if="role.permissions.allowViewTask">
                        <template #trigger>
                            <n-icon :size="22">
                                <IconEyeCheck class="doneo-cursor-help" />
                            </n-icon>
                        </template>
                        {{ t("View task allowed") }}
                    </n-tooltip>
                </td>
                <td class="doneo-text-center">
                    <n-button-group size="small">
                        <n-button @click="onUpdate(role, index)" :disabled="props.loading">
                            {{ t("Update") }}
                            <template #icon>
                                <n-icon :size="22">
                                    <IconEdit />
                                </n-icon>
                            </template>
                        </n-button>
                        <n-button @click="onConfirmDelete(role, index)" :disabled="props.loading">
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
                <td :colspan="columns.length + 1" v-if="roles.length < 1 && !props.loading">
                    <n-empty :description="t('No roles found')">
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