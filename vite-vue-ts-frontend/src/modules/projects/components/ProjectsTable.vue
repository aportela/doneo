<script setup lang="ts">
    import { computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NFlex, NEmpty, NTag, NAvatar, NButtonGroup, NButton, NIcon } from 'naive-ui';
    import { IconFilePencil } from '@tabler/icons-vue';

    import { Project } from '../models/project';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import { type SortOrder } from '../../../shared/types/common';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import TableCellHeaderSortIcon from '../../../shared/components/tables/TableCellHeaderSortIcon.vue';
    import RefreshAddActionsColumn from '../../../shared/components/tables/RefreshAddActionsColumn.vue';
    import ProjectPrioritySelector from '../../project-priorities/components/ProjectPrioritySelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import { getNaiveUITagColorProperty } from '../../../shared/composables/color';

    interface Props {
        loading: boolean;
        projects: Project[];
        sortField: string;
        sortOrder: SortOrder;
    }

    const { t } = useI18n();

    const emit = defineEmits(['refresh', 'add', 'toggleSort', 'textfilterKeydownEnter']);

    const props = defineProps<Props>();

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: t("modules.project.components.ProjectsTable.header.columns.key"),
            field: "key",
            sortable: true,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.type"),
            field: "type",
            sortable: true,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.priority"),
            field: "priority",
            sortable: true,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.status"),
            field: "status",
            sortable: true,
        },
        {
            label: t("modules.project.components.ProjectsTable.header.columns.summary"),
            field: "summary",
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

    const projectKeyFilter = defineModel<string>("projectKeyFilter", {
        default: "",
    });

    const projectSummaryFilter = defineModel<string>("projectSummaryFilter", {
        default: "",
    });

    const onToggleSort = (field: string) => {
        emit("toggleSort", field);
    };

    const onRefresh = () => {
        emit("refresh");
    };

    const onAdd = () => {
        emit("add");
    };


    const onTextFilterKeyDownEnter = () => {
        emit("textfilterKeydownEnter");
    };

</script>

<template>
    <ManageTable size="small">
        <template #thead>
            <tr class="table-header-click-action">
                <th v-for="column in columns" :key="column.field" @click="column.sortable && onToggleSort(column.field)"
                    :class="{ 'doneo-cursor-pointer': column.sortable, 'doneo-text-center': column.align === 'center' }">
                    <n-flex justify="space-between" v-if="column.sortable">
                        <span>{{ column.label }}</span>
                        <TableCellHeaderSortIcon v-if="props.sortField === column.field" :order="props.sortOrder" />
                    </n-flex>
                    <span v-else>{{ column.label }}</span>
                </th>
                <th class="doneo-table-actions-column">{{ t("shared.components.table.header.columns.actions") }}</th>
            </tr>
            <tr class="hide-mobile">
                <th>
                    <TextFilterInput clearable size="small"
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.key.placeholder')"
                        v-model:value="projectKeyFilter" @keydown-enter="onTextFilterKeyDownEnter" />
                </th>
                <th>
                    <ProjectPrioritySelector />
                </th>
                <th>
                    <ProjectPrioritySelector />
                </th>
                <th>
                    <ProjectPrioritySelector />
                </th>
                <th>
                    <TextFilterInput clearable size="small"
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.summary.placeholder')"
                        v-model:value="projectSummaryFilter" @keydown-enter="onTextFilterKeyDownEnter" />
                </th>
                <th>
                    <DateFilterSelect />
                </th>
                <th></th>
                <th class="doneo-text-center">
                    <RefreshAddActionsColumn @refresh="onRefresh" @add="onAdd" />
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="project, _index in projects" :key="project.id">
                <td>
                    {{ project.key }}
                </td>
                <td><n-tag :bordered="false" :color="getNaiveUITagColorProperty(project.type.hexColor)">{{
                    project.type.name }}</n-tag>
                </td>
                <td><n-tag :bordered="false" :color="getNaiveUITagColorProperty(project.priority.hexColor)">{{
                    project.priority.name
                        }}</n-tag></td>
                <td><n-tag :bordered="false" :color="getNaiveUITagColorProperty(project.status.hexColor)">{{
                    project.status.name }}</n-tag></td>
                <td>{{ project.summary }}</td>
                <td>{{ project.createdAt.toLocaleString() }}</td>
                <td>
                    <div class="doneo-flex-center-align" style="gap: 8px;">
                        <n-avatar v-if="project.createdBy.avatarUrl" :src="project.createdBy.avatarUrl"
                            class="avatar" />
                        {{ project.createdBy.name }}
                    </div>

                </td>
                <td class="doneo-text-center">
                    <n-button-group size="small">
                        <router-link :to="{ name: 'project', params: { id: project.id } }">
                            <n-button>
                                {{ t("shared.buttons.Open.label") }}
                                <template #icon>
                                    <n-icon :size="22">
                                        <IconFilePencil />
                                    </n-icon>
                                </template>
                            </n-button>
                        </router-link>
                    </n-button-group>
                </td>
            </tr>
            <tr>
                <td :colspan="columns.length + 1" v-if="projects.length < 1 && !props.loading">
                    <n-empty :description="t('modules.project.components.ProjectsTable.warnings.noItemsFound')">
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