<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue';

    import { NInputGroup, NInput, NButton, NSelect, NIcon, type SelectOption, type SelectSize, type SelectInst } from 'naive-ui';
    import { IconSquare, IconSquareFilled, IconAlertCircle } from '@tabler/icons-vue';

    import { useCacheStore } from '../../../stores/cache';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { taskStatusService } from '../services/task-status';
    import type { TaskStatusResponse } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';
    import { handleAPIError } from '../../../api/client/errorHandler';

    interface TaskStatusSelectorProps {
        autoFocus?: boolean;
        required?: boolean;
        placeholder?: string;
        clearable?: boolean;
        size?: SelectSize;
        hidePrefix?: boolean;
        disabled?: boolean;
        readonly?: boolean;
        setDefaultValueOnStart?: boolean;
    }

    const cacheStore = useCacheStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const selectInstRef = ref<SelectInst | null>(null)

    const isDisabled = computed(() => props.disabled || state.ajaxRunning);

    const taskStatusId = defineModel<string | null>('id');

    const taskStatuses = ref<TaskStatusResponse[]>([]);

    const emit = defineEmits(["fillEmptyStartDate", "setStartDate", "fillEmptyFinishDate", "setFinishDate", "unsetFinishDateOnLeave"]);

    const props = defineProps<TaskStatusSelectorProps>();

    const options = shallowRef<SelectOption[]>([]);

    const fillEmptyStartDateStatusId = ref<string | null>(null);
    const setStartDateStatusId = ref<string | null>(null);
    const fillEmptyFinishDateStatusId = ref<string | null>(null);
    const setFinishDateStatusId = ref<string | null>(null);
    const unsetFinishDateOnLeaveStatusId = ref<string | null>(null);

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {

            if (cacheStore.taskStatuses.length > 0) {
                taskStatuses.value = cacheStore.taskStatuses;
            } else {
                const response = await taskStatusService.searchBase();
                taskStatuses.value = response.taskStatuses;
                cacheStore.setTaskStatusesCache(taskStatuses.value);
            }
            if (taskStatusId.value) {
                selectedColor.value = taskStatuses.value.find((taskStatus) => taskStatus.id === taskStatusId.value)?.hexColor
            }
            options.value = taskStatuses.value.map((taskStatus: TaskStatusResponse) => ({ label: taskStatus.name, value: taskStatus.id }));
            if (!taskStatusId.value && props.setDefaultValueOnStart) {
                taskStatusId.value = taskStatuses.value.find((taskStatus: TaskStatusResponse) => taskStatus.flags.defaultStatusOnCreation === true)?.id;
            }
            fillEmptyStartDateStatusId.value = taskStatuses.value.find((taskStatus: TaskStatusResponse) => taskStatus.flags.fillEmptyStartDate === true)?.id ?? null;
            setStartDateStatusId.value = taskStatuses.value.find((taskStatus: TaskStatusResponse) => taskStatus.flags.setStartDate === true)?.id ?? null;
            fillEmptyFinishDateStatusId.value = taskStatuses.value.find((taskStatus: TaskStatusResponse) => taskStatus.flags.fillEmptyFinishDate === true)?.id ?? null;
            setFinishDateStatusId.value = taskStatuses.value.find((taskStatus: TaskStatusResponse) => taskStatus.flags.setFinishDate === true)?.id ?? null;
            unsetFinishDateOnLeaveStatusId.value = taskStatuses.value.find((taskStatus: TaskStatusResponse) => taskStatus.flags.unsetFinishDateOnLeave === true)?.id ?? null;
            if (props.autoFocus) {
                focus();
            }
        } catch (error: unknown) {
            options.value.length = 0;
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TaskStatusSelector.onRefresh" } });
                            break;
                        default:
                            console.error("Unhandled API error", { file: "TaskStatusSelector.vue", method: "onRefresh" });
                            break;
                    }
                },
                (fatalError) => {
                    console.error("Unhandled API error", { file: "TaskStatusSelector.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
        }
    };

    const selectedColor = ref<string | undefined>();

    watch(taskStatusId, (newValue, oldValue) => {
        selectedColor.value = taskStatuses.value.find((taskStatus) => taskStatus.id === newValue)?.hexColor
        if (oldValue && oldValue === unsetFinishDateOnLeaveStatusId.value) {
            emit("unsetFinishDateOnLeave");
        } else if (newValue) {
            switch (newValue) {
                case fillEmptyStartDateStatusId.value:
                    emit("fillEmptyStartDate");
                    break;
                case setStartDateStatusId.value:
                    emit("setStartDate");
                    break;
                case fillEmptyFinishDateStatusId.value:
                    emit("fillEmptyFinishDate");
                    break;
                case setFinishDateStatusId.value:
                    emit("setFinishDate");
                    break;
            }
        }
    });

    const readOnlyLabel = computed({
        get() {
            return taskStatuses.value.find((item) => item.id == taskStatusId.value)?.name;
        },
        set(_value) {
        }
    });

    const focus = () => {
        nextTick(() => {
            selectInstRef.value?.focus();
        });
    };

    const reset = () => {
        taskStatusId.value = null;
    }

    defineExpose({ reset });

    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("TaskStatusSelector.onRefresh")) {
                onRefresh();
            }
        });
        onRefresh();
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-input-group>
        <n-button secondary :disabled="true" class="doneo-cursor-default doneo-disable-opacity"
            v-if="!props.hidePrefix">
            <template #icon>
                <n-icon :color="selectedColor" :component="selectedColor ? IconSquareFilled : IconSquare">
                </n-icon>
            </template>
        </n-button>
        <n-select filterable ref="selectInstRef" :required="props.required" :clearable="props.clearable"
            v-model:value="taskStatusId" :options="options" :placeholder="props.placeholder" :size="props.size"
            :disabled="isDisabled" v-if="!props.readonly" />
        <n-input v-else placeholder="" v-model:value="readOnlyLabel" readonly />
        <n-button secondary :disabled="true" class="doneo-cursor-default doneo-disable-opacity" v-if="state.ajaxErrors">
            <template #icon>
                <n-icon color="red" :component="IconAlertCircle">
                </n-icon>
            </template>
        </n-button>
    </n-input-group>
</template>

<style lang="css" scoped></style>