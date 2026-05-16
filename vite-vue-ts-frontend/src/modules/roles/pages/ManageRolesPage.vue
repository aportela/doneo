<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard } from 'naive-ui';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';
    import { roleService as roleService } from '../services/role';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { RoleResponse, SearchRequest } from '../types/dto';
    import { Role } from '../models/role';
    import RolesTable from '../components/RolesTable.vue';
    import RoleForm from '../components/RoleForm.vue';
    import Pager from '../../../shared/components/tables/Pager.vue';
    import { Sort } from '../../../shared/types/models/sort';

    const { notify } = useNotify();

    const { t } = useI18n();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const roles = shallowRef<Role[]>([]);

    const sort = ref<Sort>(new Sort("name", "ASC"));

    const nameFilter = ref<string>("");
    const currentPage = ref(1);
    const pageSize = ref(10);
    const totalResults = ref(0);
    const totalPages = ref(0);

    const showUserDialog = ref<boolean>(false);
    const roleDialogMode = ref<string>("add");

    const selectedRoleId = ref<string>("");

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    watch(pageSize, () => {
        if (currentPage.value != 1) {
            currentPage.value = 1;
        } else {
            onRefresh();
        }
    });

    watch(currentPage, () => {
        onRefresh();
    });

    const onToggleSort = (field: string) => {
        sort.value.toggleSort(field);
        onRefresh();
    };

    const onShowAddForm = () => {
        roleDialogMode.value = "add";
        showUserDialog.value = true;
    };

    const onShowUpdateForm = (role: Role, _index: number) => {
        selectedRoleId.value = role.id;
        roleDialogMode.value = "update";
        showUserDialog.value = true;
    };

    const onAdd = (role: Role) => {
        showUserDialog.value = false;
        notify('success', t("roleAddedNotification", { name: role.name }));
        onRefresh();
    };

    const onUpdate = (role: Role) => {
        showUserDialog.value = false;
        notify('success', t("roleUpdatedNotification", { name: role.name }));
        onRefresh();
    };

    const onCancel = () => {
        showUserDialog.value = false;
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: SearchRequest = {
                pager: {
                    currentPage: currentPage.value,
                    resultsPage: pageSize.value,
                },
                order: {
                    field: sort.value.field,
                    sort: sort.value.order,
                },
                filter: {
                    name: nameFilter.value,
                    permissions: {
                        allowCreate: true,
                        allowUpdate: undefined,
                        allowDelete: undefined,
                        allowView: undefined,
                        allowList: undefined,
                        allowExecute: undefined,
                    }
                }
            };
            const response = await roleService.search(payload);
            totalPages.value = response.pager.totalPages;
            totalResults.value = response.pager.totalResults;
            roles.value = response.roles.map((role: RoleResponse) => new Role(role));
        } catch (error: unknown) {
            roles.value.length = 0;
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageRolesPage.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while refreshing the role list");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while refreshing the role list");
                    console.error("Unhandled API error", { file: "ManageRolesPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
        }
    };

    const onDelete = async (role: Role, _index: number) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await roleService.delete(role.id);
            notify('success', t("roleDeletedNotification", { name: role.name }));
            onRefresh();
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageRolesPage.onDelete" } });
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("We couldn’t find the specified role");
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while deleting the role");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while deleting the role");
                    console.error("Unhandled API error", { file: "ManageRolesPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageRolesPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageRolesPage.onDelete")) {
                // TODO: missing role/index param at this point
                //onDelete();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-modal v-model:show="showUserDialog">
        <RoleForm :mode="roleDialogMode == 'add' ? 'add' : 'update'" :roleId="selectedRoleId" style="width: 40%;"
            @add="onAdd" @update="onUpdate" @delete="onDelete" @cancel="onCancel" />
    </n-modal>

    <n-card :title="t('Manage roles')">
        <Pager v-model:current-page="currentPage" v-model:page-size="pageSize" :total-pages="totalPages"
            :total-results="totalResults" class="doneo-pager-container">
            <template #total-results-label="{ totalResults }">
                {{ t("TotalRolesPagerLabel", { total: totalResults }) }}
            </template>
        </Pager>
        <RolesTable :roles="roles" :loading="state.ajaxRunning" @refresh="onRefresh" @add="onShowAddForm"
            @update="onShowUpdateForm" @delete="onDelete" @textfilter-keydown-enter="onRefresh" :sort-field="sort.field"
            :sort-order="sort.order" @toggle-sort="onToggleSort" v-model:role-name-filter="nameFilter" />
    </n-card>
</template>

<style lang="css">
    .doneo-pager-container {
        margin-bottom: 4px;
    }
</style>