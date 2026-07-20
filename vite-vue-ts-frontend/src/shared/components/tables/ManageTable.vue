<script setup lang="ts" generic="T">
    import { ref, computed, useSlots } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NTable, type TableSize, NFlex, NIcon, NDrawer, NDrawerContent, NCollapse, NCollapseItem, NButton, NButtonGroup, NEmpty } from 'naive-ui';
    import { ArrowDown, ArrowDownWideNarrow, ArrowUp, ArrowUpWideNarrow, Eye, EyeOff, Funnel, FunnelX, ListRestart, Plus, Settings } from '@lucide/vue';

    import Pager from './Pager.vue';
    import { type TableHeaderColumn } from '../../types/table-header-column';
    import type { Order } from '../../types/order.ts';
    import RenderCell from './RenderCell.ts';
    import { PAGER_DEFAULT_RESULTS_PAGE } from '../../types/pager.ts';

    import { useTableSettingsStore } from '../../../stores/tableSettings.ts';

    type actionButton = "refresh" | "add" | "settings";

    interface IProps {
        id: string;
        disabled?: boolean;
        size?: TableSize;
        striped?: boolean;
        columns: TableHeaderColumn<T>[];
        rows: T[];
        rowKey: (row: T) => string;
        order: Order,
        pagerData?: {
            currentPage: number;
            resultsPage: number;
            totalPages: number;
            totalResults: number;
        };

        pagerPosition?: "top" | "bottom" | "both";

        buttons?: actionButton[];

        noItemsWarningMessage?: string;
        showNoItemsWarningMessage?: boolean;
    };

    const props = withDefaults(defineProps<IProps>(), {
        buttons: () => ["refresh", "add", "settings"],
    });

    const emit = defineEmits(['pagerChanged', 'sort', 'refresh', 'add', 'clearFilters']);

    const { t } = useI18n();
    const slots = useSlots()
    const tableSettingsStore = useTableSettingsStore();

    const showSettingsDrawer = ref(false);

    const visibleColumns = computed<TableHeaderColumn<T>[]>(() => props.columns.filter((column: TableHeaderColumn<T>) => column.visible));

    const isFiltered = computed(() => props.columns.find((column) => column.isFiltered?.() === true));

    const currentPage = computed<number>({
        get() {
            return (props.pagerData?.currentPage ?? 1);
        },
        set(value: number) {
            if (props.pagerData) {
                emit("pagerChanged", { ...props.pagerData, currentPage: value });
            }
        },
    });

    const resultsPage = computed<number>({
        get() {
            return (props.pagerData?.resultsPage ?? PAGER_DEFAULT_RESULTS_PAGE);
        },
        set(value: number) {
            if (props.pagerData) {
                emit("pagerChanged", { ...props.pagerData, resultsPage: value });
            }
        },
    });

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
            console.error("Error saving table settings: invalid table id")
        }
    };

    const onShowAllColumns = () => {
        if (props.id) {
            tableSettingsStore.showAllColumns(props.id);
        } else {
            console.error("Error saving table settings: invalid table id")
        }
    };

    const onHideAllColumns = () => {
        if (props.id) {
            tableSettingsStore.hideAllColumns(props.id);
        } else {
            console.error("Error saving table settings: invalid table id")
        }
    };

    const onToggleAllColumns = () => {
        if (props.id) {
            tableSettingsStore.toggleAllColumns(props.id);
        } else {
            console.error("Error saving table settings: invalid table id")
        }
    };

    const onMoveColumn = (field: string, direction: "up" | "down") => {
        if (props.id) {
            tableSettingsStore.moveColumn(props.id, field, direction);
        } else {
            console.error("Error saving table settings: invalid table id")
        }
    };
</script>

<template>
    <n-drawer v-model:show="showSettingsDrawer" placement="right">
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
    <Pager class="doneo-table-pager" v-if="props.pagerData && (pagerPosition === 'top' || pagerPosition === 'both')"
        :disabled="props.disabled" :total-results="props.pagerData.totalResults"
        :total-pages="props.pagerData.totalPages" v-model:current-page="currentPage" v-model:page-size="resultsPage" />
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
                                v-if="column.sortable && props.order?.field === column.field"
                                :component="props.order?.direction == 'DESC' ? ArrowDownWideNarrow : ArrowUpWideNarrow">
                            </n-icon>
                        </div>
                    </n-flex>
                </th>
                <!-- common table actions (refresh/add/settings)-->
                <th>
                    <n-button-group class="doneo-table-actions-button-group">
                        <n-button @click="onRefresh" :disabled="props.disabled" v-if="props.buttons.includes('refresh')"
                            class="doneo-table-actions-button">
                            <template #icon>
                                <n-icon :component="ListRestart" />
                            </template>
                            {{ t("shared.components.tables.ManageTable.components.buttons.refresh.label") }}
                        </n-button>
                        <n-button @click="onAdd" :disabled="props.disabled" v-if="props.buttons.includes('add')"
                            class="doneo-table-actions-button">
                            <template #icon>
                                <n-icon :component="Plus" />
                            </template>
                            {{ t("shared.components.tables.ManageTable.components.buttons.add.label") }}
                        </n-button>
                        <n-button @click="onSettings" :disabled="props.disabled"
                            v-if="props.buttons.includes('settings')" class="doneo-table-actions-button">
                            <template #icon>
                                <n-icon :component="Settings" />
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
                            <n-icon :component="FunnelX" />
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
    <Pager class="doneo-table-pager" v-if="props.pagerData && (pagerPosition === 'bottom' || pagerPosition === 'both')"
        :disabled="props.disabled" :total-results="props.pagerData.totalResults"
        :total-pages="props.pagerData.totalPages" v-model:current-page="currentPage" v-model:page-size="resultsPage" />
</template>

<style lang="css" scoped>
    .doneo-table-header-icon {
        margin-top: 4px;
    }

    .doneo-table-pager {
        margin: 4px 0px;
    }
</style>