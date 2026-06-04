<script setup lang="ts">
    import { h, ref, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NEmpty, NIcon, NTooltip } from 'naive-ui';
    import { IconEdit, IconEyeCheck, IconSquarePlus, IconTrash } from '@tabler/icons-vue';

    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { RolesTableFilters } from '../types/roles-table-filters.ts';

    import type { ReseteableComponent } from '../../../shared/types/ReseteableComponent.ts';
    import { Role } from '../models/role';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import ProjectPermissionSelect from '../../../shared/components/selectors/ProjectPermissionSelect.vue';
    import TaskPermissionSelect from '../../../shared/components/selectors/TaskPermissionSelect.vue';
    import ClearFiltersTableButton from '../../../shared/components/tables/ClearFiltersTableButton.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';

    interface Props {
        disabled: boolean;
        items: Role[];
    }

    const { t } = useI18n();
    const dialog = useDialog();

    const emit = defineEmits(['refresh', 'add', 'update', 'delete']);

    const props = defineProps<Props>();

    const projectPermissionSelectorRef = ref<ReseteableComponent | undefined>();
    const taskPermissionSelectorRef = ref<ReseteableComponent | undefined>();

    const filters = defineModel<RolesTableFilters>("filters", {
        default: () => ({
            name: "",
            projectPermission: null,
            taskPermission: null,
        })
    });

    const isFilteredByName = computed<boolean>(() => filters.value.name.length > 0);
    const isFilteredByProjectPermission = computed<boolean>(() => filters.value.projectPermission !== null);
    const isFilteredByTaskPermission = computed<boolean>(() => filters.value.taskPermission !== null);

    const hasFilters = computed<boolean>(() =>
        isFilteredByName.value || isFilteredByProjectPermission.value || isFilteredByTaskPermission.value
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
            isFiltered: () => isFilteredByProjectPermission.value,
        },
        {
            label: t("modules.role.components.RolesTable.header.columns.taskPermissions"),
            field: "taskPermissions",
            visible: true,
            sortable: false,
            align: "center",
            isFiltered: () => isFilteredByTaskPermission.value,
        },
    ]);

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
        projectPermissionSelectorRef.value?.reset();
        taskPermissionSelectorRef.value?.reset();
    };

    const permissionIconSize = 22;
</script>

<template>
    <ManageTable size="small" :columns="columns" @refresh="onRefresh" @add="onAdd">
        <template #thead>
            <tr>
                <th>
                    <TextFilterInput clearable :disabled="props.disabled" size="small"
                        :placeholder="t('modules.role.components.RolesTable.filters.name.placeholder')"
                        v-model:value="filters.name" />
                </th>
                <th>
                    <ProjectPermissionSelect v-model:permission="filters.projectPermission"
                        :placeholder="t('shared.components.selectors.ProjectPermissionSelect.placeholder')" clearable
                        ref="projectPermissionSelectorRef" />
                </th>
                <th>
                    <TaskPermissionSelect v-model:permission="filters.taskPermission"
                        :placeholder="t('shared.components.selectors.TaskPermissionSelect.placeholder')" clearable
                        ref="taskPermissionSelectorRef" />
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
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconSquarePlus"
                                :class="{ 'doneo-disabled-permission-icon': !role.permissions.allowAddTask }" />
                        </template>
                        {{ t(role.permissions.allowAddTask ?
                            "modules.role.components.RolesTable.body.columns.permissionsHints.addTaskAllowed" :
                            "modules.role.components.RolesTable.body.columns.permissionsHints.addTaskDenied") }}
                    </n-tooltip>
                </td>
                <td class="doneo-text-center">
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
                        :delete-disabled="props.disabled" :disabled="props.disabled" @update="onUpdate(role, index)"
                        @delete="onConfirmDelete(role, index)" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="!props.disabled && items.length < 1">
                    <n-empty :description="t('modules.role.components.RolesTable.warnings.noItemsFound')" />
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