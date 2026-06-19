<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard } from 'naive-ui';

    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import type { FormMode } from '../../../shared/types/form-mode';
    import type { SearchRequest, RoleResponse } from '../types/dto';
    import type { RolesTableFilters } from '../types/roles-table-filters.ts';

    import { Sort } from '../../../shared/types/models/sort';
    import { Role } from '../models/role';

    import { roleService as roleService } from '../services/role';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import RoleForm from '../components/RoleForm.vue';
    import RolesTable from '../components/RolesTable.vue';


    const { t } = useI18n();
    const { notify } = useNotify();
    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<Role[]>([]);

    const sort = reactive<Sort>(new Sort("name", "ASC"));

    const filters = reactive<RolesTableFilters>({
        name: "",
        projectPermission: null,
        taskPermission: null,
    });

    const nameFilterLowerCase = computed(() =>
        filters.name.toLowerCase()
    );

    const filteredItems = computed(() => {
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

    const showModal = ref<boolean>(false);
    const modalFormMode = ref<FormMode>("add");

    const selectedItem = ref<Role>(new Role());

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    const onShowAddForm = () => {
        modalFormMode.value = "add";
        showModal.value = true;
    };

    const onShowUpdateForm = (role: Role, _index?: number) => {
        selectedItem.value = role;
        modalFormMode.value = "update";
        showModal.value = true;
    };

    const onCancelForm = () => {
        showModal.value = false;
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: SearchRequest = {
                pager: {
                    currentPage: 1,
                    resultsPage: 0,
                },
                order: {
                    field: sort.field,
                    sort: sort.order,
                },
                filter: {
                    name: filters.name.length > 0 ? filters.name : undefined,
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
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageRolesPage.onRefresh" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.role.components.ManageRolesPage.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.role.components.ManageRolesPage.errors.refreshError");
                    console.error("Unhandled API error", { file: "ManageRolesPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onDelete = async (role: Role, _index?: number) => {
        if (role.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await roleService.delete(role.id);
                notify('success', t("modules.role.components.ManageRolesPage.notifications.roleDeleted", { name: role.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                selectedItem.value = role;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageRolesPage.onDelete" } });
                                break;
                            case 403:
                                state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.role.components.ManageRolesPage.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.role.components.ManageRolesPage.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.role.components.ManageRolesPage.errors.deleteError");
                        console.error("Unhandled API error", { file: "ManageRolesPage.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("role id not set", { file: "ManageRolesPage.vue", method: "onDelete" });
        }
    };

    const onAdded = (role: Role) => {
        showModal.value = false;
        notify('success', t("modules.role.components.ManageRolesPage.notifications.roleAdded", { name: role.name }));
        onRefresh();
    };

    const onUpdated = (role: Role) => {
        showModal.value = false;
        notify('success', t("modules.role.components.ManageRolesPage.notifications.roleUpdated", { name: role.name }));
        onRefresh();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageRolesPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageRolesPage.onDelete")) {
                onDelete(selectedItem.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-modal v-model:show="showModal">
        <RoleForm :mode="modalFormMode == 'add' ? 'add' : 'update'" :role-id="selectedItem.id" class="modal-form"
            @add="onAdded" @update="onUpdated" @cancel="onCancelForm" />
    </n-modal>

    <n-card :title="t('modules.role.components.ManageRolesPage.header.title')">
        <RolesTable :items="filteredItems" :disabled="state.ajaxRunning" @refresh="onRefresh" @add="onShowAddForm"
            @update="onShowUpdateForm" @delete="onDelete" v-model:filters="filters" />
    </n-card>
</template>

<style lang="css">
    .modal-form {
        width: 40%;
    }
</style>