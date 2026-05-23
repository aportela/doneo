<script setup lang="ts">

    import { useI18n } from "vue-i18n";

    import { NButtonGroup, NButton, NIcon } from 'naive-ui';
    import { IconRefresh, IconPlus } from '@tabler/icons-vue';

    interface RefreshAddActionsColumnProps {
        disabled?: boolean;
        iconSize?: number;
    }

    const emit = defineEmits(['refresh', 'add'])

    const props = withDefaults(defineProps<RefreshAddActionsColumnProps>(), {
        disabled: false,
        iconSize: 22,
    });

    const { t } = useI18n();

    const onRefresh = () => {
        emit("refresh");
    };

    const onAdd = () => {
        emit("add");
    };
</script>

<template>
    <n-button-group size="small">
        <n-button @click="onRefresh">
            <template #icon>
                <n-icon :size="22">
                    <IconRefresh />
                </n-icon>
            </template>
            {{ t("shared.buttons.Refresh.label") }}
        </n-button>
        <n-button @click="onAdd">
            <template #icon>
                <n-icon :size="22">
                    <IconPlus />
                </n-icon>
            </template>
            {{ t("shared.buttons.Add.label") }}
        </n-button>
    </n-button-group>
</template>

<style lang="css" scoped></style>