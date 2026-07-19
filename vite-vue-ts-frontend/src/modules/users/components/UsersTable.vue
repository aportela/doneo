<script setup lang="ts">
    import { ref, reactive, computed, h } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NIcon } from 'naive-ui';
    import { IconUserKey, IconUser, IconTrash, IconTrashOff } from '@tabler/icons-vue';

    import { useSessionStore } from '../../../stores/session';
    import { useUserSettingsStore } from '../../../stores/userSettings.ts';

    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import type { Order } from '../../../shared/types/order.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { UsersTableFilters } from '../types/users-table-filters.ts';
    import { UserPermissionFilterValue } from '../types/user-admin-permission-filter';
    import type { DateFilterSelectComponent } from './date-filter-select-component.ts';
    import { User } from '../models/user';
    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import UserPermissionsFilterSelector from '../components/UserPermissionsFilterSelector.vue';
    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';

    interface Props {
        id: string;
        disabled: boolean;
        items: User[];
        order: Order;
    }

    const { t } = useI18n();
    const sessionStore = useSessionStore();
    const dialog = useDialog();
    const userSettingsStore = useUserSettingsStore();
    const tableSettingsStore = useTableSettingsStore();

    const emit = defineEmits(['refresh', 'add', 'update', 'delete', 'undelete', 'sort']);

    const props = defineProps<Props>();

    const createdAtFilterRef = ref<DateFilterSelectComponent | undefined>();
    const updatedAtFilterRef = ref<DateFilterSelectComponent | undefined>();
    const deletedAtFilterRef = ref<DateFilterSelectComponent | undefined>();

    const filters = defineModel<UsersTableFilters>("filters", {
        default: () => ({
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
        })
    });

    const isFilteredByPermissions = computed<boolean>(() => filters.value.permissions != UserPermissionFilterValue.Any);
    const isFilteredByName = computed<boolean>(() => filters.value.name.length > 0);
    const isFilteredByEmail = computed<boolean>(() => filters.value.email.length > 0);
    const isFilteredByCreatedAt = computed<boolean>(() => filters.value.createdAt.from != null || filters.value.createdAt.to != null);
    const isFilteredByUpdatedAt = computed<boolean>(() => filters.value.updatedAt.from != null || filters.value.updatedAt.to != null);
    const isFilteredByDeletedAt = computed<boolean>(() => filters.value.deletedAt.from != null || filters.value.deletedAt.to != null);

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
                        h(NIcon, {
                            size: 16,
                            style: {
                                marginRight: "6px",
                            },
                            component: row.permissions?.isSuperUser ? IconUserKey : IconUser,
                            color: row.permissions?.isSuperUser ? "red" : undefined,
                        }),
                        t(
                            row.permissions?.isSuperUser
                                ? "modules.user.components.UsersTable.body.columns.permissions.administrator"
                                : "modules.user.components.UsersTable.body.columns.permissions.user"
                        ),
                    ]
                )
            }
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByName.value,
            render: (row: User) => {
                return h(AvatarUserName, {
                    userId: row.id,
                    userName: row.name,
                })
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
                return h("span", {}, { default: () => row.createdAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) });
            }
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.updatedAt"),
            field: "updatedAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByUpdatedAt.value,
            render: (row: User) => {
                return h("span", {}, { default: () => row.updatedAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) });
            }
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.deletedAt"),
            field: "deletedAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByDeletedAt.value,
            render: (row: User) => {
                return h("span", {}, { default: () => row.deletedAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) });
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
                isFiltered: definition?.isFiltered ?? (() => false),
                render: definition?.render ?? (() => "")
            };
        })
    );

    const onSort = (sort: Order) => {
        emit("sort", sort);
    };

    const onRefresh = () => {
        emit("refresh");
    };

    const onAdd = () => {
        emit("add");
    };

    const onUpdate = (user: User, index: number) => {
        emit("update", user, index);
    };

    const onConfirmDelete = (user: User, index: number) => {
        dialog.warning({
            title: t("modules.user.components.UsersTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash)(24),
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
                emit("delete", user, index)
            },
        });
    };

    const onConfirmUnDelete = (user: User, index: number) => {
        dialog.warning({
            title: t("modules.user.components.UsersTable.dialogs.undeleteConfirmation.title"),
            icon: renderIcon(IconTrashOff)(24),
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
                emit("undelete", user, index)
            },
        })
    };

    const onClearFilters = () => {
        filters.value.permissions = UserPermissionFilterValue.Any;
        filters.value.name = "";
        filters.value.email = "";
        createdAtFilterRef.value?.reset();
        updatedAtFilterRef.value?.reset();
        deletedAtFilterRef.value?.reset();
    };
</script>

<template>
    <ManageTable id="ManageUsers" size="small" :rows="items" :row-key="row => row.id" :columns="columns" :order="order"
        @sort="onSort" @refresh="onRefresh" @add="onAdd">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <UserPermissionsFilterSelector v-if="column.field === 'permissions'" size="small"
                    v-model:value="filters.permissions" :disabled="props.disabled" />
                <TextFilterInput v-else-if="column.field === 'name'" clearable size="small"
                    :placeholder="t('modules.user.components.UsersTable.header.filters.name.placeholder')"
                    v-model:value="filters.name" @keydown-enter="onRefresh" :disabled="props.disabled" />
                <TextFilterInput v-else-if="column.field === 'email'" clearable size="small"
                    :placeholder="t('modules.user.components.UsersTable.header.filters.email.placeholder')"
                    v-model:value="filters.email" @keydown-enter="onRefresh" :disabled="props.disabled" />
                <DateFilterSelect v-else-if="column.field === 'createdAt'" clearable v-model:range="filters.createdAt"
                    ref="createdAtFilterRef" :disabled="props.disabled" />
                <DateFilterSelect v-else-if="column.field === 'updatedAt'" clearable v-model:range="filters.updatedAt"
                    ref="updatedAtFilterRef" :disabled="props.disabled" />
                <DateFilterSelect v-else-if="column.field === 'deletedAt'" clearable v-model:range="filters.deletedAt"
                    ref="deletedAtFilterRef" :disabled="props.disabled" />
            </th>
        </template>
        <template #rowactions="{ row, index }">
            <ManageTableActionButtons show-update show-delete show-restore
                :update-disabled="props.disabled || !!row.deletedAt?.msTimestamp"
                :delete-disabled="props.disabled || sessionStore.sessionUserId === row.id || !!row.deletedAt?.msTimestamp"
                :restored-disabled="props.disabled || !row.deletedAt?.msTimestamp" :disabled="props.disabled"
                @update="onUpdate(row, index)" @delete="onConfirmDelete(row, index)"
                @restore="onConfirmUnDelete(row, index)" @clear-filters="onClearFilters"
                :no-items-warning-message="t('modules.project.components.UsersTable.warnings.noItemsFound')"
                :show-no-items-warning-message="items.length < 1 && !props.disabled" />
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>