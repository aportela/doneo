<script setup lang="ts">

    // TODO: delete component ?
    import { useI18n } from "vue-i18n";

    import { NButtonGroup, NButton, NIcon } from 'naive-ui';
    import { IconEdit, IconTrash } from '@tabler/icons-vue';

    interface Props {
        disabled?: boolean;
        iconSize?: number;
    }

    const props = withDefaults(defineProps<Props>(), {
        disabled: false,
        iconSize: 22,
    });

    const emit = defineEmits(['update', 'delete'])

    const { t } = useI18n();

    const onUpdate = () => {
        emit("update");
    };

    const onDelete = () => {
        emit("delete");
    };
</script>

<template>
    <n-button-group size="small">
        <n-button @click="onUpdate" :disabled="props.disabled">
            {{ t("shared.buttons.Update.label") }}
            <template #icon>
                <n-icon :size="props.iconSize">
                    <IconEdit />
                </n-icon>
            </template>
        </n-button>
        <n-button @click="onDelete" :disabled="props.disabled">
            {{ t("shared.buttons.Delete.label") }}
            <template #icon>
                <n-icon :size="props.iconSize">
                    <IconTrash />
                </n-icon>
            </template>
        </n-button>
    </n-button-group>
</template>

<style lang="css" scoped></style>