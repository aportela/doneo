<script setup lang="ts">
    import { h, ref, computed, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NSelect, NEmpty, type SelectOption } from 'naive-ui';
    import { IconTrash } from '@tabler/icons-vue';

    import { useUserSettingsStore } from '../../../stores/userSettings.ts';
    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { AttachmentsTableFilters } from '../types/attachments-table-filter.ts';
    import { Attachment } from '../models/attachment.ts';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import ClearFiltersTableButton from '../../../shared/components/tables/ClearFiltersTableButton.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';

    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import { formatBytes } from '../../../shared/composables/format.ts';
    import type { DateFilterSelectComponent } from '../../users/components/date-filter-select-component.ts';



    interface Props {
        disabled: boolean;
        readOnly?: boolean;
        items: Attachment[];
        projectId: string;
        errorMessage?: string | null;
    }

    const { t } = useI18n();
    const dialog = useDialog();
    const userSettingsStore = useUserSettingsStore();

    const emit = defineEmits(['refresh', 'add', 'delete', 'download', 'preview']);

    const props = defineProps<Props>();

    const createdAtFilterRef = ref<DateFilterSelectComponent | undefined>();

    const filters = defineModel<AttachmentsTableFilters>("filters", {
        default: () => ({
            name: "",
            createdByUserId: null,
            createdAt: {
                from: null,
                to: null,
            },
            contentType: null,
        })
    });


    const isFilteredByName = computed<boolean>(() => filters.value.name.length > 0);
    const isFilteredByCreator = computed<boolean>(() => filters.value.createdByUserId !== null);
    const isFilteredByCreatedAt = computed<boolean>(() => filters.value.createdAt.from != null || filters.value.createdAt.to != null);
    const isFilteredByContentType = computed<boolean>(() => filters.value.contentType !== null);

    const contentTypeOptions = ref<SelectOption[]>([]);

    watch(() => props.items, (newValue: Attachment[], oldValue: Attachment[]) => {
        if (oldValue.length === 0) {
            contentTypeOptions.value = [...new Set(newValue.map((item: Attachment) => { return (item.contentType) }))].map((contentType) => { return ({ label: contentType, value: contentType }); });
        }
    });


    const hasFilters = computed<boolean>(() =>
        isFilteredByName.value ||
        isFilteredByCreator.value ||
        isFilteredByCreatedAt.value ||
        isFilteredByContentType.value
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
            isFiltered: () => isFilteredByContentType.value,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdAt"),
            field: "createdAt",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByCreatedAt.value,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdBy"),
            field: "createdBy",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredByCreator.value,
        },
    ]);


    const onRefresh = () => {
        emit("refresh");
    };

    const onAdd = () => {
        emit("add");
    };

    const onConfirmDelete = (attachment: Attachment, index: number) => {
        dialog.warning({
            title: t("modules.projectAttachment.components.projectAttachmentsTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    t("modules.projectAttachment.components.projectAttachmentsTable.dialogs.deleteConfirmation.message", { name: attachment.name, size: formatBytes(attachment.size) }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                emit("delete", attachment, index)
            },
        });
    };

    const onDownload = (attachment: Attachment, index: number) => {
        emit("download", attachment, index)
    };

    const onPreview = (attachment: Attachment, index: number) => {
        emit("preview", attachment, index)
    };

    const onClearFilters = () => {
        filters.value.name = "";
        filters.value.createdByUserId = null;
        createdAtFilterRef.value?.reset();
        filters.value.contentType = null;
    };
</script>

<template>
    <ManageTable size="small" :columns="columns" @refresh="onRefresh" @add="onAdd" :hide-add="props.readOnly">
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
                    <n-select size="small" clearable :disabled="props.disabled" :options="contentTypeOptions"
                        v-model:value="filters.contentType"
                        :placeholder="t('modules.projectAttachment.components.projectAttachmentsTable.filters.contentType.placeholder')" />
                </th>
                <th>
                    <DateFilterSelect clearable v-model:range="filters.createdAt" ref="createdAtFilterRef"
                        :disabled="props.disabled" />
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
            <tr v-for="attachment, index in items" :key="attachment.id ?? index">
                <td>{{ attachment.name }}</td>
                <td>{{ formatBytes(attachment.size) }}</td>
                <td>{{ attachment.contentType }}</td>
                <td>{{ attachment.createdAt.toCustomMaskString(userSettingsStore.currentDatetimeMask) }}</td>
                <td>
                    <AvatarUserName :user-id="attachment.createdBy?.id ?? ''"
                        :user-name="attachment.createdBy?.name ?? ''" />
                </td>
                <td class="doneo-text-center">
                    <ManageTableActionButtons show-delete show-download show-preview
                        @delete="onConfirmDelete(attachment, index)" @download="onDownload(attachment, index)"
                        @preview="onPreview(attachment, index)" :disabled="props.disabled"
                        :delete-disabled="props.disabled" :download-disabled="props.disabled"
                        :preview-disabled="props.disabled || !(attachment.allowImagePreview() || attachment.allowAudioPreview() || attachment.allowPDFPreview())" />
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