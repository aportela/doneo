<script setup lang="ts">
    import { reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount, type CSSProperties, type Component } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard, NTimeline, NTimelineItem, NIcon } from "naive-ui";

    import { useLoadingStore } from '../../../stores/loading';
    import { appBus } from '../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import type { SearchResponse } from "../../project-history-operations/types/dto";
    import type { ProjectHistoryOperationsTableFilters } from "../../project-history-operations/types/project-history-operations-table-filters.ts";

    import { projectHistoryOperationsService } from "../../project-history-operations/services/project-history-operations";
    import { handleAPIError } from '../../../api/client/errorHandler';

    import { ProjectHistoryOperation } from "../../project-history-operations/models/project-history-operation";
    import { IconSquarePlus, IconEdit, IconDeviceUnknown, IconTrash, IconMessagePlus, IconFileUpload } from "@tabler/icons-vue";
    import AvatarUserName from "../../../shared/components/AvatarUserName.vue";
    import ProjectHistoryOperationsTable from "../../project-history-operations/components/ProjectHistoryOperationsTable.vue";

    interface ProjectHistoryOperationsTabProps {
        style?: string | CSSProperties;
        projectId: string;
    }

    const props = defineProps<ProjectHistoryOperationsTabProps>();

    const { t } = useI18n();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectHistoryOperation[]>([]);

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    const filters = reactive<ProjectHistoryOperationsTableFilters>({
        userId: null,
    });

    const filteredItems = computed(() => {
        return items.value.filter((operation: ProjectHistoryOperation) => {
            return (filters.userId === null || filters.userId == operation.createdBy.id);
        });
    });

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    watch(() => props.projectId, (newValue, oldValue) => {
        if (!oldValue && newValue) {
            onRefresh();
        }
    });

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const results: SearchResponse = await projectHistoryOperationsService.getProjectHistoryOperations(props.projectId);
            items.value = results.historyOperations.map((operation) => new ProjectHistoryOperation(operation));
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

    const OperationIcons: Record<number, Component> = {
        1: IconSquarePlus,
        2: IconEdit,
        3: IconTrash,
        4: IconMessagePlus,
        5: IconMessagePlus,
        6: IconTrash,
        7: IconFileUpload,
        8: IconTrash,
    };

    const showTimeline = false;
</script>

<template>
    <n-card bordered :style="props.style">
        <n-timeline size="large" v-if="showTimeline">
            <n-timeline-item v-for="item, index in items" :key="index" :title="item.getOperationTypeLabel()"
                :type="item.getNaiveUITimelineItemType()" :time="item.createdAt.toLocaleString()">
                <template #icon>
                    <n-icon :size="24" :component="OperationIcons[item.operationType] ?? IconDeviceUnknown" />
                </template>
                <template #default>
                    <AvatarUserName :user-id="item.createdBy.id" :user-name="item.createdBy.name" />
                </template>
            </n-timeline-item>
        </n-timeline>
        <ProjectHistoryOperationsTable v-else :project-id="props.projectId" :items="filteredItems"
            :disabled="state.ajaxRunning" v-model:filters="filters" @refresh="onRefresh" />
    </n-card>
</template>

<style lang="css" scoped></style>