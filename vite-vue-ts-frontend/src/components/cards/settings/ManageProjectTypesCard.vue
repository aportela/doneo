<script setup lang="ts">
    import { ref, onMounted, nextTick } from 'vue'
    import { NSpin, NTable, NButton, NGrid, NGridItem, NFlex, useDialog, NModal } from 'naive-ui'
    //import { v7 as uuidv7 } from 'uuid';
    import { api } from '../../../composables/api';
    import { IconDeviceFloppy, IconEdit, IconPlus, IconTrash } from '@tabler/icons-vue';
    import ProjectTypeForm from '../../forms/ProjectTypeForm.vue';

    interface ProjectTypeInterface {
        id: string;
        name: string;
    }

    class ProjectType implements ProjectTypeInterface {
        id: string;
        name: string;
        constructor(item: ProjectTypeInterface) {
            this.id = item.id;
            this.name = item.name;
        }
    }

    const tableFooter = ref<HTMLElement | null>(null);
    console.log(tableFooter);
    nextTick(() => { });

    const projectTypes = ref<ProjectType[]>([]);
    const loading = ref<boolean>(false);

    onMounted(() => {
        loading.value = true;
        api.projectTypes.search().then((successResponse: any) => {
            projectTypes.value = successResponse.data.projectTypes;
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        }).finally(() => { loading.value = false; })
    });

    const dialog = useDialog()

    /*
    const onAddProjectType = () => {
        projectTypes.value.push(
            new ProjectType({ id: uuidv7(), name: "" })
        );
        nextTick(() => {
            if (tableFooter.value) {
                tableFooter.value?.scrollIntoView({
                    behavior: 'smooth',
                    block: 'end'
                });
            }
        });
    };
    */

    const onRemoveProjectType = (projectType: ProjectType, index: number) => {
        dialog.warning({
            title: 'Please confirm: ',
            content: `Are you sure you want to remove project type: ${projectType.name}?`,
            positiveText: 'Sure',
            negativeText: 'Cancel',
            draggable: true,
            onPositiveClick: () => {
                projectTypes.value.splice(index, 1);
            },
            onNegativeClick: () => {
            }
        })
    }
    const showAddDialog = ref<boolean>(false);
</script>

<template>
    <n-spin :show="loading">
        <n-modal v-model:show="showAddDialog" @cancel="showAddDialog = false">
            <ProjectTypeForm />
        </n-modal>
        <n-table size="small">
            <caption class="table-caption">
                <n-grid :cols="2" align="center">
                    <n-grid-item style="text-align: left;">
                        <span class="table-caption-title">Project types</span>
                    </n-grid-item>
                    <n-grid-item style="display: flex; justify-content: flex-end;">
                        <n-flex>
                            <n-button @click="showAddDialog = true">
                                <template #icon>
                                    <IconPlus />
                                </template>
                                Add new
                            </n-button>
                            <n-button disabled>
                                <template #icon>
                                    <IconDeviceFloppy />
                                </template>
                                Save
                            </n-button>
                        </n-flex>
                    </n-grid-item>
                </n-grid>
            </caption>
            <thead>
                <tr>
                    <th>Name</th>
                    <th class="text-center column-action">Actions</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="projectType, index in projectTypes" :key="projectType.id">
                    <td>{{ projectType.name }}</td>
                    <td class="text-center">
                        <n-flex>
                            <n-button @click="onRemoveProjectType(projectType, index)">
                                Update
                                <template #icon>
                                    <IconEdit :size="22" />
                                </template>
                            </n-button>
                            <n-button @click="onRemoveProjectType(projectType, index)">
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

    .column-action {
        width: 12%;
        white-space: nowrap;
    }

    .text-center {
        text-align: center;
    }
</style>