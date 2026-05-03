<script setup lang="ts">
    import { ref, reactive, onMounted, nextTick, computed } from 'vue'
    import { NSpin, NTable, NButton, NGrid, NGridItem, NFlex, useDialog, NModal, NTag } from 'naive-ui'
    import { api } from '../../../composables/api';
    import { IconEdit, IconPlus, IconTrash } from '@tabler/icons-vue';
    import ProjectTypeForm from '../../forms/ProjectTypeForm.vue';
    import { getNaiveUITagColorProperty } from '../../../composables/color';
    import type { ProjectTypeInterface } from '../../../types/models/projectType';
    import { type AjaxStateInterface, defaultAjaxState } from '../../../types/ajaxState';
    import { type EntityAction } from '../../../types/common';

    const dialog = useDialog();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    if (dialog !== null) {
        dialog.destroyAll();
    }

    const tableFooter = ref<HTMLElement | null>(null);

    const projectTypes = ref<ProjectTypeInterface[]>([]);

    const onRefresh = () => {
        state.ajaxRunning = true;
        api.projectTypes.search().then((successResponse: any) => {
            projectTypes.value = successResponse.data.projectTypes;
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        }).finally(() => {
            state.ajaxRunning = false;
        })
    };

    onMounted(() => {
        onRefresh();
    });

    const selectedProjectTypeId = ref<string | undefined>(undefined);

    const onAddProjectType = () => {
        actionDialogMode.value = "add";
        nextTick(() => {
            if (tableFooter.value) {
                tableFooter.value?.scrollIntoView({
                    behavior: 'smooth',
                    block: 'end'
                });
            }
        });

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
        onRefresh();
    };

    const onUpdate = () => {
        isVisibleActionDialog.value = false;
        onRefresh();
    };

    const onDelete = () => {
        isVisibleActionDialog.value = false;
        onRefresh();
    };

    const onCancel = () => {
        isVisibleActionDialog.value = false;
    };

</script>

<template>
    <n-spin :show="state.ajaxRunning">
        <n-modal v-model:show="isVisibleActionDialog">
            <ProjectTypeForm :mode="actionDialogMode" :project-type-id="selectedProjectTypeId" style="width: 40%;"
                @add="onAdd" @update="onUpdate" @delete="onDelete" @cancel="onCancel" />
        </n-modal>
        <n-table size="small">
            <caption class="table-caption">
                <n-grid :cols="2" align="center">
                    <n-grid-item style="text-align: left;">
                        <span class="table-caption-title">Project types</span>
                    </n-grid-item>
                    <n-grid-item style="display: flex; justify-content: flex-end;">
                        <n-flex>
                            <n-button @click="onAddProjectType" :disabled="state.ajaxRunning">
                                <template #icon>
                                    <IconPlus />
                                </template>
                                Add new
                            </n-button>
                        </n-flex>
                    </n-grid-item>
                </n-grid>
            </caption>
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="projectType, index in projectTypes" :key="projectType.id">
                    <td><n-tag :color="getNaiveUITagColorProperty(projectType.hexColor)">{{ projectType.name
                    }}</n-tag></td>
                    <td class="text-center">
                        <n-flex>
                            <n-button @click="onUpdateProjectType(projectType, index)">
                                Update
                                <template #icon>
                                    <IconEdit :size="22" />
                                </template>
                            </n-button>
                            <n-button @click="onDeleteProjectType(projectType, index)">
                                Delete
                                <template #icon>
                                    <IconTrash :size="22" />
                                </template>
                            </n-button>
                        </n-flex>
                    </td>
                </tr>
            </tbody>
            <tfoot ref="tableFooter"></tfoot>
        </n-table>
    </n-spin>
</template>

<style lang="css" scoped>
    .table-caption {
        padding-bottom: 4px;
    }

    .table-caption-title {
        font-weight: 700;
        font-size: large;

    }

    .text-center {
        text-align: center;
    }
</style>