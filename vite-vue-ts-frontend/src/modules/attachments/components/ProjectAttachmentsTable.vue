<script setup lang="ts">
    import { h, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NEmpty } from 'naive-ui';
    import { IconTrash } from '@tabler/icons-vue';

    import { ProjectAttachment } from '../models/project-attachment.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import RefreshAddActionsColumn from '../../../shared/components/tables/RefreshAddActionsColumn.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import { formatBytes } from '../../../shared/composables/format.ts';
    import type { SortOrder } from '../../../shared/types/common.ts';

    interface Props {
        // TODO: refactor to disabled
        loading: boolean;
        projectId: string;
        projectAttachments: ProjectAttachment[];
        errorMessage?: string | null;
    }

    const { t } = useI18n();
    const dialog = useDialog();

    const emit = defineEmits(['refresh', 'add', 'delete', 'download', 'preview']);

    const props = defineProps<Props>();

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("modules.projectAttachment.components.projectAttachmentsTable.header.columns.name"),
            field: "name",
            sortable: true,
        },
        {
            label: t("modules.projectAttachment.components.projectAttachmentsTable.header.columns.size"),
            field: "size",
            sortable: true,
        },
        {
            label: t("modules.projectAttachment.components.projectAttachmentsTable.header.columns.contentType"),
            field: "contentType",
            sortable: true,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdAt"),
            field: "createdAt",
            sortable: true,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdBy"),
            field: "createdBy",
            sortable: true,
        },
    ]);

    const nameFilter = defineModel<string>("nameFilter", {
        default: "",
    });
    const userFilter = defineModel<string>("userFilter", {
        default: "",
    });

    const onRefresh = () => {
        emit("refresh");
    };

    const onAdd = () => {
        emit("add");
    };

    const onConfirmDelete = (projectAttachment: ProjectAttachment, index: number) => {
        dialog.warning({
            title: t("modules.projectAttachment.components.projectAttachmentsTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    t("modules.projectAttachment.components.projectAttachmentsTable.dialogs.deleteConfirmation.message", { name: projectAttachment.name, size: formatBytes(projectAttachment.size) }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                emit("delete", projectAttachment, index)
            },
        });
    };

    const onDownload = (projectAttachment: ProjectAttachment, index: number) => {
        emit("download", projectAttachment, index)
    };

    const onPreview = (projectAttachment: ProjectAttachment, index: number) => {
        emit("preview", projectAttachment, index)
    };
</script>

<template>
    <ManageTable size="small" :columns="columns" :sort-field="sortField" :sort-order="sortOrder" @sort="onToggleSort"
        @refresh="onRefresh" @add="onAdd">
        <template #thead>
            <tr>
                <th v-for="column in columns" :key="column.field"
                    :class="{ 'doneo-text-center': column.align === 'center' }">
                    {{ column.label }}
                </th>
                <th class="doneo-table-actions-column">{{ t("shared.components.table.header.columns.actions") }}</th>
            </tr>
            <tr>
                <th>
                    <TextFilterInput clearable size="small"
                        :placeholder="t('modules.projectAttachment.components.projectAttachmentsTable.filters.name.placeholder')"
                        v-model:value="nameFilter" />
                </th>
                <th></th>
                <th></th>
                <th></th>
                <th>
                    <UserSelector hideAvatar clearable v-model:id="userFilter"
                        :placeholder="t('modules.projectAttachment.components.projectAttachmentsTable.filters.user.placeholder')" />
                </th>
                <th class="doneo-text-center">
                    <RefreshAddActionsColumn @refresh="onRefresh" @add="onAdd" />
                </th>
            </tr>
        </template>
        <template #tbody v-if="!props.errorMessage">
            <tr v-for="projectAttachment, index in projectAttachments" :key="projectAttachment.id ?? index">
                <td>{{ projectAttachment.name }}</td>
                <td>{{ formatBytes(projectAttachment.size) }}</td>
                <td>{{ projectAttachment.contentType }}</td>
                <td>{{ projectAttachment.createdAt.toLocaleString() }}</td>
                <td>
                    <AvatarUserName :user-id="projectAttachment.createdBy?.id ?? ''"
                        :user-name="projectAttachment.createdBy?.name ?? ''" />
                </td>
                <td class="doneo-text-center">
                    <ManageTableActionButtons show-delete show-download
                        @delete="onConfirmDelete(projectAttachment, index)"
                        @download="onDownload(projectAttachment, index)" @preview="onPreview(projectAttachment, index)"
                        :disabled="props.loading" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="projectAttachments.length < 1 && !props.loading">
                    <n-empty
                        :description="t('modules.projectAttachment.components.projectAttachmentsTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>