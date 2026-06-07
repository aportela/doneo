<script setup lang="ts">
    import { ref, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NEmpty, NTag, NButtonGroup, NButton, NIcon } from 'naive-ui';
    import { IconFilePencil } from '@tabler/icons-vue';

    import type { Sort } from '../../../shared/types/models/sort.ts';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import type { DateFilterSelectComponent } from '../../users/components/date-filter-select-component.ts';

    import type { ProjectTasksTableFilters } from '../types/project-tasks-table-filters.ts';
    import { ProjectTask } from '../models/tasks.ts';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import ProjectPrioritySelector from '../../project-priorities/components/ProjectPrioritySelector.vue';
    import ProjectStatusSelector from '../../project-statuses/components/ProjectStatusSelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import ClearFiltersTableButton from '../../../shared/components/tables/ClearFiltersTableButton.vue';
    import { getNaiveUITagColorProperty } from '../../../shared/composables/color';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';

    interface Props {
        disabled: boolean;
        items: ProjectTask[];
        sort?: Sort;
    }

    const { t } = useI18n();
    // TODO: dialog for delete ?

    const emit = defineEmits(['refresh', 'add', 'sort']);

    const props = defineProps<Props>();

    const createdAtFilterRef = ref<DateFilterSelectComponent | undefined>();

    const filters = defineModel<ProjectTasksTableFilters>("filters", {
        default: () => ({
            slug: null,
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

    const isFilteredBySlug = computed<boolean>(() => filters.value.slug !== null);
    const isFilteredByPriority = computed<boolean>(() => filters.value.priorityId !== null);
    const isFilteredByStatus = computed<boolean>(() => filters.value.statusId !== null);
    const isFilteredBySummary = computed<boolean>(() => filters.value.summary.length > 0);
    const isFilteredByCreationDate = computed<boolean>(() => filters.value.createdAt.from != null || filters.value.createdAt.to != null);
    const isFilteredByCreator = computed<boolean>(() => filters.value.createdByUserId !== null);

    const hasFilters = computed<boolean>(() =>
        isFilteredBySlug.value ||
        isFilteredByPriority.value ||
        isFilteredByStatus.value ||
        isFilteredBySummary.value ||
        isFilteredByCreationDate.value ||
        isFilteredByCreator.value
    );

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: "Slug",
            field: "slug",
            visible: true,
            sortable: false,
            isFiltered: () => false,
        },
        {
            label: "Priority",
            field: "priority",
            visible: true,
            sortable: false,
            isFiltered: () => false,
        },
        {
            label: "Status",
            field: "status",
            visible: true,
            sortable: false,
            isFiltered: () => false,
        },
        {
            label: "Summary",
            field: "summary",
            visible: true,
            sortable: false,
            isFiltered: () => false,
        },
        {
            label: "Created at",
            field: "createdAt",
            visible: true,
            sortable: false,
            isFiltered: () => false,
        },
        {
            label: "Created by",
            field: "createdBy",
            visible: true,
            sortable: false,
            isFiltered: () => false,
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
                        :placeholder="t('modules.task.components.ProjectTasksTable.header.filters.slug.placeholder')"
                        v-model:value="filters.slug" @keydown-enter="onRefresh" />
                </th>
                <th>
                    <ProjectPrioritySelector :disabled="props.disabled" v-model:id="filters.priorityId"
                        :hide-prefix="true" clearable
                        :placeholder="t('modules.task.components.ProjectTasksTable.header.filters.priority.placeholder')" />
                </th>
                <th>
                    <ProjectStatusSelector :disabled="props.disabled" v-model:id="filters.statusId" :hide-prefix="true"
                        clearable
                        :placeholder="t('modules.task.components.ProjectTasksTable.header.filters.status.placeholder')" />
                </th>
                <th>
                    <TextFilterInput clearable :disabled="props.disabled" size="small"
                        :placeholder="t('modules.task.components.ProjectTasksTable.header.filters.summary.placeholder')"
                        v-model:value="filters.summary" @keydown-enter="onRefresh" />
                </th>
                <th>
                    <DateFilterSelect clearable ref="createdAtFilterRef" :disabled="props.disabled"
                        v-model:range="filters.createdAt" />
                </th>
                <th>
                    <UserSelector hideAvatar clearable :disabled="props.disabled" v-model:id="filters.createdByUserId"
                        :placeholder="t('modules.task.components.ProjectTasksTable.header.filters.creator.placeholder')" />
                </th>
                <th class="doneo-text-center">
                    <ClearFiltersTableButton @clear="onClearFilters" :disabled="props.disabled || !hasFilters" />
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="task, index in items" :key="task.id ?? index">
                <td>
                    {{ task.slug }}
                </td>
                <td><n-tag :bordered="false" :color="getNaiveUITagColorProperty(task.priority.hexColor ?? '#888888')">{{
                    task.priority.name
                        }}</n-tag></td>
                <td><n-tag :bordered="false" :color="getNaiveUITagColorProperty(task.status.hexColor ?? '#888888')">{{
                    task.status.name }}</n-tag></td>
                <td><router-link :to="{ name: 'projectTab', params: { id: task.id, tab: 'metadata' } }"
                        class="doneo-link-text-color-default">{{
                            task.summary
                        }}</router-link></td>
                <td>{{ task.createdAt.toLocaleString() }}</td>
                <td>
                    <AvatarUserName :user-id="task.createdBy.id" :user-name="task.createdBy.name" />
                </td>
                <td class="doneo-text-center">
                    <!-- TODO: use ManageTableActionButtons -->
                    <n-button-group size="small">
                        <router-link :to="{ name: 'projectTab', params: { id: task.id, tab: 'metadata' } }">
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
                    <n-empty :description="t('modules.task.components.ProjectTasksTable.warnings.noItemsFound')">
                    </n-empty>
                </td>
            </tr>
        </template>
    </ManageTable>
</template>

<style lang="css" scoped></style>