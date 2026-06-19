<script setup lang="ts">
    import { ref, reactive, shallowRef, watch, onMounted, onBeforeUnmount, type CSSProperties } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard, NModal } from "naive-ui";

    import { useLoadingStore } from '../../../../stores/loading';
    import { useNotify } from '../../../../shared/composables/notification';
    import { appBus } from '../../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../../shared/types/ajaxState';
    import type { SearchResponse } from "../../../time-trackings/types/dto.ts";
    import type { TimeTrackingsTableFilters } from "../../../time-trackings/types/time-trackings-table-filters.ts";

    import { timeTrackingService } from "../../../time-trackings/services/time-tracking.ts";
    import { handleAPIError } from '../../../../api/client/errorHandler';

    import TimeTrackingForm from "../../../time-trackings/components/TimeTrackingForm.vue";
    import { TimeTracking } from "../../../time-trackings/models/time-tracking.ts"
    import TimeTrackingsTable from "../../../time-trackings/components/TimeTrackingsTable.vue";

    interface TimeTrackingsTabProps {
        style?: string | CSSProperties;
        projectId: string;
        taskId: string;
    }

    const props = defineProps<TimeTrackingsTabProps>();

    const { t } = useI18n();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<TimeTracking[]>([]);

    const itemCount = defineModel<number>("itemCount", { default: 0 });

    const filters = reactive<TimeTrackingsTableFilters>({
        createdByUserId: null,
        createdAt: {
            from: null,
            to: null,
        },
        summary: "",
    });

    const showForm = ref<boolean>(false);

    const selectedItem = ref<TimeTracking>(new TimeTracking());

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    const onShowAddForm = () => {
        showForm.value = true;
    };

    const onCancelForm = () => {
        showForm.value = false;
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const results: SearchResponse = await timeTrackingService.getTaskTimeTrackings(props.projectId, props.taskId);
            items.value = results.timeTrackings.map((timeTracking) => new TimeTracking(timeTracking));
            itemCount.value = items.value?.length ?? 0;
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TrackTimeTrackingsTab.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.task.components.TimeTrackingsTab.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.task.components.TimeTrackingsTab.errors.refreshError");
                    console.error("Unhandled API error", { file: "TrackTimeTrackingsTab.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    const onDelete = async (timeTracking: TimeTracking, _index?: number) => {
        if (timeTracking.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await timeTrackingService.deleteTaskTimeTracking(props.projectId, props.taskId, timeTracking.id);
                items.value = items.value.filter((item) => item.id != timeTracking.id)
                itemCount.value = items.value?.length ?? 0;
                notify('success', t("modules.timeTracking.components.taskTimeTrackingsTab.notifications.timeTrackingDeleted", { summary: timeTracking.summary }));
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                selectedItem.value = timeTracking;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectPermissions.onDelete" } });
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.timeTracking.components.taskTimeTrackingsTab.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.timeTracking.components.taskTimeTrackingsTab.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.timeTracking.components.taskTimeTrackingsTab.errors.deleteError");
                        console.error("Unhandled API error", { file: "TimeTrackings.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
            }
        } else {
            console.error("(project permission id || project id) not set", { file: "TimeTrackings.vue", method: "onDelete" });
        }
    };
    const onAdded = (timeTracking: TimeTracking) => {
        showForm.value = false;
        items.value = [timeTracking, ...items.value]
        itemCount.value = items.value?.length ?? 0;
        notify('success', t("modules.timeTracking.components.taskTimeTrackingsTab.notifications.taskTimeTrackingAdded", { summary: timeTracking.summary }));
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("TrackTimeTrackingsTab.onRefresh")) {
                onRefresh();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-modal v-model:show="showForm">
        <TimeTrackingForm :project-id="props.projectId" :task-id="props.taskId" mode="add" style="width: 40%;"
            @add="onAdded" @cancel="onCancelForm" />
    </n-modal>
    <n-card bordered :style="props.style">
        <TimeTrackingsTable :project-id="props.projectId" :task-id="props.taskId" :items="items"
            :disabled="state.ajaxRunning" v-model:filters="filters" @refresh="onRefresh" @add="onShowAddForm"
            @delete="onDelete" />
    </n-card>
</template>

<style lang="css" scoped></style>