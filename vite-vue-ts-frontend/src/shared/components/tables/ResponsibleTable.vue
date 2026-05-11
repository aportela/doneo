<script setup lang="ts">
    import { NDataTable, NCard, type DataTableColumn } from 'naive-ui'

    type Column<T = any> = {
        key: string
        label: string
        render?: (row: T) => any
    }

    const props = defineProps<{
        data: any[]
        columns: Column[]
    }>()

    const naiveColumns = props.columns.map((col): DataTableColumn => ({
        title: col.label,
        key: col.key,
        render(row: any) {
            if (col.render) {
                return col.render(row)
            }
            return row[col.key]
        }
    }))
</script>

<template>
    <div class="table-desktop">
        <n-data-table :columns="naiveColumns" :data="data" />
    </div>
    <div class="table-mobile">
        <n-card v-for="(row, index) in data" :key="index" class="mobile-card">
            <div v-for="col in columns" :key="col.key" class="row-field">
                <strong class="label">{{ col.label }}:</strong>
                <span class="value">
                    <template v-if="col.render">
                        <component :is="col.render(row)" />
                    </template>
                    <template v-else>
                        {{ row[col.key] }}
                    </template>
                </span>
            </div>
        </n-card>
    </div>
</template>

<style scoped>
    .table-mobile {
        display: none;
        flex-direction: column;
        gap: 12px;
    }

    .mobile-card {
        padding: 12px;
    }

    .row-field {
        display: flex;
        justify-content: space-between;
        margin-bottom: 6px;
    }

    .label {
        font-weight: 600;
    }

    .value {
        text-align: right;
    }

    @media (max-width: 768px) {
        .table-desktop {
            display: none;
        }

        .table-mobile {
            display: flex;
        }
    }
</style>