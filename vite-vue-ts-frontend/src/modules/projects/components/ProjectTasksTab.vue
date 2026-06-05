<script setup lang="ts">
    import { computed, type CSSProperties } from "vue";
    //import { useI18n } from "vue-i18n";

    import { NCard } from "naive-ui";

    import { UserBase } from "../../users/models/user.ts";

    import AvatarUserName from "../../../shared/components/AvatarUserName.vue";
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import type { TableHeaderColumn } from "../../../shared/types/table-header-column.ts";

    interface ProjectTasksProps {
        style?: string | CSSProperties;
        projectId: string | null;
    }

    const props = defineProps<ProjectTasksProps>();

    //const { t } = useI18n();

    const user = new UserBase({ id: '019e5b68-f701-7af7-a694-c8146e3366a1', name: 'John Doe' });

    const itemCount = defineModel<number>("itemCount", { default: 0 });
    itemCount.value = 0;

    const columns = computed<TableHeaderColumn[]>(() => [
        {
            label: "Identifier",
            field: "id",
            visible: true,
            sortable: false,
            isFiltered: () => false,
        },
        {
            label: "Summary",
            field: "summary",
            visible: true,
            sortable: false,
            isFiltered: () => false,
        },
        {
            label: "Status",
            field: "status",
            visible: true,
            sortable: false,
            isFiltered: () => false,
        },
        {
            label: "Created at",
            field: "createdAt",
            visible: true,
            sortable: false,
            isFiltered: () => false,
        },
        {
            label: "Created by",
            field: "createdBy",
            visible: true,
            sortable: false,
            isFiltered: () => false,
        },
    ]);

    const onRefresh = () => {

    };

    const onAdd = () => {

    };
</script>

<template>
    <n-card bordered :style="props.style">
        <ManageTable size="small" :columns="columns" @refresh="onRefresh" @add="onAdd">
            <template #thead>
            </template>
            <template #tbody>
                <tr v-for="index in [1, 2, 3, 4, 5, 6, 7, 8, 9]" :key="index">
                    <td>
                        PRJ-000{{ index }}
                    </td>
                    <td>Get things done</td>
                    <td>Pending</td>
                    <td>23/02/2026</td>
                    <td>
                        <AvatarUserName :user-id="user.id" :user-name="user.name" />
                    </td>
                    <td class="doneo-text-center">
                        <ManageTableActionButtons show-update show-delete />
                    </td>
                </tr>
            </template>
        </ManageTable>
    </n-card>
</template>

<style lang="css" scoped></style>