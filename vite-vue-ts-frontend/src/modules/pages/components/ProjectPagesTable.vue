<script setup lang="ts">
    import { ref, computed } from 'vue';
    import { useRouter } from 'vue-router';
    import { useI18n } from "vue-i18n";

    import { NEmpty } from 'naive-ui';

    import { useUserSettingsStore } from '../../../stores/userSettings.ts';
    import { Page } from '../models/page.ts';

    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { ProjectPagesTableFilters } from '../types/project-pages-table-filters.ts';

    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import ClearTableFiltersButton from '../../../shared/components/buttons/ClearTableFiltersButton.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import type { DateFilterSelectComponent } from '../../users/components/date-filter-select-component.ts';

    interface IProps {
        disabled: boolean;
        items: Page[];
        projectId: string;
        errorMessage?: string | null;
    }

    const { t } = useI18n();
    const router = useRouter();

    const userSettingsStore = useUserSettingsStore();

    const emit = defineEmits(['refresh']);

    const props = defineProps<IProps>();

    const createdAtFilterRef = ref<DateFilterSelectComponent | undefined>();
    const updatedAtFilterRef = ref<DateFilterSelectComponent | undefined>();

    const filters = defineModel<ProjectPagesTableFilters>("filters", {
        default: () => ({
            title: "",
            userId: "",
            createdAt: {
                from: null,
                to: null,
            },
            updatedAt: {
                from: null,
                to: null,
            },
        })
    });

    const isFilteredByTitle = computed<boolean>(() => filters.value.title !== null);
    const isFilteredByUser = computed<boolean>(() => filters.value.userId !== null);
    const isFilteredByCreatedAt = computed<boolean>(() => filters.value.createdAt.from != null || filters.value.createdAt.to != null);
    const isFilteredByUpdatedAt = computed<boolean>(() => filters.value.updatedAt.from != null || filters.value.updatedAt.to != null);

    const hasFilters = computed<boolean>(() =>
        isFilteredByTitle.value || isFilteredByUser.value || isFilteredByCreatedAt.value || isFilteredByUpdatedAt.value
    );

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("modules.projectPage.components.ProjectPagesTable.header.columns.title"),
            field: "title",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByUser.value,
        },
        {
            label: t("modules.projectPage.components.ProjectPagesTable.header.columns.createdAt"),
            field: "createdAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByCreatedAt.value,
        },
        {
            label: t("modules.projectPage.components.ProjectPagesTable.header.columns.createdBy"),
            field: "createdBy",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByUser.value,
        },
        {
            label: t("modules.projectPage.components.ProjectPagesTable.header.columns.updatedAt"),
            field: "updatedAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByUpdatedAt.value,
        },
    ]);

    const onRefresh = () => {
        emit("refresh");
    };

    const onClearFilters = () => {
        filters.value.title = null;
        filters.value.userId = null;
        createdAtFilterRef.value?.reset();
        updatedAtFilterRef.value?.reset();
    };

    const onOpen = (projectPageId: string | null) => {
        router.push(
            { name: "projectPage", params: { projectId: props.projectId, pageId: projectPageId } }
        ).catch((e) => {
            console.error(e);
        });
    };

    const onAdd = () => {

    };
</script>

<template>
    <ManageTable id="ProjectPages" size="small" :columns="columns" @refresh="onRefresh" @add="onAdd">
        <template #thead>
            <tr>
                <th>
                    <TextFilterInput clearable :disabled="props.disabled" size="small"
                        :placeholder="t('modules.projectPage.components.ProjectPagesTable.filters.title.placeholder')"
                        v-model:value="filters.title" />
                </th>
                <th>
                    <DateFilterSelect clearable v-model:range="filters.createdAt" ref="createdAtFilterRef"
                        :disabled="props.disabled" />
                </th>
                <th>
                    <UserSelector v-model:id="filters.userId" :disabled="props.disabled" size="small" hide-avatar
                        clearable
                        :placeholder="t('modules.projectPage.components.ProjectPagesTable.filters.user.placeholder')" />
                </th>
                <th>
                    <DateFilterSelect clearable v-model:range="filters.updatedAt" ref="updatedAtFilterRef"
                        :disabled="props.disabled" />
                </th>
                <th class="doneo-text-center">
                    <ClearTableFiltersButton @clear="onClearFilters" :disabled="props.disabled || !hasFilters" />
                </th>
            </tr>
        </template>
        <template #tbody v-if="!props.errorMessage">
            <tr v-for="projectPage, index in items" :key="projectPage.id ?? index">
                <td>
                    {{ projectPage.title }}
                </td>
                <td>{{ projectPage.createdAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) }}
                </td>
                <td>
                    <AvatarUserName :user-id="projectPage.createdBy?.id ?? ''"
                        :user-name="projectPage.createdBy?.name ?? ''" />
                </td>
                <td>{{ projectPage.updatedAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) }}
                </td>
                <td class="doneo-text-center">
                    <ManageTableActionButtons show-open :disabled="props.disabled" @open="onOpen(projectPage.id)" />
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="items.length < 1 && !props.disabled">
                    <n-empty
                        :description="t('modules.projectHistoryOperation.components.ProjectHistoryOperationsTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>