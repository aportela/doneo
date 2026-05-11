<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, computed, shallowRef, h, type Component, watch } from 'vue';
    import { useI18n } from "vue-i18n";
    import { NAvatar, NInput, NSelect, NIcon, NButton, NModal, NButtonGroup, useDialog, NEmpty, NCard, NPagination, NFlex } from 'naive-ui';
    import { IconUser, IconUserKey, IconPlus, IconEdit, IconSearch, IconTrash, IconTrashOff, IconRefresh } from '@tabler/icons-vue';
    import { api } from '../composables/api';
    import { type AjaxStateInterface, defaultAjaxState } from '../types/ajaxState';
    import type { AxiosAPIError } from '../composables/axios';
    import { type EntityAction } from '../types/common';
    import { useLoadingStore } from '../stores/loading';
    import { useSessionStore } from '../stores/session';
    import { useNotify } from '../composables/notification';
    import { default as UserForm } from '../components/forms/UserForm.vue';
    import { default as ManageTable } from '../components/custom/ManageTable.vue';
    import { default as DateFilter } from '../components/forms/DateFilter.vue';
    import { useAppBus, type AppBusEvent } from '../composables/bus';
    import { userService } from '../api/services/user';
    import { User } from '../api/models/user';
    import type { UserResponse } from '../api/types/dto/user';

    const appBus = useAppBus();

    const { notify } = useNotify();

    const { t } = useI18n();

    const sessionStore = useSessionStore();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    class SearchableUser extends User {
        _searchName: string;
        _searchEmail: string;

        constructor(data: UserResponse) {
            super(data);
            this._searchName = data.name.toLowerCase();
            this._searchEmail = data.email.toLowerCase();
        }
    }

    const users = shallowRef<SearchableUser[]>([]);

    const filterUserOptions = [
        { label: 'All users', value: 0 },
        { label: 'Administrators only', value: 1 },
        { label: 'Users only', value: 2 }
    ];

    const filterByUsername = ref<string | null>(null);
    const filterByEmail = ref<string | null>(null);
    const userFilterType = ref<number | null>(0);

    const currentPage = ref(1);
    const pageSize = ref(10);
    const totalResuls = ref(0);
    const totalPages = ref(0);


    watch([filterByUsername, filterByEmail, userFilterType], () => {
        currentPage.value = 1;
    });

    watch(pageSize, () => {
        currentPage.value = 1;
    });

    watch(currentPage, () => {
        onRefresh();
    });

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
                    field: "name",
                    sort: "ASC",
                }
            };
            const response = await userService.search(payload);
            totalPages.value = response.pager.totalPages;
            totalResuls.value = response.pager.totalResults;
            users.value = response.users.map((user: UserResponse) => new SearchableUser(user));
        }
        finally {
            state.ajaxRunning = false;
            loadingStore.set(false);
        }
        /*
        api.user.search().then((successResponse: SearchUsersResponse) => {
            users.value = successResponse.data.users.map((u: UserInterface) => new UserClass(u));
        }).catch((errorResponse: AxiosAPIError) => {
            state.ajaxErrors = true;
            if (errorResponse.isAPIError) {
                state.ajaxAPIErrorDetails = errorResponse.customAPIErrorDetails;
                switch (errorResponse.status) {
                    case 401:
                        state.ajaxErrors = false;
                        appBus.emitReauthRequired("ManageUsersPage.onRefresh");
                        break;
                    default:
                        state.ajaxErrorMessage = t("There was a problem while refreshing the user list");
                        break;
                }
            } else {
                state.ajaxErrorMessage = t("There was a problem while refreshing the user list");
                console.error(errorResponse);
            }
        }).finally(() => {
            loadingStore.set(false);
            state.ajaxRunning = false;
        });
        */
    };

    const selectedUserId = ref<string | undefined>(undefined);

    const onAddUser = () => {
        actionDialogMode.value = "add";
    };

    const onUpdateUser = (user: User, _index: number) => {
        actionDialogMode.value = "update";
        selectedUserId.value = user.id;
    };

    const actionDialogMode = ref<EntityAction>("none");

    const isVisibleActionDialog = computed<boolean>({
        get: () => actionDialogMode.value !== "none",
        set: (value: boolean) => {
            if (!value) {
                actionDialogMode.value = "none";
            }
        }
    });

    const onAdd = () => {
        isVisibleActionDialog.value = false;
        notify('success', t("User added"))
        onRefresh();
    };

    const onUpdate = () => {
        isVisibleActionDialog.value = false;
        notify('success', t("User updated"))
        onRefresh();
    };

    const onCancel = () => {
        isVisibleActionDialog.value = false;
    };

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


    const onDelete = (userId: string) => {
        Object.assign(state, defaultAjaxState);
        state.ajaxRunning = true;
        api.user.delete(userId).then((_response: any) => {
            notify('success', t("User deleted"))
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

    const onUnDelete = (userId: string) => {
        Object.assign(state, defaultAjaxState);
        state.ajaxRunning = true;
        api.user.unDelete(userId).then((_response: any) => {
            notify('success', t("User undeleted"))
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

    // TODO: this is already declared on menu.ts

    const renderIcon = (icon: Component) => {
        return (size = 32) =>
            () =>
                h(
                    NIcon,
                    { size },
                    {
                        default: () => h(icon),
                    },
                );
    };


    const dialog = useDialog();

    function confirmDelete(user: User, _index: number) {
        dialog.warning({
            title: t("Delete user"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    `You are about to delete the user "${user.name}" from the system.`,
                    h('br'),
                    h('br'),
                    'Do you want to continue?',
                ]),
            positiveText: t("Delete"),
            negativeText: t("Cancel"),
            onPositiveClick: () => {
                onDelete(user.id);
            },
        })
    };

    function confirmUnDelete(user: User, _index: number) {
        dialog.warning({
            title: t("Restore user"),
            icon: renderIcon(IconTrashOff)(24),
            content: () =>
                h('div', [
                    `You are about to restore the user "${user.name}" from the system.`,
                    h('br'),
                    h('br'),
                    'Do you want to continue?',
                ]),
            positiveText: t("Restore"),
            negativeText: t("Cancel"),
            onPositiveClick: () => {
                onUnDelete(user.id);
            },
        })
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
</script>

<template>
    <n-modal v-model:show="isVisibleActionDialog">
        <UserForm :mode="actionDialogMode" :user-id="selectedUserId" style="width: 40%;" @add="onAdd" @update="onUpdate"
            @delete="onDelete" @cancel="onCancel" @undelete="onUnDelete" />
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
        <ManageTable size="small">
            <template #thead>
                <tr>
                    <th>Type</th>
                    <th>Name</th>
                    <th>Email</th>
                    <th class="hide-mobile">Created at</th>
                    <th class="hide-mobile">Updated at</th>
                    <th class="hide-mobile">Deleted at</th>
                    <th class="text-center">Operations</th>
                </tr>
                <tr class="hide-mobile">
                    <th>
                        <n-select size="small" trigger="click" :options="filterUserOptions"
                            v-model:value="userFilterType" placeholder="Search by user type">
                        </n-select>
                    </th>
                    <th>
                        <n-input size="small" placeholder="search by name..." v-model:value="filterByUsername"
                            clearable>
                            <template #prefix>
                                <n-icon>
                                    <IconSearch />
                                </n-icon>
                            </template>
                        </n-input>
                    </th>
                    <th>
                        <n-input size="small" placeholder="search by email..." v-model:value="filterByEmail" clearable>
                            <template #prefix>
                                <n-icon>
                                    <IconSearch />
                                </n-icon>
                            </template>
                        </n-input>
                    </th>
                    <th>
                        <DateFilter />
                    </th>
                    <th>
                        <DateFilter />
                    </th>
                    <th>
                        <DateFilter />
                    </th>
                    <th class="text-center">
                        <n-button-group>
                            <n-button size="small" @click="onRefresh">
                                <template #icon>
                                    <n-icon :size="22">
                                        <IconRefresh />
                                    </n-icon>
                                </template>
                                {{ t("Refresh") }}
                            </n-button>
                            <n-button size="small" @click="onAddUser">
                                <template #icon>
                                    <n-icon :size="22">
                                        <IconPlus />
                                    </n-icon>
                                </template>
                                {{ t("Add") }}
                            </n-button>
                        </n-button-group>
                    </th>
                </tr>
            </template>
            <template #tbody>
                <tr v-for="user, index in users" :key="user.id">
                    <td class="text-center">
                        <span style="display: flex; align-items: center;" v-if="user.isSuperUser">
                            <n-icon :size="16" style="margin-right: 6px;">
                                <IconUserKey color="red" />
                            </n-icon>
                            Administrator
                        </span>
                        <span style="display: flex; align-items: center;" v-else>
                            <n-icon :size="16" style="margin-right: 6px;">
                                <IconUser />
                            </n-icon>
                            User
                        </span>
                    </td>
                    <td>
                        <div style="display: flex; align-items: center; gap: 8px;">
                            <n-avatar v-if="user.avatarUrl" :src="user.avatarUrl" class="avatar" /> {{ user.name }}
                        </div>
                    </td>
                    <td><a :href="'mailto:' + user.email">{{ user.email }}</a></td>
                    <td class="hide-mobile">{{ user.createdAt.toLocaleString() }}</td>
                    <td class="hide-mobile">{{ user.updatedAt?.toLocaleString() }}</td>
                    <td class="hide-mobile">{{ user.deletedAt?.toLocaleString() }}</td>
                    <td class="text-center">
                        <n-button-group v-if="!user.deletedAt">
                            <n-button size="small" @click="onUpdateUser(user, index)">
                                {{ t("Update") }}
                                <template #icon>
                                    <n-icon :size="22">
                                        <IconEdit />
                                    </n-icon>
                                </template>
                            </n-button>
                            <n-button size="small" @click="confirmDelete(user, index)"
                                :disabled="user.id === sessionStore.sessionUserId">
                                {{ t("Delete") }}
                                <template #icon>
                                    <n-icon :size="22">
                                        <IconTrash />
                                    </n-icon>
                                </template>
                            </n-button>
                        </n-button-group>
                        <n-button size="small" @click="confirmUnDelete(user, index)" v-else>
                            {{ t("Restore") }}
                            <template #icon>
                                <n-icon :size="22">
                                    <IconTrashOff />
                                </n-icon>
                            </template>
                        </n-button>
                    </td>
                </tr>
                <tr>
                    <td colspan="7" v-if="totalResuls < 1 && !state.ajaxRunning">
                        <n-empty :description="t('No results')">
                        </n-empty>
                    </td>
                </tr>
            </template>
        </ManageTable>
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

    @media (max-width: 768px) {
        .hide-mobile {
            display: none;
        }
    }
</style>