<script setup lang="ts">
    import { onMounted, ref, computed, shallowRef } from 'vue';
    import { useI18n } from "vue-i18n";
    import { api } from '../composables/api';
    import { NAvatar, NInput, NSelect, NIcon, NButton, NButtonGroup } from 'naive-ui';
    import { default as ManageTable } from '../components/custom/ManageTable.vue';
    import { IconUser, IconUserKey, IconSearch, IconPlus, IconEdit, IconTrash } from '@tabler/icons-vue';

    import { type UserInterface, UserClass } from '../types/models/user';

    const { t } = useI18n();
    const users = shallowRef<UserClass[]>([]);

    const loading = ref<boolean>(false);

    onMounted(() => {
        loading.value = true;
        api.user.search().then((successResponse: any) => {
            users.value = successResponse.data.users.map((u: UserInterface) => new UserClass(u));
        }).catch((errorResponse: any) => {
            console.log(errorResponse);
        }).finally(() => { loading.value = false; })
    });


    const filterUserOptions = [
        { label: 'All users', value: 1 },
        { label: 'Administrators only', value: 2 },
        { label: 'Users only', value: 3 }
    ];

    const filterDateOptions = [
        { label: 'Any date', value: 0 },
        { label: 'Today', value: 1 },
        { label: 'Yesterday', value: 2 },
        { label: 'This week', value: 3 }
    ];

    const filterByUsername = ref<string | null>(null);
    const filterByEmail = ref<string | null>(null);
    const userFilterType = ref<number | null>(1);

    const createdAtFilter = ref<number | null>(0);
    const updatedAtFilter = ref<number | null>(0);
    const deletedAtFilter = ref<number | null>(0);

    const searchMappedUsers = computed(() => {
        return users.value.map(u => ({
            ...u,
            _searchName: (u.name).toLowerCase(),
            _searchEmail: (u.email).toLowerCase()
        }));
    })

    const filteredUsers = computed(() => {
        const name = filterByUsername.value?.trim().toLowerCase();
        const email = filterByEmail.value?.trim().toLowerCase();
        return searchMappedUsers.value.filter(u => {
            const matchName = !name || u._searchName.includes(name);
            const matchEmail = !email || u._searchEmail.includes(email);
            return matchName && matchEmail;
        });
    })
</script>

<template>
    <ManageTable size="small" title="Manage users">
        <template #thead>
            <tr>
                <th>Name</th>
                <th>Email</th>
                <th class="text-center">Type</th>
                <th>Created at</th>
                <th>Updated at</th>
                <th>Deleted at</th>
                <th class="text-center">Operations</th>
            </tr>
            <tr>
                <th>
                    <n-input size="small" placeholder="search by name..." v-model:value="filterByUsername" clearable>
                        <template #prefix>
                            <n-icon>
                                <IconSearch />
                            </n-icon>
                        </template>
                    </n-input>
                </th>
                <th>
                    <n-input size="small" placeholder="search by email..." v-model:value="filterByEmail" clearable>
                        <template #prefix>
                            <n-icon>
                                <IconSearch />
                            </n-icon>
                        </template>
                    </n-input>
                </th>
                <th>
                    <n-select size="small" trigger="click" :options="filterUserOptions" v-model:value="userFilterType"
                        placeholder="Search by user type">
                    </n-select>
                </th>
                <th>
                    <n-select size="small" trigger="click" :options="filterDateOptions" v-model:value="createdAtFilter"
                        placeholder="Select date filter">
                    </n-select>
                </th>
                <th>
                    <n-select size="small" trigger="click" :options="filterDateOptions" v-model:value="updatedAtFilter"
                        placeholder="Select date filter">
                    </n-select>
                </th>
                <th>
                    <n-select size="small" trigger="click" :options="filterDateOptions" v-model:value="deletedAtFilter"
                        placeholder="Select date filter">
                    </n-select>
                </th>
                <th class="text-center">
                    <n-button size="small" block>
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
            <tr v-for="user in filteredUsers" :key="user.id">
                <td>
                    <div style="display: flex; align-items: center; gap: 8px;">
                        <n-avatar :src="user.avatar" class="avatar" /> {{ user.name }}
                    </div>
                </td>
                <td>{{ user.email }}</td>
                <td class="text-center">
                    <IconUserKey v-if="user.isSuperUser" color="red" />
                    <IconUser v-else />
                </td>
                <td>{{ user.createdAt ? new Date(user.createdAt).toLocaleString() : null }}</td>
                <td>{{ user.updatedAt ? new Date(user.updatedAt).toLocaleString() : null }}</td>
                <td>{{ user.deletedAt ? new Date(user.deletedAt).toLocaleString() : null }}</td>
                <td class="text-center">
                    <n-button-group>
                        <n-button size="small">
                            {{ t("Update") }}
                            <template #icon>
                                <IconEdit :size="22" />
                            </template>
                        </n-button>
                        <n-button size="small">
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

<style lang="css">
    .avatar {
        margin-right: 4px;
    }
</style>