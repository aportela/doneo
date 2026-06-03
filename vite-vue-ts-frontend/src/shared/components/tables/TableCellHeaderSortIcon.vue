<script setup lang="ts">
    import { NIcon } from 'naive-ui';
    import { IconSortAscending, IconSortDescending } from '@tabler/icons-vue';

    import { type SortOrder } from '../../types/common';

    interface TableCellHeaderSortIconProps {
        order?: SortOrder;
    };

    const props = defineProps<TableCellHeaderSortIconProps>();
</script>

<template>
    <n-icon :size="16" class="doneo-table-header-sort-icon"
        :component="props.order == 'DESC' ? IconSortDescending : IconSortAscending" />
</template>

<style lang="css" scoped>
    .doneo-table-header-sort-icon {
        margin-top: 4px;
    }
</style>