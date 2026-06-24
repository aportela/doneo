<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue';

    import { NInputGroup, NInput, NButton, NSelect, NIcon, type SelectOption, type SelectSize, type SelectInst } from 'naive-ui';
    import { IconSquare, IconSquareFilled, IconAlertCircle } from '@tabler/icons-vue';

    import { useCacheStore } from '../../../stores/cache';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectTypeService } from '../services/project-type';
    import type { ProjectTypeResponse } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';
    import { handleAPIError } from '../../../api/client/errorHandler';

    interface ProjectTypeSelectorProps {
        autoFocus?: boolean;
        required?: boolean;
        placeholder?: string;
        clearable?: boolean;
        size?: SelectSize;
        hidePrefix?: boolean;
        disabled?: boolean;
        readonly?: boolean;
    }

    const cacheStore = useCacheStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const selectInstRef = ref<SelectInst | null>(null)

    const isDisabled = computed(() => props.disabled || state.ajaxRunning);

    const projectTypeId = defineModel<string | null>('id');

    const projectTypes = ref<ProjectTypeResponse[]>([]);

    const props = defineProps<ProjectTypeSelectorProps>();

    const options = shallowRef<SelectOption[]>([]);

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            if (cacheStore.projectTypes.length > 0) {
                projectTypes.value = cacheStore.projectTypes;
            } else {
                const response = await projectTypeService.searchBase();
                projectTypes.value = response.projectTypes;
                cacheStore.setProjectTypesCache(projectTypes.value);
            }
            if (projectTypeId.value) {
                selectedColor.value = projectTypes.value.find((projectType) => projectType.id === projectTypeId.value)?.hexColor
            }
            options.value = projectTypes.value.map((projectType: ProjectTypeResponse) => ({ label: projectType.name, value: projectType.id }));
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
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectTypeSelector.onRefresh" } });
                            break;
                        default:
                            console.error("Unhandled API error", { file: "ProjectTypeSelector.vue", method: "onRefresh" });
                            break;
                    }
                },
                (fatalError) => {
                    console.error("Unhandled API error", { file: "ProjectTypeSelector.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
        }
    };

    const selectedColor = ref<string | undefined>();

    watch(projectTypeId, (newValue) => {
        selectedColor.value = projectTypes.value.find((projectType) => projectType.id === newValue)?.hexColor
    });

    const readOnlyLabel = computed({
        get() {
            return projectTypes.value.find((item) => item.id == projectTypeId.value)?.name;
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
        projectTypeId.value = null;
    }

    defineExpose({ reset });

    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectTypeSelector.onRefresh")) {
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
            v-model:value="projectTypeId" :options="options" :placeholder="props.placeholder" :size="props.size"
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