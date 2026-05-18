<script setup lang="ts">
    import { onMounted, onBeforeUnmount, ref, reactive, shallowRef, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';
    import { projectTypeService } from '../services/project-type';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { ProjectTypeResponse } from '../types/dto';
    import { ProjectType } from '../models/project-type';
    import ProjectTypesTable from '../components/ProjectTypesTable.vue';
    //import ProjectTypeForm from '../components/forms/ProjectTypeForm.vue';
    import { Sort } from '../../../shared/types/models/sort';
    import type { FormMode } from '../types/form-mode';

    const { notify } = useNotify();

    const { t } = useI18n();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectTypes = shallowRef<ProjectType[]>([]);

    const nameFilter = ref<string>("");

    const sort = ref<Sort>(new Sort("name", "ASC"));

    const showProjectTypeDialogForm = ref<boolean>(false);
    const projectTypeDialogFormMode = ref<FormMode>("add");

    const selectedProjectTypeId = ref<string>("");

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    const onToggleSort = (field: string) => {
        sort.value.toggleSort(field);
        onRefresh();
    };

    const onShowAddForm = () => {
        projectTypeDialogFormMode.value = "add";
        showProjectTypeDialogForm.value = true;
    };

    const onShowUpdateForm = (projectType: ProjectType, _index: number) => {
        selectedProjectTypeId.value = projectType.id;
        projectTypeDialogFormMode.value = "update";
        showProjectTypeDialogForm.value = true;
    };

    /*
    const onAdd = (role: Role) => {
        showProjectTypeDialogForm.value = false;
        notify('success', t("roleAddedNotification", { name: role.name }));
        onRefresh();
    };

    const onUpdate = (role: Role) => {
        showProjectTypeDialogForm.value = false;
        notify('success', t("roleUpdatedNotification", { name: role.name }));
        onRefresh();
    };

    const onCancel = () => {
        showProjectTypeDialogForm.value = false;
    };

    const onAdd = () => {
        isVisibleActionDialog.value = false;
        notify('success', t("Project type added"))
        onRefresh();
    };

    const onUpdate = () => {
        isVisibleActionDialog.value = false;
        notify('success', t("Project type updated"))
        onRefresh();
    };

    const onDelete = () => {
        isVisibleActionDialog.value = false;
        notify('success', t("Project type deleted"))
        onRefresh();
    };

    const onCancel = () => {
        isVisibleActionDialog.value = false;
    };
    */

    const onRefresh = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response = await projectTypeService.search({});
            projectTypes.value = response.projectTypes.map((projectType: ProjectTypeResponse) => new ProjectType(projectType))
        } catch (error: unknown) {
            projectTypes.value.length = 0;
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

    const onDelete = async (projectType: ProjectType, _index: number) => {
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
                // TODO: missing role/index param at this point
                //onDelete();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <!--
    <n-modal v-model:show="isVisibleActionDialog">
        <ProjectTypeForm :mode="actionDialogMode" :project-type-id="selectedProjectTypeId" style="width: 40%;"
            @add="onAdd" @update="onUpdate" @delete="onDelete" @cancel="onCancel" />
    </n-modal>
-->

    <n-card :title="t('Manage project types')">
        <ProjectTypesTable :project-types="projectTypes" :loading="state.ajaxRunning" @refresh="onRefresh"
            @add="onShowAddForm" @update="onShowUpdateForm" @delete="onDelete" @textfilter-keydown-enter="onRefresh"
            :sort-field="sort.field" :sort-order="sort.order" @toggle-sort="onToggleSort"
            v-model:role-name-filter="nameFilter" />
    </n-card>
</template>

<style lang="css" scoped></style>