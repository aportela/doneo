<script setup lang="ts">
    import { ref, reactive, shallowRef, computed, watch, onMounted, onBeforeUnmount } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NCard } from 'naive-ui';

    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import type { FormMode } from '../../../shared/types/form-mode';
    import type { SearchRequest, ProjectTypeResponse } from '../types/dto';
    import type { ProjectTypesTableFilters } from '../types/project-types-table-filters.ts';

    import { Sort } from '../../../shared/types/models/sort';
    import { ProjectType } from '../models/project-type';

    import { projectTypeService } from '../services/project-type';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import ProjectTypeForm from '../components/ProjectTypeForm.vue';
    import ProjectTypesTable from '../components/ProjectTypesTable.vue';

    const { t } = useI18n();
    const { notify } = useNotify();
    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const items = shallowRef<ProjectType[]>([]);

    const sort = reactive<Sort>(new Sort("name", "ASC"));

    const filters = reactive<ProjectTypesTableFilters>({
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

    const selectedItem = ref<ProjectType>(new ProjectType());

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    const onShowAddForm = () => {
        modalFormMode.value = "add";
        showModal.value = true;
    };

    const onShowUpdateForm = (projectType: ProjectType, _index: number) => {
        selectedItem.value = projectType;
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
                    field: sort.field,
                    sort: sort.order,
                },
                filter: {
                    //name: filters.name.length > 0 ? filters.name : undefined,
                },
            };
            const response = await projectTypeService.search(payload);
            items.value = response.projectTypes.map((projectType: ProjectTypeResponse) => new ProjectType(projectType))
        } catch (error: unknown) {
            items.value = [];
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ManageProjectTypesPage.onRefresh" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectType.components.ManageProjectTypesPage.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectType.components.ManageProjectTypesPage.errors.refreshError");
                    console.error("Unhandled API error", { file: "ManageProjectTypesPage.vue", method: "onRefresh" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onDelete = async (projectType: ProjectType, _index?: number) => {
        if (projectType.id) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await projectTypeService.delete(projectType.id);
                notify('success', t("modules.projectType.components.ManageProjectTypesPage.notifications.projectTypeDeleted", { name: projectType.name }));
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
                                state.ajaxErrorMessage = t("modules.projectType.components.ManageProjectTypesPage.errors.notFoundError");
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectType.components.ManageProjectTypesPage.errors.deleteError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectType.components.ManageProjectTypesPage.errors.deleteError");
                        console.error("Unhandled API error", { file: "ManageProjectTypesPage.vue", method: "onRefresh" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        } else {
            console.error("project type id not set", { file: "ManageProjectTypesPage.vue", method: "onDelete" });
        }
    };

    const onAdded = (projectType: ProjectType) => {
        showModal.value = false;
        notify('success', t("modules.projectType.components.ManageProjectTypesPage.notifications.projectTypeAdded", { name: projectType.name }));
        onRefresh();
    };

    const onUpdated = (projectType: ProjectType) => {
        showModal.value = false;
        notify('success', t("modules.projectType.components.ManageProjectTypesPage.notifications.projectTypeUpdated", { name: projectType.name }));
        onRefresh();
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
    <n-modal v-model:show="showModal">
        <ProjectTypeForm :mode="modalFormMode == 'add' ? 'add' : 'update'" :project-type-id="selectedItem.id"
            class="modal-form" @add="onAdded" @update="onUpdated" @cancel="onCancelForm" />
    </n-modal>

    <n-card :title="t('modules.projectType.components.ManageProjectTypesPage.header.title')">
        <ProjectTypesTable :items="filteredItems" :disabled="state.ajaxRunning" @refresh="onRefresh"
            @add="onShowAddForm" @update="onShowUpdateForm" @delete="onDelete" v-model:filters="filters" />
    </n-card>
</template>

<style lang="css" scoped>
    .modal-form {
        width: 40%;
    }
</style>