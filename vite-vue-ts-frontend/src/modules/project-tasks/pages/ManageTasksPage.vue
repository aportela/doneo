<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NTabs, NTabPane, NIcon, NCard } from 'naive-ui';
    import { IconList, IconLayoutKanban, IconCalendarWeek } from '@tabler/icons-vue';

    import { useLoadingStore } from '../../../stores/loading';
    import { appBus } from '../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import type { SearchRequest, TaskResponse } from '../types/dto';
    import type { ProjectTasksTableFilters } from '../types/project-tasks-table-filters.ts';

    import { Sort } from '../../../shared/types/models/sort';
    import { ProjectTask } from '../models/tasks';

    import { projectTaskService } from '../services/task.ts';
    import { handleAPIError } from '../../../api/client/errorHandler.ts';

    import ProjectTasksTable from '../components/ProjectTasksTable.vue';
    import ProjectTasksKanban from '../components/ProjectTasksKanban.vue';
    import ProjectTasksCalendar from '../components/ProjectTasksCalendar.vue';
    import Pager from '../../../shared/components/tables/Pager.vue';

    const { t } = useI18n();
    const loadingStore = useLoadingStore();

    const tab = ref<string>("List");

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectTask[]>([]);

    const sort = reactive<Sort>(new Sort("createdAt", "DESC"));

    const resetPager = ref<boolean>(false);
    const currentPage = ref(1);
    const pageSize = ref(10);
    const totalResults = ref(0);
    const totalPages = ref(0);

    const filters = reactive<ProjectTasksTableFilters>({
        slug: null,
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

    let stopBusReauthListener: () => void;

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
                    summary: filters.summary.length > 0 ? filters.summary : undefined,
                    priorityId: filters.priorityId !== null ? filters.priorityId : undefined,
                    statusId: filters.statusId !== null ? filters.statusId : undefined,
                    createdAt: filters.createdAt,
                    createdByUserId: filters.createdByUserId !== null ? filters.createdByUserId : undefined,
                }
            };
            const response = await projectTaskService.search(null, payload);
            totalPages.value = response.pager.totalPages;
            totalResults.value = response.pager.totalResults;
            items.value = response.tasks.map((task: TaskResponse) => new ProjectTask(task))
        } catch (error: unknown) {
            items.value = [];
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageTasksPage.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.project.components.ManageTasksPage.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.project.components.ManageTasksPage.errors.refreshError");
                    console.error("Unhandled API error", { file: "ManageTasksPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageTasksPage.onRefresh")) {
                onRefresh();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-tabs>
        <n-tab-pane name="List" tab="List" type="segment" animated v-model="tab" display-directive="show:lazy">
            <template #tab>
                <n-icon size="16" class="tab-icon">
                    <IconList />
                </n-icon>
                List
            </template>
            <template #default>
                <n-card :title="t('modules.task.components.ManageTasksPage.header.title')">
                    <Pager v-model:current-page="currentPage" v-model:page-size="pageSize" :total-pages="totalPages"
                        :total-results="totalResults" class="doneo-pager-container">
                        <template #total-results-label="{ totalResults }">
                            {{ t("modules.task.components.ManageTasksPage.pager.totalItemsLabel", {
                                total:
                                    totalResults
                            }) }}
                        </template>
                    </Pager>
                    <ProjectTasksTable :items="items" :disabled="state.ajaxRunning" @refresh="onRefresh" :sort="sort"
                        @sort="onSort" v-model:filters="filters" :project-id="''" />
                </n-card>
            </template>
        </n-tab-pane>
        <n-tab-pane name="Kanban" tab="Kanban" display-directive="show:lazy">
            <template #tab>
                <n-icon size="16" class="tab-icon">
                    <IconLayoutKanban />
                </n-icon>
                Kanban
            </template>
            <template #default>
                <ProjectTasksKanban />
            </template>
        </n-tab-pane>
        <n-tab-pane name="Calendar" tab="Calendar" display-directive="show:lazy">
            <template #tab>
                <n-icon size="16" class="tab-icon">
                    <IconCalendarWeek />
                </n-icon>
                Calendar
            </template>
            <template #default>
                <ProjectTasksCalendar />
            </template>
        </n-tab-pane>
    </n-tabs>


</template>

<style lang="css" scoped>

    .tab-icon {
        margin-right: 4px;
    }
</style>