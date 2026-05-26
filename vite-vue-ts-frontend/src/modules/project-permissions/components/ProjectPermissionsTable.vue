<script setup lang="ts">
    import { ref, h, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NEmpty } from 'naive-ui';
    import { IconTrash } from '@tabler/icons-vue';

    import { ProjectPermission } from '../models/project-permission.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import RefreshAddActionsColumn from '../../../shared/components/tables/RefreshAddActionsColumn.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';

    interface Props {
        loading: boolean;
        projectPermissions: ProjectPermission[];
    }

    const { t } = useI18n();

    const emit = defineEmits(['refresh', 'add', 'delete']);

    const props = defineProps<Props>();

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("modules.projectPermission.components.projectPermissionsTable.header.columns.user"),
            field: "name",
            sortable: false,
        },
        {
            label: t("modules.projectPermission.components.projectPermissionsTable.header.columns.role"),
            field: "role",
            sortable: false,
        },
        {
            label: t("modules.projectPermission.components.projectPermissionsTable.header.columns.permissions"),
            field: "permissions",
            sortable: false,
        },
    ]);

    const filterByUsername = ref<string>("");
    const filterByRole = ref<string>("");

    const dialog = useDialog();

    const onRefresh = () => {
        emit("refresh");
    };

    const onAdd = () => {
        emit("add");
    };


    const onConfirmDelete = (projectPermission: ProjectPermission, index: number) => {
        dialog.warning({
            title: t("modules.projectPermission.components.projectPermissionsTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    t("modules.projectPermission.components.projectPermissionsTable.dialogs.deleteConfirmation.message", { user: projectPermission.user.name, role: projectPermission.role.name }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                emit("delete", projectPermission, index)
            },
        });
    };
</script>

<template>
    <ManageTable size="small">
        <template #thead>
            <tr>
                <th v-for="column in columns" :key="column.field"
                    :class="{ 'doneo-text-center': column.align === 'center' }">
                    {{ column.label }}
                </th>
                <th class="doneo-table-actions-column">{{ t("shared.components.table.header.columns.actions") }}</th>
            </tr>
            <tr class="hide-mobile">
                <th>
                    <TextFilterInput clearable size="small"
                        :placeholder="t('modules.projectPermission.components.projectPermissionsTable.filters.user.placeholder')"
                        v-model:value="filterByUsername" />
                </th>
                <th>
                    <TextFilterInput clearable size="small"
                        :placeholder="t('modules.projectPermission.components.projectPermissionsTable.filters.role.placeholder')"
                        v-model:value="filterByRole" />
                </th>
                <th>
                </th>
                <th class="doneo-text-center">
                    <RefreshAddActionsColumn @refresh="onRefresh" @add="onAdd" />
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="projectPermission, index in projectPermissions" :key="projectPermission.id">
                <td>
                    <AvatarUserName :user-id="projectPermission.user.id" :user-name="projectPermission.user.name" />
                </td>
                <td>{{ projectPermission.role.name }}</td>
                <td>[1][2][3][4]</td>
                <td class="doneo-text-center">
                    <ManageTableActionButtons show-delete @delete="onConfirmDelete(projectPermission, index)" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="projectPermissions.length < 1 && !props.loading">
                    <n-empty
                        :description="t('modules.projectPermission.components.projectPermissionsTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped>

    .table-header-click-action th:not(:last-of-type) .n-icon {
        margin-top: 4px;
    }

    @media (max-width: 768px) {
        .hide-mobile {
            display: none;
        }
    }
</style>