<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard } from 'naive-ui';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';
    import { projectTypeService } from '../services/project-type';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { ProjectTypeResponse, SearchRequest } from '../types/dto';
    import { ProjectType } from '../models/project-type';
    import ProjectTypesTable from '../components/ProjectTypesTable.vue';
    import ProjectTypeForm from '../components/ProjectTypeForm.vue';
    import { Sort } from '../../../shared/types/models/sort';
    import type { FormMode } from '../../../shared/types/form-mode';

    const { t } = useI18n();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectType[]>([]);

    const sort = ref<Sort>(new Sort("name", "ASC"));

    const nameFilter = ref<string>("");

    const showForm = ref<boolean>(false);
    const formMode = ref<FormMode>("add");

    const selectedItem = ref<ProjectType>(new ProjectType({
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

    const onShowUpdateForm = (projectType: ProjectType, _index: number) => {
        selectedItem.value = projectType;
        formMode.value = "update";
        showForm.value = true;
    };

    const onAdd = (projectType: ProjectType) => {
        showForm.value = false;
        notify('success', t("projectTypeAddedNotification", { name: projectType.name }));
        onRefresh();
    };

    const onUpdate = (projectType: ProjectType) => {
        showForm.value = false;
        notify('success', t("projectTypeUpdatedNotification", { name: projectType.name }));
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
            const response = await projectTypeService.search(payload);
            items.value = response.projectTypes.map((projectType: ProjectTypeResponse) => new ProjectType(projectType))
        } catch (error: unknown) {
            items.value.length = 0;
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageProjectTypesPage.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while refreshing the project type list");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while refreshing the project type list");
                    console.error("Unhandled API error", { file: "ManageProjectTypesPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
        }
    };

    const onDelete = async (projectType: ProjectType, _index?: number) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await projectTypeService.delete(projectType.id);
            notify('success', t("projectTypeDeletedNotification", { name: projectType.name }));
            onRefresh();
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            selectedItem.value = projectType;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageProjectTypesPage.onDelete" } });
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("We couldn’t find the specified project type");
                            break;
                        default:
                            state.ajaxErrorMessage = t("There was a problem while deleting the project type");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("There was a problem while deleting the project type");
                    console.error("Unhandled API error", { file: "ManageProjectTypesPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onRefresh();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ManageProjectTypesPage.onRefresh")) {
                onRefresh();
            } else if (payload.to.includes("ManageProjectTypesPage.onDelete")) {
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
        <ProjectTypeForm :mode="formMode == 'add' ? 'add' : 'update'" :project-type-id="selectedItem.id"
            style="width: 40%;" @add="onAdd" @update="onUpdate" @cancel="onCancel" />
    </n-modal>

    <n-card :title="t('Manage project types')">
        <ProjectTypesTable :project-types="items" :loading="state.ajaxRunning" @refresh="onRefresh" @add="onShowAddForm"
            @update="onShowUpdateForm" @delete="onDelete" @textfilter-keydown-enter="onRefresh" :sort-field="sort.field"
            :sort-order="sort.order" @toggle-sort="onToggleSort" v-model:project-type-name-filter="nameFilter" />
    </n-card>
</template>

<style lang="css" scoped></style>