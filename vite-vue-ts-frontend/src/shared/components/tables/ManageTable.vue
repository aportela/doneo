<script setup lang="ts">
    import { NTable, type TableSize } from 'naive-ui';

    interface ManageTableProps {
        size?: TableSize;
        striped?: boolean;
    };

    const props = defineProps<ManageTableProps>();
</script>

<template>
    <n-table :size="size" :striped="striped" class="doneo-table">
        <thead>
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

<style lang="css" scoped></style>