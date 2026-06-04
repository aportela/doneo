<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard } from 'naive-ui';

    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import type { FormMode } from '../../../shared/types/form-mode';
    import type { SearchRequest, ProjectPriorityResponse } from '../types/dto';
    import type { ProjectPrioritiesTableFilters } from '../types/project-priorities-table-filters.ts';

    import { Sort } from '../../../shared/types/models/sort';
    import { ProjectPriority } from '../models/project-priority';

    import { projectPriorityService } from '../services/project-priority';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import ProjectPriorityForm from '../components/ProjectPriorityForm.vue';
    import ProjectPrioritiesTable from '../components/ProjectPrioritiesTable.vue';

    const { t } = useI18n();
    const { notify } = useNotify();
    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectPriority[]>([]);

    const sort = ref<Sort>(new Sort("name", "ASC"));

    const filters = reactive<ProjectPrioritiesTableFilters>({
        name: "",
    });

    const nameFilterLowerCase = computed(() =>
        filters.name.toLowerCase()
    );

    const filteredItems = computed(() => {
        return items.value.filter((projectPriority) => {
            const name = projectPriority.name?.toLowerCase();
            return ((!name || name?.includes(nameFilterLowerCase.value))
            );
        });
    });

    const showModal = ref<boolean>(false);
    const modalFormMode = ref<FormMode>("add");

    const selectedItem = ref<ProjectPriority>(new ProjectPriority());

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    const onSort = (newSort: Sort) => {
        sort.field = newSort.field;
        sort.order = newSort.order;
        onRefresh();
    };

    const onShowAddForm = () => {
        modalFormMode.value = "add";
        showModal.value = true;
    };

    const onShowUpdateForm = (projectPriority: ProjectPriority, _index: number) => {
        selectedItem.value = projectPriority;
        modalFormMode.value = "update";
        showModal.value = true;
    };

    const onCancelForm = () => {
        showModal.value = false;
    };

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
                    //name: filters.name.length > 0 ? filters.name : undefined,
                },
            };
            const response = await projectPriorityService.search(payload);
            items.value = response.projectPriorities.map((projectPriority: ProjectPriorityResponse) => new ProjectPriority(projectPriority))
        } catch (error: unknown) {
            items.value = [];
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageProjectPrioritiesPage.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectPriority.components.ManageProjectPrioritiesPage.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectPriority.components.ManageProjectPrioritiesPage.errors.refreshError");
                    console.error("Unhandled API error", { file: "ManageProjectPrioritiesPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onDelete = async (projectPriority: ProjectPriority, _index?: number) => {
        if (projectPriority.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await projectPriorityService.delete(projectPriority.id);
                notify('success', t("modules.projectPriority.components.ManageProjectPrioritiesPage.notifications.projectPriorityDeleted", { name: projectPriority.name }));
                onRefresh();
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                selectedItem.value = projectPriority;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageProjectPrioritiesPage.onDelete" } });
                                break;
                            case 404:
                                state.ajaxErrorMessage = t("modules.projectPriority.components.ManageProjectPrioritiesPage.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectPriority.components.ManageProjectPrioritiesPage.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectPriority.components.ManageProjectPrioritiesPage.errors.deleteError");
                        console.error("Unhandled API error", { file: "ManageProjectPrioritiesPage.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("project priority id not set", { file: "ManageProjectPrioritiesPage.vue", method: "onDelete" });
        }
    };

    const onAdded = (projectPriority: ProjectPriority) => {
        showModal.value = false;
        notify('success', t("modules.projectPriority.components.ManageProjectPrioritiesPage.notifications.projectPriorityAdded", { name: projectPriority.name }));
        onRefresh();
    };

    const onUpdated = (projectPriority: ProjectPriority) => {
        showModal.value = false;
        notify('success', t("modules.projectPriority.components.ManageProjectPrioritiesPage.notifications.projectPriorityUpdated", { name: projectPriority.name }));
        onRefresh();
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageProjectPrioritiesPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageProjectPrioritiesPage.onDelete")) {
                onDelete(selectedItem.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-modal v-model:show="showModal">
        <ProjectPriorityForm :mode="modalFormMode == 'add' ? 'add' : 'update'" :project-priority-id="selectedItem.id"
            class="modal-form" @add="onAdded" @update="onUpdated" @cancel="onCancelForm" />
    </n-modal>

    <n-card :title="t('modules.projectPriority.components.ManageProjectPrioritiesPage.header.title')">
        <ProjectPrioritiesTable :items="filteredItems" :disabled="state.ajaxRunning" @refresh="onRefresh"
            @add="onShowAddForm" @update="onShowUpdateForm" @delete="onDelete" :sort="sort" @sort="onSort"
            v-model:filters="filters" />
    </n-card>
</template>

<style lang="css" scoped>
    .modal-form {
        width: 40%;
    }
</style>