<script setup lang="ts">
    import { ref, reactive, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NEmpty, NTag, NButtonGroup, NButton, NIcon } from 'naive-ui';
    import { IconFilePencil } from '@tabler/icons-vue';

    import { useUserSettingsStore } from '../../../stores/userSettings.ts';
    import type { Sort } from '../../../shared/types/models/sort.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { ProjectsTableFilters } from '../types/projects-table-filters.ts';
    import type { DateFilterSelectComponent } from '../../users/components/date-filter-select-component.ts';
    import { Project } from '../models/project';

    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/form-blocks/TextFilterInput.vue';
    import ProjectPrioritySelector from '../../project-priorities/components/ProjectPrioritySelector.vue';
    import ProjectTypeSelector from '../../project-types/components/ProjectTypeSelector.vue';
    import ProjectStatusSelector from '../../project-statuses/components/ProjectStatusSelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import ClearTableFiltersButton from '../../../shared/components/buttons/ClearTableFiltersButton.vue';
    import { getNaiveUITagColorProperty } from '../../../shared/composables/color';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import ChangeProjectStatusDropdown from '../../../shared/components/dropdowns/ChangeProjectStatusDropdown.vue';
    import type { ProjectStatus } from '../../project-statuses/models/project-status.ts';
    import { DEFAULT_BUTTON_SIZE } from '../../../constants.ts';
    import ProjectResumeFloatingCard from './ProjectResumeFloatingCard.vue';

    interface Props {
        disabled: boolean;
        readOnly?: boolean;
        items: Project[];
        sort?: Sort;
    }

    const { t } = useI18n();
    const userSettingsStore = useUserSettingsStore();
    // TODO: dialog for delete ?

    const emit = defineEmits(['refresh', 'add', 'sort', 'statusChanged']);

    const props = defineProps<Props>();

    const createdAtFilterRef = ref<DateFilterSelectComponent | undefined>();

    const filters = defineModel<ProjectsTableFilters>("filters", {
        default: () => ({
            slug: "",
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

    const isFilteredBySlug = computed<boolean>(() => filters.value.slug.length > 0);
    const isFilteredByType = computed<boolean>(() => filters.value.typeId !== null);
    const isFilteredByPriority = computed<boolean>(() => filters.value.priorityId !== null);
    const isFilteredByStatus = computed<boolean>(() => filters.value.statusId !== null);
    const isFilteredBySummary = computed<boolean>(() => filters.value.summary.length > 0);
    const isFilteredByCreationDate = computed<boolean>(() => filters.value.createdAt.from != null || filters.value.createdAt.to != null);
    const isFilteredByCreator = computed<boolean>(() => filters.value.createdByUserId !== null);

    const hasFilters = computed<boolean>(() =>
        isFilteredBySlug.value ||
        isFilteredByType.value ||
        isFilteredByPriority.value ||
        isFilteredByStatus.value ||
        isFilteredBySummary.value ||
        isFilteredByCreationDate.value ||
        isFilteredByCreator.value
    );

    const columnDefinitions = reactive<Record<string, TableHeaderColumn>>({
        slug: {
            label: t("modules.project.components.ProjectsTable.header.columns.slug"),
            field: "slug",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredBySlug.value,
        },
        type: {
            label: t("modules.project.components.ProjectsTable.header.columns.type"),
            field: "type",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByType.value,
        },
        priority: {
            label: t("modules.project.components.ProjectsTable.header.columns.priority"),
            field: "priority",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByPriority.value,
        },
        status: {
            label: t("modules.project.components.ProjectsTable.header.columns.status"),
            field: "status",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByStatus.value,
        },
        summary: {
            label: t("modules.project.components.ProjectsTable.header.columns.summary"),
            field: "summary",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredBySummary.value,
        },
        createdAt: {
            label: t("modules.project.components.ProjectsTable.header.columns.createdAt"),
            field: "createdAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByCreationDate.value,
        },
        createdBy: {
            label: t("modules.project.components.ProjectsTable.header.columns.createdBy"),
            field: "createdBy",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByCreator.value,
        },
    });

    const columns = computed<TableHeaderColumn[]>(() => Object.values(columnDefinitions))

    const onSort = (sort: Sort) => {
        emit("sort", sort);
    };

    const onRefresh = () => {
        emit("refresh");
    };

    const onAdd = () => {
        emit("add");
    };

    const onStatusChange = (project: Project, status: ProjectStatus) => {
        emit("statusChanged", project, status);
    }

    const onClearFilters = () => {
        filters.value.slug = "";
        filters.value.typeId = null;
        filters.value.priorityId = null;
        filters.value.statusId = null;
        filters.value.summary = "";
        createdAtFilterRef.value?.reset();
        filters.value.createdByUserId = null;
    };

    const showDrawer = ref<boolean>(false);

    const currentProject = ref<Project | null>(null);

    const onShowProjectResume = (project: Project) => {
        showDrawer.value = true;
        currentProject.value = project;
    };

    const onShowColumn = (column: TableHeaderColumn) => {
        columnDefinitions[column.field].visible = true;
    };

    const onHideColumn = (column: TableHeaderColumn) => {
        columnDefinitions[column.field].visible = false;
    };
</script>

<template>
    <ProjectResumeFloatingCard v-if="showDrawer && currentProject?.id" v-model:show="showDrawer"
        :project-id="currentProject?.id" />
    <ManageTable size="small" :columns="columns" :current-sort="sort" @sort="onSort" @refresh="onRefresh" @add="onAdd"
        @show-column="onShowColumn" @hide-column="onHideColumn">
        <template #thead>
            <tr>
                <th v-if="columnDefinitions.slug.visible">
                    <TextFilterInput clearable :disabled="props.disabled" size="small"
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.slug.placeholder')"
                        v-model:value="filters.slug" @keydown-enter="onRefresh" />
                </th>
                <th v-if="columnDefinitions.type.visible">
                    <ProjectTypeSelector :disabled="props.disabled" size="small" v-model:id="filters.typeId"
                        :hide-prefix="true" clearable
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.type.placeholder')" />
                </th>
                <th v-if="columnDefinitions.priority.visible">
                    <ProjectPrioritySelector :disabled="props.disabled" size="small" v-model:id="filters.priorityId"
                        :hide-prefix="true" clearable
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.priority.placeholder')" />
                </th>
                <th v-if="columnDefinitions.status.visible">
                    <ProjectStatusSelector :disabled="props.disabled" size="small" v-model:id="filters.statusId"
                        :hide-prefix="true" clearable
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.status.placeholder')" />
                </th>
                <th v-if="columnDefinitions.summary.visible">
                    <TextFilterInput clearable :disabled="props.disabled" size="small"
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.summary.placeholder')"
                        v-model:value="filters.summary" @keydown-enter="onRefresh" />
                </th>
                <th v-if="columnDefinitions.createdAt.visible">
                    <DateFilterSelect clearable ref="createdAtFilterRef" size="small" :disabled="props.disabled"
                        v-model:range="filters.createdAt" />
                </th>
                <th v-if="columnDefinitions.createdBy.visible">
                    <UserSelector hideAvatar clearable :disabled="props.disabled" size="small"
                        v-model:id="filters.createdByUserId"
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.creator.placeholder')" />
                </th>
                <th class="doneo-text-center">
                    <ClearTableFiltersButton @clear="onClearFilters" :disabled="props.disabled || !hasFilters" />
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="project, index in items" :key="project.id ?? index"
                :class="{ 'tr-archived-project': project.archivedAt.hasValue() }">
                <td v-if="columnDefinitions.slug.visible">
                    {{ project.slug }}
                </td>
                <td v-if="columnDefinitions.type.visible"><n-tag :bordered="false"
                        :color="getNaiveUITagColorProperty(project.type.hexColor ?? '#888888')">{{
                            project.type.name }}</n-tag>
                </td>
                <td v-if="columnDefinitions.priority.visible"><n-tag :bordered="false"
                        :color="getNaiveUITagColorProperty(project.priority.hexColor ?? '#888888')">{{
                            project.priority.name
                        }}</n-tag></td>
                <td v-if="columnDefinitions.status.visible"><n-tag :bordered="false"
                        :color="getNaiveUITagColorProperty(project.status.hexColor ?? '#888888')">{{
                            project.status.name }}</n-tag></td>
                <td v-if="columnDefinitions.summary.visible">
                    <router-link :to="{ name: 'projectTab', params: { projectId: project.id, tab: 'metadata' } }"
                        class="doneo-link-text-color-default">{{
                            project.summary
                        }}</router-link>
                </td>
                <td v-if="columnDefinitions.createdAt.visible">{{
                    project.createdAt.toCustomMaskString(userSettingsStore.currentDatetimeMask) }}</td>
                <td v-if="columnDefinitions.createdBy.visible">
                    <AvatarUserName :user-id="project.createdBy.id" :user-name="project.createdBy.name" />
                </td>
                <td class="doneo-text-center">
                    <!-- TODO: use ManageTableActionButtons -->
                    <n-button-group class="doneo-table-actions-button-group" :size="DEFAULT_BUTTON_SIZE">
                        <n-button class="doneo-table-actions-button" :disabled="props.disabled"
                            :size="DEFAULT_BUTTON_SIZE" @click="onShowProjectResume(project)">
                            {{ t("shared.buttons.Open.label") }}
                            <template #icon>
                                <n-icon :size="22" :component="IconFilePencil" />
                            </template>
                        </n-button>
                        <ChangeProjectStatusDropdown className="doneo-table-actions-button" :disabled="props.disabled"
                            :read-only="props.readOnly" :current-status="project.status"
                            @change="(status: ProjectStatus) => onStatusChange(project, status)" />
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

<style lang="css" scoped>
    .tr-archived-project td {
        opacity: 0.5;
    }
</style>