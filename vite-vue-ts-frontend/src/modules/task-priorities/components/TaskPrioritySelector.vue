<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, watch, onMounted, onBeforeUnmount } from 'vue';

    import { NInputGroup, NInput, NButton, NSelect, NIcon, type SelectOption } from 'naive-ui';

    import { useCacheStore } from '../../../stores/cache';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { taskPriorityService } from '../services/task-priority';
    import type { TaskPriorityResponse } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import { DONEO_ICON_ALERT, DONEO_ICON_SQUARE, DONEO_ICON_SQUARE_FILLED } from '../../../shared/types/icons';

    interface Props {
        clearable?: boolean;
        disabled?: boolean;
        placeholder?: string;
        readOnly?: boolean;
        showPrefixIcon?: boolean;
    };

    const props = defineProps<Props>();

    const cacheStore = useCacheStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const isDisabled = computed(() => props.disabled || state.ajaxRunning);

    const taskPriorityId = defineModel<string | null>('id');

    const taskPriorities = ref<TaskPriorityResponse[]>([]);

    const options = shallowRef<SelectOption[]>([]);

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            if (cacheStore.taskPriorities.length === 0) {
                const response = await taskPriorityService.searchBase();
                taskPriorities.value = response.taskPriorities;
                cacheStore.setTaskPrioritiesCache(taskPriorities.value);
            }
            taskPriorities.value = cacheStore.taskPriorities;
            if (taskPriorityId.value) {
                selectedColor.value = taskPriorities.value.find((taskPriority) => taskPriority.id === taskPriorityId.value)?.hexColor
            }
            options.value = taskPriorities.value.map((taskPriority: TaskPriorityResponse) => ({ label: taskPriority.name, value: taskPriority.id }));
        } catch (error: unknown) {
            options.value.length = 0;
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TaskPrioritySelector.onRefresh" } });
                            break;
                        default:
                            console.error("Unhandled API error", { file: "TaskPrioritySelector.vue", method: "onRefresh" });
                            break;
                    }
                },
                (fatalError) => {
                    console.error("Unhandled API error", { file: "TaskPrioritySelector.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
        }
    };

    const selectedColor = ref<string | undefined>();

    watch(taskPriorityId, (newValue) => {
        selectedColor.value = taskPriorities.value.find((taskPriority) => taskPriority.id === newValue)?.hexColor
    });

    const readOnlyLabel = computed({
        get() {
            return taskPriorities.value.find((item) => item.id == taskPriorityId.value)?.name;
        },
        set(_value) {
        }
    });

    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("TaskPrioritySelector.onRefresh")) {
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
        <n-button secondary disabled class="doneo-cursor-default doneo-disable-opacity" v-if="props.showPrefixIcon">
            <template #icon>
                <n-icon :color="selectedColor"
                    :component="selectedColor ? DONEO_ICON_SQUARE_FILLED : DONEO_ICON_SQUARE" />
            </template>
        </n-button>
        <n-select filterable ref="selectInstRef" :clearable="props.clearable" v-model:value="taskPriorityId"
            :options="options" :placeholder="props.placeholder" :disabled="isDisabled" v-if="!props.readOnly" />
        <n-input v-else placeholder="" v-model:value="readOnlyLabel" readonly />
        <n-button secondary disabled class="doneo-cursor-default doneo-disable-opacity" v-if="state.ajaxErrors">
            <template #icon>
                <n-icon color="red" :component="DONEO_ICON_ALERT" />
            </template>
        </n-button>
    </n-input-group>
</template>

<style lang="css" scoped></style>