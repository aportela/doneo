<script setup lang="ts">
    import { ref, reactive, watch, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NInput, NEmpty, type InputInst } from 'naive-ui';

    import { useDebounceFn } from '@vueuse/core';

    import { appBus } from '../../../shared/composables/bus';

    import { searchService } from '../../../modules/search/services/search';
    import { SearchResultItem } from '../../../modules/search/models/search';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from "../../../shared/types/ajaxState";

    const { t } = useI18n();

    const show = defineModel<boolean>("show");

    const bodyStyle = {
        width: '600px'
    };

    const segmented = {
        content: 'soft',
        footer: 'soft'
    } as const;


    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const showWarningNoResults = ref<boolean>(false);

    const searchQuery = ref<string>("");

    const searchQueryRef = ref<InputInst | null>(null);

    const debouncedSearch = useDebounceFn(() => {
        onSearch(searchQuery.value);
    }, 300);

    watch(searchQuery, () => {
        debouncedSearch();
    });

    const results = ref<SearchResultItem[]>([]);

    const onSearch = async (query: string) => {
        showWarningNoResults.value = false;
        if (query) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                results.value = (await searchService.search({})).results?.map((item) => new SearchResultItem(item));
                showWarningNoResults.value = results.value?.length === 0;
            } catch (error) {
                // TODO
                console.error(error);
            } finally {
                state.ajaxRunning = false;
                await nextTick();
                searchQueryRef.value?.focus();
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        }
    };
</script>

<template>
    <n-modal v-model:show="show" :title="t('shared.components.modals.SearchModal.title')" :closable="true" preset="card"
        size="medium" :bordered="true" :segmented="segmented" :style="bodyStyle">
        <n-input :placeholder="t('shared.components.modals.SearchModal.inputs.query.placeholder')" autofocus
            :loading="state.ajaxRunning" :disabled="state.ajaxRunning" v-model:value="searchQuery"
            ref="searchQueryRef" />
        <n-empty v-if="showWarningNoResults" :description="t('shared.components.modals.SearchModal.labels.noResults')"
            style="margin-top: 16px;" />
        <template #footer>
            <div class="shortcut">
                <kbd>↵</kbd><span>Select</span>
                <kbd>↑</kbd><kbd>↓</kbd><span>Toggle</span>
                <kbd>Del</kbd><span>Delete</span>
                <kbd>Esc</kbd><span>Close</span>
            </div>
        </template>
    </n-modal>
</template>

<style lang="css" scoped>
    .shortcut {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 12px;
    }

    kbd {
        padding: 2px 8px;
        border-radius: 4px;
        border: 1px solid #ccc;
        background: #f5f5f5;
        font-family: monospace;
        font-size: 14px;
    }
</style>