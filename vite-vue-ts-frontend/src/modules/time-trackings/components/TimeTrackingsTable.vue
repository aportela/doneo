<script setup lang="ts">
    import { h, ref, computed, } from 'vue';
    import { useI18n } from "vue-i18n";

    import { useDialog, NEmpty } from 'naive-ui';
    import { IconTrash } from '@tabler/icons-vue';

    import { useUserSettingsStore } from '../../../stores/userSettings.ts';
    import { renderIcon } from '../../../shared/composables/naive-ui-icon';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { TimeTrackingsTableFilters } from '../types/time-trackings-table-filters.ts';
    import type { TimeTracking } from '../models/time-tracking.ts';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import ClearFiltersTableButton from '../../../shared/components/tables/ClearFiltersTableButton.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';

    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import type { DateFilterSelectComponent } from '../../users/components/date-filter-select-component.ts';

    interface TimeTrackingsTableProps {
        disabled: boolean;
        items: TimeTracking[];
        projectId: string;
        taskId: string; // TODO: required ?
        errorMessage?: string | null;
    }

    const { t } = useI18n();
    const dialog = useDialog();
    const userSettingsStore = useUserSettingsStore();

    const emit = defineEmits(['refresh', 'add', 'delete', 'download', 'preview']);

    const props = defineProps<TimeTrackingsTableProps>();

    const createdAtFilterRef = ref<DateFilterSelectComponent | undefined>();

    const filters = defineModel<TimeTrackingsTableFilters>("filters", {
        default: () => ({
            createdByUserId: null,
            createdAt: {
                from: null,
                to: null,
            },
            summary: null,
        })
    });


    const isFilteredByCreator = computed<boolean>(() => filters.value.createdByUserId !== null);
    const isFilteredByCreatedAt = computed<boolean>(() => filters.value.createdAt.from != null || filters.value.createdAt.to != null);
    const isFilteredBySummary = computed<boolean>(() => filters.value.summary.length > 0);



    const hasFilters = computed<boolean>(() =>
        isFilteredByCreator.value ||
        isFilteredByCreatedAt.value ||
        isFilteredBySummary.value
    );

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("modules.timeTracking.components.TimeTrackingsTable.header.columns.summary"),
            field: "name",
            visible: true,
            sortable: false,
            isFiltered: () => isFilteredBySummary.value,
        },
        {
            label: t("modules.timeTracking.components.TimeTrackingsTable.header.columns.spentTime"),
            field: "spentTime",
            visible: true,
            sortable: false,
            isFiltered: () => false,
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

    const onConfirmDelete = (timeTracking: TimeTracking, index: number) => {
        dialog.warning({
            title: t("modules.timeTracking.components.TimeTrackingsTable.dialogs.deleteConfirmation.title"),
            icon: renderIcon(IconTrash)(24),
            content: () =>
                h('div', [
                    t("modules.timeTracking.components.TimeTrackingsTable.dialogs.deleteConfirmation.message", { summary: timeTracking.summary }),
                    h('br'),
                    h('br'),
                    t("shared.components.dialogs.confirmation.continueMessage"),
                ]),
            positiveText: t("shared.buttons.Delete.label"),
            negativeText: t("shared.buttons.Cancel.label"),
            onPositiveClick: () => {
                emit("delete", timeTracking, index)
            },
        });
    };

    const onClearFilters = () => {
        filters.value.createdByUserId = null;
        createdAtFilterRef.value?.reset();
        filters.value.summary = "";

    };
</script>

<template>
    <ManageTable size="small" :columns="columns" @refresh="onRefresh" @add="onAdd">
        <template #thead>
            <tr>
                <th>
                    <TextFilterInput clearable :disabled="props.disabled" size="small"
                        :placeholder="t('modules.timeTracking.components.TimeTrackingsTable.filters.summary.placeholder')"
                        v-model:value="filters.summary" />
                </th>
                <th></th>
                <th>
                    <DateFilterSelect clearable v-model:range="filters.createdAt" ref="createdAtFilterRef"
                        :disabled="props.disabled" />
                </th>
                <th>
                    <UserSelector hideAvatar clearable :disabled="props.disabled" v-model:id="filters.createdByUserId"
                        :placeholder="t('modules.timeTracking.components.TimeTrackingsTable.filters.user.placeholder')" />
                </th>
                <th class="doneo-text-center">
                    <ClearFiltersTableButton @clear="onClearFilters" :disabled="props.disabled || !hasFilters" />
                </th>
            </tr>
        </template>
        <template #tbody v-if="!props.errorMessage">
            <tr v-for="timeTracking, index in items" :key="timeTracking.id ?? index">
                <td>{{ timeTracking.summary }}</td>
                <td>
                    {{
                        timeTracking.geti18nTimeParts()
                            .map(({ key, count }) => `${count} ${t(key, count)}`)
                            .join(", ")
                    }}
                </td>
                <td>{{ timeTracking.createdAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) }}</td>
                <td>
                    <AvatarUserName :user-id="timeTracking.createdBy?.id ?? ''"
                        :user-name="timeTracking.createdBy?.name ?? ''" />
                </td>
                <td class="doneo-text-center">
                    <ManageTableActionButtons show-delete show-download show-preview
                        @delete="onConfirmDelete(timeTracking, index)" :disabled="props.disabled"
                        :delete-disabled="props.disabled" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="items.length < 1 && !props.disabled">
                    <n-empty
                        :description="t('modules.timeTracking.components.TimeTrackingsTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>