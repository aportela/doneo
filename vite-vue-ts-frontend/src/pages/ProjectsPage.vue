<script setup lang="ts">
    import { onMounted, h, ref, shallowRef } from 'vue';
    import { api } from '../composables/api';
    import { NDataTable, NTag, NTable } from 'naive-ui';
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

    interface ProjectStatusInterface {
        id: string;
        name: string;
    }

    class ProjectStatus implements ProjectStatusInterface {
        id: string;
        name: string;
        constructor(item: ProjectStatusInterface) {
            this.id = item.id;
            this.name = item.name;
        }
    }

    interface ProjectPriorityInterface {
        id: string;
        name: string;
    }

    class ProjectPriority implements ProjectPriorityInterface {
        id: string;
        name: string;
        constructor(item: ProjectPriority) {
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
        status: ProjectStatus;
        priority: ProjectPriority;
        summary: string;
        taskCount: number;
        createdBy: UserBase;
        createdAt: number;
    };

    class Project implements ProjectInterface {
        id: string;
        key: string;
        type: ProjectType;
        status: ProjectStatus;
        priority: ProjectPriority;
        summary: string;
        taskCount: number;
        createdBy: UserBase;
        createdAt: number;

        constructor(item: ProjectInterface) {
            this.id = item.id;
            this.key = item.key;
            this.summary = item.summary;
            this.type = item.type;
            this.status = item.status;
            this.priority = item.priority;
            this.taskCount = item.taskCount;
            this.createdBy = item.createdBy;
            this.createdAt = item.createdAt;
        }

    }

    const columns: DataTableColumns<ProjectInterface> = [
        {
            title: 'Key',
            key: 'key',
            width: 100,
            minWidth: 100,
        },
        {
            title: 'Type',
            key: 'type',
            render(row) {
                return row.type.name
            }
        },
        {
            title: 'Priority',
            key: 'priority',
            align: 'center',
            render(row) {
                return h(
                    NTag,
                    {
                        style: {
                            marginRight: '6px'
                        },
                        type: 'info',
                        bordered: false
                    },
                    {
                        default: () => row.priority.name
                    }
                )
            }
        },
        {
            title: 'Status',
            key: 'status',
            render(row) {
                return h(
                    NTag,
                    {
                        style: {
                            marginRight: '6px'
                        },
                        type: 'success',
                        bordered: false
                    },
                    {
                        default: () => row.status.name
                    }
                )
            }
        },
        {
            title: 'Summary',
            key: 'summary',
            /*
            width: 1200,
            ellipsis: {
                tooltip: true
            }
            */
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
            render(row) {
                return row.createdBy.name
            }
        },
    ];

    const projects = shallowRef<Project[]>([]);

    const loading = ref<boolean>(false);

    onMounted(() => {
        loading.value = true;
        api.project.search().then((successResponse: any) => {
            projects.value = successResponse.data.projects;
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        }).finally(() => { loading.value = false; })
    });

    const pagination = false as const

    const simpleTable = true;
</script>

<template>
    <h1>Manage projects</h1>

    <n-table :bordered="true" size="small" :striped="false" v-if="simpleTable">
        <thead>
            <tr>
                <th>Key</th>
                <th>Type</th>
                <th>Priority</th>
                <th>Status</th>
                <th>Summary</th>
                <th>Created at</th>
                <th>Creator</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="project in projects" :key="project.id">
                <td>{{ project.key }}</td>
                <td>{{ project.type.name }}</td>
                <td><n-tag type="info">{{ project.priority.name }}</n-tag></td>
                <td><n-tag type="error">{{ project.status.name }}</n-tag></td>
                <td>{{ project.summary }}</td>
                <td>{{ project.createdAt }}</td>
                <td>{{ project.createdBy.name }}</td>
            </tr>
        </tbody>
    </n-table>
    <n-data-table size="small" :columns="columns" :data="projects" :pagination="pagination" :bordered="false"
        :loading="loading" :style="{ height: `80vh` }" flex-height v-else />
</template>