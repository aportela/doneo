<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NFlex, NPagination } from 'naive-ui';

    interface PagerProps {
        totalResults: number;
        totalPages: number;
    };

    const props = defineProps<PagerProps>();

    const { t } = useI18n();
    const pageSizes = [
        {
            label: 'All results',
            value: 0
        },
        {
            label: '5 results/page',
            value: 5
        },
        {
            label: '10 results/page',
            value: 10
        },
        {
            label: '20 results/page',
            value: 20
        },
        {
            label: '50 results/page',
            value: 50
        },
        {
            label: '100 results/page',
            value: 100
        },
        {
            label: '200 results/page',
            value: 200
        },
    ];

    const currentPage = defineModel<number>("currentPage");

    const pageSize = defineModel<number>("pageSize");

</script>

<template>
    <n-flex justify="space-between" class="doneo-pagination-container">
        <div class="doneo-pagination-total-results-container">
            <span>Total results: {{ props.totalResults }}</span>
        </div>
        <n-pagination v-model:page="currentPage" v-model:page-size="pageSize" :page-count="totalPages"
            :page-sizes="pageSizes" show-size-picker :page-slot="8">
            <template #prefix="{ page, pageCount }">
                {{ t("labelPageOfTotalPages", { currentPage: page, totalPages: pageCount }) }}
            </template>
        </n-pagination>
    </n-flex>
</template>

<style lang="css" scoped>
    .doneo-pagination-container {
        padding: 4px;
        background-color: rgba(250, 250, 252, 1);
        margin-bottom: 8px;
        border: 1px solid;
        border-color: rgb(239, 239, 245);
        border-radius: 3px;
    }

    .doneo-doneo-pagination-total-results-container {
        padding-top: 2px;
        padding-left: 2px;
    }
</style>