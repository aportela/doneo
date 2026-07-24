<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, useDialog, NIcon, NButton, NButtonGroup } from 'naive-ui';

    import { useLoadingStore } from '../../../stores/loading';
    import { useCacheStore } from '../../../stores/cache.ts';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import type { Order } from '../../../shared/types/order.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';

    import { ProjectPermission } from '../models/project-permission.ts';

    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectPermissionService } from '../services/project-permission.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import RoleSelector from '../../roles/components/RoleSelector.vue';
    import ProjectPermissionSelect from '../../../shared/components/selectors/ProjectPermissionSelect.vue';
    import type { ProjectPermissionSelectValue } from '../../../shared/types/project-permission-select-value.ts';
    import type { TaskPermissionSelectValue } from '../../../shared/types/task-permission-select-value.ts';
    import TaskPermissionSelect from '../../../shared/components/selectors/TaskPermissionSelect.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';

    import { renderIcon, renderLabel, renderProjectPermissionIcons, renderTaskPermissionIcons } from '../../../shared/composables/naive-ui-helpers.ts';
    import { DONEO_ICON_ACTION_DELETE, DONEO_ICON_ACTION_EDIT } from '../../../shared/types/icons.ts';
    import type { ProjectPermissionResponse, SearchResponse } from '../types/dto.ts';
    import ProjectPermissionForm from './ProjectPermissionForm.vue';

    interface Props {
        id?: string;
        readOnly?: boolean;
        projectId: string;
    }

    const props = withDefaults(defineProps<Props>(), { id: "ProjectPermissionsTable" });;

    const emit = defineEmits(['add', 'delete']);

    const itemCount = defineModel<number>("itemCount", { default: 0 });

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

    const items = shallowRef<ProjectPermission[]>([]);

    const tmpItem = ref<ProjectPermission>(new ProjectPermission());

    const showNoItemsWarningMessage = ref<boolean>(false);

    const currentOrder = reactive<Order>({ field: "name", direction: "ASC" });

    const onSort = (newOrder: Order) => {
        currentOrder.field = newOrder.field;
        currentOrder.direction = newOrder.direction;
        // we have all results, use local sorting for avoiding server load
        if (currentOrder.direction === "ASC") {
            items.value = [...items.value].sort((a, b) =>
                currentOrder.field === "user" ? a.user.name.localeCompare(b.user.name) : String(a.role.name).localeCompare(String(b.role.name))
            );
        } else {
            items.value = [...items.value].sort((a, b) =>
                currentOrder.field === "user" ? b.user.name.localeCompare(a.user.name) : String(b.role.name).localeCompare(String(a.role.name))
            );
        }
    };

    const projectPermissionSelectorRef = ref<InstanceType<typeof ProjectPermissionSelect>[] | null>(null);
    const taskPermissionSelectorRef = ref<InstanceType<typeof ProjectPermissionSelect>[] | null>(null);

    interface ProjectPermissionsTableFilters {
        userId: string | null;
        roleId: string | null;
        projectPermission: ProjectPermissionSelectValue | null;
        taskPermission: TaskPermissionSelectValue | null;
    };

    const filters = reactive<ProjectPermissionsTableFilters>(
        {
            userId: null,
            roleId: null,
            projectPermission: null,
            taskPermission: null,
        }
    );

    const isFilteredByUser = computed<boolean>(() => filters.userId !== null);
    const isFilteredByRole = computed<boolean>(() => filters.roleId !== null);
    const isFilteredByProjectPermission = computed<boolean>(() => filters.projectPermission !== null);
    const isFilteredByTaskPermission = computed<boolean>(() => filters.taskPermission !== null);

    const onClearFilters = () => {
        filters.userId = null;
        filters.roleId = null;
        if (projectPermissionSelectorRef.value) {
            projectPermissionSelectorRef.value[0]?.reset();
        }
        if (taskPermissionSelectorRef.value) {
            taskPermissionSelectorRef.value[0]?.reset();
        }
    };

    const localFilteredItems = computed(() => {
        return items.value.filter((permission: ProjectPermission) => {
            return (
                (filters.userId === null || filters.userId === permission.user.id) &&
                (filters.roleId === null || filters.roleId === permission.role.id) &&
                (filters.projectPermission === null || (filters.projectPermission !== null && (
                    (filters.projectPermission === "updateProjectAllowed" && permission.role.permissions.allowUpdateProject) ||
                    (filters.projectPermission === "updateProjectDenied" && !permission.role.permissions.allowUpdateProject) ||
                    (filters.projectPermission === "deleteProjectAllowed" && permission.role.permissions.allowDeleteProject) ||
                    (filters.projectPermission === "deleteProjectDenied" && !permission.role.permissions.allowDeleteProject) ||
                    (filters.projectPermission === "viewProjectAllowed" && permission.role.permissions.allowViewProject) ||
                    (filters.projectPermission === "viewProjectDenied" && !permission.role.permissions.allowViewProject) ||
                    (filters.projectPermission === "addTaskAllowed" && permission.role.permissions.allowAddTask) ||
                    (filters.projectPermission === "addTaskDenied" && !permission.role.permissions.allowAddTask)
                ))
                ) &&
                (filters.taskPermission === null || (filters.taskPermission !== null && (
                    (filters.taskPermission === "updateTaskAllowed" && permission.role.permissions.allowUpdateTask) ||
                    (filters.taskPermission === "updateTaskDenied" && !permission.role.permissions.allowUpdateTask) ||
                    (filters.taskPermission === "deleteTaskAllowed" && permission.role.permissions.allowDeleteTask) ||
                    (filters.taskPermission === "deleteTaskDenied" && !permission.role.permissions.allowDeleteTask) ||
                    (filters.taskPermission === "viewTaskAllowed" && permission.role.permissions.allowViewTask) ||
                    (filters.taskPermission === "viewTaskDenied" && !permission.role.permissions.allowViewTask)
                ))
                )
            );
        });
    });

    const columnDefinitions = reactive<TableHeaderColumn<ProjectPermission>[]>([
        {
            label: t("modules.projectPermission.components.projectPermissionsTable.header.columns.user"),
            field: "user",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByUser.value,
            render: (row: ProjectPermission) => h(AvatarUserName, { userId: row.user.id, userName: row.user.name }),
        },
        {
            label: t("modules.projectPermission.components.projectPermissionsTable.header.columns.role"),
            field: "role",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByRole.value,
            render: (row: ProjectPermission) => renderLabel(row.role.name),
        },
        {
            label: t("modules.projectPermission.components.projectPermissionsTable.header.columns.projectPermissions"),
            field: "projectPermissions",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByProjectPermission.value,
            align: "center",
            render: (row: ProjectPermission) => renderProjectPermissionIcons(row.role, t),
        },
        {
            label: t("modules.projectPermission.components.projectPermissionsTable.header.columns.taskPermissions"),
            field: "taskPermissions",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByTaskPermission.value,
            align: "center",
            render: (row: ProjectPermission) => renderTaskPermissionIcons(row.role, t),
        },
    ]);

    // create (if not found) default settings for this table (column order & visibility)
    tableSettingsStore.register(props.id, { columns: columnDefinitions.map((column) => { return { field: column.field, visible: column.visible } }) ?? [] });

    // restore previous settings
    const tableSettings = tableSettingsStore.get(props.id);

    // build columns based on saved order visibility settings
    const columns = computed<TableHeaderColumn<ProjectPermission>[]>(() =>
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

    const onDelete = async (projectPermission: ProjectPermission) => {
        if (projectPermission.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await projectPermissionService.delete(props.projectId, projectPermission.id);
                items.value = items.value.filter((item) => item.id != projectPermission.id)
                itemCount.value = items.value?.length ?? 0;
                notify('success', t("modules.projectPermission.components.projectPermissionsTab.notifications.projectPermissionDeleted", { user: projectPermission.user.name, role: projectPermission.role.name }));
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                tmpItem.value = projectPermission;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectPermissions.onDelete" } });
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissionsTab.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissionsTab.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissionsTab.errors.deleteError");
                        console.error("Unhandled API error", { file: "ProjectPermissions.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("(project permission id || project id) not set", { file: "ProjectPermissions.vue", method: "onDelete" });
        }
    };

    const onConfirmDelete = (projectPermission: ProjectPermission) => {
        dialog.warning({
            title: t("modules.projectPermission.components.projectPermissionsTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(DONEO_ICON_ACTION_DELETE, { size: 24 }),
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
                onDelete(projectPermission);
            },
        });
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        showNoItemsWarningMessage.value = false;
        try {
            const results: SearchResponse = await projectPermissionService.search(props.projectId);
            items.value = results.projectPermissions.map((permission) => new ProjectPermission(permission));
            itemCount.value = items.value?.length ?? 0;
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectPermissions.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissionsTab.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissionsTab.errors.refreshError");
                    console.error("Unhandled API error", { file: "ProjectPermissions.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const showFormModal = ref<boolean>(false);

    const onAdd = () => {
        tmpItem.value = new ProjectPermission();
        showFormModal.value = true;
    };

    const onProjectPermissionAdded = (projectPermission: ProjectPermissionResponse) => {
        cacheStore.clearUsersCache();
        showFormModal.value = false;
        tmpItem.value = new ProjectPermission();
        notify('success', t("modules.projectPermission.components.projectPermissionsTab.notifications.projectPermissionAdded", { user: projectPermission.user.name, role: projectPermission.role.name }));
        onRefresh();
    };

    const hideFormModal = () => {
        showFormModal.value = false;
        tmpItem.value = new ProjectPermission();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectPermissionsTable.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ProjectPermissionsTable.onDelete")) {
                onDelete(tmpItem.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-modal v-model:show="showFormModal" v-if="showFormModal">
        <ProjectPermissionForm class="project-permission-form" :project-id="props.projectId"
            @add="onProjectPermissionAdded" @cancel="hideFormModal" />
    </n-modal>
    <ManageTable :id="props.id" size="small" :disabled="state.ajaxRunning" :rows="localFilteredItems"
        :row-key="row => row.id" :columns="columns" :order="currentOrder"
        :show-no-items-warning-message="showNoItemsWarningMessage || (items.length > 0 && localFilteredItems.length === 0)"
        :no-items-warning-message="t('modules.projectPermission.components.projectPermissionsTable.warnings.noItemsFound')"
        @sort="onSort" @refresh="onRefresh" @add="onAdd" @clear-filters="onClearFilters"
        :buttons="props.readOnly ? ['refresh', 'settings'] : ['refresh', 'add', 'settings']">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <UserSelector v-if="column.field === 'user'" hideAvatar clearable :disabled="state.ajaxRunning"
                    v-model:id="filters.userId"
                    :placeholder="t('modules.projectPermission.components.projectPermissionsTable.filters.user.placeholder')" />
                <RoleSelector v-if="column.field === 'role'" clearable :disabled="state.ajaxRunning"
                    v-model:id="filters.roleId"
                    :placeholder="t('modules.projectPermission.components.projectPermissionsTable.filters.role.placeholder')" />
                <ProjectPermissionSelect v-if="column.field === 'projectPermission'"
                    v-model:permission="filters.projectPermission"
                    :placeholder="t('shared.components.selectors.ProjectPermissionSelect.placeholder')" clearable
                    ref="projectPermissionSelectorRef" />
                <TaskPermissionSelect v-if="column.field === 'taskPermission'"
                    v-model:permission="filters.taskPermission"
                    :placeholder="t('shared.components.selectors.TaskPermissionSelect.placeholder')" clearable
                    ref="taskPermissionSelectorRef" />
            </th>
        </template>
        <template #rowactions="{ row }">
            <n-button-group class="doneo-table-actions-button-group" size="small">
                <n-button @click="onConfirmDelete(row)" :disabled="state.ajaxRunning || props.readOnly"
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
    .project-permission-form {
        width: 95%;
        max-width: 640px;
    }
</style>