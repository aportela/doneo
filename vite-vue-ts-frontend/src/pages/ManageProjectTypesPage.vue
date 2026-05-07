<script setup lang="ts">
    import { ref, reactive, shallowRef, onMounted, computed } from 'vue'
    import { useI18n } from "vue-i18n";
    import { NButton, NButtonGroup, NModal, NTag, NInput, NIcon } from 'naive-ui'
    import { api } from '../composables/api';
    import { IconEdit, IconPlus, IconTrash, IconSearch } from '@tabler/icons-vue';
    import { getNaiveUITagColorProperty } from '../composables/color';
    import type { ProjectTypeInterface } from '../types/models/projectType';
    import type { SearchProjectTypesResponse } from '../types/apiResponses';
    import { type AjaxStateInterface, defaultAjaxState } from '../types/ajaxState';
    import { type EntityAction } from '../types/common';
    import { useLoadingStore } from '../stores/loading';
    import { useNotify } from '../composables/notification';
    import { default as ProjectTypeForm } from '../components/forms/ProjectTypeForm.vue';
    import { default as ManageTable } from '../components/custom/ManageTable.vue';

    const { notify } = useNotify();

    const { t } = useI18n();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectTypes = shallowRef<ProjectTypeInterface[]>([]);

    const onRefresh = () => {
        state.ajaxRunning = true;
        loadingStore.set(true);
        api.projectType.search().then((successResponse: SearchProjectTypesResponse) => {
            projectTypes.value = [...successResponse.data.projectTypes];
        }).catch((errorResponse: any) => {
            // TODO
            console.log(errorResponse);
        }).finally(() => {
            loadingStore.set(false);
            state.ajaxRunning = false;
        });
    };

    const selectedProjectTypeId = ref<string | undefined>(undefined);

    const onAddProjectType = () => {
        actionDialogMode.value = "add";
    };

    const onUpdateProjectType = (projectType: ProjectTypeInterface, _index: number) => {
        actionDialogMode.value = "update";
        selectedProjectTypeId.value = projectType.id;
    };
    const onDeleteProjectType = (projectType: ProjectTypeInterface, _index: number) => {
        actionDialogMode.value = "delete";
        selectedProjectTypeId.value = projectType.id;
    };

    const actionDialogMode = ref<EntityAction>("none");

    const isVisibleActionDialog = computed<boolean>({
        get: () => actionDialogMode.value !== "none",
        set: (value: boolean) => {
            if (!value) {
                actionDialogMode.value = "none";
            }
        }
    });

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

    const filterByName = ref<string | null>(null);

    onMounted(() => {
        onRefresh();
    });

</script>

<template>
    <n-modal v-model:show="isVisibleActionDialog">
        <ProjectTypeForm :mode="actionDialogMode" :project-type-id="selectedProjectTypeId" style="width: 40%;"
            @add="onAdd" @update="onUpdate" @delete="onDelete" @cancel="onCancel" />
    </n-modal>

    <ManageTable size="small" title="Manage project types">
        <template #thead>
            <tr>
                <th style="width: auto;">Name</th>
                <th class="text-center" style="width: 1%;white-space: nowrap;">Actions</th>
            </tr>
            <tr>
                <th>
                    <n-input size="small" placeholder="search by name..." v-model:value="filterByName" clearable>
                        <template #prefix>
                            <n-icon>
                                <IconSearch />
                            </n-icon>
                        </template>
                    </n-input>
                </th>
                <th class="text-center">
                    <n-button size="small" block @click="onAddProjectType">
                        <template #icon>
                            <NIcon>
                                <IconPlus />
                            </NIcon>
                        </template>
                        Add
                    </n-button>
                </th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="projectType, index in projectTypes" :key="projectType.id">
                <td>
                    <n-tag :color="getNaiveUITagColorProperty(projectType.hexColor)">{{ projectType.name }}</n-tag>
                </td>
                <td class="text-center">
                    <n-button-group>
                        <n-button size="small" @click="onUpdateProjectType(projectType, index)">
                            {{ t("Update") }}
                            <template #icon>
                                <IconEdit :size="22" />
                            </template>
                        </n-button>
                        <n-button size="small" @click="onDeleteProjectType(projectType, index)">
                            {{ t("Delete") }}
                            <template #icon>
                                <IconTrash :size="22" />
                            </template>
                        </n-button>
                    </n-button-group>
                </td>
            </tr>
        </template>

    </ManageTable>
</template>

<style lang="css" scoped></style>