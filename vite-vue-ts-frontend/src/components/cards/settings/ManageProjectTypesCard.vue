<script setup lang="ts">
    import { ref, reactive, shallowRef, onMounted, computed } from 'vue'
    import { useI18n } from "vue-i18n";
    import { NTable, NButton, NButtonGroup, NGrid, NGridItem, NModal, NTag } from 'naive-ui'
    import { api } from '../../../composables/api';
    import { IconEdit, IconPlus, IconTrash } from '@tabler/icons-vue';
    import { getNaiveUITagColorProperty } from '../../../composables/color';
    import type { ProjectTypeInterface } from '../../../types/models/projectType';
    import { type AjaxStateInterface, defaultAjaxState } from '../../../types/ajaxState';
    import { type EntityAction } from '../../../types/common';
    import { useLoadingStore } from '../../../stores/loading';
    import { useNotify } from '../../../composables/notification';
    import { default as ProjectTypeForm } from '../../forms/ProjectTypeForm.vue';
    import { default as ManageTable } from '../../custom/ManageTable.vue';

    const { notify } = useNotify();

    const { t } = useI18n();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectTypes = shallowRef<ProjectTypeInterface[]>([]);

    const onRefresh = () => {
        state.ajaxRunning = true;
        loadingStore.set(true);
        api.projectTypes.search().then((successResponse: any) => {
            projectTypes.value = [...successResponse.data.projectTypes];
        }).catch((errorResponse: any) => {
            // TODO
            console.log(errorResponse);
        }).finally(() => {
            loadingStore.set(false);
            state.ajaxRunning = false;
        })
    };

    const selectedProjectTypeId = ref<string | undefined>(undefined);

    const onAddProjectType = () => {
        actionDialogMode.value = "add";
    };

    const onUpdateProjectType = (_projectType: ProjectTypeInterface, _index: number) => {
        actionDialogMode.value = "update";
        selectedProjectTypeId.value = _projectType.id;
    };
    const onDeleteProjectType = (_projectType: ProjectTypeInterface, _index: number) => {
        actionDialogMode.value = "delete";
        selectedProjectTypeId.value = _projectType.id;
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

    onMounted(() => {
        onRefresh();
    });

</script>

<template>
    <n-modal v-model:show="isVisibleActionDialog">
        <ProjectTypeForm :mode="actionDialogMode" :project-type-id="selectedProjectTypeId" style="width: 40%;"
            @add="onAdd" @update="onUpdate" @delete="onDelete" @cancel="onCancel" />
    </n-modal>
    <ManageTable size="small" :title="t('Project types')">
        <template #caption-extra>
            <n-button @click="onAddProjectType" :disabled="state.ajaxRunning">
                <template #icon>
                    <IconPlus />
                </template>
                {{ t("Add") }}
            </n-button>
        </template>
        <template #thead>
            <tr>
                <th>{{ t("Name") }}</th>
                <th class="text-center">{{ t("Actions") }}</th>
            </tr>
        </template>
        <template #tbody>
            <tr v-for="projectType, index in projectTypes" :key="projectType.id">
                <td><n-tag :color="getNaiveUITagColorProperty(projectType.hexColor)">{{ projectType.name
                        }}</n-tag><span style="display: inline-block; width: 100%;">&nbsp;</span></td>
                <td class="text-center">
                    <n-button-group>
                        <n-button @click="onUpdateProjectType(projectType, index)">
                            {{ t("Update") }}
                            <template #icon>
                                <IconEdit :size="22" />
                            </template>
                        </n-button>
                        <n-button @click="onDeleteProjectType(projectType, index)">
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
    <n-table size="small" v-if="false">
        <caption class="table-caption">
            <n-grid :cols="2" align="center">
                <n-grid-item style="text-align: left;">
                    <span class="table-caption-title">{{ t("Project types") }}</span>
                </n-grid-item>
                <n-grid-item style="display: flex; justify-content: flex-end;">
                    <n-button @click="onAddProjectType" :disabled="state.ajaxRunning">
                        <template #icon>
                            <IconPlus />
                        </template>
                        {{ t("Add") }}
                    </n-button>
                </n-grid-item>
            </n-grid>
        </caption>
        <thead>
            <tr>
                <th>{{ t("Name") }}</th>
                <th class="text-center">{{ t("Actions") }}</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="projectType, index in projectTypes" :key="projectType.id">
                <td><n-tag :color="getNaiveUITagColorProperty(projectType.hexColor)">{{ projectType.name
                        }}</n-tag></td>
                <td class="text-center">
                    <n-button-group>
                        <n-button @click="onUpdateProjectType(projectType, index)">
                            {{ t("Update") }}
                            <template #icon>
                                <IconEdit :size="22" />
                            </template>
                        </n-button>
                        <n-button @click="onDeleteProjectType(projectType, index)">
                            {{ t("Delete") }}
                            <template #icon>
                                <IconTrash :size="22" />
                            </template>
                        </n-button>
                    </n-button-group>
                </td>
            </tr>
        </tbody>
        <tfoot ref="tableFooter"></tfoot>
    </n-table>
</template>

<style lang="css" scoped>
    .table-caption {
        padding-bottom: 4px;
    }

    .table-caption-title {
        font-weight: 700;
        font-size: large;

    }
</style>