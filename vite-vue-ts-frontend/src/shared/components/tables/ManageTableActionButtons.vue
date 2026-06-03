<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButtonGroup, NButton, NIcon } from 'naive-ui';
    import { IconEdit, IconTrash, IconDownload, IconEyeSearch } from '@tabler/icons-vue';


    interface UpdateDeleteActionsColumnProps {
        disabled?: boolean;
        iconSize?: number;
        showUpdate?: boolean;
        updateDisabled?: boolean;
        showDelete?: boolean;
        deleteDisabled?: boolean;
        showDownload?: boolean;
        downloadDisabled?: boolean;
        showPreview?: boolean;
        previewDisabled?: boolean;
    }

    const emit = defineEmits(['update', 'delete', 'download', 'preview'])

    const props = withDefaults(defineProps<UpdateDeleteActionsColumnProps>(), {
        disabled: false,
        iconSize: 22,
        showUpdate: false,
        updateDisabled: false,
        showDelete: false,
        deleteDisabled: false,
        showDownload: false,
        downloadDisabled: false,
        showPreview: false,
        previewDisabled: false,
    });

    const { t } = useI18n();

    const onUpdate = () => {
        emit("update");
    };

    const onDelete = () => {
        emit("delete");
    };

    const onDownload = () => {
        emit("download");
    };

    const onPreview = () => {
        emit("preview");
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
        <n-button @click.prevent="onDownload" :disabled="props.disabled || props.downloadDisabled" v-if="showDownload">
            {{ t("shared.buttons.Download.label") }}
            <template #icon>
                <n-icon :size="props.iconSize">
                    <IconDownload />
                </n-icon>
            </template>
        </n-button>
        <n-button @click.prevent="onPreview" :disabled="props.disabled || props.previewDisabled" v-if="showPreview">
            {{ t("shared.buttons.Preview.label") }}
            <template #icon>
                <n-icon :size="props.iconSize">
                    <IconEyeSearch />
                </n-icon>
            </template>
        </n-button>
    </n-button-group>
</template>

<style lang="css" scoped></style>