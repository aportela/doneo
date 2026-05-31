<script setup lang="ts">
    import { ref, type CSSProperties } from "vue";
    import { useI18n } from "vue-i18n";

    import { NCard, NButton } from "naive-ui";

    import { UserBase } from "../../users/models/user.ts";
    import ManageTable from '../../../shared/components/tables/ManageTable.vue';
    import ManageTableActionButtons from '../../../shared/components/tables/ManageTableActionButtons.vue';
    import AvatarUserName from "../../../shared/components/AvatarUserName.vue";
    import UploadDialog from "../../attachments/components/UploadDialog.vue";

    interface ProjectAttachmentsProps {
        style?: string | CSSProperties;
        projectId: string | null;
    }

    const props = defineProps<ProjectAttachmentsProps>();

    const user = new UserBase({ id: '019e5b68-f701-7af7-a694-c8146e3366a1', name: 'John Doe' });

    const { t } = useI18n();

    const itemCount = defineModel<number>("itemCount", { default: 0 });
    itemCount.value = 0;

    const showUploadDialog = ref<boolean>(false);
</script>

<template>

    <UploadDialog v-model:show="showUploadDialog" />
    <n-button @click="showUploadDialog = true">add attachment</n-button>
    <n-card bordered :style="props.style">
        <ManageTable size="small">
            <template #thead>
                <tr>
                    <th>Filename</th>
                    <th>Size</th>
                    <th>Mime</th>
                    <th>Uploaded at</th>
                    <th>Uploaded by</th>
                    <th class="doneo-table-actions-column">{{
                        t("shared.components.table.header.columns.actions") }}</th>
                </tr>
            </template>
            <template #tbody>
                <tr v-for="index in [1, 2, 3, 4, 5]" :key="index">
                    <td>
                        test.pdf
                    </td>
                    <td>8230Kb</td>
                    <td>application/pdf</td>
                    <td>23/02/2026</td>
                    <td>
                        <AvatarUserName :user-id="user.id" :user-name="user.name" />
                    </td>
                    <td class="doneo-text-center">
                        <ManageTableActionButtons show-delete show-download />
                    </td>
                </tr>
            </template>
        </ManageTable>
    </n-card>
</template>

<style lang="css" scoped></style>