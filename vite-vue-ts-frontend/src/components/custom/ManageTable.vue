<script setup lang="ts">
    import { type CSSProperties } from 'vue';
    import { NTable, type TableSize } from 'naive-ui';


    interface ManageTableProps {
        size: TableSize;
        title: string;
        style?: string | CSSProperties;
    };

    const props = defineProps<ManageTableProps>();
</script>

<template>
    <div class="table-header">
        <div>
            <span class="table-caption-title">{{ title }}</span>
        </div>
        <div>
            <slot name="caption-extra" />
        </div>
    </div>
    <n-table :size="size" :style="style">
        <slot name="thead">
            <thead></thead>
        </slot>
        <slot name="tbody">
            <tbody></tbody>
        </slot>
        <slot name="tfoot">
            <tfoot></tfoot>
        </slot>
    </n-table>
</template>

<style lang="css" scoped>
    .table-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 4px;
    }

    .table-caption-title {
        font-weight: 700;
        font-size: large;
    }
</style>