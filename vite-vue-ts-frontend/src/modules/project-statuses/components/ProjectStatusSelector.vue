<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, watch, onMounted, onBeforeUnmount } from 'vue';

    import { NInputGroup, NInput, NButton, NSelect, NIcon, type SelectOption } from 'naive-ui';

    import { useCacheStore } from '../../../stores/cache';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectStatusService } from '../services/project-status';
    import type { ProjectStatusResponse } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import { DONEO_ICON_ALERT, DONEO_ICON_SQUARE, DONEO_ICON_SQUARE_FILLED } from '../../../shared/types/icons';

    interface Props {
        clearable?: boolean;
        disabled?: boolean;
        placeholder?: string;
        readOnly?: boolean;
        setDefaultValueOnStart?: boolean;
        showPrefixIcon?: boolean;
    };

    const props = defineProps<Props>();

    const cacheStore = useCacheStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const isDisabled = computed(() => props.disabled || state.ajaxRunning);

    const projectStatusId = defineModel<string | null>('id');

    const projectStatuses = ref<ProjectStatusResponse[]>([]);

    const emit = defineEmits(["fillEmptyStartDate", "setStartDate", "fillEmptyFinishDate", "setFinishDate", "unsetFinishDateOnLeave"]);

    const options = shallowRef<SelectOption[]>([]);

    const fillEmptyStartDateStatusId = ref<string | null>(null);
    const setStartDateStatusId = ref<string | null>(null);
    const fillEmptyFinishDateStatusId = ref<string | null>(null);
    const setFinishDateStatusId = ref<string | null>(null);
    const unsetFinishDateOnLeaveStatusId = ref<string | null>(null);

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            if (cacheStore.projectStatuses.length === 0) {
                const response = await projectStatusService.searchBase();
                projectStatuses.value = response.projectStatuses;
                cacheStore.setProjectStatusesCache(projectStatuses.value);
            }
            projectStatuses.value = cacheStore.projectStatuses;
            if (projectStatusId.value) {
                selectedColor.value = projectStatuses.value.find((projectStatus) => projectStatus.id === projectStatusId.value)?.hexColor
            }
            options.value = projectStatuses.value.map((projectStatus: ProjectStatusResponse) => ({ label: projectStatus.name, value: projectStatus.id }));
            if (!projectStatusId.value && props.setDefaultValueOnStart) {
                projectStatusId.value = projectStatuses.value.find((projectStatus: ProjectStatusResponse) => projectStatus.flags.defaultStatusOnCreation === true)?.id;
            }
            fillEmptyStartDateStatusId.value = projectStatuses.value.find((projectStatus: ProjectStatusResponse) => projectStatus.flags.fillEmptyStartDate === true)?.id ?? null;
            setStartDateStatusId.value = projectStatuses.value.find((projectStatus: ProjectStatusResponse) => projectStatus.flags.setStartDate === true)?.id ?? null;
            fillEmptyFinishDateStatusId.value = projectStatuses.value.find((projectStatus: ProjectStatusResponse) => projectStatus.flags.fillEmptyFinishDate === true)?.id ?? null;
            setFinishDateStatusId.value = projectStatuses.value.find((projectStatus: ProjectStatusResponse) => projectStatus.flags.setFinishDate === true)?.id ?? null;
            unsetFinishDateOnLeaveStatusId.value = projectStatuses.value.find((projectStatus: ProjectStatusResponse) => projectStatus.flags.unsetFinishDateOnLeave === true)?.id ?? null;
        } catch (error: unknown) {
            options.value.length = 0;
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectStatusSelector.onRefresh" } });
                            break;
                        default:
                            console.error("Unhandled API error", { file: "ProjectStatusSelector.vue", method: "onRefresh" });
                            break;
                    }
                },
                (fatalError) => {
                    console.error("Unhandled API error", { file: "ProjectStatusSelector.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
        }
    };

    const selectedColor = ref<string | undefined>();

    watch(projectStatusId, (newValue, oldValue) => {
        selectedColor.value = projectStatuses.value.find((projectStatus) => projectStatus.id === newValue)?.hexColor
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
            return projectStatuses.value.find((item) => item.id == projectStatusId.value)?.name;
        },
        set(_value) {
        }
    });

    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectStatusSelector.onRefresh")) {
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
        <n-select filterable ref="selectInstRef" :clearable="props.clearable" v-model:value="projectStatusId"
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