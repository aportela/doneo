<script setup lang="ts" generic="T">
    import { ref, computed, useSlots } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NTable, type TableSize, NFlex, NIcon, NDrawer, NDrawerContent, NCollapse, NCollapseItem, NButton, NButtonGroup, NEmpty } from 'naive-ui';

    import { DONEO_ICON_ACTION_ADD, DONEO_ICON_ACTION_CLEAR_FILTERS, DONEO_ICON_ACTION_DOWN, DONEO_ICON_ACTION_HIDE, DONEO_ICON_ACTION_REFRESH, DONEO_ICON_ACTION_SETTINGS, DONEO_ICON_ACTION_SHOW, DONEO_ICON_ACTION_UP, DONEO_ICON_FILTERED, DONEO_ICON_TOGGLE_SORT_ASCENDING, DONEO_ICON_TOGGLE_SORT_DESCENDING } from '../../types/icons.ts';

    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';

    import { type TableHeaderColumn } from '../../types/table-header-column';
    import type { Order } from '../../types/order.ts';
    import RenderCell from './RenderCell.ts';
    import { type Pagination } from '../../types/pager.ts';

    import Pager from './Pager.vue';

    type ActionButton = "refresh" | "add" | "settings";

    interface Props {
        id?: string;
        disabled?: boolean;
        size?: TableSize;
        striped?: boolean;
        columns: TableHeaderColumn<T>[];
        rows: T[];
        rowKey: (row: T) => string;
        order: Order;
        pagerData?: Pagination;

        pagerPosition?: "top" | "bottom" | "both";

        buttons?: ActionButton[];

        noItemsWarningMessage?: string;
        showNoItemsWarningMessage?: boolean;
    };

    const props = withDefaults(defineProps<Props>(), {
        buttons: () => ["refresh", "add", "settings"],
    });

    const emit = defineEmits(['pagerChanged', 'sort', 'refresh', 'add', 'clearFilters']);

    const { t } = useI18n();
    const slots = useSlots()
    const tableSettingsStore = useTableSettingsStore();

    const showSettingsDrawer = ref(false);

    const visibleColumns = computed<TableHeaderColumn<T>[]>(() => props.columns.filter((column: TableHeaderColumn<T>) => column.visible));

    const isFiltered = computed(() => props.columns.find((column) => column.isFiltered?.() === true));

    const onUpdateCurrentPageIndex = (currentPageIndex: number) => {
        emit("pagerChanged", { ...props.pagerData, currentPage: currentPageIndex });
    };

    const onUpdateResultsPage = (resultsPage: number) => {
        emit("pagerChanged", { ...props.pagerData, currentPage: 1, resultsPage: resultsPage });
    };

    const onToggleSort = (column: TableHeaderColumn<T>) => {
        if (!props.disabled && props.order && column.sortable) {
            if (props.order.field !== column.field) {
                emit("sort", { field: column.field, direction: "ASC" });
            } else {
                emit("sort", { field: props.order.field, direction: props.order.direction === "ASC" ? "DESC" : "ASC" });
            }
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
            showSettingsDrawer.value = true;
        }
    };

    const onClearFilters = () => {
        emit("clearFilters");
    };

    const onToggleVisibleColumn = (field: string) => {
        if (props.id) {
            tableSettingsStore.toggleVisibleColumn(props.id, field);
        } else {
            console.error("Error saving table settings: missing table id")
        }
    };

    const onShowAllColumns = () => {
        if (props.id) {
            tableSettingsStore.showAllColumns(props.id);
        } else {
            console.error("Error saving table settings: missing table id")
        }
    };

    const onHideAllColumns = () => {
        if (props.id) {
            tableSettingsStore.hideAllColumns(props.id);
        } else {
            console.error("Error saving table settings: missing table id")
        }
    };

    const onToggleAllColumns = () => {
        if (props.id) {
            tableSettingsStore.toggleAllColumns(props.id);
        } else {
            console.error("Error saving table settings: missing table id")
        }
    };

    const onMoveColumn = (field: string, direction: "up" | "down") => {
        if (props.id) {
            tableSettingsStore.moveColumn(props.id, field, direction);
        } else {
            console.error("Error saving table settings: missing table id")
        }
    };
</script>

<template>
    <n-drawer v-if="props.id" v-model:show="showSettingsDrawer" placement="right">
        <n-drawer-content :title="t('shared.components.tables.ManageTable.components.settingsDrawer.title')">
            <n-collapse accordion default-expanded-names="columnVisibility">
                <n-collapse-item title="Column settings" key="columnVisibility">
                    <n-button-group size="tiny">
                        <n-button @click="onShowAllColumns">{{
                            t("shared.components.tables.ManageTable.components.settingsDrawer.buttons.showAllColumns.label")
                            }}</n-button>
                        <n-button @click="onHideAllColumns">{{
                            t("shared.components.tables.ManageTable.components.settingsDrawer.buttons.HideAllColumns.label")
                        }}</n-button>
                        <n-button @click="onToggleAllColumns">{{
                            t("shared.components.tables.ManageTable.components.settingsDrawer.buttons.ToggleColumns.label")
                            }}</n-button>
                    </n-button-group>
                    <p v-for="column, index in props.columns" class="doneo-cursor-pointer doneo-flex-center-align">
                        <n-button-group size="tiny" style="margin-right: 8px;">
                            <n-button @click="onMoveColumn(column.field, 'up')" :disabled="index < 1">
                                <template #icon>
                                    <n-icon :component="DONEO_ICON_ACTION_UP" />
                                </template>
                            </n-button>
                            <n-button @click="onMoveColumn(column.field, 'down')"
                                :disabled="index >= props.columns.length - 1">
                                <template #icon>
                                    <n-icon :component="DONEO_ICON_ACTION_DOWN" />
                                </template>
                            </n-button>
                            <n-button @click="onToggleVisibleColumn(column.field)">
                                <template #icon>
                                    <n-icon :color="column.visible ? 'green' : 'red'"
                                        :component="column.visible ? DONEO_ICON_ACTION_SHOW : DONEO_ICON_ACTION_HIDE" />
                                </template>
                            </n-button>
                        </n-button-group>
                        {{ column.label }}
                    </p>
                </n-collapse-item>
            </n-collapse>
        </n-drawer-content>
    </n-drawer>
    <Pager class="doneo-table-pager" v-if="props.pagerData && (pagerPosition === 'top' || pagerPosition === 'both')"
        :disabled="props.disabled" :pagination="props.pagerData" @update-current-page-index="onUpdateCurrentPageIndex"
        @update-results-page="onUpdateResultsPage" />
    <div class="doneo-table-container" role="region" aria-labelledby="table-caption" tabindex="0">
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
                                <n-icon :component="DONEO_ICON_FILTERED" class="doneo-table-header-icon"
                                    v-if="column.isFiltered?.() ?? false" />
                                <n-icon class="doneo-table-header-icon"
                                    v-if="column.sortable && props.order?.field === column.field"
                                    :component="props.order?.direction == 'DESC' ? DONEO_ICON_TOGGLE_SORT_DESCENDING : DONEO_ICON_TOGGLE_SORT_ASCENDING">
                                </n-icon>
                            </div>
                        </n-flex>
                    </th>
                    <!-- common table actions (refresh/add/settings)-->
                    <th>
                        <n-button-group class="doneo-table-actions-button-group">
                            <n-button @click="onRefresh" :disabled="props.disabled"
                                v-if="props.buttons.includes('refresh')" class="doneo-table-actions-button">
                                <template #icon>
                                    <n-icon :component="DONEO_ICON_ACTION_REFRESH" />
                                </template>
                                {{ t("shared.components.tables.ManageTable.components.buttons.refresh.label") }}
                            </n-button>
                            <n-button @click="onAdd" :disabled="props.disabled" v-if="props.buttons.includes('add')"
                                class="doneo-table-actions-button">
                                <template #icon>
                                    <n-icon :component="DONEO_ICON_ACTION_ADD" />
                                </template>
                                {{ t("shared.components.tables.ManageTable.components.buttons.add.label") }}
                            </n-button>
                            <n-button @click="onSettings" :disabled="props.disabled"
                                v-if="props.buttons.includes('settings') && props.id"
                                class="doneo-table-actions-button">
                                <template #icon>
                                    <n-icon :component="DONEO_ICON_ACTION_SETTINGS" />
                                </template>
                                {{ t("shared.components.tables.ManageTable.components.buttons.settings.label") }}
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
                            :disabled="props.disabled || !isFiltered">
                            <template #icon>
                                <n-icon :component="DONEO_ICON_ACTION_CLEAR_FILTERS" />
                            </template>
                            <!-- TODO: remove component & change label -->
                            {{ t("shared.components.tables.ManageTable.components.buttons.clearFilters.label") }}
                        </n-button>
                    </th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="row, index in props.rows" :key="props.rowKey(row)">
                    <!-- row content -->
                    <td v-for="column in visibleColumns" :key="String(column.field)">
                        <RenderCell :render="column.render" :row="row" />
                    </td>
                    <!-- row actions -->
                    <td class="doneo-text-center">
                        <slot name="rowactions" :row="row" :index="index" />
                    </td>
                </tr>
                <tr v-if="rows.length == 0 && props.noItemsWarningMessage && props.showNoItemsWarningMessage">
                    <td :colspan="visibleColumns.length + 1">
                        <n-empty :description="props.noItemsWarningMessage" />
                    </td>
                </tr>
            </tbody>
        </n-table>
    </div>
    <Pager class="doneo-table-pager" v-if="props.pagerData && (pagerPosition === 'bottom' || pagerPosition === 'both')"
        :disabled="props.disabled" :pagination="props.pagerData" @update-current-page-index="onUpdateCurrentPageIndex"
        @update-results-page="onUpdateResultsPage" />
</template>

<style lang="css" scoped>
    .doneo-table-container {
        overflow-x: auto;
        width: 100%;
    }

    .doneo-table-header-icon {
        margin-top: 4px;
    }

    .doneo-table-pager {
        margin: 4px 0px;
    }

    th {
        min-width: 10em;
    }

</style>