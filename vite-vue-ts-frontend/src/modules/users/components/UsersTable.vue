<script setup lang="ts">
    import { h, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NButtonGroup, NButton, NFlex, NEmpty, NAvatar, NIcon } from 'naive-ui';
    import { IconUserKey, IconUser, IconEdit, IconPlus, IconRefresh, IconTrash, IconTrashOff } from '@tabler/icons-vue';

    import { User } from '../models/user';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import { type SortOrder } from '../../../shared/types/common';
    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import { type UserTypeFilter, UserTypeFilterValue } from '../types/filter-user-type';
    import FilterUserTypeSelector from '../components/FilterUserTypeSelector.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import DateFilter from '../../../shared/components/forms/DateFilter.vue';
    import TableCellHeaderSortIcon from '../../../shared/components/tables/TableCellHeaderSortIcon.vue';
    import { useSessionStore } from '../../../stores/session';

    interface Props {
        loading: boolean;
        users: User[];
        sortField: string;
        sortOrder: SortOrder;
    }

    const { t } = useI18n();

    const sessionStore = useSessionStore();

    const emit = defineEmits(['refresh', 'add', 'update', 'delete', 'undelete', 'toggleSort', 'textfilterKeydownEnter']);

    const props = defineProps<Props>();

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

    const userTypeFiter = defineModel<UserTypeFilter>("userTypeFilter", {
        default: UserTypeFilterValue.None,
    });

    const userNameFilter = defineModel<string>("userNameFilter", {
        default: "",
    });

    const emailFilter = defineModel<string>("emailFilter", {
        default: "",
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
            title: t("Delete user"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    t("deleteUserConfirmation", { name: user.name }),
                    h('br'),
                    h('br'),
                    t("Do you want to continue ?"),
                ]),
            positiveText: t("Delete"),
            negativeText: t("Cancel"),
            onPositiveClick: () => {
                emit("delete", user, index)
            },
        });
    };

    const onConfirmUnDelete = (user: User, index: number) => {
        dialog.warning({
            title: t("Restore user"),
            icon: renderIcon(IconTrashOff)(24),
            content: () =>
                h('div', [
                    t("restoreUserConfirmation", { name: user.name }),
                    h('br'),
                    h('br'),
                    t("Do you want to continue ?"),
                ]),
            positiveText: t("Restore"),
            negativeText: t("Cancel"),
            onPositiveClick: () => {
                emit("undelete", user, index)
            },
        })
    };

    const onTextFilterKeyDownEnter = () => {
        emit("textfilterKeydownEnter");
    };

    const onCreatedAtFilterChanged = (e: any) => {
        console.log(e);
    }

    const onUpdatedAtFilterChanged = (e: any) => {
        console.log(e);
    }

    const onDeletedAtFilterChanged = (e: any) => {
        console.log(e);
    }
</script>

<template>
    <ManageTable size="small">
        <template #thead>
            <tr class="table-header-click-action">
                <th v-for="column in columns" :key="column.field" @click="onToggleSort(column.field)">
                    <n-flex justify="space-between">
                        <span>{{ column.label }}</span>
                        <TableCellHeaderSortIcon v-if="props.sortField === column.field" :order="props.sortOrder" />
                    </n-flex>
                </th>
                <th class="text-center">{{ t("Actions") }}</th>
            </tr>
            <tr class="hide-mobile">
                <th>
                    <FilterUserTypeSelector size="small" v-model:value="userTypeFiter" />
                </th>
                <th>
                    <TextFilterInput clearable size="small" :placeholder="t('searchByNameDefaultPlaceholder')"
                        v-model:value="userNameFilter" @keydown-enter="onTextFilterKeyDownEnter" />
                </th>
                <th>
                    <TextFilterInput clearable size="small" :placeholder="t('searchByEmailDefaultPlaceholder')"
                        v-model:value="emailFilter" @keydown-enter="onTextFilterKeyDownEnter" />
                </th>
                <th>
                    <DateFilter @timestamp-range-change="onCreatedAtFilterChanged" />
                </th>
                <th>
                    <DateFilter @timestamp-range-change="onUpdatedAtFilterChanged" />
                </th>
                <th>
                    <DateFilter @timestamp-range-change="onDeletedAtFilterChanged" />
                </th>
                <th class="text-center">
                    <n-button-group size="small">
                        <n-button @click="onRefresh">
                            <template #icon>
                                <n-icon :size="22">
                                    <IconRefresh />
                                </n-icon>
                            </template>
                            {{ t("Refresh") }}
                        </n-button>
                        <n-button @click="onAdd">
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
                    <!-- TODO: hide icon /label on small screens ? -->
                    <span class="doneo-flex-center-align">
                        <n-icon :size="16" style="margin-right: 6px;"
                            :component="user.isSuperUser ? IconUserKey : IconUser"
                            :color="user.isSuperUser ? 'red' : undefined">
                        </n-icon>
                        {{ t(user.isSuperUser ? "AdminTypeValue" : "UserTypeValue") }}
                    </span>
                </td>
                <td>
                    <div class="doneo-flex-center-align" style="gap: 8px;">
                        <n-avatar v-if="user.avatarUrl" :src="user.avatarUrl" class="avatar" />
                        {{ user.name }}
                    </div>
                </td>
                <td><a :href="'mailto:' + user.email">{{ user.email }}</a></td>
                <td class="hide-mobile">{{ user.createdAt.toLocaleString() }}</td>
                <td class="hide-mobile">{{ user.updatedAt?.toLocaleString() }}</td>
                <td class="hide-mobile">{{ user.deletedAt?.toLocaleString() }}</td>
                <td class="text-center">
                    <n-button-group v-if="!user.deletedAt" size="small">
                        <n-button @click="onUpdate(user, index)" :disabled="props.loading">
                            {{ t("Update") }}
                            <template #icon>
                                <n-icon :size="22">
                                    <IconEdit />
                                </n-icon>
                            </template>
                        </n-button>
                        <n-button @click="onConfirmDelete(user, index)"
                            :disabled="user.id === sessionStore.sessionUserId || props.loading">
                            {{ t("Delete") }}
                            <template #icon>
                                <n-icon :size="22">
                                    <IconTrash />
                                </n-icon>
                            </template>
                        </n-button>
                    </n-button-group>
                    <n-button size="small" @click="onConfirmUnDelete(user, index)" :disabled="props.loading" v-else>
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
                <td :colspan="columns.length + 1" v-if="users.length < 1 && !props.loading">
                    <n-empty :description="t('No users found')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped>
    .avatar {
        margin-right: 4px;
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