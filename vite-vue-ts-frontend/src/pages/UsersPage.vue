<script setup lang="ts">
    import { h, onMounted, shallowRef } from 'vue';
    import { api } from '../composables/api';
    import { NDataTable, NAvatar } from 'naive-ui';
    import type { DataTableColumns } from 'naive-ui'

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

    const columns: DataTableColumns<UserInterface> = [
        {
            title: 'Avatar',
            key: 'avatar',
            render(row) { return h(NAvatar, { src: "https://i.pravatar.cc/32?u=" + row.id, },) },
        },
        {
            title: 'Name',
            key: 'name'
        },
        {
            title: 'Email',
            key: 'email'
        },
        {
            title: 'Admin',
            key: 'isAdministrator',
        },
        {
            title: 'Created At',
            key: 'createdAt',
            render(row) {
                return new Date(row.createdAt).toLocaleString()
            }
        },
        {
            title: 'Last Updated',
            key: 'lastUpdateAt',
            render(row) {
                if (row.lastUpdateAt) {
                    return new Date(row.lastUpdateAt).toLocaleString()
                } else {
                    return null;
                }
            }
        }
    ];

    const users = shallowRef<User[]>([]);

    onMounted(() => {
        api.user.search().then((successResponse: any) => {
            users.value = successResponse.data.users;
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        });
    });

    const pagination = false as const

</script>

<template>
    <n-data-table :columns="columns" :data="users" :pagination="pagination" :bordered="false" />
</template>

<style lang="css"></style>