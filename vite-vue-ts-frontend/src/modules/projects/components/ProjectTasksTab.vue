<script setup lang="ts">
    import type { CSSProperties } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard } from "naive-ui";


    import { UserBase } from "../../users/models/user.ts";

    import AvatarUserName from "../../../shared/components/AvatarUserName.vue";
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';

    interface ProjectTasksProps {
        style?: string | CSSProperties;
        projectId: string | null;
    }

    const props = defineProps<ProjectTasksProps>();

    const { t } = useI18n();

    const user = new UserBase({ id: '019e5b68-f701-7af7-a694-c8146e3366a1', name: 'John Doe' });

</script>

<template>
    <n-card bordered :style="props.style">
        <ManageTable size="small">
            <template #thead>
                <tr>
                    <th>Identifier</th>
                    <th>Summary</th>
                    <th>Status</th>
                    <th>Created at</th>
                    <th>Created by</th>
                    <th class="doneo-table-actions-column">{{
                        t("shared.components.table.header.columns.actions") }}</th>
                </tr>
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