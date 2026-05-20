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
    import { Sort } from '../../../shared/types/models/sort';
    import type { FormMode } from '../../../shared/types/form-mode';

    const { t } = useI18n();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectPriority[]>([]);

    const sort = ref<Sort>(new Sort("name", "ASC"));

    const nameFilter = ref<string>("");

    const showForm = ref<boolean>(false);
    const formMode = ref<FormMode>("add");

    const selectedItem = ref<ProjectPriority>(new ProjectPriority({
        id: "",
        name: "",
        hexColor: "",
    }));

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    const onToggleSort = (field: string) => {
        sort.value.toggleSort(field);
        onRefresh();
    };

    const onShowAddForm = () => {
        formMode.value = "add";
        showForm.value = true;
    };

    const onShowUpdateForm = (projectPriority: ProjectPriority, _index: number) => {
        selectedItem.value = projectPriority;
        formMode.value = "update";
        showForm.value = true;
    };

    const onAdd = (projectPriority: ProjectPriority) => {
        showForm.value = false;
        notify('success', t("projectStatusAddedNotification", { name: projectPriority.name }));
        onRefresh();
    };

    const onUpdate = (projectPriority: ProjectPriority) => {
        showForm.value = false;
        notify('success', t("projectStatusUpdatedNotification", { name: projectPriority.name }));
        onRefresh();
    };

    const onCancel = () => {
        showForm.value = false;
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
                    name: nameFilter.value,
                }
            };
            const response = await projectPriorityService.search(payload);
            items.value = response.projectPriorities.map((projectPriority: ProjectPriorityResponse) => new ProjectPriority(projectPriority))
        } catch (error: unknown) {
            items.value.length = 0;
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
                            selectedItem.value = projectPriority;
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
                onDelete(selectedItem.value);
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-modal v-model:show="showForm">
        <ProjectPriorityForm :mode="formMode == 'add' ? 'add' : 'update'" :project-priority-id="selectedItem.id"
            style="width: 40%;" @add="onAdd" @update="onUpdate" @cancel="onCancel" />
    </n-modal>

    <n-card :title="t('Manage project priorities')">
        <ProjectPrioritiesTable :projectPriorities="items" :loading="state.ajaxRunning" @refresh="onRefresh"
            @add="onShowAddForm" @update="onShowUpdateForm" @delete="onDelete" @textfilter-keydown-enter="onRefresh"
            :sort-field="sort.field" :sort-order="sort.order" @toggle-sort="onToggleSort"
            v-model:project-priority-name-filter="nameFilter" />
    </n-card>
</template>

<style lang="css" scoped></style>