<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h, type Component } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, useDialog, NIcon, NTooltip } from 'naive-ui';
    import { IconTrash } from '@tabler/icons-vue';

    import { useLoadingStore } from '../../../stores/loading';
    import { useCacheStore } from '../../../stores/cache.ts';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import type { Order } from '../../../shared/types/order.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';

    import { Role } from '../models/role';

    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { roleService } from '../services/role.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';
    import ProjectPermissionSelect from '../../../shared/components/selectors/ProjectPermissionSelect.vue';
    import TaskPermissionSelect from '../../../shared/components/selectors/TaskPermissionSelect.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import type { ProjectPermissionSelectValue } from '../../../shared/types/project-permission-select-value.ts';
    import type { TaskPermissionSelectValue } from '../../../shared/types/task-permission-select-value.ts';
    import type { RoleResponse, SearchRequest } from '../types/dto.ts';
    import RoleForm from './RoleForm.vue';
    import { DONEO_ICON_ACTION_ADD, DONEO_ICON_ACTION_DELETE, DONEO_ICON_ACTION_EDIT, DONEO_ICON_ACTION_SHOW } from '../../../shared/types/icons.ts';
    import { renderIcon, renderLabel } from '../../../shared/composables/naive-ui-helpers.ts';

    interface Props {
        id?: string;
    }

    const props = withDefaults(defineProps<Props>(), { id: "RolesTable" });;

    const { t } = useI18n();
    const dialog = useDialog();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();
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
        // we have all results, use local sorting for avoiding server load
        if (order.direction === "ASC") {
            items.value = [...items.value].sort((a, b) =>
                a.name.localeCompare(b.name)
            );
        } else {
            items.value = [...items.value].sort((a, b) =>
                b.name.localeCompare(a.name)
            );
        }
    };

    const projectPermissionFilterRef = ref<InstanceType<typeof ProjectPermissionSelect>[] | null>(null);
    const taskPermissionFilterRef = ref<InstanceType<typeof TaskPermissionSelect>[] | null>(null);

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
        if (projectPermissionFilterRef.value) {
            projectPermissionFilterRef.value[0]?.reset();
        }
        if (taskPermissionFilterRef.value) {
            taskPermissionFilterRef.value[0]?.reset();
        }
    };


    const nameFilterLowerCase = computed(() => filters.name.toLowerCase());

    // we have all results, use local filtering for avoiding server load
    const localFilteredItems = computed<Role[]>(() => {
        return items.value.filter((role: Role) => {
            const name = role.name?.toLowerCase();
            return (
                (!name || name?.includes(nameFilterLowerCase.value)) &&
                (filters.projectPermission === null || (filters.projectPermission !== null && (
                    (filters.projectPermission === "updateProjectAllowed" && role.permissions.allowUpdateProject) ||
                    (filters.projectPermission === "updateProjectDenied" && !role.permissions.allowUpdateProject) ||
                    (filters.projectPermission === "deleteProjectAllowed" && role.permissions.allowDeleteProject) ||
                    (filters.projectPermission === "deleteProjectDenied" && !role.permissions.allowDeleteProject) ||
                    (filters.projectPermission === "viewProjectAllowed" && role.permissions.allowViewProject) ||
                    (filters.projectPermission === "viewProjectDenied" && !role.permissions.allowViewProject) ||
                    (filters.projectPermission === "addTaskAllowed" && role.permissions.allowAddTask) ||
                    (filters.projectPermission === "addTaskDenied" && !role.permissions.allowAddTask)
                ))
                ) &&
                (filters.taskPermission === null || (filters.taskPermission !== null && (
                    (filters.taskPermission === "updateTaskAllowed" && role.permissions.allowUpdateTask) ||
                    (filters.taskPermission === "updateTaskDenied" && !role.permissions.allowUpdateTask) ||
                    (filters.taskPermission === "deleteTaskAllowed" && role.permissions.allowDeleteTask) ||
                    (filters.taskPermission === "deleteTaskDenied" && !role.permissions.allowDeleteTask) ||
                    (filters.taskPermission === "viewTaskAllowed" && role.permissions.allowViewTask) ||
                    (filters.taskPermission === "viewTaskDenied" && !role.permissions.allowViewTask)
                ))
                )
            );
        });
    });

    type PermissionIcon = {
        allowed: boolean;
        icon: Component;
        allowedKey: string;
        deniedKey: string;
    };


    // return rendered permissions array
    const renderPermissionIcons = (permissions: PermissionIcon[]) =>
        h(
            "div",
            { class: "doneo-flex doneo-gap-2" },
            permissions.map((permission) =>
                h(
                    NTooltip,
                    { trigger: "hover" },
                    {
                        trigger: () =>
                            h(NIcon, {
                                size: 20,
                                component: permission.icon,
                                class: [
                                    "doneo-cursor-help",
                                    { "doneo-disabled-icon": !permission.allowed },
                                ],
                            }),
                        default: () =>
                            t(permission.allowed ? permission.allowedKey : permission.deniedKey),
                    }
                )
            )
        );

    // get project permissions array from role row
    const renderProjectPermissionColumn = (row: Role) => {
        return renderPermissionIcons([
            {
                allowed: row.permissions.allowUpdateProject,
                icon: DONEO_ICON_ACTION_EDIT,
                allowedKey:
                    "modules.role.components.RolesTable.body.columns.permissionsHints.updateProjectAllowed",
                deniedKey:
                    "modules.role.components.RolesTable.body.columns.permissionsHints.updateProjectDenied",
            },
            {
                allowed: row.permissions.allowDeleteProject,
                icon: DONEO_ICON_ACTION_DELETE,
                allowedKey:
                    "modules.role.components.RolesTable.body.columns.permissionsHints.deleteProjectAllowed",
                deniedKey:
                    "modules.role.components.RolesTable.body.columns.permissionsHints.deleteProjectDenied",
            },
            {
                allowed: row.permissions.allowViewProject,
                icon: DONEO_ICON_ACTION_SHOW,
                allowedKey:
                    "modules.role.components.RolesTable.body.columns.permissionsHints.viewProjectAllowed",
                deniedKey:
                    "modules.role.components.RolesTable.body.columns.permissionsHints.viewProjectDenied",
            },
            {
                allowed: row.permissions.allowAddTask,
                icon: DONEO_ICON_ACTION_ADD,
                allowedKey:
                    "modules.role.components.RolesTable.body.columns.permissionsHints.addTaskAllowed",
                deniedKey:
                    "modules.role.components.RolesTable.body.columns.permissionsHints.addTaskDenied",
            },
        ]);
    };

    // get task permissions array from role row
    const renderTaskPermissionColumn = (row: Role) => {
        return renderPermissionIcons([
            {
                allowed: row.permissions.allowUpdateTask,
                icon: DONEO_ICON_ACTION_EDIT,
                allowedKey: "modules.role.components.RolesTable.body.columns.permissionsHints.updateTaskAllowed",
                deniedKey: "modules.role.components.RolesTable.body.columns.permissionsHints.updateTaskDenied",
            },
            {
                allowed: row.permissions.allowDeleteTask,
                icon: DONEO_ICON_ACTION_DELETE,
                allowedKey: "modules.role.components.RolesTable.body.columns.permissionsHints.deleteTaskAllowed",
                deniedKey: "modules.role.components.RolesTable.body.columns.permissionsHints.deleteTaskDenied",
            },
            {
                allowed: row.permissions.allowViewTask,
                icon: DONEO_ICON_ACTION_SHOW,
                allowedKey: "modules.role.components.RolesTable.body.columns.permissionsHints.viewTaskAllowed",
                deniedKey: "modules.role.components.RolesTable.body.columns.permissionsHints.viewTaskDenied",
            },
        ]);
    };

    const columnDefinitions = reactive<TableHeaderColumn<Role>[]>([
        {
            label: t("modules.role.components.RolesTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByName.value,
            render: (row: Role) => renderLabel(row.name),
        },
        {
            label: t("modules.role.components.RolesTable.header.columns.projectPermissions"),
            field: "projectPermissions",
            visible: true,
            sortable: false,
            align: "center",
            isFiltered: () => isFilteredByProjectPermission.value,
            render: (row: Role) => renderProjectPermissionColumn(row),
        },
        {
            label: t("modules.role.components.RolesTable.header.columns.taskPermissions"),
            field: "taskPermissions",
            visible: true,
            sortable: false,
            align: "center",
            isFiltered: () => isFilteredByTaskPermission.value,
            render: (row: Role) => renderTaskPermissionColumn(row),
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
                align: definition?.align,
                isFiltered: definition?.isFiltered ?? (() => false),
                render: definition?.render ?? (() => "")
            };
        })
    );

    const onDelete = async (role: Role) => {
        if (role.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await roleService.delete(role.id);
                notify('success', t("modules.role.components.RolesTable.notifications.roleDeleted", { name: role.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                tmpRole.value = role;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "RolesTable.onDelete" } });
                                break;
                            case 403:
                                state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.role.components.RolesTable.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.role.components.RolesTable.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.role.components.RolesTable.errors.deleteError");
                        console.error("Unhandled API error", { file: "RolesTable.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("role id not set", { file: "RolesTable.vue", method: "onDelete" });
        }
    };

    const onConfirmDelete = (role: Role) => {
        dialog.warning({
            title: t("modules.role.components.RolesTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash, { size: 24 }),
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
                onDelete(role);
            },
        });
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: SearchRequest = {
                pager: {
                    enabled: false,
                    currentPage: 1,
                    resultsPage: 0,
                },
                order: {
                    field: order.field,
                    direction: order.direction,
                },
                filter: {
                    //name: filters.name.length > 0 ? filters.name : undefined,
                }
            };
            const response = await roleService.search(payload);
            items.value = response.roles.map((role: RoleResponse) => new Role(role));
        } catch (error: unknown) {
            items.value = [];
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "RolesTable.onRefresh" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.role.components.RolesTable.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.role.components.RolesTable.errors.refreshError");
                    console.error("Unhandled API error", { file: "RolesTable.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

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
        notify('success', t("modules.role.components.RolesTable.notifications.roleAdded", { name: role.name }));
        onRefresh();
    };

    const onRoleUpdated = (role: RoleResponse) => {
        cacheStore.clearUsersCache();
        showFormModal.value = false;
        tmpRole.value = new Role();
        notify('success', t("modules.role.components.RolesTable.notifications.roleUpdated", { name: role.name }));
        onRefresh();
    };

    const hideFormModal = () => {
        showFormModal.value = false;
        tmpRole.value = new Role();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("RolesTable.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("RolesTable.onDelete")) {
                onDelete(tmpRole.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });

</script>

<template>
    <n-modal v-model:show="showFormModal" v-if="showFormModal">
        <RoleForm class="role-form" :role-id="tmpRole.id" @add="onRoleAdded" @update="onRoleUpdated"
            @cancel="hideFormModal" />
    </n-modal>
    <ManageTable :id="props.id" size="small" :disabled="state.ajaxRunning" :rows="localFilteredItems"
        :row-key="row => row.id" :columns="columns" :order="order" @sort="onSort" @refresh="onRefresh" @add="onAdd"
        @clear-filters="onClearFilters">
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
    </ManageTable>
</template>

<style lang="css" scoped>
    .role-form {
        width: 95%;
        max-width: 640px;
    }
</style>