<script setup lang="ts">
    import { ref, reactive, shallowRef, watch, onMounted, onBeforeUnmount } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard } from 'naive-ui';

    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import type { FormMode } from '../../../shared/types/form-mode';
    import type { SearchRequest, UserResponse } from '../types/dto';
    import type { UsersTableFilters } from '../types/users-table-filters.ts';
    import { UserPermissionFilterValue } from '../types/user-admin-permission-filter';

    import { Sort } from '../../../shared/types/models/sort';
    import { User } from '../models/user';

    import { userService } from '../services/user';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import Pager from '../../../shared/components/tables/Pager.vue';
    import UserForm from '../components/UserForm.vue';
    import UsersTable from '../components/UsersTable.vue';

    const { t } = useI18n();
    const { notify } = useNotify();
    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<User[]>([]);

    const sort = reactive<Sort>(new Sort("name", "ASC"));

    const resetPager = ref<boolean>(false);
    const currentPage = ref(1);
    const pageSize = ref(10);
    const totalResults = ref(0);
    const totalPages = ref(0);

    const filters = reactive<UsersTableFilters>({
        permissions: UserPermissionFilterValue.Any,
        name: "",
        email: "",
        createdAt: {
            from: null,
            to: null,
        },
        updatedAt: {
            from: null,
            to: null,
        },
        deletedAt: {
            from: null,
            to: null,
        }
    });

    const showModal = ref<boolean>(false);
    const modalFormMode = ref<FormMode>("add");

    const selectedItem = ref<User>(new User());

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    watch(() => filters, () => {
        resetPager.value = true;
    }, { deep: true });

    watch(pageSize, () => {
        if (currentPage.value != 1) {
            resetPager.value = true;
        } else {
            onRefresh();
        }
    });

    watch(currentPage, () => {
        onRefresh();
    });

    const onSort = (newSort: Sort) => {
        sort.field = newSort.field;
        sort.order = newSort.order;
        onRefresh();
    };

    const onShowAddForm = () => {
        modalFormMode.value = "add";
        showModal.value = true;
    };

    const onShowUpdateForm = (user: User, _index: number) => {
        selectedItem.value = user;
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
                    currentPage: resetPager.value ? 1 : currentPage.value,
                    resultsPage: pageSize.value,
                },
                order: {
                    field: sort.field,
                    sort: sort.order,
                },
                filter: {
                    name: filters.name,
                    email: filters.email,
                    permissions: {
                        isSuperUser: filters.permissions == UserPermissionFilterValue.Any ? undefined : (filters.permissions === UserPermissionFilterValue.OnlyAdministrators ? true : false),
                    },
                    createdAt: filters.createdAt,
                    updatedAt: filters.updatedAt,
                    deletedAt: filters.deletedAt,
                }
            };
            const response = await userService.search(payload);
            totalPages.value = response.pager.totalPages;
            totalResults.value = response.pager.totalResults;
            currentPage.value = response.pager.currentPage;
            items.value = response.users.map((user: UserResponse) => new User(user));
            resetPager.value = false;
        } catch (error: unknown) {
            items.value = [];
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageUsersPage.onRefresh" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.user.components.ManageUsersPage.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.user.components.ManageUsersPage.errors.refreshError");
                    console.error("Unhandled API error", { file: "ManageUsersPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onDelete = async (user: User, _index?: number) => {
        if (user.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await userService.delete(user.id);
                notify('success', t("modules.user.components.ManageUsersPage.notifications.userDeleted", { name: user.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                selectedItem.value = user;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageUsersPage.onDelete" } });
                                break;
                            case 403:
                                state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.user.components.ManageUsersPage.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.user.components.ManageUsersPage.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.user.components.ManageUsersPage.errors.deleteError");
                        console.error("Unhandled API error", { file: "ManageUsersPage.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("user id not set", { file: "ManageUsersPage.vue", method: "onDelete" });
        }
    };

    const onUnDelete = async (user: User, _index?: number) => {
        if (user.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await userService.unDelete(user.id);
                notify('success', t("modules.user.components.ManageUsersPage.notifications.userRestored", { name: user.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                selectedItem.value = user;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageUsersPage.onUnDelete" } });
                                break;
                            case 403:
                                state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.user.components.ManageUsersPage.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.user.components.ManageUsersPage.errors.restoreError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.user.components.ManageUsersPage.errors.restoreError");
                        console.error("Unhandled API error", { file: "ManageUsersPage.vue", method: "onUnDelete" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("user id not set", { file: "ManageUsersPage.vue", method: "onUnDelete" });
        }
    };

    const onAdded = (user: User) => {
        showModal.value = false;
        notify('success', t("modules.user.components.ManageUsersPage.notifications.userAdded", { name: user.name }));
        onRefresh();
    };

    const onUpdated = (user: User) => {
        showModal.value = false;
        notify('success', t("modules.user.components.ManageUsersPage.notifications.userUpdated", { name: user.name }));
        onRefresh();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageUsersPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageUsersPage.onDelete")) {
                onDelete(selectedItem.value)
            } else if (payload.to.includes("ManageUsersPage.onUnDelete")) {
                onUnDelete(selectedItem.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-modal v-model:show="showModal">
        <UserForm :mode="modalFormMode == 'add' ? 'add' : 'update'" :user-id="selectedItem.id" class="modal-form"
            @add="onAdded" @update="onUpdated" @cancel="onCancelForm" />
    </n-modal>

    <n-card :title="t('modules.user.components.ManageUsersPage.header.title')">
        <Pager v-model:current-page="currentPage" v-model:page-size="pageSize" :total-pages="totalPages"
            :total-results="totalResults" class="doneo-pager-container">
            <template #total-results-label="{ totalResults }">
                {{ t("modules.user.components.ManageUsersPage.pager.totalItemsLabel", { total: totalResults }) }}
            </template>
        </Pager>
        <UsersTable :items="items" :disabled="state.ajaxRunning" @refresh="onRefresh" @add="onShowAddForm"
            @update="onShowUpdateForm" @delete="onDelete" @undelete="onUnDelete" :sort="sort" @sort="onSort"
            v-model:filters="filters" />
    </n-card>
</template>

<style lang="css" scoped>
    .modal-form {
        width: 40%;
    }
</style>