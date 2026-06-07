<script setup lang="ts">
    import { ref, reactive, shallowRef, watch, onMounted, onBeforeUnmount, type CSSProperties } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard } from "naive-ui";

    import { useLoadingStore } from '../../../stores/loading';
    import { appBus } from '../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import type { SearchResponse } from "../../project-tasks/types/dto.ts";
    import type { ProjectTasksTableFilters } from "../../project-tasks/types/project-tasks-table-filters.ts";

    import { projectTaskService } from "../../project-tasks/services/task.ts";
    import { handleAPIError } from '../../../api/client/errorHandler';

    import { Sort } from '../../../shared/types/models/sort';
    import { ProjectTask } from "../../project-tasks/models/tasks.ts";
    import Pager from '../../../shared/components/tables/Pager.vue';
    import type { SearchRequest } from "../types/dto.ts";
    import ProjectTasksTable from "../../project-tasks/components/ProjectTasksTable.vue";

    interface ProjectTasksProps {
        style?: string | CSSProperties;
        projectId: string;
    }

    const props = defineProps<ProjectTasksProps>();

    const { t } = useI18n();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectTask[]>([]);

    const itemCount = defineModel<number>("itemCount", { default: 0 });
    itemCount.value = 0;

    const sort = reactive<Sort>(new Sort("createdAt", "DESC"));

    const resetPager = ref<boolean>(false);
    const currentPage = ref(1);
    const pageSize = ref(10);
    const totalResults = ref(0);
    const totalPages = ref(0);

    const filters = reactive<ProjectTasksTableFilters>({
        priorityId: null,
        statusId: null,
        summary: "",
        createdAt: {
            from: null,
            to: null,
        },
        createdByUserId: null,
    });

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    watch(() => filters, () => {
        resetPager.value = true;
    }, { deep: true });

    watch(pageSize, () => {
        if (currentPage.value != 1) {
            currentPage.value = 1;
        } else {
            onRefresh();
        }
    });

    watch(currentPage, () => {
        onRefresh();
    });

    const onSort = (newSort: Sort) => {
        sort.field = newSort.field;
        sort.order = newSort.order;
        onRefresh();
    };

    watch(() => props.projectId, (newValue, oldValue) => {
        if (!oldValue && newValue) {
            onRefresh();
        }
    });

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: SearchRequest = {
                pager: {
                    currentPage: currentPage.value,
                    resultsPage: pageSize.value,
                },
                order: {
                    field: sort.field,
                    sort: sort.order,
                },
                filter: {
                }
            };
            const results: SearchResponse = await projectTaskService.search(props.projectId, payload);
            items.value = results.tasks.map((task) => new ProjectTask(task));
            itemCount.value = items.value?.length ?? 0;
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProjectTasksTab.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectPermission.components.projectPermissions.errors.refreshError");
                    console.error("Unhandled API error", { file: "ProjectTasksTab.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProjectTasksTab.onRefresh")) {
                onRefresh();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-card bordered :style="props.style">
        <Pager v-model:current-page="currentPage" v-model:page-size="pageSize" :total-pages="totalPages"
            :total-results="totalResults" class="doneo-pager-container">
            <template #total-results-label="{ totalResults }">
                {{ t("modules.project.components.ManageProjectsPage.pager.totalItemsLabel", { total: totalResults }) }}
            </template>
        </Pager>
        <ProjectTasksTable :items="items" :disabled="state.ajaxRunning" @refresh="onRefresh" :sort="sort" @sort="onSort"
            v-model:filters="filters" />
    </n-card>
</template>

<style lang="css" scoped></style>