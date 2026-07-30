<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, watch, onMounted, onBeforeUnmount } from 'vue';

    import { NInputGroup, NInput, NButton, NSelect, NIcon, type SelectOption } from 'naive-ui';

    import { useCacheStore } from '../../../stores/cache';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectPriorityService } from '../services/project-priority';
    import type { ProjectPriorityResponse } from '../types/dto';
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

    const projectPriorityId = defineModel<string | null>('id');

    const projectPriorities = ref<ProjectPriorityResponse[]>([]);

    const options = shallowRef<SelectOption[]>([]);

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            if (cacheStore.projectPriorities.length === 0) {
                const response = await projectPriorityService.searchBase();
                projectPriorities.value = response.projectPriorities;
                cacheStore.setProjectPrioritiesCache(projectPriorities.value);
            }
            projectPriorities.value = cacheStore.projectPriorities;
            if (projectPriorityId.value) {
                selectedColor.value = projectPriorities.value.find((projectPriority) => projectPriority.id === projectPriorityId.value)?.hexColor
            }
            options.value = projectPriorities.value.map((projectPriority: ProjectPriorityResponse) => ({ label: projectPriority.name, value: projectPriority.id }));
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
        <n-button secondary disabled class="doneo-cursor-default doneo-disable-opacity" v-if="props.showPrefixIcon">
            <template #icon>
                <n-icon :color="selectedColor"
                    :component="selectedColor ? DONEO_ICON_SQUARE_FILLED : DONEO_ICON_SQUARE" />
            </template>
        </n-button>
        <n-select filterable ref="selectInstRef" :clearable="props.clearable" v-model:value="projectPriorityId"
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