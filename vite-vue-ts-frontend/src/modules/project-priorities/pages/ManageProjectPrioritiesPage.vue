<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard } from 'naive-ui';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';
    import { projectPriorityService } from '../services/project-priority';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { ProjectPriorityResponse, SearchRequest } from '../types/dto';
    import { ProjectPriority } from '../models/project-priority';
    import ProjectPrioritiesTable from '../components/ProjectPrioritiesTable.vue';
    import ProjectPriorityForm from '../components/ProjectPriorityForm.vue';
    import Pager from '../../../shared/components/tables/Pager.vue';
    import { Sort } from '../../../shared/types/models/sort';
    import type { FormMode } from '../types/form-mode';

    const { notify } = useNotify();

    const { t } = useI18n();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectPriorities = shallowRef<ProjectPriority[]>([]);

    const sort = ref<Sort>(new Sort("index", "ASC"));

    const nameFilter = ref<string>("");

    const currentPage = ref(1);
    const pageSize = ref(0);
    const totalResults = ref(0);
    const totalPages = ref(0);


    const showProjectPriorityDialogForm = ref<boolean>(false);
    const projectPriorityDialogFormMode = ref<FormMode>("add");

    const selectedProjectPriority = ref<ProjectPriority>(new ProjectPriority({
        id: "",
        name: "",
        hexColor: "",
    }));

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

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

    const onToggleSort = (field: string) => {
        sort.value.toggleSort(field);
        onRefresh();
    };

    const onShowAddForm = () => {
        projectPriorityDialogFormMode.value = "add";
        showProjectPriorityDialogForm.value = true;
    };

    const onShowUpdateForm = (projectPriority: ProjectPriority, _index: number) => {
        selectedProjectPriority.value = projectPriority;
        projectPriorityDialogFormMode.value = "update";
        showProjectPriorityDialogForm.value = true;
    };

    const onAdd = (projectPriority: ProjectPriority) => {
        showProjectPriorityDialogForm.value = false;
        notify('success', t("projectStatusAddedNotification", { name: projectPriority.name }));
        onRefresh();
    };

    const onUpdate = (projectPriority: ProjectPriority) => {
        showProjectPriorityDialogForm.value = false;
        notify('success', t("projectStatusUpdatedNotification", { name: projectPriority.name }));
        onRefresh();
    };

    const onCancel = () => {
        showProjectPriorityDialogForm.value = false;
    };

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: SearchRequest = {
                pager: {
                    currentPage: currentPage.value,
                    resultsPage: pageSize.value,
                },
                order: {
                    field: sort.value.field,
                    sort: sort.value.order,
                },
                filter: {
                    name: nameFilter.value,
                }
            };
            const response = await projectPriorityService.search(payload);
            totalPages.value = response.pager.totalPages;
            totalResults.value = response.pager.totalResults;
            projectPriorities.value = response.projectPriorities.map((projectPriority: ProjectPriorityResponse) => new ProjectPriority(projectPriority))
        } catch (error: unknown) {
            projectPriorities.value.length = 0;
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageProjectPrioritiesPage.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while refreshing the project priority list");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while refreshing the project priority list");
                    console.error("Unhandled API error", { file: "ManageProjectPrioritiesPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
        }
    };

    const onDelete = async (projectPriority: ProjectPriority, _index?: number) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await projectPriorityService.delete(projectPriority.id);
            notify('success', t("projectPriorityIndexMovedNotification", { name: projectPriority.name }));
            onRefresh();
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            selectedProjectPriority.value = projectPriority;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageProjectPrioritiesPage.onDelete" } });
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("We couldn’t find the specified project priority");
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while deleting the project priority");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while deleting the project priority");
                    console.error("Unhandled API error", { file: "ManageProjectPrioritiesPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageProjectPrioritiesPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageProjectPrioritiesPage.onDelete")) {
                onDelete(selectedProjectPriority.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-modal v-model:show="showProjectPriorityDialogForm">
        <ProjectPriorityForm :mode="projectPriorityDialogFormMode == 'add' ? 'add' : 'update'"
            :project-priority-id="selectedProjectPriority.id" style="width: 40%;" @add="onAdd" @update="onUpdate"
            @cancel="onCancel" />
    </n-modal>

    <n-card :title="t('Manage project priorities')">
        <Pager v-model:current-page="currentPage" v-model:page-size="pageSize" :total-pages="totalPages"
            :total-results="totalResults" class="doneo-pager-container">
            <template #total-results-label="{ totalResults }">
                {{ t("TotalProjectPrioritiesPagerLabel", { total: totalResults }) }}
            </template>
        </Pager>
        <ProjectPrioritiesTable :projectPriorities="projectPriorities" :loading="state.ajaxRunning" @refresh="onRefresh"
            @add="onShowAddForm" @update="onShowUpdateForm" @delete="onDelete" @textfilter-keydown-enter="onRefresh"
            :sort-field="sort.field" :sort-order="sort.order" @toggle-sort="onToggleSort"
            v-model:project-priority-name-filter="nameFilter" />
    </n-card>
</template>

<style lang="css" scoped>
    .doneo-pager-container {
        margin-bottom: 4px;
    }
</style>