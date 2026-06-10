<script setup lang="ts">
    import { ref } from 'vue';

    import { NTable, type TableSize, NFlex, NIcon, NModal } from 'naive-ui';
    import { IconEye, IconEyeOff, IconFilter, IconSortAscending, IconSortDescending } from '@tabler/icons-vue';

    import { type TableHeaderColumn } from '../../types/table-header-column';
    import RefreshAddActionsColumn from './RefreshAddActionsColumn.vue';
    import { Sort } from '../../types/models/sort.ts';

    interface ManageTableProps {
        disabled?: boolean;
        size?: TableSize;
        striped?: boolean;
        columns: TableHeaderColumn[];
        currentSort?: Sort,
        hideAdd?: boolean;
    };

    const emit = defineEmits(['sort', 'refresh', 'add']);

    const props = withDefaults(defineProps<ManageTableProps>(), {
        disabled: false,
        hideAdd: false,
    });

    const showModal = ref<boolean>(false);

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
            showModal.value = true;
        }
    };

    const tableHeaderIconSize = 16;
</script>

<template>
    <n-modal v-model:show="showModal" preset="card" style="width: 50%;">
        <div>
            Table column visibility settings
            <p v-for="column in columns" class="doneo-cursor-pointer doneo-flex-center-align">
                <n-icon :size="22" :component="column.visible ? IconEye : IconEyeOff" style="margin-right: 4px;" />
                {{ column.label }}
            </p>
        </div>
    </n-modal>
    <n-table :size="size" :striped="striped" class="doneo-table" :single-line="false" :single-column="false">
        <thead>
            <tr v-if="props.columns && props.columns?.length > 0">
                <th v-for="column in props.columns" :key="column.field" @click="onToggleSort(column)"
                    :class="{ 'doneo-cursor-pointer': column.sortable }">
                    <n-flex justify="space-between">
                        <span v-if="column.align === 'center'"> </span>
                        <span>{{ column.label }}</span>
                        <div>
                            <n-icon :size="tableHeaderIconSize" :component="IconFilter" class="doneo-table-header-icon"
                                v-if="column.isFiltered?.() ?? false" />
                            <n-icon :size="tableHeaderIconSize" class="doneo-table-header-icon"
                                v-if="column.sortable && props.currentSort?.field === column.field">
                                <IconSortDescending v-if="props.currentSort?.order == 'DESC'" />
                                <IconSortAscending v-else />
                            </n-icon>
                        </div>
                    </n-flex>
                </th>
                <th>
                    <RefreshAddActionsColumn :columns="columns" :disabled="props.disabled" :hide-add="props.hideAdd"
                        @refresh="onRefresh" @add="onAdd" @settings="onSettings" />
                </th>
            </tr>
            <slot name="thead" />
        </thead>
        <tbody>
            <slot name="tbody" />
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