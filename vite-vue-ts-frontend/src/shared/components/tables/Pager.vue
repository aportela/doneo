<script setup lang="ts">
    import { computed } from "vue";
    import { useI18n } from "vue-i18n";

    import { NPagination } from 'naive-ui';
    import type { PaginationSize, PaginationSizeOption } from "naive-ui";
    import { type Pagination } from "../../types/pager";

    interface Props {
        disabled?: boolean;
        size?: PaginationSize
        simple?: boolean;
        pagination: Pagination;
    };

    const props = defineProps<Props>();

    const emit = defineEmits(["updateCurrentPageIndex", "updateResultsPage"]);

    const { t } = useI18n();

    const pageSizes = computed<PaginationSizeOption[]>(() => [
        {
            label: t("shared.components.pager.selector.options.allResults"),
            value: 0
        },
        {
            label: t("shared.components.pager.selector.options.nnnResultsPage", { number: 5 }),
            value: 5
        },
        {
            label: t("shared.components.pager.selector.options.nnnResultsPage", { number: 10 }),
            value: 10
        },
        {
            label: t("shared.components.pager.selector.options.nnnResultsPage", { number: 20 }),
            value: 20
        },
        {
            label: t("shared.components.pager.selector.options.nnnResultsPage", { number: 50 }),
            value: 50
        },
        {
            label: t("shared.components.pager.selector.options.nnnResultsPage", { number: 100 }),
            value: 100
        },
        {
            label: t("shared.components.pager.selector.options.nnnResultsPage", { number: 200 }),
            value: 200
        },
    ]);

    const page = computed(
        {
            get() {
                return props.pagination.currentPage;
            },
            set(currentPageIndex: number) { emit("updateCurrentPageIndex", currentPageIndex); },
        }
    );

    const pageSize = computed(
        {
            get() {
                return props.pagination.resultsPage;
            },
            set(resultsPage: number) { emit("updateResultsPage", resultsPage); },
        }
    );
</script>

<template>
    <div class="doneo-pagination-container">
        <div class="doneo-pagination-total-results-label">
            <slot name="total-results-label" :total-results="props.pagination.totalResults">
                {{ t("shared.components.pager.labels.totalResults") }} {{ props.pagination.totalResults }}
            </slot>
        </div>
        <!-- TODO: simple property on small screens ? -->
        <n-pagination :size="props.size" :simple="props.simple" :disabled="props.disabled" v-model:page="page"
            v-model:page-size="pageSize" :page-count="props.pagination.totalPages" :page-sizes="pageSizes"
            show-size-picker :page-slot="8">
            <template #prefix="{ page, pageCount }">
                {{
                    t("shared.components.pager.labels.currentPageOfTotal", { currentPage: page, totalPages: pageCount })
                }}
            </template>
        </n-pagination>
    </div>
</template>

<style lang="css" scoped>
    .doneo-pagination-container {
        justify-content: space-between;
        display: flex;
        align-items: center;
        padding: 4px;
        background-color: var(--n-th-color);
        border: 1px solid var(--n-border-color);
        border-radius: var(--n-border-radius);
    }

    .doneo-pagination-total-results-label {
        padding-top: 2px;
        padding-left: 2px;
    }
</style>