<script setup lang="ts">
    import { h, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NEmpty, NIcon, NTooltip, NSelect, type SelectOption } from 'naive-ui';
    import { IconEdit, IconEyeCheck, IconSquarePlus, IconTrash } from '@tabler/icons-vue';

    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import type { Sort } from '../../../shared/types/models/sort.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { RolesTableFilters } from '../types/roles-table-filters.ts';
    import { Role } from '../models/role';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import ClearFiltersTableButton from '../../../shared/components/tables/ClearFiltersTableButton.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';

    interface Props {
        disabled: boolean;
        items: Role[];
        sort?: Sort;
    }

    const { t } = useI18n();
    const dialog = useDialog();

    const emit = defineEmits(['refresh', 'add', 'update', 'delete', 'sort']);

    const props = defineProps<Props>();

    const filters = defineModel<RolesTableFilters>("filters", {
        default: () => ({
            name: "",
            allowedProjectPermissions: [],
            allowedTaskPermissions: [],
        })
    });

    const isFilteredByName = computed<boolean>(() => filters.value.name.length > 0);
    const isFilteredByAllowedProjectPermissions = computed<boolean>(() => filters.value.allowedProjectPermissions.length > 0);
    const isFilteredByAllowedTaskPermissions = computed<boolean>(() => filters.value.allowedTaskPermissions.length > 0);


    const hasFilters = computed<boolean>(() =>
        isFilteredByName.value || isFilteredByAllowedProjectPermissions.value || isFilteredByAllowedTaskPermissions.value
    );

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("modules.role.components.RolesTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByName.value,
        },
        {
            label: t("modules.role.components.RolesTable.header.columns.projectPermissions"),
            field: "projectPermissions",
            visible: true,
            sortable: false,
            align: "center",
            isFiltered: () => isFilteredByAllowedProjectPermissions.value,
        },
        {
            label: t("modules.role.components.RolesTable.header.columns.taskPermissions"),
            field: "taskPermissions",
            visible: true,
            sortable: false,
            align: "center",
            isFiltered: () => isFilteredByAllowedTaskPermissions.value,
        },
    ]);

    // TODO:
    const projectPermissionsSelectorOptions = computed<SelectOption[]>(() => [
        {
            label: t("modules.role.components.RolesTable.body.columns.permissionsHints.updateProjectAllowed"),
            value: 1
        },
        {
            label: t("modules.role.components.RolesTable.body.columns.permissionsHints.deleteProjectAllowed"),
            value: 2
        },
        {
            label: t("modules.role.components.RolesTable.body.columns.permissionsHints.viewProjectAllowed"),
            value: 3
        },
    ]);

    // TODO:
    const taskPermissionsSelectorOptions = computed<SelectOption[]>(() => [
        {
            label: t("modules.role.components.RolesTable.body.columns.permissionsHints.addTaskAllowed"),
            value: 1
        },
        {
            label: t("modules.role.components.RolesTable.body.columns.permissionsHints.updateTaskAllowed"),
            value: 2
        },
        {
            label: t("modules.role.components.RolesTable.body.columns.permissionsHints.deleteTaskAllowed"),
            value: 3
        },
        {
            label: t("modules.role.components.RolesTable.body.columns.permissionsHints.viewTaskAllowed"),
            value: 4
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

    const onClearFilters = () => {
        filters.value.name = "";
        filters.value.allowedProjectPermissions = [];
        filters.value.allowedTaskPermissions = [];
    };

    const permissionIconSize = 22;
</script>

<template>
    <ManageTable size="small" :columns="columns" :current-sort="sort" @sort="onSort" @refresh="onRefresh" @add="onAdd">
        <template #thead>
            <tr>
                <th>
                    <TextFilterInput clearable size="small"
                        :placeholder="t('modules.role.components.RolesTable.filters.name.placeholder')"
                        v-model:value="filters.name" />
                </th>
                <th>
                    <n-select size="small" clearable multiple :options="projectPermissionsSelectorOptions"
                        v-model:value="filters.allowedProjectPermissions" disabled />
                </th>
                <th>
                    <n-select size="small" clearable multiple :options="taskPermissionsSelectorOptions"
                        v-model:value="filters.allowedTaskPermissions" disabled />
                </th>
                <th class="doneo-text-center">
                    <ClearFiltersTableButton @clear="onClearFilters" :disabled="props.disabled || !hasFilters" />
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="role, index in items" :key="role.id ?? index">
                <td>
                    <div class="doneo-flex-center-align" style="gap: 8px;">
                        {{ role.name }}
                    </div>
                </td>
                <td class="doneo-text-center">
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconEdit"
                                :class="{ 'doneo-disabled-permission-icon': !role.permissions.allowUpdateProject }" />
                        </template>
                        {{ t(role.permissions.allowUpdateProject ?
                            "modules.role.components.RolesTable.body.columns.permissionsHints.updateProjectAllowed" :
                            "modules.role.components.RolesTable.body.columns.permissionsHints.updateProjectDenied") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconTrash"
                                :class="{ 'doneo-disabled-permission-icon': !role.permissions.allowDeleteProject }" />
                        </template>
                        {{ t(role.permissions.allowDeleteProject ?
                            "modules.role.components.RolesTable.body.columns.permissionsHints.deleteProjectAllowed" :
                            "modules.role.components.RolesTable.body.columns.permissionsHints.deleteProjectDenied") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconEyeCheck"
                                :class="{ 'doneo-disabled-permission-icon': !role.permissions.allowViewProject }" />
                        </template>
                        {{ t(role.permissions.allowViewProject ?
                            "modules.role.components.RolesTable.body.columns.permissionsHints.viewProjectAllowed" :
                            "modules.role.components.RolesTable.body.columns.permissionsHints.viewProjectDenied") }}
                    </n-tooltip>
                </td>
                <td class="doneo-text-center">
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconSquarePlus"
                                :class="{ 'doneo-disabled-permission-icon': !role.permissions.allowAddTask }" />
                        </template>
                        {{ t(role.permissions.allowAddTask ?
                            "modules.role.components.RolesTable.body.columns.permissionsHints.addTaskAllowed" :
                            "modules.role.components.RolesTable.body.columns.permissionsHints.addTaskDenied") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconEdit"
                                :class="{ 'doneo-disabled-permission-icon': !role.permissions.allowUpdateTask }" />
                        </template>
                        {{ t(role.permissions.allowUpdateTask ?
                            "modules.role.components.RolesTable.body.columns.permissionsHints.updateTaskAllowed" :
                            "modules.role.components.RolesTable.body.columns.permissionsHints.updateTaskDenied") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconTrash"
                                :class="{ 'doneo-disabled-permission-icon': !role.permissions.allowDeleteTask }" />
                        </template>
                        {{ t(role.permissions.allowDeleteTask ?
                            "modules.role.components.RolesTable.body.columns.permissionsHints.deleteTaskAllowed" :
                            "modules.role.components.RolesTable.body.columns.permissionsHints.deleteTaskDenied") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconEyeCheck"
                                :class="{ 'doneo-disabled-permission-icon': !role.permissions.allowViewTask }" />
                        </template>
                        {{ t(role.permissions.allowViewTask ?
                            "modules.role.components.RolesTable.body.columns.permissionsHints.viewTaskAllowed" :
                            "modules.role.components.RolesTable.body.columns.permissionsHints.viewTaskDenied") }}
                    </n-tooltip>
                </td>
                <td class="doneo-text-center">
                    <ManageTableActionButtons show-update show-delete :update-disabled="props.disabled"
                        :delete-disabled="props.disabled" @update="onUpdate(role, index)"
                        @delete="onConfirmDelete(role, index)" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="!props.disabled && items.length < 1">
                    <n-empty :description="t('modules.role.components.RolesTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped>
    .doneo-disabled-permission-icon {
        opacity: 0.1;
    }
</style>