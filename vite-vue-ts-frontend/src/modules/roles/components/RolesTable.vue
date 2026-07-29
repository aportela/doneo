<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, useDialog, NIcon, NButton, NButtonGroup } from 'naive-ui';

    import { DONEO_ICON_ACTION_DELETE, DONEO_ICON_ACTION_EDIT } from '../../../shared/types/icons.ts';

    import { useLoadingStore } from '../../../stores/loading';
    import { useCacheStore } from '../../../stores/cache.ts';
    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';


    import { Role } from '../models/role';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { PAGER_DEFAULT_RESULTS_PAGE_NO_PAGINATION, type Pagination } from '../../../shared/types/pager.ts';
    import type { Order } from '../../../shared/types/order.ts';
    import { roleService } from '../services/role.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { RoleResponse, SearchRequest } from '../types/dto.ts';

    import RoleForm from './RoleForm.vue';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';
    import ProjectPermissionSelect from '../../../shared/components/selectors/ProjectPermissionSelect.vue';
    import TaskPermissionSelect from '../../../shared/components/selectors/TaskPermissionSelect.vue';
    import type { ProjectPermissionSelectValue } from '../../../shared/types/project-permission-select-value.ts';
    import type { TaskPermissionSelectValue } from '../../../shared/types/task-permission-select-value.ts';
    import { renderIcon, renderLabel, renderProjectPermissionIcons, renderTaskPermissionIcons } from '../../../shared/composables/naive-ui-helpers.ts';

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

    const tmpItem = ref<Role>(new Role());

    const showNoItemsWarningMessage = ref<boolean>(false);

    const currentOrder = reactive<Order>({ field: "name", direction: "ASC" });

    const onSort = (newOrder: Order) => {
        currentOrder.field = newOrder.field;
        currentOrder.direction = newOrder.direction;
        // we have all results, use local sorting for avoiding server load
        if (currentOrder.direction === "ASC") {
            items.value = [...items.value].sort((a, b) =>
                a.name.localeCompare(b.name)
            );
        } else {
            items.value = [...items.value].sort((a, b) =>
                b.name.localeCompare(a.name)
            );
        }
    };

    const currentPagination = reactive<Pagination>({ enabled: false, currentPage: 1, resultsPage: PAGER_DEFAULT_RESULTS_PAGE_NO_PAGINATION, totalPages: 1, totalResults: 0 });

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
        projectPermissionSelectorRef.value = null;
        taskPermissionSelectorRef.value = null;
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
            render: (row: Role) => renderProjectPermissionIcons(row, t),
        },
        {
            label: t("modules.role.components.RolesTable.header.columns.taskPermissions"),
            field: "taskPermissions",
            visible: true,
            sortable: false,
            align: "center",
            isFiltered: () => isFilteredByTaskPermission.value,
            render: (row: Role) => renderTaskPermissionIcons(row, t),
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
                                tmpItem.value = role;
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
            icon: renderIcon(DONEO_ICON_ACTION_DELETE, { size: 24 }),
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
        showNoItemsWarningMessage.value = false;
        try {
            const payload: SearchRequest = {
                pager: currentPagination,
                order: currentOrder,
                filter: {
                    // disabled (local filtering)
                    //name: filters.name !== "" ? filters.name : undefined,
                }
            };
            const response = await roleService.search(payload);
            items.value = response.roles.map((role: RoleResponse) => new Role(role));
            showNoItemsWarningMessage.value = items.value.length === 0;
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
        tmpItem.value = new Role();
        showFormModal.value = true;
    };

    const onUpdate = (role: Role) => {
        tmpItem.value = role;
        showFormModal.value = true;
    };

    const onRoleAdded = (role: RoleResponse) => {
        cacheStore.clearUsersCache();
        showFormModal.value = false;
        tmpItem.value = new Role();
        notify('success', t("modules.role.components.RolesTable.notifications.roleAdded", { name: role.name }));
        onRefresh();
    };

    const onRoleUpdated = (role: RoleResponse) => {
        cacheStore.clearUsersCache();
        showFormModal.value = false;
        tmpItem.value = new Role();
        notify('success', t("modules.role.components.RolesTable.notifications.roleUpdated", { name: role.name }));
        onRefresh();
    };

    const hideFormModal = () => {
        showFormModal.value = false;
        tmpItem.value = new Role();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("RolesTable.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("RolesTable.onDelete")) {
                onDelete(tmpItem.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });

</script>

<template>
    <n-modal v-model:show="showFormModal">
        <RoleForm class="doneo-role-form" :role-id="tmpItem.id" @add="onRoleAdded" @update="onRoleUpdated"
            @cancel="hideFormModal" v-if="showFormModal" />
    </n-modal>
    <ManageTable :id="props.id" size="small" :disabled="state.ajaxRunning" :rows="localFilteredItems"
        :row-key="row => row.id" :columns="columns" :order="currentOrder"
        :show-no-items-warning-message="showNoItemsWarningMessage || (items.length > 0 && localFilteredItems.length === 0)"
        :no-items-warning-message="t('modules.role.components.RolesTable.warnings.noItemsFound')" @sort="onSort"
        @refresh="onRefresh" @add="onAdd" @clear-filters="onClearFilters">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <TextFilterInput v-if="column.field === 'name'" clearable :disabled="state.ajaxRunning"
                    :placeholder="t('modules.role.components.RolesTable.filters.name.placeholder')"
                    v-model:value="filters.name" />
                <ProjectPermissionSelect v-else-if="column.field == 'projectPermissions'" clearable
                    :disabled="state.ajaxRunning"
                    :placeholder="t('shared.components.selectors.ProjectPermissionSelect.placeholder')"
                    v-model:permission="filters.projectPermission" ref="projectPermissionSelectorRef" />
                <TaskPermissionSelect v-else-if="column.field === 'taskPermissions'" clearable
                    :disabled="state.ajaxRunning"
                    :placeholder="t('shared.components.selectors.TaskPermissionSelect.placeholder')"
                    v-model:permission="filters.taskPermission" ref="taskPermissionSelectorRef" />
            </th>
        </template>
        <template #rowactions="{ row }">
            <n-button-group class="doneo-table-actions-button-group" size="small">
                <n-button @click="onUpdate(row)" :disabled="state.ajaxRunning" class="doneo-table-actions-button">
                    {{ t("shared.buttons.Edit.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_EDIT" />
                    </template>
                </n-button>
                <n-button @click="onConfirmDelete(row)" :disabled="state.ajaxRunning"
                    class="doneo-table-actions-button">
                    {{ t("shared.buttons.Delete.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_DELETE" />
                    </template>
                </n-button>
            </n-button-group>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped>
    .doneo-role-form {
        width: 95%;
        max-width: 640px;
    }
</style>