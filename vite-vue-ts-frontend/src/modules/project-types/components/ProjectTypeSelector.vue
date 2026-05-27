<script setup lang="ts">
    import { ref, shallowRef, reactive, computed, watch, onMounted, onBeforeUnmount } from 'vue';

    import { NInputGroup, NButton, NSelect, NIcon, type SelectOption, type SelectSize } from 'naive-ui';
    import { IconSquare, IconSquareFilled, IconAlertCircle } from '@tabler/icons-vue';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectTypeService } from '../services/project-type';
    import type { SearchRequest, ProjectTypeResponse } from '../types/dto';
    import { Sort } from '../../../shared/types/models/sort';
    import { appBus } from '../../../shared/composables/bus';
    import { handleAPIError } from '../../../api/client/errorHandler';

    interface ProjectTypeSelectorProps {
        placeholder?: string;
        size?: SelectSize;
        hidePrefix?: boolean;
    }

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const isDisabled = computed(() => state.ajaxRunning);

    const projectTypeId = defineModel<string | null>('id');

    const projectTypes = ref<ProjectTypeResponse[]>([]);

    const props = defineProps<ProjectTypeSelectorProps>();

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
            const response = await projectTypeService.search(payload);
            projectTypes.value = response.projectTypes;
            if (projectTypeId.value) {
                selectedColor.value = projectTypes.value.find((projectType) => projectType.id === projectTypeId.value)?.hexColor
            }
            options.value = response.projectTypes.map((projectType: ProjectTypeResponse) => ({ label: projectType.name, value: projectType.id }));
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
        <n-select filterable v-model:value="projectTypeId" :options="options" :placeholder="props.placeholder"
            :size="props.size" :disabled="isDisabled" />
        <n-button secondary :disabled="true" class="doneo-cursor-default doneo-disable-opacity" v-if="state.ajaxErrors">
            <template #icon>
                <n-icon color="red" :component="IconAlertCircle">
                </n-icon>
            </template>
        </n-button>
    </n-input-group>
</template>

<style lang="css" scoped></style>