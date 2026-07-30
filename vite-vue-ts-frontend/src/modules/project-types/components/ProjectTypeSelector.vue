<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, watch, onMounted, onBeforeUnmount } from 'vue';

    import { NInputGroup, NInput, NButton, NSelect, NIcon, type SelectOption } from 'naive-ui';

    import { useCacheStore } from '../../../stores/cache';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectTypeService } from '../services/project-type';
    import type { ProjectTypeResponse } from '../types/dto';
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

    const projectTypeId = defineModel<string | null>('id');

    const projectTypes = ref<ProjectTypeResponse[]>([]);

    const options = shallowRef<SelectOption[]>([]);

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            if (cacheStore.projectTypes.length === 0) {
                const response = await projectTypeService.searchBase();
                projectTypes.value = response.projectTypes;
                cacheStore.setProjectTypesCache(projectTypes.value);
            }
            projectTypes.value = cacheStore.projectTypes;
            if (projectTypeId.value) {
                selectedColor.value = projectTypes.value.find((projectType) => projectType.id === projectTypeId.value)?.hexColor
            }
            options.value = projectTypes.value.map((projectType: ProjectTypeResponse) => ({ label: projectType.name, value: projectType.id }));
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
        <n-button secondary disabled class="doneo-cursor-default doneo-disable-opacity" v-if="props.showPrefixIcon">
            <template #icon>
                <n-icon :color="selectedColor"
                    :component="selectedColor ? DONEO_ICON_SQUARE_FILLED : DONEO_ICON_SQUARE" />
            </template>
        </n-button>
        <n-select filterable :clearable="props.clearable" v-model:value="projectTypeId" :options="options"
            :placeholder="props.placeholder" :disabled="isDisabled" v-if="!props.readOnly" />
        <n-input v-else placeholder="" v-model:value="readOnlyLabel" readonly />
        <n-button secondary disabled class="doneo-cursor-default doneo-disable-opacity" v-if="state.ajaxErrors">
            <template #icon>
                <n-icon color="red" :component="DONEO_ICON_ALERT" />
            </template>
        </n-button>
    </n-input-group>
</template>

<style lang="css" scoped></style>