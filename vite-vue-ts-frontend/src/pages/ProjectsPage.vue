<script setup lang="ts">
    import { onMounted, ref, shallowRef } from 'vue';
    import { api } from '../composables/api';
    import { NDataTable } from 'naive-ui';
    import type { DataTableColumns } from 'naive-ui'

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

    interface UserBaseInterface {
        id: string;
        name: string;
    }

    class UserBase implements UserBaseInterface {
        id: string;
        name: string;
        constructor(item: UserBaseInterface) {
            this.id = item.id;
            this.name = item.name;
        }
    }

    interface ProjectInterface {
        id: string;
        key: string;
        type: ProjectType;
        summary: string;
        status: string;
        taskCount: number;
        createdBy: UserBase;
        createdAt: number;
    };

    class Project implements ProjectInterface {
        id: string;
        key: string;
        type: ProjectType;
        summary: string;
        status: string;
        taskCount: number;
        createdBy: UserBase;
        createdAt: number;

        constructor(item: ProjectInterface) {
            this.id = item.id;
            this.key = item.key;
            this.type = item.type;
            this.summary = item.summary;
            this.status = item.status;
            this.taskCount = item.taskCount;
            this.createdBy = item.createdBy;
            this.createdAt = item.createdAt;
        }

    }

    const columns: DataTableColumns<ProjectInterface> = [
        {
            title: 'Key',
            key: 'key'
        },
        {
            title: 'Type',
            key: 'type'
        },
        {
            title: 'Summary',
            key: 'summary',
        },
        {
            title: 'Status',
            key: 'status',
        },
        {
            title: 'Created At',
            key: 'createdAt',
            render(row) {
                return new Date(row.createdAt).toLocaleString()
            }
        },
        {
            title: 'Creator',
            key: 'createdBy',
        },
    ];

    const projects = shallowRef<Project[]>([]);

    const loading = ref<boolean>(false);

    onMounted(() => {
        loading.value = true;
        api.project.search().then((successResponse: any) => {
            projects.value = successResponse.data;
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        }).finally(() => { loading.value = false; })
    });

    const pagination = false as const
</script>

<template>
    <h1>Manage projects</h1>
    <n-data-table size="small" :columns="columns" :data="projects" :pagination="pagination" :bordered="false"
        :loading="loading" :style="{ height: `80vh` }" flex-height />
</template>