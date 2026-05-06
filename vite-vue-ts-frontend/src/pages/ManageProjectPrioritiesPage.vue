<script setup lang="ts">
    import { ref, reactive, shallowRef, onMounted, computed } from 'vue'
    import { useI18n } from "vue-i18n";
    import { NButton, NButtonGroup, NModal, NTag, NInput, NIcon } from 'naive-ui'
    import { api } from '../composables/api';
    import { IconEdit, IconPlus, IconTrash, IconSearch, IconArrowUp, IconArrowDown } from '@tabler/icons-vue';
    import { getNaiveUITagColorProperty } from '../composables/color';
    import type { ProjectPriorityInterface } from '../types/models/projectPriority';
    import type { SearchProjectPrioritiesResponse } from '../types/apiResponses';
    import { type AjaxStateInterface, defaultAjaxState } from '../types/ajaxState';
    import { type EntityAction } from '../types/common';
    import { useLoadingStore } from '../stores/loading';
    import { useNotify } from '../composables/notification';
    import { default as ProjectPriorityForm } from '../components/forms/ProjectTypeForm.vue';
    import { default as ManageTable } from '../components/custom/ManageTable.vue';

    const { notify } = useNotify();

    const { t } = useI18n();

    const loadingStore = useLoadingStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const projectPriorities = shallowRef<ProjectPriorityInterface[]>([]);

    const onRefresh = () => {
        state.ajaxRunning = true;
        loadingStore.set(true);
        api.projectPriorities.search().then((successResponse: SearchProjectPrioritiesResponse) => {
            projectPriorities.value = [...successResponse.data.projectPriorities];
        }).catch((errorResponse: any) => {
            // TODO
            console.log(errorResponse);
        }).finally(() => {
            loadingStore.set(false);
            state.ajaxRunning = false;
        })
    };

    const selectedProjectPriorityId = ref<string | undefined>(undefined);

    const onAddProjectPriority = () => {
        actionDialogMode.value = "add";
    };

    const onUpdateProjectPriority = (_projectType: ProjectPriorityInterface, _index: number) => {
        actionDialogMode.value = "update";
        selectedProjectPriorityId.value = _projectType.id;
    };
    const onDeleteProjectPriority = (_projectType: ProjectPriorityInterface, _index: number) => {
        actionDialogMode.value = "delete";
        selectedProjectPriorityId.value = _projectType.id;
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
        notify('success', t("Project priority added"))
        onRefresh();
    };

    const onUpdate = () => {
        isVisibleActionDialog.value = false;
        notify('success', t("Project priority updated"))
        onRefresh();
    };

    const onDelete = () => {
        isVisibleActionDialog.value = false;
        notify('success', t("Project priority deleted"))
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
        <ProjectPriorityForm :mode="actionDialogMode" :project-type-id="selectedProjectPriorityId" style="width: 40%;"
            @add="onAdd" @update="onUpdate" @delete="onDelete" @cancel="onCancel" />
    </n-modal>

    <ManageTable size="small" title="Manage project priorities">
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
                    <n-button size="small" block @click="onAddProjectPriority">
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
            <tr v-for="projectPriority, index in projectPriorities" :key="projectPriority.id">
                <td>
                    <n-tag :color="getNaiveUITagColorProperty(projectPriority.hexColor)">{{ projectPriority.name
                        }}</n-tag>
                </td>
                <td class="text-center">
                    <n-button-group>
                        <n-button size="small">
                            Up
                            <template #icon>
                                <IconArrowUp />
                            </template>
                        </n-button>
                        <n-button size="small">
                            Down
                            <template #icon>
                                <IconArrowDown />
                            </template>
                        </n-button>
                        <n-button size="small" @click="onUpdateProjectPriority(projectPriority, index)">
                            {{ t("Update") }}
                            <template #icon>
                                <IconEdit :size="22" />
                            </template>
                        </n-button>
                        <n-button size="small" @click="onDeleteProjectPriority(projectPriority, index)">
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