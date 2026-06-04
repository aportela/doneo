<script setup lang="ts">
    import { h, ref, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NEmpty, NIcon } from 'naive-ui';
    import { IconUserKey, IconUser, IconTrash, IconTrashOff } from '@tabler/icons-vue';

    import { useSessionStore } from '../../../stores/session';

    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import type { Sort } from '../../../shared/types/models/sort.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { UsersTableFilters } from '../types/users-table-filters.ts';
    import { UserPermissionFilterValue } from '../types/user-admin-permission-filter';
    import type { DateFilterSelectComponent } from './date-filter-select-component.ts';
    import { User } from '../models/user';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import UserPermissionsFilterSelector from '../components/UserPermissionsFilterSelector.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import ClearFiltersTableButton from '../../../shared/components/tables/ClearFiltersTableButton.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';

    interface Props {
        disabled: boolean;
        items: User[];
        sort?: Sort;
    }

    const { t } = useI18n();
    const sessionStore = useSessionStore();
    const dialog = useDialog();

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

    const hasFilters = computed<boolean>(() =>
        isFilteredByPermissions.value ||
        isFilteredByName.value ||
        isFilteredByEmail.value ||
        isFilteredByCreatedAt.value ||
        isFilteredByUpdatedAt.value ||
        isFilteredByDeletedAt.value
    );

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("modules.user.components.UsersTable.header.columns.permissions"),
            field: "permissions",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByPermissions.value,

        },
        {
            label: t("modules.user.components.UsersTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByName.value,
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.email"),
            field: "email",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByEmail.value,
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.createdAt"),
            field: "createdAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByCreatedAt.value,
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.updatedAt"),
            field: "updatedAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByUpdatedAt.value,
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.deletedAt"),
            field: "deletedAt",
            visible: true,
            sortable: true,
        },
    ]);

    const onSort = (sort: Sort) => {
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
    <ManageTable size="small" :columns="columns" :current-sort="sort" @sort="onSort" @refresh="onRefresh" @add="onAdd">
        <template #thead>
            <tr>
                <th>
                    <UserPermissionsFilterSelector size="small" v-model:value="filters.permissions"
                        :disabled="props.disabled" />
                </th>
                <th>
                    <TextFilterInput clearable size="small"
                        :placeholder="t('modules.user.components.UsersTable.header.filters.name.placeholder')"
                        v-model:value="filters.name" @keydown-enter="onRefresh" :disabled="props.disabled" />
                </th>
                <th>
                    <TextFilterInput clearable size="small"
                        :placeholder="t('modules.user.components.UsersTable.header.filters.email.placeholder')"
                        v-model:value="filters.email" @keydown-enter="onRefresh" :disabled="props.disabled" />
                </th>
                <th>
                    <DateFilterSelect v-model:range="filters.createdAt" ref="createdAtFilterRef"
                        :disabled="props.disabled" />
                </th>
                <th>
                    <DateFilterSelect v-model:range="filters.updatedAt" ref="updatedAtFilterRef"
                        :disabled="props.disabled" />
                </th>
                <th>
                    <DateFilterSelect v-model:range="filters.deletedAt" ref="deletedAtFilterRef"
                        :disabled="props.disabled" />
                </th>
                <th>
                    <ClearFiltersTableButton @clear="onClearFilters" :disabled="props.disabled || !hasFilters" />
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="user, index in items" :key="user.id ?? index">
                <td class="doneo-text-center">
                    <span class="doneo-flex-center-align">
                        <n-icon :size="16" style="margin-right: 6px;"
                            :component="user.permissions?.isSuperUser ? IconUserKey : IconUser"
                            :color="user.permissions?.isSuperUser ? 'red' : undefined">
                        </n-icon>
                        {{ t(user.permissions?.isSuperUser ?
                            "modules.user.components.UsersTable.body.columns.permissions.administrator" :
                            "modules.user.components.UsersTable.body.columns.permissions.user") }}
                    </span>
                </td>
                <td>
                    <AvatarUserName :user-id="user.id" :user-name="user.name" />
                </td>
                <td><a :href="'mailto:' + user.email">{{ user.email }}</a></td>
                <td>{{ user.createdAt?.toLocaleString() }}</td>
                <td>{{ user.updatedAt?.toLocaleString() }}</td>
                <td>{{ user.deletedAt?.toLocaleString() }}</td>
                <td class="doneo-text-center">
                    <ManageTableActionButtons show-update show-delete show-restore
                        :update-disabled="props.disabled || !!user.deletedAt?.msTimestamp"
                        :delete-disabled="props.disabled || sessionStore.sessionUserId === user.id || !!user.deletedAt?.msTimestamp"
                        :restored-disabled="props.disabled || !user.deletedAt?.msTimestamp" :disabled="props.disabled"
                        @update="onUpdate(user, index)" @delete="onConfirmDelete(user, index)"
                        @restore="onConfirmUnDelete(user, index)" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="!props.disabled && items.length < 1">
                    <n-empty :description="t('modules.user.components.UsersTable.warnings.noItemsFound')" />
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>