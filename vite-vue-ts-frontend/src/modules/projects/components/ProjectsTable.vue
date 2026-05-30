<script setup lang="ts">
    import { computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NFlex, NEmpty, NTag, NButtonGroup, NButton, NIcon } from 'naive-ui';
    import { IconFilePencil } from '@tabler/icons-vue';

    import { Project } from '../models/project';
    import type { TableHeaderColumn } from '../../../shared/types/table-header-column';
    import { type SortOrder } from '../../../shared/types/common';
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import TextFilterInput from '../../../shared/components/TextFilterInput.vue';
    import TableCellHeaderSortIcon from '../../../shared/components/tables/TableCellHeaderSortIcon.vue';
    import RefreshAddActionsColumn from '../../../shared/components/tables/RefreshAddActionsColumn.vue';
    import ProjectPrioritySelector from '../../project-priorities/components/ProjectPrioritySelector.vue';
    import ProjectTypeSelector from '../../project-types/components/ProjectTypeSelector.vue';
    import ProjectStatusSelector from '../../project-statuses/components/ProjectStatusSelector.vue';
    import DateFilterSelect from '../../../shared/components/selectors/DateFilterSelect.vue';
    import UserSelector from '../../users/components/UserSelector.vue';
    import AvatarUserName from '../../../shared/components/AvatarUserName.vue';
    import { getNaiveUITagColorProperty } from '../../../shared/composables/color';
    import type { TimestampRange } from '../../../shared/composables/timestamps.ts';

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

    const keyFilter = defineModel<string | null>("keyFilter", {
        default: "",
    });

    const summaryFilter = defineModel<string | null>("summaryFilter", {
        default: "",
    });

    const typeIdFilter = defineModel<string | null>("typeIdFilter", {
        default: "",
    });

    const priorityIdFilter = defineModel<string | null>("priorityIdFilter", {
        default: "",
    });

    const statusIdFilter = defineModel<string | null>("statusIdFilter", {
        default: "",
    });

    const createdAtFilter = defineModel<TimestampRange>("createdAtFilter", {
        default: {
            from: null,
            to: null
        }
    });

    const createdByUserIdFilter = defineModel<string | null>("createdByUserIdFilter", {
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
            <tr class="doneo-table-header-click-action">
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
            <tr>
                <th>
                    <TextFilterInput clearable size="small"
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.key.placeholder')"
                        v-model:value="keyFilter" @keydown-enter="onTextFilterKeyDownEnter" />
                </th>
                <th>
                    <ProjectTypeSelector v-model:id="typeIdFilter" :hide-prefix="true" clearable />
                </th>
                <th>
                    <ProjectPrioritySelector v-model:id="priorityIdFilter" :hide-prefix="true" clearable />
                </th>
                <th>
                    <ProjectStatusSelector v-model:id="statusIdFilter" :hide-prefix="true" clearable />
                </th>
                <th>
                    <TextFilterInput clearable size="small"
                        :placeholder="t('modules.project.components.ProjectsTable.header.filters.summary.placeholder')"
                        v-model:value="summaryFilter" @keydown-enter="onTextFilterKeyDownEnter" />
                </th>
                <th>
                    <DateFilterSelect v-model:range="createdAtFilter" />
                </th>
                <th>
                    <UserSelector hideAvatar clearable v-model:id="createdByUserIdFilter" />
                </th>
                <th class="doneo-text-center">
                    <RefreshAddActionsColumn @refresh="onRefresh" @add="onAdd" />
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="project, index in projects" :key="project.id ?? index">
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
                    <n-button-group size="small">
                        <router-link :to="{ name: 'projectTab', params: { id: project.id, tab: 'metadata' } }">
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

<style lang="css" scoped></style>