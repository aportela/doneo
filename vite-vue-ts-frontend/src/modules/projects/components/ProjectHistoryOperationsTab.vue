<script setup lang="ts">
    import { shallowRef, reactive, onMounted, onBeforeUnmount, watch, type CSSProperties, type Component } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard, NTimeline, NTimelineItem, NIcon } from "naive-ui";

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { useLoadingStore } from '../../../stores/loading';
    import { appBus } from '../../../shared/composables/bus';

    import { projectHistoryOperationsService } from "../../project-history-operations/services/project-history-operations";
    import { handleAPIError } from '../../../api/client/errorHandler';

    import type { SearchResponse } from "../../project-history-operations/types/dto";
    import { ProjectHistoryOperation } from "../../project-history-operations/models/project-history-operation";
    import { IconSquarePlus, IconEdit, IconDeviceUnknown } from "@tabler/icons-vue";
    import AvatarUserName from "../../../shared/components/AvatarUserName.vue";

    interface ProjectNotesProps {
        style?: string | CSSProperties;
        projectId: string;
    }

    const props = defineProps<ProjectNotesProps>();

    const { t } = useI18n();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectHistoryOperation[]>([]);

    const itemCount = defineModel<number>("itemCount", { default: 0 });

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
        2: IconEdit
    };

</script>

<template>
    <n-card bordered :style="props.style">
        <n-timeline size="large">
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
    </n-card>
</template>

<style lang="css" scoped></style>