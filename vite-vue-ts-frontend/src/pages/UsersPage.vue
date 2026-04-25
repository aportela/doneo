<script setup lang="ts">
    import { ref, onMounted, shallowRef, reactive } from 'vue';
    import { api } from '../composables/api';
    import { IconUser, IconUserKey } from '@tabler/icons-vue';

    interface UserInterface {
        id: string;
        name: string;
        email: string;
        isAdministrator: boolean;
        createdAt: number;
        lastUpdateAt: number;
        avatar: string;
    };

    class User implements UserInterface {
        id: string;
        name: string;
        email: string;
        isAdministrator: boolean;
        createdAt: number;
        lastUpdateAt: number;
        avatar: string;

        constructor(item: UserInterface) {
            this.id = item.id;
            this.name = item.name;
            this.email = item.email;
            this.isAdministrator = item.isAdministrator;
            this.createdAt = item.createdAt;
            this.lastUpdateAt = item.lastUpdateAt;
            this.avatar = item.avatar;
        }
    }

    const users = shallowRef<User[]>([]);

    onMounted(() => {
        api.user.search().then((successResponse: any) => {
            users.value = successResponse.data;
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
                    <h3 class="card-title">Users</h3>
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
                                    Type
                                    <!-- Download SVG icon from http://tabler.io/icons/icon/chevron-up -->
                                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                                        fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                                        stroke-linejoin="round" class="icon icon-sm icon-thick icon-2">
                                        <path d="M6 15l6 -6l6 6"></path>
                                    </svg>
                                </th>
                                <th class="text-center">Avatar</th>
                                <th>Name</th>
                                <th>Email</th>
                                <th>Created at</th>
                                <th>Last update</th>
                                <th></th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="user in users" v-bind:key="user.id">
                                <td><span class="text-secondary">
                                        {{ user.isAdministrator ? 'Administrator' : 'User'
                                        }}
                                    </span></td>
                                <td class="text-center">
                                    <span class="avatar" :style="'background-image: url(' + user.avatar + ')'"></span>
                                </td>
                                <td><span class="text-secondary">{{ user.name }}</span></td>
                                <td>{{ user.email }}</td>
                                <td>{{ new Date(user.createdAt).toDateString() }}</td>
                                <td>{{ user.lastUpdateAt ? new Date(user.lastUpdateAt).toDateString() : null }}</td>
                                <td class="text-end">
                                    <span class="dropdown">
                                        <button class="btn dropdown-toggle align-text-top" data-bs-boundary="viewport"
                                            data-bs-toggle="dropdown">Actions</button>
                                        <div class="dropdown-menu dropdown-menu-end">
                                            <a class="dropdown-item" href="#"> Action </a>
                                            <a class="dropdown-item" href="#"> Another action </a>
                                        </div>
                                    </span>
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

<style lang="css"></style>