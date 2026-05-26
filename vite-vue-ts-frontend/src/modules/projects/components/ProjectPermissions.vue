<script setup lang="ts">
    import { ref, shallowRef, reactive, onMounted, onBeforeUnmount, watch, type CSSProperties } from "vue";
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
    import { Role } from "../../roles/models/role.ts";
    import { UserBase } from "../../users/models/user.ts";

    interface ProjectPermissionsProps {
        style?: string | CSSProperties;
        projectId: string;
    }

    const props = defineProps<ProjectPermissionsProps>();

    const emit = defineEmits(['delete']);

    const { t } = useI18n();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectPermission[]>([]);

    const showForm = ref<boolean>(false);

    const selectedItem = ref<ProjectPermission>(new ProjectPermission({
        id: "",
        user: new UserBase({ id: "", name: "" }),
        role: new Role({
            id: "", name: "",
            permissions: {
                allowUpdateProject: false,
                allowDeleteProject: false,
                allowViewProject: false,
                allowAddTask: false,
                allowUpdateTask: false,
                allowDeleteTask: false,
                allowViewTask: false,
            },
        }),
    }));


    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    const onShowAddForm = () => {
        showForm.value = true;
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const results: SearchResponse = await projectPermissionService.search(props.projectId);
            items.value = results.projectPermissions.map((permission) => new ProjectPermission(permission));
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
    };

    const onDelete = async (projectPermission: ProjectPermission, _index?: number) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await projectPermissionService.delete(props.projectId, projectPermission.id);
            notify('success', t("modules.projectPermission.components.projectPermissions.notifications.projectPermissionDeleted", { user: projectPermission.user.name, role: projectPermission.role.name }));
            onRefresh();
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
    };


    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
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
        TODO: add permission dialog
    </n-modal>
    <n-card bordered :style="props.style">
        <ProjectPermissionsTable :project-permissions="items" :loading="state.ajaxRunning" @refresh="onRefresh"
            @add="onShowAddForm" @delete="onDelete" />
    </n-card>
</template>

<style lang="css" scoped></style>