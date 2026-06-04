<script setup lang="ts">
    import { h, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NEmpty } from 'naive-ui';
    import { IconTrash } from '@tabler/icons-vue';

    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { ProjectAttachmentsTableFilters } from '../types/project-attachments-table-filter.ts';
    import { ProjectAttachment } from '../models/project-attachment.ts';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import ClearFiltersTableButton from '../../../shared/components/tables/ClearFiltersTableButton.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import { formatBytes } from '../../../shared/composables/format.ts';

    interface Props {
        disabled: boolean;
        items: ProjectAttachment[];
        projectId: string;
        errorMessage?: string | null;
    }

    const { t } = useI18n();
    const dialog = useDialog();

    const emit = defineEmits(['refresh', 'add', 'delete', 'download', 'preview']);

    const props = defineProps<Props>();

    const filters = defineModel<ProjectAttachmentsTableFilters>("filters", {
        default: () => ({
            name: "",
            createdByUserId: null,
        })
    });

    const isFilteredByName = computed<boolean>(() => filters.value.name.length > 0);
    const isFilteredByCreator = computed<boolean>(() => filters.value.createdByUserId !== null);

    const hasFilters = computed<boolean>(() =>
        isFilteredByName.value ||
        isFilteredByCreator.value
    );

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("modules.projectAttachment.components.projectAttachmentsTable.header.columns.name"),
            field: "name",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByName.value,
        },
        {
            label: t("modules.projectAttachment.components.projectAttachmentsTable.header.columns.size"),
            field: "size",
            visible: true,
            sortable: false,
            isFiltered: () => false,
        },
        {
            label: t("modules.projectAttachment.components.projectAttachmentsTable.header.columns.contentType"),
            field: "contentType",
            visible: true,
            sortable: false,
            isFiltered: () => false,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdAt"),
            field: "createdAt",
            visible: true,
            sortable: false,
            isFiltered: () => false,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdBy"),
            field: "createdBy",
            visible: true,
            sortable: false,
            isFiltered: () => false,
        },
    ]);

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

    const onClearFilters = () => {
        filters.value.name = "";
        filters.value.createdByUserId = null;
    };
</script>

<template>
    <ManageTable size="small" :columns="columns" @refresh="onRefresh" @add="onAdd">
        <template #thead>
            <tr>
                <th>
                    <TextFilterInput clearable :disabled="props.disabled" size="small"
                        :placeholder="t('modules.projectAttachment.components.projectAttachmentsTable.filters.name.placeholder')"
                        v-model:value="filters.name" />
                </th>
                <th>
                    <!-- TODO: size filter -->
                </th>
                <th>
                    <!-- TODO: content type filter -->
                </th>
                <th>
                    <!-- TODO: created at filter -->
                </th>
                <th>
                    <UserSelector hideAvatar clearable :disabled="props.disabled" v-model:id="filters.createdByUserId"
                        :placeholder="t('modules.projectAttachment.components.projectAttachmentsTable.filters.user.placeholder')" />
                </th>
                <th class="doneo-text-center">
                    <ClearFiltersTableButton @clear="onClearFilters" :disabled="props.disabled || !hasFilters" />
                </th>
            </tr>
        </template>
        <template #tbody v-if="!props.errorMessage">
            <tr v-for="projectAttachment, index in items" :key="projectAttachment.id ?? index">
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
                        :disabled="props.disabled" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="items.length < 1 && !props.disabled">
                    <n-empty
                        :description="t('modules.projectAttachment.components.projectAttachmentsTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>