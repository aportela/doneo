<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard, NPagination, NFlex } from 'naive-ui';
    //import { IconTrash, IconTrashOff, } from '@tabler/icons-vue';

    import { api } from '../../../shared/composables/api';

    import { type AjaxStateInterface, defaultAjaxState } from '../../../shared/types/ajaxState';
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
    const totalResuls = ref(0);
    const totalPages = ref(0);

    const showUserDialog = ref<boolean>(false);
    const userDialogMode = ref<string>("add");

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

    const onRefresh = async () => {
        state.ajaxRunning = true;
        loadingStore.set(true);
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
            totalResuls.value = response.pager.totalResults;
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
            loadingStore.set(false);
        }
    };

    const selectedUserId = ref<string | undefined>(undefined);

    const onAddUser = () => {
        userDialogMode.value = "add";
        showUserDialog.value = true;
    };

    const onUpdateUser = (user: User, _index: number) => {
        selectedUserId.value = user.id;
        userDialogMode.value = "update";
        showUserDialog.value = true;
    };


    const onAdd = () => {
        showUserDialog.value = false;
        notify('success', t("User added"))
        onRefresh();
    };

    const onUpdate = () => {
        showUserDialog.value = false;
        notify('success', t("User updated"))
        onRefresh();
    };

    const onCancel = () => {
        showUserDialog.value = false;
    };

    const onDelete = (user: User, _index: number) => {
        Object.assign(state, defaultAjaxState);
        state.ajaxRunning = true;
        api.user.delete(user.id).then((_response: any) => {
            notify('success', t(`User ${user.name} has been deleted`))
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
                        state.ajaxErrorMessage = t("There was a problem deleting the user data");
                        break;
                }
            } else {
                state.ajaxErrorMessage = t("There was a problem deleting the user data");
                console.error(errorResponse);
            }
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };

    const onUnDelete = (user: User, _index: number) => {
        Object.assign(state, defaultAjaxState);
        state.ajaxRunning = true;
        api.user.unDelete(user.id).then((_response: any) => {
            notify('success', t(`User ${user.name} has been restored`))
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
                        state.ajaxErrorMessage = t("There was a problem undeleting the user data");
                        break;
                }
            } else {
                state.ajaxErrorMessage = t("There was a problem undeleting the user data");
                console.error(errorResponse);
            }
        }).finally(() => {
            state.ajaxRunning = false;
        });
    };

    const pageSizes = [
        {
            label: 'All results',
            value: 0
        },
        {
            label: '5 results/page',
            value: 5
        },
        {
            label: '10 results/page',
            value: 10
        },
        {
            label: '20 results/page',
            value: 20
        },
        {
            label: '50 results/page',
            value: 50
        },
        {
            label: '100 results/page',
            value: 100
        },
        {
            label: '200 results/page',
            value: 200
        },
    ];

    const columns = [
        {
            label: "Type",
            field: "isSuperUser",
        },
        {
            label: "Name",
            field: "name",
        },
        {
            label: "Email",
            field: "email",
        },
        {
            label: "Created at",
            field: "createdAt"
        },
        {
            label: "Updated at",
            field: "updatedAt"
        },
        {
            label: "Deleted at",
            field: "deletedAt"
        },
    ];

    onMounted(() => {
        onRefresh();
        appBus.on((event: AppBusEvent) => {
            if (event.type == "reauthValidNotify" && event.to.includes("ManageUsersPage.onRefresh")) {
                onRefresh();
            }
        });
    });

    onBeforeUnmount(() => {
        // TODO: manage this
        //bus.off("reAuthSucess");
    });

</script>

<template>
    <n-modal v-model:show="showUserDialog">
        <UserForm :mode="userDialogMode == 'add' ? 'add' : 'update'" :user-id="selectedUserId" style="width: 40%;"
            @add="onAdd" @update="onUpdate" @delete="onDelete" @cancel="onCancel" @undelete="onUnDelete" />
    </n-modal>

    <n-card :title="t('Manage users')">
        <n-flex justify="space-between" class="pagination-container">
            <div style="padding-top: 2px; padding-left: 2px;">
                <span>Total users: {{ totalResuls }}</span>
            </div>
            <n-pagination v-model:page="currentPage" :page-count="totalPages" v-model:page-size="pageSize"
                show-size-picker :page-sizes="pageSizes" :page-slot="8">
                <template #prefix="{ page, pageCount }">
                    Page {{ page }} of {{ pageCount }}
                </template>
            </n-pagination>
        </n-flex>
        <UsersTable :users="users" :columns="columns" :loading="state.ajaxRunning" @refresh="onRefresh" @add="onAddUser"
            @update="onUpdateUser" @delete="onDelete" @undelete="onUnDelete" :sort-field="sortField"
            :sort-order="sortOrder" @toggle-sort="onToggleSort" />
    </n-card>
</template>

<style lang="css">
    .avatar {
        margin-right: 4px;
    }

    .pagination-container {
        padding: 4px;
        background-color: rgba(250, 250, 252, 1);
        margin-bottom: 8px;
        border: 1px solid;
        border-color: rgb(239, 239, 245);
        border-radius: 3px;
    }

    .table-header-click-action th:not(:last-of-type) {
        cursor: pointer;
    }

    .table-header-click-action th:not(:last-of-type) .n-icon {
        margin-top: 4px;
    }

    @media (max-width: 768px) {
        .hide-mobile {
            display: none;
        }
    }
</style>