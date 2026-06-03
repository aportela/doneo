<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButtonGroup, NButton, NIcon } from 'naive-ui';
    import { IconEdit, IconTrash, IconTrashOff, IconDownload, IconEyeSearch } from '@tabler/icons-vue';


    interface UpdateDeleteActionsColumnProps {
        disabled?: boolean;
        buttonSize?: ButtonSize;
        iconSize?: number;
        showUpdate?: boolean;
        updateDisabled?: boolean;
        showDelete?: boolean;
        deleteDisabled?: boolean;
        showRestore?: boolean;
        restoredDisabled?: boolean;
        showDownload?: boolean;
        downloadDisabled?: boolean;
        showPreview?: boolean;
        previewDisabled?: boolean;
    }

    const emit = defineEmits(['update', 'delete', 'restore', 'download', 'preview'])

    const props = withDefaults(defineProps<UpdateDeleteActionsColumnProps>(), {
        disabled: false,
        buttonSize: "small",
        iconSize: 22,
        showUpdate: false,
        updateDisabled: false,
        showDelete: false,
        deleteDisabled: false,
        showRestore: false,
        restoredDisabled: false,
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

    const onRestore = () => {
        emit("restore");
    };

    const onDownload = () => {
        emit("download");
    };

    const onPreview = () => {
        emit("preview");
    };
</script>

<template>
    <n-button-group :size="props.buttonSize">
        <n-button @click="onUpdate" :disabled="props.disabled || props.updateDisabled" v-if="showUpdate">
            {{ t("shared.buttons.Update.label") }}
            <template #icon>
                <n-icon :size="props.iconSize" :component="IconEdit" />
            </template>
        </n-button>
        <n-button @click="onDelete" :disabled="props.disabled || props.deleteDisabled" v-if="showDelete">
            {{ t("shared.buttons.Delete.label") }}
            <template #icon>
                <n-icon :size="props.iconSize" :component="IconTrash" />
            </template>
        </n-button>
        <n-button @click="onRestore" :disabled="props.disabled || props.restoredDisabled" v-if="showRestore">
            {{ t("shared.buttons.Restore.label") }}
            <template #icon>
                <n-icon :size="props.iconSize" :component="IconTrashOff" />
            </template>
        </n-button>
        <n-button @click.prevent="onDownload" :disabled="props.disabled || props.downloadDisabled" v-if="showDownload">
            {{ t("shared.buttons.Download.label") }}
            <template #icon>
                <n-icon :size="props.iconSize" :component="IconDownload" />
            </template>
        </n-button>
        <n-button @click.prevent="onPreview" :disabled="props.disabled || props.previewDisabled" v-if="showPreview">
            {{ t("shared.buttons.Preview.label") }}
            <template #icon>
                <n-icon :size="props.iconSize" :component="IconEyeSearch" />
            </template>
        </n-button>
    </n-button-group>
</template>

<style lang="css" scoped></style>