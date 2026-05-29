<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, onMounted, onBeforeUnmount, watch, type CSSProperties } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard, NModal } from "naive-ui";

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';

    import { appBus } from '../../../shared/composables/bus';


    import { projectPermissionService } from "../../project-permissions/services/project-permission.ts";
    import { handleAPIError } from '../../../api/client/errorHandler';

    import type { SearchResponse } from "../../project-permissions/types/dto.ts";
    import { ProjectPermission } from "../../project-permissions/models/project-permission.ts";

    import ProjectPermissionsTable from "../../project-permissions/components/ProjectPermissionsTable.vue";
    import ProjectPermissionForm from "../../project-permissions/components/ProjectPermissionForm.vue";

    interface ProjectPermissionsProps {
        style?: string | CSSProperties;
        projectId: string | null;
    }

    const props = defineProps<ProjectPermissionsProps>();

    const emit = defineEmits(['delete']);

    const { t } = useI18n();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectPermission[]>([]);

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    const filterByUser = ref<string>("");
    const filterByRole = ref<string>("");

    const userFilter = computed(() =>
        filterByUser.value?.toLowerCase() ?? ''
    );

    const roleFilter = computed(() =>
        filterByRole.value?.toLowerCase() ?? ''
    );

    const filteredPermissions = computed(() => {
        return items.value.filter((permission) => {
            const userName = permission.user?.name?.toLowerCase();
            const roleName = permission.role?.name?.toLowerCase();

            return (
                (!userFilter.value || userName?.includes(userFilter.value)) &&
                (!roleFilter.value || roleName?.includes(roleFilter.value))
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

    const onAdd = (_projectPermission: ProjectPermission) => {
        showForm.value = false;
        // TODO:
        //notify('success', t("modules.projectStatus.components.ManageProjectStatusesPage.notifications.projectStatusAdded", { name: projectPermission.name }));
        onRefresh();
    };

    const onCancel = () => {
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
                                state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
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
                notify('success', t("modules.projectPermission.components.projectPermissions.notifications.projectPermissionDeleted", { user: projectPermission.user.name, role: projectPermission.role.name }));
                onRefresh();
                // TODO: not refreshing permission count on parent component
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
                                state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.deleteError");
                        console.error("Unhandled API error", { file: "ProjectPermissions.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
            }
        } else {
            console.error("(project permission id || project id) not set", { file: "ProjectPermissions.vue", method: "onDelete" });
        }
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
        <ProjectPermissionForm :project-id="props.projectId" mode="add" style="width: 40%;" @add="onAdd"
            @cancel="onCancel" />
    </n-modal>
    <n-card bordered :style="props.style">
        <ProjectPermissionsTable :project-permissions="filteredPermissions" :loading="state.ajaxRunning"
            @refresh="onRefresh" @add="onShowAddForm" @delete="onDelete" v-model:user-filter="filterByUser"
            v-model:role-filter="filterByRole" />
    </n-card>
</template>

<style lang="css" scoped></style>