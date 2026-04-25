<script setup lang="ts">
    import { onMounted, shallowRef } from 'vue';
    import { api } from '../composables/api';
    import { IconEdit, IconTrash } from '@tabler/icons-vue';

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

    const projects = shallowRef<Project[]>([]);
    onMounted(() => {
        api.project.search().then((successResponse: any) => {
            projects.value = successResponse.data;
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        });
    });

</script>

<template>
    <div class="row row-deck row-cards">
        <div class="col-12">
            <div class="card">
                <div class="card-header">
                    <h3 class="card-title">Projects</h3>
                </div>
                <div class="card-body border-bottom py-3">
                    <div class="d-flex">
                        <div class="text-secondary">
                            Show
                            <div class="mx-2 d-inline-block">
                                <input type="text" class="form-control form-control-sm" value="8" size="3"
                                    aria-label="Invoices count">
                            </div>
                            entries
                        </div>
                        <div class="ms-auto text-secondary">
                            Search:
                            <div class="ms-2 d-inline-block">
                                <input type="text" class="form-control form-control-sm" aria-label="Search invoice">
                            </div>
                        </div>
                    </div>
                </div>
                <div class="table-responsive">
                    <table class="table table-selectable card-table table-vcenter text-nowrap datatable">
                        <thead>
                            <tr>
                                <th class="w-1">
                                    KEY
                                    <!-- Download SVG icon from http://tabler.io/icons/icon/chevron-up -->
                                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                                        fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                                        stroke-linejoin="round" class="icon icon-sm icon-thick icon-2">
                                        <path d="M6 15l6 -6l6 6"></path>
                                    </svg>
                                </th>
                                <th>Type</th>
                                <th>Summary</th>
                                <th>Creator</th>
                                <th>Created</th>
                                <!--
                                <th>Last update</th>
                                <th>Status</th>
                                -->
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="project in projects" v-bind:key="project.id">
                                <td><span class="text-secondary">{{ project.key }}</span></td>
                                <td><span class="text-secondary">{{ project.type.name }}</span></td>
                                <td><span class="text-secondary">{{ project.summary }}</span></td>
                                <td>{{ project.createdBy.name }}</td>
                                <td><span class="text-secondary">{{ new Date(project.createdAt).toDateString() }}</span>
                                </td>
                                <!--
                                <td>{{ project.taskCount }}</td>
                                <td><span class="badge bg-success me-1"></span> {{ project.status }}</td>
                                -->
                                <td>
                                    <div class="btn-actions">
                                        <a href="#" class="btn btn-action btn-icon" aria-label="Button">
                                            <IconEdit /> Edit
                                        </a>
                                        <a href="#" class="btn btn-action btn-icon ms-2" aria-label="Button">
                                            <IconTrash /> Remove
                                        </a>
                                    </div>
                                </td>
                            </tr>

                        </tbody>
                    </table>
                </div>
                <div class="card-footer">
                    <div class="row g-2 justify-content-center justify-content-sm-between">
                        <div class="col-auto d-flex align-items-center">
                            <p class="m-0 text-secondary">Showing <strong>1 to 8</strong> of <strong>16 entries</strong>
                            </p>
                        </div>
                        <div class="col-auto">
                            <ul class="pagination m-0 ms-auto">
                                <li class="page-item disabled">
                                    <a class="page-link" href="#" tabindex="-1" aria-disabled="true">
                                        <!-- Download SVG icon from http://tabler.io/icons/icon/chevron-left -->
                                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                            viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                                            stroke-linecap="round" stroke-linejoin="round" class="icon icon-1">
                                            <path d="M15 6l-6 6l6 6"></path>
                                        </svg>
                                    </a>
                                </li>
                                <li class="page-item">
                                    <a class="page-link" href="#">1</a>
                                </li>
                                <li class="page-item">
                                    <a class="page-link" href="#">2</a>
                                </li>
                                <li class="page-item active">
                                    <a class="page-link" href="#">3</a>
                                </li>
                                <li class="page-item">
                                    <a class="page-link" href="#">4</a>
                                </li>
                                <li class="page-item">
                                    <a class="page-link" href="#">5</a>
                                </li>
                                <li class="page-item">
                                    <a class="page-link" href="#">
                                        <!-- Download SVG icon from http://tabler.io/icons/icon/chevron-right -->
                                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                            viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                                            stroke-linecap="round" stroke-linejoin="round" class="icon icon-1">
                                            <path d="M9 6l6 6l-6 6"></path>
                                        </svg>
                                    </a>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>