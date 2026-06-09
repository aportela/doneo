<script setup lang="ts">
    import { h, ref, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NEmpty, NTooltip, NIcon } from 'naive-ui';
    import { IconTrash, IconEdit, IconEyeCheck, IconSquarePlus } from '@tabler/icons-vue';

    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { ProjectPermissionsTableFilters } from '../types/project-permissions-table-filter.ts';

    import type { ReseteableComponent } from '../../../shared/types/ReseteableComponent.ts';
    import { ProjectPermission } from '../models/project-permission.ts';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import ClearFiltersTableButton from '../../../shared/components/tables/ClearFiltersTableButton.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import RoleSelector from '../../roles/components/RoleSelector.vue';
    import ProjectPermissionSelect from '../../../shared/components/selectors/ProjectPermissionSelect.vue';
    import TaskPermissionSelect from '../../../shared/components/selectors/TaskPermissionSelect.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';

    interface Props {
        disabled: boolean;
        items: ProjectPermission[];
        projectId: string;
        errorMessage?: string | null;
    }

    const { t } = useI18n();
    const dialog = useDialog();

    const emit = defineEmits(['refresh', 'add', 'delete']);

    const props = defineProps<Props>();

    const projectPermissionSelectorRef = ref<ReseteableComponent | undefined>();
    const taskPermissionSelectorRef = ref<ReseteableComponent | undefined>();

    const filters = defineModel<ProjectPermissionsTableFilters>("filters", {
        default: () => ({
            userId: null,
            roleId: null,
            projectPermission: null,
            taskPermission: null,
        })
    });

    const isFilteredByUser = computed<boolean>(() => filters.value.userId !== null);
    const isFilteredByRole = computed<boolean>(() => filters.value.roleId !== null);
    const isFilteredByProjectPermission = computed<boolean>(() => filters.value.projectPermission !== null);
    const isFilteredByTaskPermission = computed<boolean>(() => filters.value.taskPermission !== null);

    const hasFilters = computed<boolean>(() =>
        isFilteredByUser.value || isFilteredByRole.value || isFilteredByProjectPermission.value || isFilteredByTaskPermission.value
    );

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("modules.projectPermission.components.projectPermissionsTable.header.columns.user"),
            field: "user",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByUser.value,
        },
        {
            label: t("modules.projectPermission.components.projectPermissionsTable.header.columns.role"),
            field: "role",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByRole.value,
        },
        {
            label: t("modules.projectPermission.components.projectPermissionsTable.header.columns.projectPermissions"),
            field: "projectPermissions",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByProjectPermission.value,
            align: "center",
        },
        {
            label: t("modules.projectPermission.components.projectPermissionsTable.header.columns.taskPermissions"),
            field: "taskPermissions",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByTaskPermission.value,
            align: "center",
        },
    ]);

    const onRefresh = () => {
        emit("refresh");
    };

    const onAdd = () => {
        emit("add");
    };

    const onConfirmDelete = (projectPermission: ProjectPermission, index: number) => {
        dialog.warning({
            title: t("modules.projectPermission.components.projectPermissionsTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    t("modules.projectPermission.components.projectPermissionsTable.dialogs.deleteConfirmation.message", { user: projectPermission.user.name, role: projectPermission.role.name }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                emit("delete", projectPermission, index)
            },
        });
    };

    const onClearFilters = () => {
        filters.value.userId = null;
        filters.value.roleId = null;
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
                    <UserSelector hideAvatar clearable :disabled="props.disabled" v-model:id="filters.userId"
                        :placeholder="t('modules.projectPermission.components.projectPermissionsTable.filters.user.placeholder')" />
                </th>
                <th>
                    <RoleSelector clearable :disabled="props.disabled" v-model:id="filters.roleId"
                        :placeholder="t('modules.projectPermission.components.projectPermissionsTable.filters.role.placeholder')" />
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
        <template #tbody v-if="!props.errorMessage">
            <tr v-for="projectPermission, index in items" :key="projectPermission.id ?? index">
                <td>
                    <AvatarUserName :user-id="projectPermission.user?.id ?? ''"
                        :user-name="projectPermission.user?.name ?? ''" />
                </td>
                <td>{{ projectPermission.role.name }}</td>
                <td class="doneo-text-center">
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconEdit"
                                :class="{ 'doneo-disabled-icon': !projectPermission.role.permissions.allowUpdateProject }" />
                        </template>
                        {{ t(projectPermission.role.permissions.allowUpdateProject ?
                            "modules.projectPermission.components.projectPermissionsTable.body.columns.permissionsHints.updateProjectAllowed"
                            :
                            "modules.projectPermission.components.projectPermissionsTable.body.columns.permissionsHints.updateProjectDenied")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconTrash"
                                :class="{ 'doneo-disabled-icon': !projectPermission.role.permissions.allowDeleteProject }" />
                        </template>
                        {{ t(projectPermission.role.permissions.allowDeleteProject ?
                            "modules.projectPermission.components.projectPermissionsTable.body.columns.permissionsHints.deleteProjectAllowed"
                            :
                            "modules.projectPermission.components.projectPermissionsTable.body.columns.permissionsHints.deleteProjectDenied")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconEyeCheck"
                                :class="{ 'doneo-disabled-icon': !projectPermission.role.permissions.allowViewProject }" />
                        </template>
                        {{ t(projectPermission.role.permissions.allowViewProject ?
                            "modules.projectPermission.components.projectPermissionsTable.body.columns.permissionsHints.viewProjectAllowed"
                            :
                            "modules.projectPermission.components.projectPermissionsTable.body.columns.permissionsHints.viewProjectDenied")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconSquarePlus"
                                :class="{ 'doneo-disabled-icon': !projectPermission.role.permissions.allowAddTask }" />
                        </template>
                        {{ t(projectPermission.role.permissions.allowAddTask ?
                            "modules.projectPermission.components.projectPermissionsTable.body.columns.permissionsHints.addTaskAllowed"
                            :
                            "modules.projectPermission.components.projectPermissionsTable.body.columns.permissionsHints.addTaskDenied")
                        }}
                    </n-tooltip>
                </td>
                <td class="doneo-text-center">
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconEdit"
                                :class="{ 'doneo-disabled-icon': !projectPermission.role.permissions.allowUpdateTask }" />
                        </template>
                        {{ t(projectPermission.role.permissions.allowUpdateTask ?
                            "modules.projectPermission.components.projectPermissionsTable.body.columns.permissionsHints.updateTaskAllowed"
                            :
                            "modules.projectPermission.components.projectPermissionsTable.body.columns.permissionsHints.updateTaskDenied")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconTrash"
                                :class="{ 'doneo-disabled-icon': !projectPermission.role.permissions.allowDeleteTask }" />
                        </template>
                        {{ t(projectPermission.role.permissions.allowDeleteTask ?
                            "modules.projectPermission.components.projectPermissionsTable.body.columns.permissionsHints.deleteTaskAllowed"
                            :
                            "modules.projectPermission.components.projectPermissionsTable.body.columns.permissionsHints.deleteTaskDenied")
                        }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconEyeCheck"
                                :class="{ 'doneo-disabled-icon': !projectPermission.role.permissions.allowViewTask }" />
                        </template>
                        {{ t(projectPermission.role.permissions.allowViewTask ?
                            "modules.projectPermission.components.projectPermissionsTable.body.columns.permissionsHints.viewTaskAllowed"
                            :
                            "modules.projectPermission.components.projectPermissionsTable.body.columns.permissionsHints.viewTaskDenied")
                        }}
                    </n-tooltip>
                </td>
                <td class="doneo-text-center">
                    <ManageTableActionButtons show-delete @delete="onConfirmDelete(projectPermission, index)"
                        :disabled="props.disabled" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="items.length < 1 && !props.disabled">
                    <n-empty
                        :description="t('modules.projectPermission.components.projectPermissionsTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>