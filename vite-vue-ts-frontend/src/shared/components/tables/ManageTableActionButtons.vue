<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButtonGroup, NButton, NIcon } from 'naive-ui';
    import { IconEdit, IconTrash, IconDownload } from '@tabler/icons-vue';


    interface UpdateDeleteActionsColumnProps {
        disabled?: boolean;
        iconSize?: number;
        showUpdate?: boolean;
        updateDisabled?: boolean;
        deleteDisabled?: boolean;
        showDelete?: boolean;
        showDownload?: boolean;
        downloadDisabled?: boolean;

    }

    const emit = defineEmits(['update', 'delete', 'donwload'])

    const props = withDefaults(defineProps<UpdateDeleteActionsColumnProps>(), {
        disabled: false,
        iconSize: 22,
        showUpdate: false,
        updateDisabled: false,
        showDelete: false,
        deleteDisabled: false,
        showDownload: false,
        downloadDisabled: false,
    });

    const { t } = useI18n();

    const onUpdate = () => {
        emit("update");
    };

    const onDelete = () => {
        emit("delete");
    };

    const onDownload = () => {
        emit("donwload");
    };
</script>

<template>
    <n-button-group size="small">
        <n-button @click="onUpdate" :disabled="props.disabled || props.updateDisabled" v-if="showUpdate">
            {{ t("shared.buttons.Update.label") }}
            <template #icon>
                <n-icon :size="props.iconSize">
                    <IconEdit />
                </n-icon>
            </template>
        </n-button>
        <n-button @click="onDelete" :disabled="props.disabled || props.deleteDisabled" v-if="showDelete">
            {{ t("shared.buttons.Delete.label") }}
            <template #icon>
                <n-icon :size="props.iconSize">
                    <IconTrash />
                </n-icon>
            </template>
        </n-button>
        <n-button @click="onDownload" :disabled="props.disabled || props.downloadDisabled" v-if="showDownload">
            {{ t("shared.buttons.Download.label") }}
            <template #icon>
                <n-icon :size="props.iconSize">
                    <IconDownload />
                </n-icon>
            </template>
        </n-button>
    </n-button-group>
</template>

<style lang="css" scoped></style>