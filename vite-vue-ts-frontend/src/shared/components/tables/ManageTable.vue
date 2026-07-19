<script setup lang="ts" generic="T">
    import { ref, computed, useSlots } from 'vue';
    import { useI18n } from "vue-i18n";


    import { NTable, type TableSize, NFlex, NIcon, NDrawer, NDrawerContent, NCollapse, NCollapseItem, NButton, NButtonGroup, NEmpty } from 'naive-ui';
    import { ArrowDown, ArrowDownWideNarrow, ArrowUp, ArrowUpWideNarrow, Eye, EyeOff, Funnel, FunnelX, ListRestart, Plus, Settings } from '@lucide/vue';

    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';

    import { type TableHeaderColumn } from '../../types/table-header-column';
    import { Order } from '../../types/models/sort.ts';
    import Pager from './Pager.vue';
    import RenderCell from './RenderCell.ts';


    interface IProps {
        id: string;
        disabled?: boolean;
        size?: TableSize;
        striped?: boolean;
        columns: TableHeaderColumn<T>[];
        rows: T[];
        rowKey: (row: T) => string;
        currentSort?: Order,

        pagination?: {
            currentPageIndex: number;
            totalPages: number;
            totalResults: number;
        },
        pagerPosition?: "top" | "bottom" | "both";

        hideRefresh?: boolean;
        hideAdd?: boolean;
        hideSettings?: boolean;

        noItemsWarningMessage?: string;
        showNoItemsWarningMessage?: boolean;
    };

    const emit = defineEmits(['sort', 'refresh', 'add', 'clearFilters']);

    const props = withDefaults(defineProps<IProps>(), {
        disabled: false,
        hideRefresh: false,
        hideAdd: false,
        hideSettings: false,
    });

    const { t } = useI18n();
    const slots = useSlots()

    const tableSettingsStore = useTableSettingsStore();

    const isPaginationEnabled = computed<boolean>(() => tableSettingsStore.tables[props.id]?.pagination.enabled ?? true)
    const currentPageIndex = computed(() => tableSettingsStore.tables[props.id]?.pagination.page ?? 1);
    const currentPageSize = computed(() => tableSettingsStore.tables[props.id]?.pagination.pageSize ?? 1);

    const visibleColumns = computed<TableHeaderColumn<T>[]>(() => props.columns.filter((column: TableHeaderColumn<T>) => column.visible));

    const hasColumnsWithFilter = computed(() => props.columns.find((column) => column.isFiltered?.() === true))

    const showDrawerSettings = ref(false);

    const showTopPager = computed(() => isPaginationEnabled.value && props.pagerPosition === "top" || props.pagerPosition === "both");
    const showBottomPager = computed(() => isPaginationEnabled.value && props.pagerPosition === "bottom" || props.pagerPosition === "both");

    const onToggleSort = (column: TableHeaderColumn<T>) => {
        if (!props.disabled && props.currentSort && column.sortable) {
            const newSort = new Order(props.currentSort?.field, props.currentSort?.sort);
            newSort.toggleSort(column.field);
            emit("sort", newSort);
        }
    };

    const onRefresh = () => {
        if (!props.disabled) {
            emit("refresh");
        }
    };

    const onAdd = () => {
        if (!props.disabled) {
            emit("add");
        }
    };

    const onSettings = () => {
        if (!props.disabled) {
            showDrawerSettings.value = true;
        }
    };

    const onClearFilters = () => {
        emit("clearFilters");
    };

    const onToggleVisibleColumn = (field: string) => {
        tableSettingsStore.toggleVisibleColumn(props.id, field);
    };

    const onShowAllColumns = () => {
        tableSettingsStore.showAllColumns(props.id);
    };

    const onHideAllColumns = () => {
        tableSettingsStore.hideAllColumns(props.id);

    };
    const onToggleAllColumns = () => {
        tableSettingsStore.toggleAllColumns(props.id);
    };

    const onMoveColumn = (field: string, direction: "up" | "down") => {
        tableSettingsStore.moveColumn(props.id, field, direction);
    };

</script>

<template>
    <n-drawer v-model:show="showDrawerSettings" placement="right">
        <n-drawer-content title="Table settings">
            <n-collapse accordion default-expanded-names="columnVisibility">
                <n-collapse-item title="Column settings" key="columnVisibility">
                    <n-button-group size="tiny">
                        <n-button @click="onShowAllColumns">Show all</n-button>
                        <n-button @click="onHideAllColumns">Hide all</n-button>
                        <n-button @click="onToggleAllColumns">Toggle values</n-button>
                    </n-button-group>
                    <p v-for="column, index in props.columns" class="doneo-cursor-pointer doneo-flex-center-align">
                        <n-button-group size="tiny" style="margin-right: 8px;">
                            <n-button @click="onMoveColumn(column.field, 'up')" :disabled="index < 1">
                                <template #icon>
                                    <n-icon :component="ArrowUp" />
                                </template>
                            </n-button>
                            <n-button @click="onMoveColumn(column.field, 'down')"
                                :disabled="index >= props.columns.length - 1">
                                <template #icon>
                                    <n-icon :component="ArrowDown" />
                                </template>
                            </n-button>
                            <n-button @click="onToggleVisibleColumn(column.field)">
                                <template #icon>
                                    <n-icon :color="column.visible ? 'green' : 'red'"
                                        :component="column.visible ? Eye : EyeOff" />
                                </template>
                            </n-button>
                        </n-button-group>
                        {{ column.label }}
                    </p>
                </n-collapse-item>
            </n-collapse>
        </n-drawer-content>
    </n-drawer>
    <Pager v-if="showTopPager" v-model:page-size="currentPageSize" v-model:current-page="currentPageIndex"
        :total-results="0" :total-pages="1" class="doneo-table-pager" />
    <n-table :size="size" :striped="striped" class="doneo-table" :single-line="false" :single-column="false">
        <thead>
            <tr>
                <!-- column header labels -->
                <th v-for="column in visibleColumns" :key="column.field" @click="onToggleSort(column)"
                    :class="{ 'doneo-cursor-pointer': column.sortable }">
                    <n-flex align="center" justify="space-between">
                        <span v-if="column.align === 'center'"></span>
                        <span>{{ column.label }}</span>
                        <div>
                            <n-icon :component="Funnel" class="doneo-table-header-icon"
                                v-if="column.isFiltered?.() ?? false" />
                            <n-icon class="doneo-table-header-icon"
                                v-if="column.sortable && props.currentSort?.field === column.field"
                                :component="props.currentSort?.sort == 'DESC' ? ArrowDownWideNarrow : ArrowUpWideNarrow">
                            </n-icon>
                        </div>
                    </n-flex>
                </th>
                <!-- common table actions (refresh/add/settings)-->
                <th>
                    <n-button-group class="doneo-table-actions-button-group">
                        <n-button @click="onRefresh" :disabled="props.disabled" v-if="!props.hideRefresh"
                            class="doneo-table-actions-button">
                            <template #icon>
                                <n-icon :component="ListRestart" />
                            </template>
                            {{ t("shared.buttons.Refresh.label") }}
                        </n-button>
                        <n-button @click="onAdd" :disabled="props.disabled" v-if="!props.hideAdd"
                            class="doneo-table-actions-button">
                            <template #icon>
                                <n-icon :component="Plus" />
                            </template>
                            {{ t("shared.buttons.Add.label") }}
                        </n-button>
                        <n-button @click="onSettings" :disabled="props.disabled" v-if="!props.hideSettings"
                            class="doneo-table-actions-button">
                            <template #icon>
                                <n-icon :component="Settings" />
                            </template>
                            {{ t("shared.buttons.Settings.label") }}
                        </n-button>
                    </n-button-group>
                </th>
            </tr>
            <tr v-if="slots['thead-column-filters']">
                <!-- slot for extra header column filters -->
                <slot name="thead-column-filters" :columns="visibleColumns" />
                <!-- clear filters button -->
                <th class="doneo-text-center">
                    <n-button :size="props.size" block @click="onClearFilters"
                        :disabled="props.disabled || !hasColumnsWithFilter">
                        <template #icon>
                            <n-icon :component="FunnelX" />
                        </template>
                        <!-- TODO: remove component & change label -->
                        {{ t("shared.components.table.filters.button.clearFilters.label") }}
                    </n-button>
                </th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="row in props.rows" :key="props.rowKey(row)">
                <!-- row content -->
                <td v-for="column in visibleColumns" :key="String(column.field)">
                    <RenderCell :render="column.render" :row="row" />
                </td>
                <!-- row actions -->
                <td class="doneo-text-center">
                    <slot name="rowactions" :row="row" />
                </td>
            </tr>
            <tr v-if="rows.length == 0 && props.noItemsWarningMessage && props.showNoItemsWarningMessage">
                <td :colspan="visibleColumns.length + 1">
                    <n-empty :description="props.noItemsWarningMessage" />
                </td>
            </tr>
        </tbody>
    </n-table>
    <Pager v-if="showBottomPager" :totalResults="0" :total-pages="1" class="doneo-table-pager" />
</template>

<style lang="css" scoped>
    .doneo-table-header-icon {
        margin-top: 4px;
    }

    .doneo-table-pager {
        margin: 4px 0px;
    }


</style>