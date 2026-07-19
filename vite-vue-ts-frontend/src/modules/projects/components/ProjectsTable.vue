<script setup lang="ts">
    import { ref, reactive, computed, h } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NTag, NButtonGroup, NButton, NIcon } from 'naive-ui';
    import { IconFilePencil } from '@tabler/icons-vue';

    import { useUserSettingsStore } from '../../../stores/userSettings.ts';
    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';
    import type { Order } from '../../../shared/types/order.ts';
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
    import { getNaiveUITagColorProperty } from '../../../shared/composables/color';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import ChangeProjectStatusDropdown from '../../../shared/components/dropdowns/ChangeProjectStatusDropdown.vue';
    import type { ProjectStatus } from '../../project-statuses/models/project-status.ts';
    import { DEFAULT_BUTTON_SIZE } from '../../../constants.ts';
    import ProjectResumeFloatingCard from './ProjectResumeFloatingCard.vue';

    interface Props {
        id: string;
        disabled: boolean;
        readOnly?: boolean;
        items: Project[];
        order: Order;
    }

    const { t } = useI18n();
    const userSettingsStore = useUserSettingsStore();
    const tableSettingsStore = useTableSettingsStore();

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

    const columnDefinitions = reactive<TableHeaderColumn<Project>[]>([
        {
            label: t("modules.project.components.ProjectsTable.header.columns.slug"),
            field: "slug",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredBySlug.value,
            render: (row: Project) => {
                return h("span", {}, { default: () => row.slug });
            }
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.type"),
            field: "type",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByType.value,
            render: (row: Project) => {
                return h(
                    NTag,
                    {
                        bordered: false,
                        color: getNaiveUITagColorProperty(row.type.hexColor ?? "#888888"),
                    },
                    {
                        default: () => row.type.name,
                    }
                );
            }
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.priority"),
            field: "priority",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByPriority.value,
            render: (row: Project) => {
                return h(
                    NTag,
                    {
                        bordered: false,
                        color: getNaiveUITagColorProperty(row.priority.hexColor ?? "#888888"),
                    },
                    {
                        default: () => row.priority.name,
                    }
                );
            }
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.status"),
            field: "status",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByStatus.value,
            render: (row: Project) => {
                return h(
                    NTag,
                    {
                        bordered: false,
                        color: getNaiveUITagColorProperty(row.status.hexColor ?? "#888888"),
                    },
                    {
                        default: () => row.status.name,
                    }
                );
            }
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.summary"),
            field: "summary",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredBySummary.value,
            render: (row: Project) => {
                return h("span", {}, { default: () => row.summary });
            }
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdAt"),
            field: "createdAt",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByCreationDate.value,
            render: (row: Project) => {
                return h("span", {}, { default: () => row.createdAt.toCustomMaskString(userSettingsStore.currentDatetimeMask) });
            }
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.createdBy"),
            field: "createdBy",
            visible: true,
            sortable: true,
            isFiltered: () => isFilteredByCreator.value,
            render: (row: Project) => {
                return h(AvatarUserName, { userId: row.createdBy.id, userName: row.createdBy.name });
            }
        },
    ]);

    const TABLE_ID = "ManageProjects";

    tableSettingsStore.register(
        TABLE_ID,
        {
            columns: columnDefinitions.map((column) => { return { field: column.field, visible: column.visible } }) ?? []
        }
    );

    const tableSettings = tableSettingsStore.get(TABLE_ID);

    const columns = computed<TableHeaderColumn<Project>[]>(() =>
        tableSettings.columns.map((column) => {
            const definition = columnDefinitions.find((c) => c.field === column.field);
            return {
                label: definition!.label,
                field: column.field,
                visible: column.visible,
                sortable: definition!.sortable,
                isFiltered: definition!.isFiltered ?? (() => false),
                render: definition!.render,
            };
        }
        )
    );

    const onSort = (sort: Order) => {
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
        if (createdAtFilterRef.value?.reset) { // TODO: fix
            createdAtFilterRef.value?.reset();
        }
        filters.value.createdByUserId = null;
    };

    const showDrawer = ref<boolean>(false);

    const currentProject = ref<Project | null>(null);

    const onShowProjectResume = (project: Project) => {
        showDrawer.value = true;
        currentProject.value = project;
    };
</script>

<template>
    <ProjectResumeFloatingCard v-if="showDrawer && currentProject?.id" v-model:show="showDrawer"
        :project-id="currentProject?.id" />
    <ManageTable id="ManageProjects" size="small" :rows="items" :row-key="row => row.id ?? ''" :columns="columns"
        :order="order" @sort="onSort" @refresh="onRefresh" @add="onAdd" @clear-filters="onClearFilters"
        :no-items-warning-message="t('modules.project.components.ProjectsTable.warnings.noItemsFound')"
        :show-no-items-warning-message="items.length < 1 && !props.disabled" pagerPosition="top">
        <template #thead-column-filters="{ columns }">
            <th v-for="column in columns">
                <TextFilterInput v-if="column.field === 'slug'" clearable :disabled="props.disabled" size="small"
                    :placeholder="t('modules.project.components.ProjectsTable.header.filters.slug.placeholder')"
                    v-model:value="filters.slug" @keydown-enter="onRefresh" />
                <ProjectTypeSelector v-else-if="column.field === 'type'" :disabled="props.disabled" size="small"
                    v-model:id="filters.typeId" :hide-prefix="true" clearable
                    :placeholder="t('modules.project.components.ProjectsTable.header.filters.type.placeholder')" />
                <ProjectPrioritySelector v-else-if="column.field === 'priority'" :disabled="props.disabled" size="small"
                    v-model:id="filters.priorityId" :hide-prefix="true" clearable
                    :placeholder="t('modules.project.components.ProjectsTable.header.filters.priority.placeholder')" />
                <ProjectStatusSelector v-else-if="column.field === 'status'" :disabled="props.disabled" size="small"
                    v-model:id="filters.statusId" :hide-prefix="true" clearable
                    :placeholder="t('modules.project.components.ProjectsTable.header.filters.status.placeholder')" />
                <TextFilterInput v-else-if="column.field === 'summary'" clearable :disabled="props.disabled"
                    size="small"
                    :placeholder="t('modules.project.components.ProjectsTable.header.filters.summary.placeholder')"
                    v-model:value="filters.summary" @keydown-enter="onRefresh" />
                <DateFilterSelect v-else-if="column.field === 'createdAt'" clearable ref="createdAtFilterRef"
                    size="small" :disabled="props.disabled" v-model:range="filters.createdAt" />
                <UserSelector v-else-if="column.field === 'createdBy'" hideAvatar clearable :disabled="props.disabled"
                    size="small" v-model:id="filters.createdByUserId"
                    :placeholder="t('modules.project.components.ProjectsTable.header.filters.creator.placeholder')" />
            </th>
        </template>
        <template #rowactions="{ row }">
            <!-- TODO: use ManageTableActionButtons -->
            <n-button-group class="doneo-table-actions-button-group" :size="DEFAULT_BUTTON_SIZE">
                <n-button class="doneo-table-actions-button" :disabled="props.disabled" :size="DEFAULT_BUTTON_SIZE"
                    @click="onShowProjectResume(row)">
                    {{ t("shared.buttons.Open.label") }}
                    <template #icon>
                        <n-icon :size="22" :component="IconFilePencil" />
                    </template>
                </n-button>
                <ChangeProjectStatusDropdown className="doneo-table-actions-button" :disabled="props.disabled"
                    :read-only="props.readOnly" :current-status="row.status"
                    @change="(status: ProjectStatus) => onStatusChange(row, status)" />
            </n-button-group>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>