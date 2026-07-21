<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, h } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, useDialog, NButtonGroup, NButton, NIcon } from 'naive-ui';

    import { DONEO_ICON_ACTION_DELETE, DONEO_ICON_ACTION_EDIT, DONEO_ICON_ADMIN_USER, DONEO_ICON_ACTIONRESTORE, DONEO_ICON_USER } from '../../../shared/types/icons.ts';

    import { useLoadingStore } from '../../../stores/loading';
    import { useCacheStore } from '../../../stores/cache.ts';
    import { useSessionStore } from '../../../stores/session';
    import { useUserSettingsStore } from '../../../stores/userSettings.ts';
    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { renderIcon } from '../../../shared/composables/naive-ui-icon';

    import { User } from '../models/user';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { PAGER_DEFAULT_RESULTS_PAGE, type Pagination } from '../../../shared/types/pager.ts';
    import type { Order } from '../../../shared/types/order.ts';
    import { userService } from '../services/user';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { SearchRequest, UserResponse } from '../types/dto.ts';
    import type { TimestampRange } from '../../../shared/composables/timestamps.ts';
    import { UserPermissionFilterValue, type UserPermissionFilter } from '../types/user-admin-permission-filter';

    import UserForm from './UserForm.vue';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import UserPermissionsFilterSelector from '../components/UserPermissionsFilterSelector.vue';
    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';

    interface Props {
        id?: string;
    }

    const props = withDefaults(defineProps<Props>(), { id: "UsersTable" });;

    const { t } = useI18n();
    const dialog = useDialog();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();
    const sessionStore = useSessionStore();
    const userSettingsStore = useUserSettingsStore();
    const tableSettingsStore = useTableSettingsStore();
    const cacheStore = useCacheStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    watch(
        () => state.ajaxRunning,
        (ajaxRunning) => {
            loadingStore.set(ajaxRunning);
        }
    );

    const items = shallowRef<User[]>([]);
    const tmpItem = ref<User>(new User());

    const order = reactive<Order>({ field: "name", direction: "ASC" });

    const onSort = (newOrder: Order) => {
        order.field = newOrder.field;
        order.direction = newOrder.direction;
        onRefresh();
    };

    const pagination = reactive<Pagination>({ enabled: true, currentPage: 1, resultsPage: PAGER_DEFAULT_RESULTS_PAGE, totalPages: 1, totalResults: 0 });
    const resetPager = ref<boolean>(false);

    const onPagerChanged = (newPagination: Pagination) => {
        pagination.enabled = newPagination.enabled;
        pagination.currentPage = newPagination.currentPage;
        pagination.resultsPage = newPagination.resultsPage;
        onRefresh();
    };

    const createdAtFilterRef = ref<InstanceType<typeof DateFilterSelect>[] | null>(null);
    const updatedAtFilterRef = ref<InstanceType<typeof DateFilterSelect>[] | null>(null);
    const deletedAtFilterRef = ref<InstanceType<typeof DateFilterSelect>[] | null>(null);

    interface UsersTableFilters {
        permissions: UserPermissionFilter;
        name: string;
        email: string;
        createdAt: TimestampRange;
        updatedAt: TimestampRange;
        deletedAt: TimestampRange;
    }

    const filters = reactive<UsersTableFilters>(
        {
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
            },
        }
    );

    watch(
        () => [
            filters.permissions,
            filters.name,
            filters.email,
            filters.createdAt.from,
            filters.createdAt.to,
            filters.updatedAt.from,
            filters.updatedAt.to,
            filters.deletedAt.from,
            filters.deletedAt.to,
        ],
        () => {
            resetPager.value = true;
        },
    );

    const isFilteredByPermissions = computed<boolean>(() => filters.permissions != UserPermissionFilterValue.Any);
    const isFilteredByName = computed<boolean>(() => filters.name !== "");
    const isFilteredByEmail = computed<boolean>(() => filters.email !== "");
    const isFilteredByCreatedAt = computed<boolean>(() => filters.createdAt.from != null || filters.createdAt.to != null);
    const isFilteredByUpdatedAt = computed<boolean>(() => filters.updatedAt.from != null || filters.updatedAt.to != null);
    const isFilteredByDeletedAt = computed<boolean>(() => filters.deletedAt.from != null || filters.deletedAt.to != null);

    const onClearFilters = () => {
        filters.permissions = UserPermissionFilterValue.Any;
        filters.name = "";
        filters.email = "";
        if (createdAtFilterRef.value) {
            createdAtFilterRef.value[0]?.reset();
        }
        if (updatedAtFilterRef.value) {
            updatedAtFilterRef.value[0]?.reset();
        }
        if (deletedAtFilterRef.value) {
            deletedAtFilterRef.value[0]?.reset();
        }
    };

    const columnDefinitions = reactive<TableHeaderColumn<User>[]>([
        {
            label: t("modules.user.components.UsersTable.header.columns.permissions"),
            field: "permissions",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByPermissions.value,
            render: (row: User) => {
                return h(
                    "span",
                    {
                        class: "doneo-flex-center-align",
                    },
                    [
                        h(
                            NIcon,
                            {
                                size: 16,
                                style: {
                                    marginRight: "6px",
                                },
                                component: row.permissions?.isSuperUser ? DONEO_ICON_ADMIN_USER : DONEO_ICON_USER,
                                color: row.permissions?.isSuperUser ? "red" : undefined,
                            }
                        ),
                        t(
                            row.permissions?.isSuperUser
                                ? "modules.user.components.UsersTable.body.columns.permissions.administrator"
                                : "modules.user.components.UsersTable.body.columns.permissions.user"
                        ),
                    ]
                );
            }
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByName.value,
            render: (row: User) => {
                return h(
                    AvatarUserName,
                    {
                        userId: row.id,
                        userName: row.name
                    }
                );
            }
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.email"),
            field: "email",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByEmail.value,
            render: (row: User) => {
                return h(
                    "a",
                    {
                        href: `mailto:${row.email}`,
                    },
                    row.email
                );
            }
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.createdAt"),
            field: "createdAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByCreatedAt.value,
            render: (row: User) => {
                return h(
                    "span",
                    {},
                    {
                        default: () => row.createdAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask)

                    }
                );
            }
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.updatedAt"),
            field: "updatedAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByUpdatedAt.value,
            render: (row: User) => {
                return h(
                    "span",
                    {},
                    {
                        default: () => row.updatedAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask)
                    }
                );
            }
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.deletedAt"),
            field: "deletedAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByDeletedAt.value,
            render: (row: User) => {
                return h(
                    "span",
                    {},
                    {
                        default: () => row.deletedAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask)
                    }
                );
            }
        },
    ]);

    // create (if not found) default settings for this table (column order & visibility)
    tableSettingsStore.register(props.id, { columns: columnDefinitions.map((column) => { return { field: column.field, visible: column.visible } }) ?? [] });

    // restore previous settings
    const tableSettings = tableSettingsStore.get(props.id);

    // build columns based on saved order visibility settings
    const columns = computed<TableHeaderColumn<User>[]>(() =>
        tableSettings.columns.map((column) => { // get saved ordered columns
            const definition = columnDefinitions.find((c) => c.field === column.field);
            return {
                label: definition?.label ?? "",
                field: column.field,
                visible: column.visible,
                sortable: definition!.sortable,
                align: definition?.align,
                isFiltered: definition?.isFiltered ?? (() => false),
                render: definition?.render ?? (() => "")
            };
        })
    );

    const onDelete = async (user: User) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await userService.delete(user.id);
            cacheStore.clearUsersCache();
            notify('success', t("modules.user.components.UsersTable.notifications.userDeleted", { name: user.name }));
            onRefresh();
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            tmpItem.value = user;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "UsersTable.onDelete" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.user.components.UsersTable.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.user.components.UsersTable.errors.deleteError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.user.components.UsersTable.errors.deleteError");
                    console.error("Unhandled API error", { file: "UsersTable.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onConfirmDelete = (user: User) => {
        dialog.warning({
            title: t("modules.user.components.UsersTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(DONEO_ICON_ACTION_DELETE)(24),
            content: () =>
                h('div', [
                    t("modules.user.components.UsersTable.dialogs.deleteConfirmation.message", { name: user.name }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                onDelete(user);
            },
        });
    };

    const onUnDelete = async (user: User) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await userService.unDelete(user.id);
            cacheStore.clearUsersCache();
            notify('success', t("modules.user.components.UsersTable.notifications.userRestored", { name: user.name }));
            onRefresh();
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            tmpItem.value = user;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "UsersTable.onUnDelete" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.user.components.UsersTable.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.user.components.UsersTable.errors.restoreError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.user.components.UsersTable.errors.restoreError");
                    console.error("Unhandled API error", { file: "UsersTable.vue", method: "onUnDelete" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onConfirmUnDelete = (user: User) => {
        dialog.warning({
            title: t("modules.user.components.UsersTable.dialogs.undeleteConfirmation.title"),
            icon: renderIcon(DONEO_ICON_ACTIONRESTORE)(24),
            content: () =>
                h('div', [
                    t("modules.user.components.UsersTable.dialogs.undeleteConfirmation.message", { name: user.name }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Restore.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                onUnDelete(user);
            },
        })
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: SearchRequest = {
                pager: {
                    enabled: pagination.enabled,
                    currentPage: resetPager.value ? 1 : pagination.currentPage,
                    resultsPage: pagination.resultsPage,
                },
                order: {
                    field: order.field,
                    direction: order.direction,
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
            pagination.enabled = response.pager.enabled;
            pagination.totalPages = response.pager.totalPages;
            pagination.totalResults = response.pager.totalResults;
            pagination.currentPage = response.pager.currentPage;
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
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "UsersTable.onRefresh" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.user.components.UsersTable.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.user.components.UsersTable.errors.refreshError");
                    console.error("Unhandled API error", { file: "UsersTable.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const showFormModal = ref<boolean>(false);

    const onAdd = () => {
        tmpItem.value = new User();
        showFormModal.value = true;
    };

    const onUpdate = (user: User) => {
        tmpItem.value = user;
        showFormModal.value = true;
    };

    const onUserAdded = (user: UserResponse) => {
        cacheStore.clearUsersCache();
        showFormModal.value = false;
        tmpItem.value = new User();
        notify('success', t("modules.user.components.UsersTable.notifications.userAdded", { name: user.name }));
        onRefresh();
    };

    const onUserUpdated = (user: UserResponse) => {
        cacheStore.clearUsersCache();
        showFormModal.value = false;
        tmpItem.value = new User();
        notify('success', t("modules.user.components.UsersTable.notifications.userUpdated", { name: user.name }));
        onRefresh();
    };

    const hideFormModal = () => {
        showFormModal.value = false;
        tmpItem.value = new User();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("UsersTable.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("UsersTable.onDelete")) {
                onDelete(tmpItem.value);
            } else if (payload.to.includes("UsersTable.onUnDelete")) {
                onUnDelete(tmpItem.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-modal v-model:show="showFormModal" v-if="showFormModal">
        <UserForm class="user-form" :user-id="tmpItem.id" @add="onUserAdded" @update="onUserUpdated"
            @cancel="hideFormModal" />
    </n-modal>
    <ManageTable :id="props.id" size="small" :disabled="state.ajaxRunning" :rows="items" :row-key="row => row.id"
        :columns="columns" :order="order" :pager-data="pagination" pager-position="both" @sort="onSort"
        @refresh="onRefresh" @add="onAdd" @pager-changed="onPagerChanged" @clear-filters="onClearFilters">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <UserPermissionsFilterSelector v-if="column.field === 'permissions'" size="small"
                    v-model:value="filters.permissions" :disabled="state.ajaxRunning" />
                <TextFilterInput v-else-if="column.field === 'name'" clearable size="small"
                    :placeholder="t('modules.user.components.UsersTable.header.filters.name.placeholder')"
                    v-model:value="filters.name" @keydown-enter="onRefresh" :disabled="state.ajaxRunning" />
                <TextFilterInput v-else-if="column.field === 'email'" clearable size="small"
                    :placeholder="t('modules.user.components.UsersTable.header.filters.email.placeholder')"
                    v-model:value="filters.email" @keydown-enter="onRefresh" :disabled="state.ajaxRunning" />
                <DateFilterSelect v-else-if="column.field === 'createdAt'" clearable v-model:range="filters.createdAt"
                    ref="createdAtFilterRef" :disabled="state.ajaxRunning" />
                <DateFilterSelect v-else-if="column.field === 'updatedAt'" clearable v-model:range="filters.updatedAt"
                    ref="updatedAtFilterRef" :disabled="state.ajaxRunning" />
                <DateFilterSelect v-else-if="column.field === 'deletedAt'" clearable v-model:range="filters.deletedAt"
                    ref="deletedAtFilterRef" :disabled="state.ajaxRunning" />
            </th>
        </template>
        <template #rowactions="{ row }">
            <n-button-group class="doneo-table-actions-button-group" size="small">
                <n-button @click="onUpdate(row)" :disabled="state.ajaxRunning || row.deletedAt?.hasValue()"
                    class="doneo-table-actions-button">
                    {{ t("shared.buttons.Edit.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_EDIT" />
                    </template>
                </n-button>
                <n-button @click="onConfirmDelete(row)"
                    :disabled="state.ajaxRunning || row.deletedAt?.hasValue() || sessionStore.sessionUserId === row.id"
                    class="doneo-table-actions-button">
                    {{ t("shared.buttons.Delete.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTION_DELETE" />
                    </template>
                </n-button>
                <n-button @click="onConfirmUnDelete(row)" :disabled="state.ajaxRunning || !row.deletedAt?.hasValue()"
                    class="doneo-table-actions-button">
                    {{ t("shared.buttons.Restore.label") }}
                    <template #icon>
                        <n-icon :component="DONEO_ICON_ACTIONRESTORE" />
                    </template>
                </n-button>
            </n-button-group>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped>
    .user-form {
        width: 95%;
        max-width: 640px;
    }
</style>