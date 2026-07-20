<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NEmpty, NIcon, NTooltip } from 'naive-ui';
    import { IconEdit, IconEyeCheck, IconSquarePlus, IconTrash } from '@tabler/icons-vue';

    import { useLoadingStore } from '../../../stores/loading';
    import { useCacheStore } from '../../../stores/cache.ts';
    import { useSessionStore } from '../../../stores/session';
    import { useUserSettingsStore } from '../../../stores/userSettings.ts';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import type { Order } from '../../../shared/types/order.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';

    import type { ReseteableComponent } from '../../../shared/types/ReseteableComponent.ts';
    import { Role } from '../models/role';

    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { roleService } from '../services/role.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';
    import ProjectPermissionSelect from '../../../shared/components/selectors/ProjectPermissionSelect.vue';
    import TaskPermissionSelect from '../../../shared/components/selectors/TaskPermissionSelect.vue';
    import ClearTableFiltersButton from '../../../shared/components/buttons/ClearTableFiltersButton.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import type { ProjectPermissionSelectValue } from '../../../shared/types/project-permission-select-value.ts';
    import type { TaskPermissionSelectValue } from '../../../shared/types/task-permission-select-value.ts';
    import type { RoleResponse } from '../types/dto.ts';
    import { ro } from '@nuxt/ui/runtime/locale/index.js';

    interface Props {
        id?: string;
    }

    const props = withDefaults(defineProps<Props>(), { id: "RolesTable" });;

    const { t } = useI18n();
    const dialog = useDialog();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();
    const sessionStore = useSessionStore();
    const userSettingsStore = useUserSettingsStore();
    const tableSettingsStore = useTableSettingsStore();
    const cacheStore = useCacheStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    watch(
        () => state.ajaxRunning,
        (ajaxRunning) => {
            loadingStore.set(ajaxRunning);
        }
    );

    const items = shallowRef<Role[]>([]);
    const tmpRole = ref<Role>(new Role());

    const order = reactive<Order>({ field: "name", direction: "ASC" });

    const onSort = (newOrder: Order) => {
        order.field = newOrder.field;
        order.direction = newOrder.direction;
        onRefresh();
    };

    const projectPermissionSelectorRef = ref<InstanceType<typeof ProjectPermissionSelect>[] | null>(null);
    const taskPermissionSelectorRef = ref<InstanceType<typeof TaskPermissionSelect>[] | null>(null);

    interface RolesTableFilters {
        name: string;
        projectPermission: ProjectPermissionSelectValue | null;
        taskPermission: TaskPermissionSelectValue | null;
    }

    const filters = reactive<RolesTableFilters>(
        {
            name: "",
            projectPermission: null,
            taskPermission: null,
        }
    );

    const isFilteredByName = computed<boolean>(() => filters.name.length > 0);
    const isFilteredByProjectPermission = computed<boolean>(() => filters.projectPermission !== null);
    const isFilteredByTaskPermission = computed<boolean>(() => filters.taskPermission !== null);

    const onClearFilters = () => {
        filters.name = "";
        if (projectPermissionSelectorRef.value) {
            projectPermissionSelectorRef.value[0]?.reset();
        }
        if (taskPermissionSelectorRef.value) {
            taskPermissionSelectorRef.value[0]?.reset();
        }
    };


    const columnDefinitions = reactive<TableHeaderColumn<Role>[]>([
        {
            label: t("modules.role.components.RolesTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByName.value,
            render: (row: Role) => {
                return h(
                    "span",
                    {},
                    {
                        default: () => row.name

                    }
                );
            }
        },
        {
            label: t("modules.role.components.RolesTable.header.columns.projectPermissions"),
            field: "projectPermissions",
            visible: true,
            sortable: false,
            align: "center",
            isFiltered: () => isFilteredByProjectPermission.value,
            render: (row: Role) => { }
        },
        {
            label: t("modules.role.components.RolesTable.header.columns.taskPermissions"),
            field: "taskPermissions",
            visible: true,
            sortable: false,
            align: "center",
            isFiltered: () => isFilteredByTaskPermission.value,
            render: (row: Role) => { }
        },
    ]);

    // create (if not found) default settings for this table (column order & visibility)
    tableSettingsStore.register(props.id, { columns: columnDefinitions.map((column) => { return { field: column.field, visible: column.visible } }) ?? [] });

    // restore previous settings
    const tableSettings = tableSettingsStore.get(props.id);

    // build columns based on saved order visibility settings
    const columns = computed<TableHeaderColumn<Role>[]>(() =>
        tableSettings.columns.map((column) => { // get saved ordered columns
            const definition = columnDefinitions.find((c) => c.field === column.field);
            return {
                label: definition?.label ?? "",
                field: column.field,
                visible: column.visible,
                sortable: definition!.sortable,
                isFiltered: definition?.isFiltered ?? (() => false),
                render: definition?.render ?? (() => "")
            };
        })
    );

    const onConfirmDelete = (role: Role) => {
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

    const onRefresh = async () => { };

    const showFormModal = ref<boolean>(false);

    const onAdd = () => {
        tmpRole.value = new Role();
        showFormModal.value = true;
    };

    const onUpdate = (role: Role) => {
        tmpRole.value = role;
        showFormModal.value = true;
    };

    const onRoleAdded = (role: RoleResponse) => {
        cacheStore.clearUsersCache();
        showFormModal.value = false;
        tmpRole.value = new Role();
        notify('success', t("modules.user.components.UsersTable.notifications.userAdded", { name: user.name }));
        onRefresh();
    };

    const onRoleUpdated = (role: RoleResponse) => {
        cacheStore.clearUsersCache();
        showFormModal.value = false;
        tmpRole.value = new Role();
        notify('success', t("modules.user.components.UsersTable.notifications.userUpdated", { name: user.name }));
        onRefresh();
    };

    const hideRoleForm = () => {
        showFormModal.value = false;
        tmpRole.value = new Role();
    };


    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageUsersPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageUsersPage.onDelete")) {
                onDelete(tmpUser.value);
            } else if (payload.to.includes("ManageUsersPage.onUnDelete")) {
                onUnDelete(tmpUser.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <ManageTable :id="props.id" size="small" :columns="columnDefinitions" @refresh="onRefresh" @add="onAdd">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <TextFilterInput v-if="column.field === 'name'" clearable :disabled="state.ajaxRunning" size="small"
                    :placeholder="t('modules.role.components.RolesTable.filters.name.placeholder')"
                    v-model:value="filters.name" />
                <ProjectPermissionSelect v-if="column.field == 'projectPermissions'"
                    v-model:permission="filters.projectPermission"
                    :placeholder="t('shared.components.selectors.ProjectPermissionSelect.placeholder')" clearable
                    ref="projectPermissionSelectorRef" />
                <TaskPermissionSelect v-if="column.field === 'taskPermissions'"
                    v-model:permission="filters.taskPermission"
                    :placeholder="t('shared.components.selectors.TaskPermissionSelect.placeholder')" clearable
                    ref="taskPermissionSelectorRef" />
            </th>
        </template>
        <template #rowactions="{ row }">
            <ManageTableActionButtons show-update show-delete :update-disabled="state.ajaxRunning"
                :delete-disabled="state.ajaxRunning" :disabled="state.ajaxRunning" @update="onUpdate(row)"
                @delete="onConfirmDelete(row)" />
        </template>
        <!--
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
                                :class="{ 'doneo-disabled-icon': !role.permissions.allowUpdateProject }" />
                        </template>
                        {{ t(role.permissions.allowUpdateProject ?
                            "modules.role.components.RolesTable.body.columns.permissionsHints.updateProjectAllowed" :
                            "modules.role.components.RolesTable.body.columns.permissionsHints.updateProjectDenied") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconTrash"
                                :class="{ 'doneo-disabled-icon': !role.permissions.allowDeleteProject }" />
                        </template>
                        {{ t(role.permissions.allowDeleteProject ?
                            "modules.role.components.RolesTable.body.columns.permissionsHints.deleteProjectAllowed" :
                            "modules.role.components.RolesTable.body.columns.permissionsHints.deleteProjectDenied") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconEyeCheck"
                                :class="{ 'doneo-disabled-icon': !role.permissions.allowViewProject }" />
                        </template>
                        {{ t(role.permissions.allowViewProject ?
                            "modules.role.components.RolesTable.body.columns.permissionsHints.viewProjectAllowed" :
                            "modules.role.components.RolesTable.body.columns.permissionsHints.viewProjectDenied") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconSquarePlus"
                                :class="{ 'doneo-disabled-icon': !role.permissions.allowAddTask }" />
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
                                :class="{ 'doneo-disabled-icon': !role.permissions.allowUpdateTask }" />
                        </template>
                        {{ t(role.permissions.allowUpdateTask ?
                            "modules.role.components.RolesTable.body.columns.permissionsHints.updateTaskAllowed" :
                            "modules.role.components.RolesTable.body.columns.permissionsHints.updateTaskDenied") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconTrash"
                                :class="{ 'doneo-disabled-icon': !role.permissions.allowDeleteTask }" />
                        </template>
                        {{ t(role.permissions.allowDeleteTask ?
                            "modules.role.components.RolesTable.body.columns.permissionsHints.deleteTaskAllowed" :
                            "modules.role.components.RolesTable.body.columns.permissionsHints.deleteTaskDenied") }}
                    </n-tooltip>
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-icon :size="permissionIconSize" class="doneo-cursor-help" :component="IconEyeCheck"
                                :class="{ 'doneo-disabled-icon': !role.permissions.allowViewTask }" />
                        </template>
                        {{ t(role.permissions.allowViewTask ?
                            "modules.role.components.RolesTable.body.columns.permissionsHints.viewTaskAllowed" :
                            "modules.role.components.RolesTable.body.columns.permissionsHints.viewTaskDenied") }}
                    </n-tooltip>
                </td>
                <td class="doneo-text-center">
                    <ManageTableActionButtons show-update show-delete :update-disabled="state.ajaxRunning"
                        :delete-disabled="state.ajaxRunning" :disabled="state.ajaxRunning" @update="onUpdate(role, index)"
                        @delete="onConfirmDelete(role, index)" />
                </td>
            </tr>
            <tr>
                <td :colspan="columnDefinitions.length + 1" v-if="!state.ajaxRunning && items.length < 1">
                    <n-empty :description="t('modules.role.components.RolesTable.warnings.noItemsFound')" />
                </td>
            </tr>
        </template>
        -->
    </ManageTable>
</template>

<style lang="css" scoped></style>