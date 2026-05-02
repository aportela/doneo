<script setup lang="ts">
    import { ref, onMounted, shallowRef } from 'vue'
    import { NDialogProvider, NCard, NTabs, NTabPane } from 'naive-ui'
    import { api } from '../composables/api';
    import { default as ManageProjectTypesCard } from '../components/cards/settings/ManageProjectTypesCard.vue';
    import { default as ManageProjectStatusesCard } from '../components/cards/settings/ManageProjectStatusesCard.vue';
    import { default as ManageProjectPrioritiesCard } from '../components/cards/settings/ManageProjectPrioritiesCard.vue';


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

    interface ProjectStatusInterface {
        id: string;
        name: string;
        index: number;
        hexColor: string;
    }

    class ProjectStatus implements ProjectStatusInterface {
        id: string;
        name: string;
        index: number;
        hexColor: string;
        constructor(item: ProjectStatusInterface) {
            this.id = item.id;
            this.name = item.name;
            this.index = item.index;
            this.hexColor = item.hexColor;
        }
    }

    interface ProjectPriorityInterface {
        id: string;
        name: string;
        index: number;
        hexColor: string;
    }

    const projectTypes = shallowRef<ProjectType[]>([]);
    const projectStatuses = shallowRef<ProjectStatus[]>([]);
    const projectPriorities = shallowRef<ProjectPriorityInterface[]>([]);

    const loading = ref<boolean>(false);

    onMounted(() => {
        loading.value = true;
        api.projectTypes.search().then((successResponse: any) => {
            projectTypes.value = successResponse.data.projectTypes;
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        }).finally(() => { loading.value = false; })
        api.projectStatuses.search().then((successResponse: any) => {
            projectStatuses.value = successResponse.data.projectTypes;
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        }).finally(() => { loading.value = false; })
        api.projectPriorities.search().then((successResponse: any) => {
            projectPriorities.value = successResponse.data.projectTypes;
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        }).finally(() => { loading.value = false; })
    });

</script>

<template>
    <n-card style="margin: 1em;">
        <n-dialog-provider>
            <n-tabs placement="left" type="line" :bar-width="48">
                <n-tab-pane name="Project types" tab="Project types">
                    <ManageProjectTypesCard />
                </n-tab-pane>
                <n-tab-pane name="Project priorities" tab="Project priorities">
                    <ManageProjectPrioritiesCard />
                </n-tab-pane>
                <n-tab-pane name="Project statuses" tab="Project statuses">
                    <ManageProjectStatusesCard />
                </n-tab-pane>
            </n-tabs>
        </n-dialog-provider>
    </n-card>


</template>

<style lang="css" scoped></style>