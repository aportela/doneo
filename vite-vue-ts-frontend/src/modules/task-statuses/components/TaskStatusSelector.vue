<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue';

    import { NInputGroup, NButton, NSelect, NIcon, type SelectOption, type SelectSize, type SelectInst } from 'naive-ui';
    import { IconSquare, IconSquareFilled, IconAlertCircle } from '@tabler/icons-vue';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { taskStatusService } from '../services/task-status';
    import type { SearchRequest, TaskStatusResponse } from '../types/dto';
    import { Sort } from '../../../shared/types/models/sort';
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
        setDefaultValueOnStart?: boolean;
    }

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const selectInstRef = ref<SelectInst | null>(null)

    const isDisabled = computed(() => props.disabled || state.ajaxRunning);

    const taskStatusId = defineModel<string | null>('id');

    const taskStatuses = ref<TaskStatusResponse[]>([]);

    const emit = defineEmits([]);

    const props = defineProps<TaskStatusSelectorProps>();

    const sort = ref<Sort>(new Sort("name", "ASC"));

    const options = shallowRef<SelectOption[]>([]);

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: SearchRequest = {
                pager: {
                    currentPage: 1,
                    resultsPage: 0,
                },
                order: {
                    field: sort.value.field,
                    sort: sort.value.order,
                },
                filter: {
                    name: undefined,
                }
            };
            const response = await taskStatusService.search(payload);
            taskStatuses.value = response.taskStatuses;
            if (taskStatusId.value) {
                selectedColor.value = taskStatuses.value.find((taskStatus) => taskStatus.id === taskStatusId.value)?.hexColor
            }
            options.value = response.taskStatuses.map((taskStatus: TaskStatusResponse) => ({ label: taskStatus.name, value: taskStatus.id }));
            if (!taskStatusId.value && props.setDefaultValueOnStart) {
                taskStatusId.value = response.taskStatuses.find((taskStatus: TaskStatusResponse) => taskStatus.flags.defaultStatusOnCreation === true)?.id;
            }
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

    watch(taskStatusId, (newValue) => {
        selectedColor.value = taskStatuses.value.find((taskStatus) => taskStatus.id === newValue)?.hexColor
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
            :disabled="isDisabled" />
        <n-button secondary :disabled="true" class="doneo-cursor-default doneo-disable-opacity" v-if="state.ajaxErrors">
            <template #icon>
                <n-icon color="red" :component="IconAlertCircle">
                </n-icon>
            </template>
        </n-button>
    </n-input-group>
</template>

<style lang="css" scoped></style>