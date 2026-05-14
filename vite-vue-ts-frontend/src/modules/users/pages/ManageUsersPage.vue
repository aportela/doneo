<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard } from 'naive-ui';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';
    import { userService } from '../services/user';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { UserResponse } from '../types/dto';
    import { User } from '../models/user';
    import { type UserTypeFilter, UserTypeFilterValue } from '../types/filter-user-type';
    import UsersTable from '../components/UsersTable.vue';
    import UserForm from '../components/UserForm.vue';
    import Pager from '../../../shared/components/tables/Pager.vue';
    import { Sort } from '../../../shared/types/models/sort';
    import type { TimestampRange } from '../../../shared/composables/timestamps';

    const { notify } = useNotify();

    const { t } = useI18n();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const users = shallowRef<User[]>([]);

    const sort = ref<Sort>(new Sort("name", "ASC"));

    const nameFilter = ref<string>("");
    const emailFilter = ref<string>("");
    const typeFilter = ref<UserTypeFilter>(UserTypeFilterValue.None);

    const createdAtFilter = ref<TimestampRange>({ from: null, to: null });
    const updatedAtFilter = ref<TimestampRange>({ from: null, to: null });
    const deletedAtFilter = ref<TimestampRange>({ from: null, to: null });

    const currentPage = ref(1);
    const pageSize = ref(10);
    const totalResults = ref(0);
    const totalPages = ref(0);

    const showUserDialog = ref<boolean>(false);
    const userDialogMode = ref<string>("add");

    const selectedUserId = ref<string>("");

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    watch([nameFilter, emailFilter, typeFilter], () => {
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
        sort.value.toggleSort(field);
        onRefresh();
    };

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
                    field: sort.value.field,
                    sort: sort.value.order,
                },
                filter: {
                    name: nameFilter.value,
                    email: emailFilter.value,
                    administratorFlag: typeFilter.value == UserTypeFilterValue.None ? undefined : (typeFilter.value === UserTypeFilterValue.OnlyAdministrators ? true : false),
                    createdAt: createdAtFilter.value,
                    updatedAt: updatedAtFilter.value,
                    deletedAt: deletedAtFilter.value
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
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageUsersPage.onRefresh" } });
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

    const onDelete = async (user: User, _index: number) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await userService.delete(user.id);
            notify('success', t("userDeletedNotification", { name: user.name }));
            onRefresh();
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageUsersPage.onDelete" } });
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("We couldn’t find the specified user");
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while deleting the user");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while deleting the user");
                    console.error("Unhandled API error", { file: "ManageUsersPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    const onUnDelete = async (user: User, _index: number) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await userService.unDelete(user.id);
            notify('success', t("userRestoredNotification", { name: user.name }));
            onRefresh();
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageUsersPage.onUnDelete" } });
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("We couldn’t find the specified user");
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while restoring the user");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while restoring the user");
                    console.error("Unhandled API error", { file: "ManageUsersPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageUsersPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageUsersPage.onRefresh")) {
                // TODO: missing user/index param at this point
                //onDelete();
            } else if (payload.to.includes("ManageUsersPage.onRefresh")) {
                // TODO: missing user/index param at this point
                //onUnDelete();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
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
        <UsersTable :users="users" :loading="state.ajaxRunning" @refresh="onRefresh" @add="onShowAddForm"
            @update="onShowUpdateForm" @delete="onDelete" @undelete="onUnDelete" @textfilter-keydown-enter="onRefresh"
            :sort-field="sort.field" :sort-order="sort.order" @toggle-sort="onToggleSort"
            v-model:user-type-filter="typeFilter" v-model:user-name-filter="nameFilter"
            v-model:email-filter="emailFilter" v-model:created-at-filter="createdAtFilter"
            v-model:updated-at-filter="updatedAtFilter" v-model:deleted-at-filter="deletedAtFilter" />
    </n-card>
</template>

<style lang="css"></style>