<script setup lang="ts">
    import { h, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NFlex, NEmpty, NIcon, NTooltip } from 'naive-ui';
    import { IconEdit, IconEyeCheck, IconSquarePlus, IconTrash } from '@tabler/icons-vue';

    import { Role } from '../models/role';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import { type SortOrder } from '../../../shared/types/common';
    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import TableCellHeaderSortIcon from '../../../shared/components/tables/TableCellHeaderSortIcon.vue';
    import UpdateDeleteActionsColumn from '../../../shared/components/tables/UpdateDeleteActionsColumn.vue';
    import RefreshAddActionsColumn from '../../../shared/components/tables/RefreshAddActionsColumn.vue';

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
            label: t("modules.role.components.RolesTable.header.columns.name"),
            field: "name",
            sortable: true,
        },
        {
            label: t("modules.role.components.RolesTable.header.columns.projectPermissions"),
            field: "projectPermissions",
            sortable: false,
            align: "center",
        },
        {
            label: t("modules.role.components.RolesTable.header.columns.taskPermissions"),
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
            title: t("modules.role.components.RolesTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    t("modules.role.components.RolesTable.dialogs.deleteConfirmation.message", { name: role.name }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
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
                <th class="doneo-table-actions-column">{{ t("shared.components.table.header.columns.actions") }}</th>
            </tr>
            <tr class="hide-mobile">
                <th>
                    <TextFilterInput clearable size="small"
                        :placeholder="t('modules.role.components.RolesTable.filters.name.placeholder')"
                        v-model:value="roleNameFilter" @keydown-enter="onTextFilterKeyDownEnter" />
                </th>
                <th></th>
                <th></th>
                <th class="doneo-text-center">
                    <RefreshAddActionsColumn @refresh="onRefresh" @add="onAdd" />
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
                        {{ t("modules.role.components.RolesTable.body.columns.permissionsHints.updateProjectAllowed") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover" v-if="role.permissions.allowDeleteProject">
                        <template #trigger>
                            <n-icon :size="22">
                                <IconTrash class="doneo-cursor-help" />
                            </n-icon>
                        </template>
                        {{ t("modules.role.components.RolesTable.body.columns.permissionsHints.deleteProjectAllowed") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover" v-if="role.permissions.allowViewProject">
                        <template #trigger>
                            <n-icon :size="22">
                                <IconEyeCheck class="doneo-cursor-help" />
                            </n-icon>
                        </template>
                        {{ t("modules.role.components.RolesTable.body.columns.permissionsHints.viewProjectAllowed") }}
                    </n-tooltip>
                </td>
                <td class="doneo-text-center">
                    <n-tooltip trigger="hover" v-if="role.permissions.allowAddTask">
                        <template #trigger>
                            <n-icon :size="22">
                                <IconSquarePlus class="doneo-cursor-help" />
                            </n-icon>
                        </template>
                        {{ t("modules.role.components.RolesTable.body.columns.permissionsHints.addTaskAllowed") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover" v-if="role.permissions.allowUpdateTask">
                        <template #trigger>
                            <n-icon :size="22">
                                <IconEdit class="doneo-cursor-help" />
                            </n-icon>
                        </template>
                        {{ t("modules.role.components.RolesTable.body.columns.permissionsHints.updateTaskAllowed") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover" v-if="role.permissions.allowDeleteTask">
                        <template #trigger>
                            <n-icon :size="22">
                                <IconTrash class="doneo-cursor-help" />
                            </n-icon>
                        </template>
                        {{ t("modules.role.components.RolesTable.body.columns.permissionsHints.deleteTaskAllowed") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover" v-if="role.permissions.allowViewTask">
                        <template #trigger>
                            <n-icon :size="22">
                                <IconEyeCheck class="doneo-cursor-help" />
                            </n-icon>
                        </template>
                        {{ t("modules.role.components.RolesTable.body.columns.permissionsHints.viewTaskAllowed") }}
                    </n-tooltip>
                </td>
                <td class="doneo-text-center">
                    <UpdateDeleteActionsColumn @update="onUpdate(role, index)" @delete="onConfirmDelete(role, index)" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="roles.length < 1 && !props.loading">
                    <n-empty :description="t('modules.role.components.RolesTable.warnings.noItemsFound')">
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