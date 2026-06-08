<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, onMounted, onBeforeUnmount, watch, type CSSProperties } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard, NModal } from "naive-ui";

    import { useLoadingStore } from '../../../../stores/loading';
    import { useNotify } from '../../../../shared/composables/notification';
    import { appBus } from '../../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../../shared/types/ajaxState';
    import type { SearchResponse } from "../../../project-permissions/types/dto.ts";

    import { ProjectPermission } from "../../../project-permissions/models/project-permission.ts";

    import { projectPermissionService } from "../../../project-permissions/services/project-permission.ts";
    import { handleAPIError } from '../../../../api/client/errorHandler';

    import ProjectPermissionForm from "../../../project-permissions/components/ProjectPermissionForm.vue";
    import ProjectPermissionsTable from "../../../project-permissions/components/ProjectPermissionsTable.vue";
    import type { ProjectPermissionsTableFilters } from "../../../project-permissions/types/project-permissions-table-filter.ts";

    interface ProjectPermissionsProps {
        style?: string | CSSProperties;
        projectId: string;
    }

    const { t } = useI18n();
    const { notify } = useNotify();
    const loadingStore = useLoadingStore();

    const props = defineProps<ProjectPermissionsProps>();

    const emit = defineEmits(['delete']);

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectPermission[]>([]);

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    const filters = reactive<ProjectPermissionsTableFilters>({
        userId: null,
        roleId: null,
        projectPermission: null,
        taskPermission: null,
    });

    const filteredPermissions = computed(() => {
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

    const showForm = ref<boolean>(false);

    const selectedItem = ref<ProjectPermission>(new ProjectPermission());

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    watch(() => props.projectId, (newValue, oldValue) => {
        if (!oldValue && newValue) {
            onRefresh();
        }
    });

    const onShowAddForm = () => {
        showForm.value = true;
    };

    const onCancelForm = () => {
        showForm.value = false;
    };

    const onRefresh = async () => {
        if (props.projectId) {
            Object.assign(state, defaultAjaxStateRunning);
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
            }
        } else {
            console.error("project id not set", { file: "ProjectPermissions.vue", method: "onRefresh" });
        }
    };

    const onDelete = async (projectPermission: ProjectPermission, _index?: number) => {
        if (props.projectId && projectPermission.id) {
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
                                selectedItem.value = projectPermission;
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
            }
        } else {
            console.error("(project permission id || project id) not set", { file: "ProjectPermissions.vue", method: "onDelete" });
        }
    };

    const onAdded = (projectPermission: ProjectPermission) => {
        showForm.value = false;
        items.value = [projectPermission, ...items.value]
        itemCount.value = items.value?.length ?? 0;
        notify('success', t("modules.projectPermission.components.projectPermissionsTab.notifications.projectPermissionAdded", { user: projectPermission.user.name, role: projectPermission.role.name }));
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        if (props.projectId) {
            onRefresh();
        }
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectPermissions.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ProjectPermissions.onDelete")) {
                onDelete(selectedItem.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });

</script>

<template>
    <n-modal v-model:show="showForm">
        <ProjectPermissionForm :project-id="props.projectId" mode="add" style="width: 40%;" @add="onAdded"
            @cancel="onCancelForm" />
    </n-modal>
    <n-card bordered :style="props.style">
        <ProjectPermissionsTable :project-id="props.projectId" :items="filteredPermissions"
            :disabled="state.ajaxRunning" v-model:filters="filters" @refresh="onRefresh" @add="onShowAddForm"
            @delete="onDelete" />
    </n-card>
</template>

<style lang="css" scoped></style>