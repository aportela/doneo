<script setup lang="ts">
    import { onMounted, h, ref, shallowRef } from 'vue';
    import { api } from '../composables/api';
    import { NDataTable, NTag, NTable, NAvatar } from 'naive-ui';
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

    class ProjectPriority implements ProjectPriorityInterface {
        id: string;
        name: string;
        index: number;
        hexColor: string;
        constructor(item: ProjectPriority) {
            this.id = item.id;
            this.name = item.name;
            this.index = item.index;
            this.hexColor = item.hexColor;
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

    const hexToRgba = (hex: string, alphaOverride?: number) => {
        if (!hex) return `rgba(0,0,0,1)`

        let h = hex.replace('#', '')

        let r, g, b, a = 1

        if (h.length === 8) {
            // RRGGBBAA
            r = parseInt(h.slice(0, 2), 16)
            g = parseInt(h.slice(2, 4), 16)
            b = parseInt(h.slice(4, 6), 16)
            a = parseInt(h.slice(6, 8), 16) / 255
        } else {
            // RRGGBB
            r = parseInt(h.slice(0, 2), 16)
            g = parseInt(h.slice(2, 4), 16)
            b = parseInt(h.slice(4, 6), 16)
        }

        const alpha = alphaOverride ?? a

        return `rgba(${r}, ${g}, ${b}, ${alpha})`
    }

    const tagColor = (base: string) => {
        console.log(base);
        console.log(hexToRgba(base, 1));
        console.log({
            color: hexToRgba(base, 0.2),
            textColor: hexToRgba(base, 1),
            borderColor: hexToRgba(base, 0.5)
        });
        return {
            color: hexToRgba(base, 0.2),
            textColor: hexToRgba(base, 1),
            borderColor: hexToRgba(base, 0.5)
        }
    }
</script>

<template>
    <h1>Manage projects</h1>

    <n-table :bordered="true" size="small" :striped="false" v-if="simpleTable">
        <thead>
            <tr>
                <th style="text-align: center">Task Key</th>
                <th>Summary</th>
                <th>Type</th>
                <th>Creator</th>
                <th>Created at</th>
                <th>Priority</th>
                <th>Status</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="project in projects" :key="project.id">
                <td style="text-align: center"><router-link :to="{ name: 'project', params: { id: project.id } }">{{
                    project.key }}</router-link>
                </td>
                <td><router-link :to="{ name: 'project', params: { id: project.id } }">{{ project.summary
                }}</router-link></td>
                <td>{{ project.type.name }}</td>
                <td><n-avatar :src="'https://i.pravatar.cc/32?u=' + project.createdBy.id"></n-avatar><span
                        class="creator-name">{{
                            project.createdBy.name }}</span></td>
                <td>{{ new Date(project.createdAt).toLocaleString() }}</td>
                <td><n-tag :color="tagColor(project.priority.hexColor)" class="clickable_tag"
                        title="Filter by this value">{{
                            project.priority.name
                        }}</n-tag>
                </td>
                <td><n-tag :color="tagColor(project.status.hexColor)" class="clickable_tag"
                        title="Filter by this value">{{ project.status.name
                        }}</n-tag></td>
            </tr>
        </tbody>
    </n-table>
    <n-data-table size="small" :columns="columns" :data="projects" :pagination="pagination" :bordered="false"
        :loading="loading" :style="{ height: `80vh` }" flex-height v-else />
</template>

<style lang="css" scoped>

    .creator-name {
        margin-left: 8px;
        vertical-align: bottom;
    }

    .clickable_tag {
        cursor: pointer;
    }
</style>