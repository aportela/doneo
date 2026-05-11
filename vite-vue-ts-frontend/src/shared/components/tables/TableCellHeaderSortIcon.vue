<script setup lang="ts">
    import { NIcon } from 'naive-ui';
    import { IconSortAscending, IconSortDescending } from '@tabler/icons-vue';
    import { type SortOrder } from '../../types/common';

    interface TableCellHeaderSortIconProps {
        order: SortOrder;
    };

    const props = defineProps<TableCellHeaderSortIconProps>();
</script>

<template>
    <n-icon :size="16">
        <IconSortAscending v-if="order == 'ASC'" />
        <IconSortDescending v-else />
    </n-icon>
</template>

<style lang="css" scoped></style>