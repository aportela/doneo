<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard } from 'naive-ui';

    import { api } from '../../../shared/composables/api';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import type { AxiosAPIError } from '../../../shared/composables/axios';
    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { useAppBus, type AppBusEvent } from '../../../shared/composables/bus';
    import { userService } from '../services/user';
    import UserForm from '../components/UserForm.vue';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { UserResponse } from '../types/dto';
    import type { SortOrder } from '../../../shared/types/common';
    import { User } from '../models/user';

    import UsersTable from '../components/UsersTable.vue';
    import Pager from '../../../shared/components/tables/Pager.vue';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';

    const appBus = useAppBus();

    const { notify } = useNotify();

    const { t } = useI18n();


    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const users = shallowRef<User[]>([]);


    const sortField = ref<string>("name");
    const sortOrder = ref<SortOrder>("ASC");

    const filterByUsername = ref<string | null>(null);
    const filterByEmail = ref<string | null>(null);
    const userFilterType = ref<number | null>(0);

    const currentPage = ref(1);
    const pageSize = ref(10);
    const totalResults = ref(0);
    const totalPages = ref(0);

    const showUserDialog = ref<boolean>(false);
    const userDialogMode = ref<string>("add");

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    watch([filterByUsername, filterByEmail, userFilterType], () => {
        currentPage.value = 1;
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
        if (field !== sortField.value) {
            sortField.value = field;
            sortOrder.value = "ASC";
        } else {
            sortOrder.value = sortOrder.value === "ASC" ? "DESC" : "ASC";
        }
        onRefresh();
    }

    const selectedUserId = ref<string | undefined>(undefined);

    const onShowAddForm = () => {
        userDialogMode.value = "add";
        showUserDialog.value = true;
    };

    const onShowUpdateForm = (user: User, _index: number) => {
        selectedUserId.value = user.id;
        userDialogMode.value = "update";
        showUserDialog.value = true;
    };

    const onAdd = (user: User) => {
        showUserDialog.value = false;
        notify('success', t("userAddedNotification", { name: user.name }));
        onRefresh();
    };

    const onUpdate = (user: User) => {
        showUserDialog.value = false;
        notify('success', t("userUpdatedNotification", { name: user.name }));
        onRefresh();
    };

    const onCancel = () => {
        showUserDialog.value = false;
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload = {
                pager: {
                    currentPage: currentPage.value,
                    resultsPage: pageSize.value,
                },
                order: {
                    field: sortField.value,
                    sort: sortOrder.value,
                }
            };
            const response = await userService.search(payload);
            totalPages.value = response.pager.totalPages;
            totalResults.value = response.pager.totalResults;
            users.value = response.users.map((user: UserResponse) => new User(user));
        } catch (error: unknown) {
            users.value.length = 0;
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emitReauthRequired("ManageUsersPage.onRefresh");
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while refreshing the user list");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while refreshing the user list");
                    console.error("Unhandled API error", { file: "ManageUsersPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
        }
    };

    const onDelete = (user: User, _index: number) => {
        Object.assign(state, defaultAjaxStateRunning);
        api.user.delete(user.id).then((_response: any) => {
            notify('success', t("userDeletedNotification", { name: user.name }));
            onRefresh();
        }).catch((errorResponse: AxiosAPIError) => {
            state.ajaxErrors = true;
            if (errorResponse.isAPIError) {
                state.ajaxAPIErrorDetails = errorResponse.customAPIErrorDetails;
                switch (errorResponse.status) {
                    case 401:
                        state.ajaxErrors = false;
                        appBus.emitReauthRequired("ManageUsersPage.onDelete");
                        break;
                    case 404:
                        state.ajaxErrorMessage = t("We couldn’t find the specified user");
                        break;
                    default:
                        state.ajaxErrorMessage = t("There was a problem while deleting the user");
                        break;
                }
            } else {
                state.ajaxErrorMessage = t("There was a problem while deleting the user");
                console.error(errorResponse);
            }
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };

    const onUnDelete = (user: User, _index: number) => {
        Object.assign(state, defaultAjaxStateRunning);
        api.user.unDelete(user.id).then((_response: any) => {
            notify('success', t("userRestoredNotification", { name: user.name }));
            onRefresh();
        }).catch((errorResponse: AxiosAPIError) => {
            state.ajaxErrors = true;
            if (errorResponse.isAPIError) {
                state.ajaxAPIErrorDetails = errorResponse.customAPIErrorDetails;
                switch (errorResponse.status) {
                    case 401:
                        state.ajaxErrors = false;
                        appBus.emitReauthRequired("ManageUsersPage.onUnDelete");
                        break;
                    case 404:
                        state.ajaxErrorMessage = t("We couldn’t find the specified user");
                        break;
                    default:
                        state.ajaxErrorMessage = t("There was a problem while restoring the user");
                        break;
                }
            } else {
                state.ajaxErrorMessage = t("There was a problem while restoring the user");
                console.error(errorResponse);
            }
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };


    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("UserTypeTableHeader"),
            field: "isSuperUser",
        },
        {
            label: t("UsernameTableHeader"),
            field: "name",
        },
        {
            label: t("EmailTableHeader"),
            field: "email",
        },
        {
            label: t("CreatedAtTableHeader"),
            field: "createdAt"
        },
        {
            label: t("UpdatedAtTableHeader"),
            field: "updatedAt"
        },
        {
            label: t("DeletedAtTableHeader"),
            field: "deletedAt"
        },
    ]);

    onMounted(() => {
        onRefresh();
        appBus.on((event: AppBusEvent) => {
            if (event.type == "reauthValidNotify" && event.to.includes("ManageUsersPage.onRefresh")) {
                onRefresh();
            }
        });
    });

    onBeforeUnmount(() => {
        appBus.reset();
    });


</script>

<template>
    <n-modal v-model:show="showUserDialog">
        <UserForm :mode="userDialogMode == 'add' ? 'add' : 'update'" :user-id="selectedUserId" style="width: 40%;"
            @add="onAdd" @update="onUpdate" @delete="onDelete" @cancel="onCancel" @undelete="onUnDelete" />
    </n-modal>

    <n-card :title="t('Manage users')">
        <Pager v-model:current-page="currentPage" v-model:page-size="pageSize" :total-pages="totalPages"
            :total-results="totalResults">
            <template #total-results-label="{ totalResults }">
                {{ t("TotalUsersPagerLabel", { total: totalResults }) }}
            </template>
        </Pager>
        <UsersTable :users="users" :columns="columns" :loading="state.ajaxRunning" @refresh="onRefresh"
            @add="onShowAddForm" @update="onShowUpdateForm" @delete="onDelete" @undelete="onUnDelete"
            :sort-field="sortField" :sort-order="sortOrder" @toggle-sort="onToggleSort" />
    </n-card>
</template>

<style lang="css"></style>