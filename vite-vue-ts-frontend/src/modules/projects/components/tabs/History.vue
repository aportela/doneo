<script setup lang="ts">
    import { reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, type CSSProperties } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard } from "naive-ui";

    import { useLoadingStore } from '../../../../stores/loading';
    import { appBus } from '../../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../../shared/types/ajaxState';
    import type { SearchResponse } from "../../../history-operations/types/dto";
    import type { HistoryOperationsTableFilters } from "../../../history-operations/types/history-operations-table-filters.ts";

    import { historyOperationsService } from "../../../history-operations/services/history-operations";
    import { handleAPIError } from '../../../../api/client/errorHandler';

    import { HistoryOperation } from "../../../history-operations/models/history-operation";
    import HistoryOperationsTable from "../../../history-operations/components/HistoryOperationsTable.vue";

    interface ProjectHistoryOperationsTabProps {
        style?: string | CSSProperties;
        projectId: string;
    }

    const props = defineProps<ProjectHistoryOperationsTabProps>();

    const { t } = useI18n();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<HistoryOperation[]>([]);

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    const filters = reactive<HistoryOperationsTableFilters>({
        userId: null,
        createdAt: {
            from: null,
            to: null,
        },
        operationType: null,
    });

    const filteredItems = computed(() => {
        return items.value.filter((operation: HistoryOperation) => {
            return (
                (filters.userId === null || filters.userId == operation.createdBy.id) &&
                (filters.operationType === null || filters.operationType == operation.operationType) &&
                ((filters.createdAt.from === null && filters.createdAt.to === null) || (operation.createdAt.msTimestamp != null && filters.createdAt.from != null && filters.createdAt.from <= operation.createdAt.msTimestamp && filters.createdAt.to != null && filters.createdAt.to >= operation.createdAt.msTimestamp))
            );
        });
    });

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const results: SearchResponse = await historyOperationsService.getProjectHistoryOperations(props.projectId);
            items.value = results.historyOperations.map((operation) => new HistoryOperation(operation));
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
        <HistoryOperationsTable :project-id="props.projectId" :items="filteredItems" :disabled="state.ajaxRunning"
            v-model:filters="filters" @refresh="onRefresh" />
    </n-card>
</template>

<style lang="css" scoped></style>