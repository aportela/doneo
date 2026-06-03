<script setup lang="ts">

    import { useI18n } from "vue-i18n";

    import { NButtonGroup, NButton, NIcon } from 'naive-ui';
    import { IconRefresh, IconPlus, IconSettings } from '@tabler/icons-vue';

    interface RefreshAddActionsColumnProps {
        disabled?: boolean;
        iconSize?: number;
        hideRefresh?: boolean;
        hideAdd?: boolean;
        hideSettings?: boolean;

    }

    const emit = defineEmits(['refresh', 'add', 'settings'])

    const props = withDefaults(defineProps<RefreshAddActionsColumnProps>(), {
        disabled: false,
        iconSize: 22,
        hideRefresh: false,
        hideAdd: false,
        hideSettings: false,
    });

    const { t } = useI18n();

    const onRefresh = () => {
        emit("refresh");
    };

    const onAdd = () => {
        emit("add");
    };

    const onSettings = () => {
        emit("settings");
    };
</script>

<template>
    <n-button-group size="small" style="width: 100%">
        <n-button @click="onRefresh" v-if="!props.hideRefresh" style="flex: 1;">
            <template #icon>
                <n-icon :size="22" :component="IconRefresh" />
            </template>
            {{ t("shared.buttons.Refresh.label") }}
        </n-button>
        <n-button @click="onAdd" v-if="!props.hideAdd" style="flex: 1;">
            <template #icon>
                <n-icon :size="22" :component="IconPlus" />
            </template>
            {{ t("shared.buttons.Add.label") }}
        </n-button>
        <n-button @click="onSettings" v-if="!props.hideSettings" style="flex: 1;">
            <template #icon>
                <n-icon :size="22" :component="IconSettings" />
            </template>
            {{ t("shared.buttons.Settings.label") }}
        </n-button>
    </n-button-group>
</template>

<style lang="css" scoped></style>