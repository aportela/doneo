<script setup lang="ts">
    import { h, onMounted, ref, shallowRef } from 'vue';
    import { api } from '../composables/api';
    import { NDataTable, NAvatar } from 'naive-ui';
    import type { DataTableColumns } from 'naive-ui'
    import { IconUser, IconUserKey } from '@tabler/icons-vue';

    interface UserInterface {
        id: string;
        name: string;
        email: string;
        isSuperUser: boolean;
        createdAt: number;
        updatedAt: number;
        deletedAt: number;
        avatar: string;
    };

    class User implements UserInterface {
        id: string;
        name: string;
        email: string;
        isSuperUser: boolean;
        createdAt: number;
        updatedAt: number;
        deletedAt: number;
        avatar: string;

        constructor(item: UserInterface) {
            this.id = item.id;
            this.name = item.name;
            this.email = item.email;
            this.isSuperUser = item.isSuperUser;
            this.createdAt = item.createdAt;
            this.updatedAt = item.updatedAt;
            this.deletedAt = item.deletedAt;
            this.avatar = item.avatar;
        }
    }

    const columns: DataTableColumns<UserInterface> = [
        {
            title: 'Avatar',
            key: 'avatar',
            render(row) {
                if (row.avatar) {
                    return h(NAvatar, { src: row.avatar })
                } else {
                    return null
                }
            },
        },
        {
            title: 'Name',
            key: 'name',
            sorter: 'default',
        },
        {
            title: 'Email',
            key: 'email',
            sorter: 'default',
        },
        {
            title: 'Type',
            key: 'type',
            render(row) {
                if (row.isSuperUser) {
                    return h(IconUserKey, { color: "red" });
                } else {
                    return h(IconUser, { color: "green" });
                }
            },
            sorter: 'default',
            defaultFilterOptionValues: ['admin', 'user'],
            filterOptions: [
                {
                    label: 'Administrator',
                    value: 'admin'
                },
                {
                    label: 'User',
                    value: 'user'
                }
            ],
            filter(value, row) {
                if (value === 'admin') return row.isSuperUser === true
                if (value === 'user') return row.isSuperUser === false
                return false
            }
        },
        {
            title: 'Created At',
            key: 'createdAt',
            render(row) {
                return new Date(row.createdAt).toLocaleString()
            },
            sorter: (a, b) => {
                if (a.createdAt === null && b.createdAt === null) return 0
                if (a.createdAt === null) return 1
                if (b.createdAt === null) return -1
                return a.createdAt - b.createdAt
            }
        },
        {
            title: 'Updated at',
            key: 'updatedAt',
            render(row) {
                if (row.updatedAt) {
                    return new Date(row.updatedAt).toLocaleString()
                } else {
                    return null;
                }
            },
            sorter: (a, b) => {
                if (a.updatedAt === null && b.updatedAt === null) return 0
                if (a.updatedAt === null) return 1
                if (b.updatedAt === null) return -1
                return a.updatedAt - b.updatedAt
            },
            defaultFilterOptionValues: ['updated', 'notUpdated'],
            filterOptions: [
                {
                    label: 'Updated',
                    value: 'updated'
                },
                {
                    label: 'Not updated',
                    value: 'notUpdated'
                }
            ],
            filter(value, row) {
                if (value === 'updated') return row.updatedAt !== null
                if (value === 'notUpdated') return row.updatedAt === null
                return false
            }
        },
        {
            title: 'Deleted at',
            key: 'deletedAt',
            render(row) {
                if (row.deletedAt) {
                    return new Date(row.deletedAt).toLocaleString()
                } else {
                    return null;
                }
            },
            sorter: (a, b) => {
                if (a.deletedAt === null && b.deletedAt === null) return 0
                if (a.deletedAt === null) return 1
                if (b.deletedAt === null) return -1
                return a.deletedAt - b.deletedAt
            },
            defaultFilterOptionValues: ['deleted', 'notDeleted'],
            filterOptions: [
                {
                    label: 'Deleted',
                    value: 'deleted'
                },
                {
                    label: 'Not deleted',
                    value: 'notDeleted'
                }
            ],
            filter(value, row) {
                if (value === 'deleted') return row.deletedAt !== null
                if (value === 'notDeleted') return row.deletedAt === null
                return false
            }
        }
    ];

    const users = shallowRef<User[]>([]);

    const loading = ref<boolean>(false);

    onMounted(() => {
        loading.value = true;
        api.user.search().then((successResponse: any) => {
            users.value = successResponse.data.users;
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        }).finally(() => { loading.value = false; })
    });

    const pagination = false as const

</script>

<template>
    <h1>Manage users</h1>
    <n-data-table size="small" :columns="columns" :data="users" :pagination="pagination" :bordered="false"
        :loading="loading" :style="{ height: `80vh` }" flex-height />
</template>

<style lang="css"></style>