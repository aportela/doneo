<script setup lang="ts">
    import { ref, computed } from 'vue';
    import { useI18n } from "vue-i18n";


    import { NTable, type TableSize, NFlex, NIcon, NDrawer, NDrawerContent, NCollapse, NCollapseItem, NButton, NButtonGroup, NEmpty } from 'naive-ui';

    import { type TableHeaderColumn } from '../../types/table-header-column';
    import { Sort } from '../../types/models/sort.ts';
    import { ArrowDown, ArrowDownWideNarrow, ArrowUp, ArrowUpWideNarrow, Eye, EyeOff, Funnel, ListRestart, Plus, Settings } from '@lucide/vue';

    interface ManageTableProps {
        disabled?: boolean;
        size?: TableSize;
        striped?: boolean;
        columns: TableHeaderColumn[];
        currentSort?: Sort,
        hideRefresh?: boolean;
        hideAdd?: boolean;
        hideSettings?: boolean;
        noItemsWarningMessage?: string;
        showNoItemsWarningMessage?: boolean;
    };

    const emit = defineEmits(['sort', 'refresh', 'add', 'showColumn', 'hideColumn', 'moveColumn']);

    const props = withDefaults(defineProps<ManageTableProps>(), {
        disabled: false,
        hideRefresh: false,
        hideAdd: false,
        hideSettings: false,
    });

    const { t } = useI18n();

    const visibleColumns = computed<TableHeaderColumn[]>(() => props.columns.filter((column: TableHeaderColumn) => column.visible));

    const TABLE_HEADER_ICON_SIZE = 16;

    const showDrawerSettings = ref(false);

    const onToggleSort = (column: TableHeaderColumn) => {
        if (!props.disabled && props.currentSort && column.sortable) {
            const newSort = new Sort(props.currentSort?.field, props.currentSort?.order);
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

    const onToggleColumnVisibility = (column: TableHeaderColumn) => {
        emit(column.visible ? "hideColumn" : "showColumn", column);
    };

    const onShowAllColumns = () => {
        props.columns.forEach((column) => {
            emit("showColumn", column);
        });
    };
    const onHideAllColumns = () => {
        props.columns.forEach((column) => {
            emit("hideColumn", column);
        });
    };
    const onToggleAllColumns = () => {
        props.columns.forEach((column) => {
            emit(column.visible ? "hideColumn" : "showColumn", column);
        });
    };

    const onSortColumn = (index: number, direction: "up" | "down") => {
        emit("moveColumn", props.columns[index], direction);
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
                            <n-button @click="onSortColumn(index, 'up')" :disabled="index < 1">
                                <template #icon>
                                    <n-icon :component="ArrowUp" />
                                </template>
                            </n-button>
                            <n-button @click="onSortColumn(index, 'down')"
                                :disabled="index >= props.columns.length - 1">
                                <template #icon>
                                    <n-icon :component="ArrowDown" />
                                </template>
                            </n-button>
                            <n-button @click="onToggleColumnVisibility(column)">
                                <template #icon>
                                    <n-icon :color="column.visible ? 'green' : 'red'"
                                        :component="column.visible ? Eye : EyeOff" style="margin-right: 4px;" />
                                </template>
                            </n-button>
                        </n-button-group>
                        {{ column.label }}
                    </p>
                </n-collapse-item>
            </n-collapse>
        </n-drawer-content>
    </n-drawer>
    <n-table :size="size" :striped="striped" class="doneo-table" :single-line="false" :single-column="false">
        <thead>
            <tr>
                <th v-for="column in visibleColumns" :key="column.field" @click="onToggleSort(column)"
                    :class="{ 'doneo-cursor-pointer': column.sortable }">
                    <n-flex align="center" justify="space-between">
                        <span v-if="column.align === 'center'"></span>
                        <span>{{ column.label }}</span>
                        <div>
                            <n-icon :size="TABLE_HEADER_ICON_SIZE" :component="Funnel" class="doneo-table-header-icon"
                                v-if="column.isFiltered?.() ?? false" />
                            <n-icon :size="TABLE_HEADER_ICON_SIZE" class="doneo-table-header-icon"
                                v-if="column.sortable && props.currentSort?.field === column.field"
                                :component="props.currentSort?.order == 'DESC' ? ArrowDownWideNarrow : ArrowUpWideNarrow">
                            </n-icon>
                        </div>
                    </n-flex>
                </th>
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
            <slot name="thead" :columns="visibleColumns" />
        </thead>
        <tbody>
            <slot name="tbody" />
            <tr v-if="props.noItemsWarningMessage && props.showNoItemsWarningMessage">
                <td :colspan="props.columns.length + 1">
                    <n-empty :description="props.noItemsWarningMessage">
                    </n-empty>
                </td>
            </tr>
        </tbody>
        <tfoot>
            <slot name="tfoot" />
        </tfoot>
    </n-table>
</template>

<style lang="css" scoped>
    .doneo-table-header-icon {
        margin-top: 4px;
    }
</style>