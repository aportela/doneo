<script setup lang="ts">
    import { h, computed, ref } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NButtonGroup, NButton, NFlex, NEmpty, NIcon } from 'naive-ui';
    import { IconUserKey, IconUser, IconEdit, IconTrash, IconTrashOff, IconFilterOff, IconFilter } from '@tabler/icons-vue';

    import { User } from '../models/user';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import { type SortOrder } from '../../../shared/types/common';
    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import { type UserPermissionFilter, UserPermissionFilterValue } from '../types/user-admin-permission-filter';
    import UserPermissionsFilterSelector from '../components/UserPermissionsFilterSelector.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import TableCellHeaderSortIcon from '../../../shared/components/tables/TableCellHeaderSortIcon.vue';
    import RefreshAddActionsColumn from '../../../shared/components/tables/RefreshAddActionsColumn.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import { useSessionStore } from '../../../stores/session';
    import type { TimestampRange } from '../../../shared/composables/timestamps';

    interface Props {
        loading: boolean;
        users: User[];
        sortField: string;
        sortOrder: SortOrder;
    }

    const { t } = useI18n();

    const sessionStore = useSessionStore();

    const emit = defineEmits(['refresh', 'add', 'update', 'delete', 'undelete', 'toggleSort', 'textfilterKeydownEnter', 'createdAtFilterChange', 'updatedAtFilterChange', 'deletedAtFilterChange']);

    const props = defineProps<Props>();

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("modules.user.components.UsersTable.header.columns.permissions"),
            field: "isSuperUser",
            sortable: true,
            isFiltered: () => userPermissionsFilter.value != UserPermissionFilterValue.Any,
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.name"),
            field: "name",
            sortable: true,
            isFiltered: () => userNameFilter.value?.length > 0,
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.email"),
            field: "email",
            sortable: true,
            isFiltered: () => emailFilter.value?.length > 0,
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.createdAt"),
            field: "createdAt",
            sortable: true,
            isFiltered: () => createdAtFilter.value?.from != null || createdAtFilter.value?.to != null,
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.updatedAt"),
            field: "updatedAt",
            sortable: true,
            isFiltered: () => updatedAtFilter.value?.from != null || updatedAtFilter.value?.to != null,
        },
        {
            label: t("modules.user.components.UsersTable.header.columns.deletedAt"),
            field: "deletedAt",
            sortable: true,
            isFiltered: () => deletedAtFilter.value?.from != null || deletedAtFilter.value?.to != null,
        },
    ]);

    const userPermissionsFilter = defineModel<UserPermissionFilter>("userTypeFilter", {
        default: UserPermissionFilterValue.Any,
    });

    const userNameFilter = defineModel<string>("userNameFilter", {
        default: "",
    });

    const emailFilter = defineModel<string>("emailFilter", {
        default: "",
    });

    const createdAtFilter = defineModel<TimestampRange>("createdAtFilter", {
        default: {
            from: null,
            to: null
        }
    });

    const updatedAtFilter = defineModel<TimestampRange>("updatedAtFilter", {
        default: {
            from: null,
            to: null
        }
    });

    const deletedAtFilter = defineModel<TimestampRange>("deletedAtFilter", {
        default: {
            from: null,
            to: null
        }
    });

    const dialog = useDialog();

    const onToggleSort = (field: string) => {
        emit("toggleSort", field);
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

    const onTextFilterKeyDownEnter = () => {
        emit("textfilterKeydownEnter");
    };

    interface DateFilterSelectComponent {
        reset: () => void
    };

    const createdAtFilterRef = ref<DateFilterSelectComponent | undefined>();
    const updatedAtFilterRef = ref<DateFilterSelectComponent | undefined>();
    const deletedAtFilterRef = ref<DateFilterSelectComponent | undefined>();

    const onClearFilters = () => {
        userPermissionsFilter.value = UserPermissionFilterValue.Any;
        userNameFilter.value = "";
        emailFilter.value = "";
        createdAtFilterRef.value?.reset();
        updatedAtFilterRef.value?.reset();
        deletedAtFilterRef.value?.reset();
    };

    const hasFilters = computed<boolean>(() =>
        userPermissionsFilter.value != UserPermissionFilterValue.Any ||
        userNameFilter.value?.length > 0 ||
        emailFilter.value?.length > 0 ||
        createdAtFilter.value?.from != null || createdAtFilter.value?.to != null ||
        updatedAtFilter.value?.from != null || updatedAtFilter.value?.to != null ||
        deletedAtFilter.value?.from != null || deletedAtFilter.value?.to != null
    );

</script>

<template>
    <ManageTable size="small">
        <template #thead>
            <tr>
                <th v-for="column in columns" :key="column.field" @click="onToggleSort(column.field)"
                    class="doneo-cursor-pointer">
                    <n-flex justify="space-between">
                        <span>{{ column.label }}</span>
                        <div>
                            <n-icon :size="16" :component="IconFilter" style="margin-top: 4px;"
                                v-if="column.isFiltered?.() ?? false" />
                            <TableCellHeaderSortIcon v-if="props.sortField === column.field" :order="props.sortOrder" />
                        </div>
                    </n-flex>
                </th>
                <!--
                <th class="doneo-table-actions-column">{{ t("shared.components.table.header.columns.actions") }}</th>
                -->
                <th>
                    <RefreshAddActionsColumn @refresh="onRefresh" @add="onAdd" />
                </th>
            </tr>
            <tr>
                <th>
                    <UserPermissionsFilterSelector size="small" v-model:value="userPermissionsFilter" />
                </th>
                <th>
                    <TextFilterInput clearable size="small"
                        :placeholder="t('modules.user.components.UsersTable.header.filters.name.placeholder')"
                        v-model:value="userNameFilter" @keydown-enter="onTextFilterKeyDownEnter" />
                </th>
                <th>
                    <TextFilterInput clearable size="small"
                        :placeholder="t('modules.user.components.UsersTable.header.filters.email.placeholder')"
                        v-model:value="emailFilter" @keydown-enter="onTextFilterKeyDownEnter" />
                </th>
                <th>
                    <DateFilterSelect v-model:range="createdAtFilter" ref="createdAtFilterRef" />
                </th>
                <th>
                    <DateFilterSelect v-model:range="updatedAtFilter" ref="updatedAtFilterRef" />
                </th>
                <th>
                    <DateFilterSelect v-model:range="deletedAtFilter" ref="deletedAtFilterRef" />
                </th>
                <th>
                    <n-button size="small" block @click="onClearFilters" :disabled="props.loading || !hasFilters">
                        <template #icon>
                            <n-icon :component="IconFilterOff" />
                        </template>
                        {{ t("shared.components.table.filters.button.clearFilters.label") }}
                    </n-button>
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="user, index in users" :key="user.id ?? index">
                <td class="doneo-text-center">
                    <!-- TODO: hide icon /label on small screens ? -->
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
                <td class="hide-mobile">{{ user.createdAt?.toLocaleString() }}</td>
                <td class="hide-mobile">{{ user.updatedAt?.toLocaleString() }}</td>
                <td class="hide-mobile">{{ user.deletedAt?.toLocaleString() }}</td>
                <td class="doneo-text-center">
                    <n-button-group v-if="!user.deletedAt" size="small">
                        <n-button @click="onUpdate(user, index)" :disabled="props.loading">
                            {{ t("shared.buttons.Update.label") }}
                            <template #icon>
                                <n-icon :size="22">
                                    <IconEdit />
                                </n-icon>
                            </template>
                        </n-button>
                        <n-button @click="onConfirmDelete(user, index)"
                            :disabled="user.id === sessionStore.sessionUserId || props.loading">
                            {{ t("shared.buttons.Delete.label") }}
                            <template #icon>
                                <n-icon :size="22">
                                    <IconTrash />
                                </n-icon>
                            </template>
                        </n-button>
                    </n-button-group>
                    <n-button size="small" @click="onConfirmUnDelete(user, index)" :disabled="props.loading" v-else>
                        {{ t("shared.buttons.Restore.label") }}
                        <template #icon>
                            <n-icon :size="22">
                                <IconTrashOff />
                            </n-icon>
                        </template>
                    </n-button>
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="users.length < 1 && !props.loading">
                    <n-empty :description="t('modules.user.components.UsersTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>