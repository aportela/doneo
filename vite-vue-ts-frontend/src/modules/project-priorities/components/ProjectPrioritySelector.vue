<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue';

    import { NInputGroup, NInput, NButton, NSelect, NIcon, type SelectOption, type SelectSize, type SelectInst } from 'naive-ui';
    import { IconSquare, IconSquareFilled, IconAlertCircle } from '@tabler/icons-vue';

    import { useCacheStore } from '../../../stores/cache';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectPriorityService } from '../services/project-priority';
    import type { ProjectPriorityResponse } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import { DEFAULT_BUTTON_ICON_SIZE, DEFAULT_SELECTOR_SIZE } from '../../../constants';

    interface Props {
        autoFocus?: boolean;
        required?: boolean;
        placeholder?: string;
        clearable?: boolean;
        size?: SelectSize;
        iconSize?: number;
        hidePrefix?: boolean;
        disabled?: boolean;
        readOnly?: boolean;
    };

    const props = withDefaults(defineProps<Props>(), {
        autoFocus: false,
        required: false,
        clearable: false,
        size: DEFAULT_SELECTOR_SIZE,
        iconSize: DEFAULT_BUTTON_ICON_SIZE,
        disabled: false,
        readOnly: false,
    });

    const cacheStore = useCacheStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const selectInstRef = ref<SelectInst | null>(null)

    const isDisabled = computed(() => props.disabled || state.ajaxRunning);

    const projectPriorityId = defineModel<string | null>('id');

    const projectPriorities = ref<ProjectPriorityResponse[]>([]);

    const options = shallowRef<SelectOption[]>([]);

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            if (cacheStore.projectPriorities.length > 0) {
                projectPriorities.value = cacheStore.projectPriorities;
            } else {
                const response = await projectPriorityService.searchBase();
                projectPriorities.value = response.projectPriorities;
                cacheStore.setProjectPrioritiesCache(projectPriorities.value);
            }
            if (projectPriorityId.value) {
                selectedColor.value = projectPriorities.value.find((projectPriority) => projectPriority.id === projectPriorityId.value)?.hexColor
            }
            options.value = projectPriorities.value.map((projectPriority: ProjectPriorityResponse) => ({ label: projectPriority.name, value: projectPriority.id }));
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
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectPrioritySelector.onRefresh" } });
                            break;
                        default:
                            console.error("Unhandled API error", { file: "ProjectPrioritySelector.vue", method: "onRefresh" });
                            break;
                    }
                },
                (fatalError) => {
                    console.error("Unhandled API error", { file: "ProjectPrioritySelector.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
        }
    };

    const selectedColor = ref<string | undefined>();

    watch(projectPriorityId, (newValue) => {
        selectedColor.value = projectPriorities.value.find((projectPriority) => projectPriority.id === newValue)?.hexColor
    });

    const readOnlyLabel = computed({
        get() {
            return projectPriorities.value.find((item) => item.id == projectPriorityId.value)?.name;
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
        projectPriorityId.value = null;
    };

    defineExpose({ reset });

    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectPrioritySelector.onRefresh")) {
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
        <n-button secondary :size="props.size" :disabled="true" class="doneo-cursor-default doneo-disable-opacity"
            v-if="!props.hidePrefix">
            <template #icon v-if="!state.ajaxErrors">
                <n-icon :size="props.iconSize" :color="selectedColor"
                    :component="selectedColor ? IconSquareFilled : IconSquare" />
            </template>
        </n-button>
        <n-select filterable ref="selectInstRef" :required="props.required" :clearable="props.clearable"
            v-model:value="projectPriorityId" :options="options" :placeholder="props.placeholder" :size="props.size"
            :disabled="isDisabled" v-if="!props.readOnly" />
        <n-input v-else :size="props.size" placeholder="" v-model:value="readOnlyLabel" readonly />
        <n-button secondary :size="props.size" :disabled="true" class="doneo-cursor-default doneo-disable-opacity"
            v-if="state.ajaxErrors">
            <template #icon>
                <n-icon :size="props.iconSize" color="red" :component="IconAlertCircle" />
            </template>
        </n-button>
    </n-input-group>
</template>

<style lang="css" scoped></style>