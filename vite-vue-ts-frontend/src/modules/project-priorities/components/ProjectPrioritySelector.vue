<script setup lang="ts">

    import { ref, shallowRef, reactive, computed, watch, onMounted, onBeforeUnmount } from 'vue';

    import { NInputGroup, NButton, NSelect, NIcon, type SelectOption, type SelectSize } from 'naive-ui';
    import { IconSquare, IconSquareFilled } from '@tabler/icons-vue';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectPriorityService } from '../services/project-priority';
    import type { SearchRequest, ProjectPriorityResponse } from '../types/dto';
    import { Sort } from '../../../shared/types/models/sort';
    import { appBus } from '../../../shared/composables/bus';
    import { handleAPIError } from '../../../api/client/errorHandler';


    interface ProjectPrioritySelectorProps {
        placeholder?: string;
        size?: SelectSize;
        hidePrefix?: boolean;
    }

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const isDisabled = computed(() => state.ajaxRunning);

    const projectPriorityId = defineModel<string | null>('id');

    const projectPriorities = ref<ProjectPriorityResponse[]>([]);

    const props = defineProps<ProjectPrioritySelectorProps>();

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
            const response = await projectPriorityService.search(payload);
            projectPriorities.value = response.projectPriorities;
            if (projectPriorityId.value) {
                selectedColor.value = projectPriorities.value.find((projectPriority) => projectPriority.id === projectPriorityId.value)?.hexColor
            }
            options.value = response.projectPriorities.map((projectPriority: ProjectPriorityResponse) => ({ label: projectPriority.name, value: projectPriority.id }));
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
                            //state.ajaxErrorMessage = t("modules.projectPriority.components.ManageProjectPrioritiesPage.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    //state.ajaxErrorMessage = t("modules.projectPriority.components.ManageProjectPrioritiesPage.errors.refreshError");
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
        <n-button secondary :disabled="true" class="doneo-cursor-default doneo-disable-opacity"
            v-if="!props.hidePrefix">
            <template #icon>
                <n-icon :color="selectedColor" :component="selectedColor ? IconSquareFilled : IconSquare">
                </n-icon>
            </template>
        </n-button>
        <n-select v-model:value="projectPriorityId" :options="options" :placeholder="props.placeholder"
            :size="props.size" :disabled="isDisabled" />
    </n-input-group>

</template>

<style lang="css" scoped></style>