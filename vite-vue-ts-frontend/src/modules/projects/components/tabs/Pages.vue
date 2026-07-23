<script setup lang="ts">
    import { reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, type CSSProperties } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard } from "naive-ui";

    import { useLoadingStore } from '../../../../stores/loading';
    import { appBus } from '../../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../../shared/types/ajaxState';
    import type { SearchResponse } from "../../../pages/types/dto";
    import type { ProjectPagesTableFilters } from "../../../pages/types/project-pages-table-filters.ts";

    import { pageService } from "../../../pages/services/page.ts";
    import { handleAPIError } from '../../../../api/client/errorHandler';

    import { Page } from "../../../pages/models/page.ts";
    import ProjectPagesTable from "../../../pages/components/ProjectPagesTable.vue";

    interface ProjectHistoryOperationsTabProps {
        style?: string | CSSProperties;
        projectId: string;
    }

    const props = defineProps<ProjectHistoryOperationsTabProps>();

    const { t } = useI18n();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<Page[]>([]);

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    const filters = reactive<ProjectPagesTableFilters>({
        title: null,
        userId: null,
        createdAt: {
            from: null,
            to: null,
        },
        updatedAt: {
            from: null,
            to: null,
        },
    });

    const filteredItems = computed(() => {
        return items.value.filter((page: Page) => {
            return (
                (filters.userId === null || filters.userId == page.createdBy.id) &&
                ((filters.createdAt.from === null && filters.createdAt.to === null) || (page.createdAt?.msTimestamp != null && filters.createdAt.from != null && filters.createdAt.from <= page.createdAt.msTimestamp && filters.createdAt.to != null && filters.createdAt.to >= page.createdAt.msTimestamp))
            );
        });
    });

    watch(
        () => state.ajaxRunning,
        (ajaxRunning) => {
            loadingStore.set(ajaxRunning);
        }
    );

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const results: SearchResponse = await pageService.getProjectPages(props.projectId);
            items.value = results.pages.map((page) => new Page(page));
            itemCount.value = items.value?.length ?? 0;
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectAttachmentsTab.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                    console.error("Unhandled API error", { file: "ProjectAttachmentsTab.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectHistoryOperationsTab.onRefresh")) {
                onRefresh();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-card bordered :style="props.style">
        <ProjectPagesTable :project-id="props.projectId" :items="filteredItems" :disabled="state.ajaxRunning"
            v-model:filters="filters" @refresh="onRefresh" />
    </n-card>
</template>

<style lang="css" scoped></style>