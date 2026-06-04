<script setup lang="ts">
    import { ref, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NEmpty, NTag, NButtonGroup, NButton, NIcon } from 'naive-ui';
    import { IconFilePencil } from '@tabler/icons-vue';

    import type { Sort } from '../../../shared/types/models/sort.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { ProjectsTableFilters } from '../types/projects-table-filters.ts';
    import type { DateFilterSelectComponent } from '../../users/components/date-filter-select-component.ts';
    import { Project } from '../models/project';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import ProjectPrioritySelector from '../../project-priorities/components/ProjectPrioritySelector.vue';
    import ProjectTypeSelector from '../../project-types/components/ProjectTypeSelector.vue';
    import ProjectStatusSelector from '../../project-statuses/components/ProjectStatusSelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import ClearFiltersTableButton from '../../../shared/components/tables/ClearFiltersTableButton.vue';
    import { getNaiveUITagColorProperty } from '../../../shared/composables/color';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';

    interface Props {
        disabled: boolean;
        items: Project[];
        sort?: Sort;
    }

    const { t } = useI18n();
    // TODO: dialog for delete ?

    const emit = defineEmits(['refresh', 'add', 'sort']);

    const props = defineProps<Props>();

    const createdAtFilterRef = ref<DateFilterSelectComponent | undefined>();

    const filters = defineModel<ProjectsTableFilters>("filters", {
        default: () => ({
            key: "",
            typeId: null,
            priorityId: null,
            statusId: null,
            summary: "",
            createdAt: {
                from: null,
                to: null,
            },
            createdByUserId: "",
        })
    });

    const isFilteredByKey = computed<boolean>(() => filters.value.key.length > 0);
    const isFilteredByType = computed<boolean>(() => filters.value.typeId !== null);
    const isFilteredByPriority = computed<boolean>(() => filters.value.priorityId !== null);
    const isFilteredByStatus = computed<boolean>(() => filters.value.statusId !== null);
    const isFilteredBySummary = computed<boolean>(() => filters.value.summary.length > 0);
    const isFilteredByCreationDate = computed<boolean>(() => filters.value.createdAt.from != null || filters.value.createdAt.to != null);
    const isFilteredByCreator = computed<boolean>(() => filters.value.createdByUserId !== null);

    const hasFilters = computed<boolean>(() =>
        isFilteredByKey.value ||
        isFilteredByType.value ||
        isFilteredByPriority.value ||
        isFilteredByStatus.value ||
        isFilteredBySummary.value ||
        isFilteredByCreationDate.value ||
        isFilteredByCreator.value
    );

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("modules.project.components.ProjectsTable.header.columns.key"),
            field: "key",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByKey.value,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.type"),
            field: "type",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByType.value,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.priority"),
            field: "priority",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByPriority.value,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.status"),
            field: "status",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByStatus.value,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.summary"),
            field: "summary",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredBySummary.value,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdAt"),
            field: "createdAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByCreationDate.value,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdBy"),
            field: "createdBy",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByCreator.value,
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

    const onClearFilters = () => {
        filters.value.key = "";
        filters.value.typeId = null;
        filters.value.priorityId = null;
        filters.value.statusId = null;
        filters.value.summary = "";
        createdAtFilterRef.value?.reset();
        filters.value.createdByUserId = null;
    };

</script>

<template>
    <ManageTable size="small" :columns="columns" :current-sort="sort" @sort="onSort" @refresh="onRefresh" @add="onAdd">
        <template #thead>
            <tr>
                <th>
                    <TextFilterInput clearable :disabled="props.disabled" size="small"
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.key.placeholder')"
                        v-model:value="filters.key" @keydown-enter="onRefresh" />
                </th>
                <th>
                    <ProjectTypeSelector :disabled="props.disabled" v-model:id="filters.typeId" :hide-prefix="true"
                        clearable
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.type.placeholder')" />
                </th>
                <th>
                    <ProjectPrioritySelector :disabled="props.disabled" v-model:id="filters.priorityId"
                        :hide-prefix="true" clearable
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.priority.placeholder')" />
                </th>
                <th>
                    <ProjectStatusSelector :disabled="props.disabled" v-model:id="filters.statusId" :hide-prefix="true"
                        clearable
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.status.placeholder')" />
                </th>
                <th>
                    <TextFilterInput clearable :disabled="props.disabled" size="small"
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.summary.placeholder')"
                        v-model:value="filters.summary" @keydown-enter="onRefresh" />
                </th>
                <th>
                    <DateFilterSelect :disabled="props.disabled" v-model:range="filters.createdAt" />
                </th>
                <th>
                    <UserSelector hideAvatar clearable :disabled="props.disabled" v-model:id="filters.createdByUserId"
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.creator.placeholder')" />
                </th>
                <th class="doneo-text-center">
                    <ClearFiltersTableButton @clear="onClearFilters" :disabled="props.disabled || !hasFilters" />
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="project, index in items" :key="project.id ?? index">
                <td>
                    {{ project.key }}
                </td>
                <td><n-tag :bordered="false" :color="getNaiveUITagColorProperty(project.type.hexColor ?? '#888888')">{{
                    project.type.name }}</n-tag>
                </td>
                <td><n-tag :bordered="false"
                        :color="getNaiveUITagColorProperty(project.priority.hexColor ?? '#888888')">{{
                            project.priority.name
                        }}</n-tag></td>
                <td><n-tag :bordered="false"
                        :color="getNaiveUITagColorProperty(project.status.hexColor ?? '#888888')">{{
                            project.status.name }}</n-tag></td>
                <td><router-link :to="{ name: 'projectTab', params: { id: project.id, tab: 'metadata' } }"
                        class="doneo-link-text-color-default">{{
                            project.summary
                        }}</router-link></td>
                <td>{{ project.createdAt.toLocaleString() }}</td>
                <td>
                    <AvatarUserName :user-id="project.createdBy.id" :user-name="project.createdBy.name" />
                </td>
                <td class="doneo-text-center">
                    <!-- TODO: use ManageTableActionButtons -->
                    <n-button-group size="small">
                        <router-link :to="{ name: 'projectTab', params: { id: project.id, tab: 'metadata' } }">
                            <n-button :disabled="props.disabled">
                                {{ t("shared.buttons.Open.label") }}
                                <template #icon>
                                    <n-icon :size="22" :component="IconFilePencil" />
                                </template>
                            </n-button>
                        </router-link>
                    </n-button-group>
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="items.length < 1 && !props.disabled">
                    <n-empty :description="t('modules.project.components.ProjectsTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>