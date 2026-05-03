<script setup lang="ts">
    import { useI18n } from "vue-i18n";
    import { NButton, NTooltip } from 'naive-ui';
    import { IconLayoutSidebarLeftExpand, IconLayoutNavbarExpand } from '@tabler/icons-vue';
    import { useUserSettingsStore } from '../../stores/userSettings';

    const { t } = useI18n();

    interface SwitchFluidLayoutButtonProps {
        size?: number,
    };

    withDefaults(defineProps<SwitchFluidLayoutButtonProps>(), {
        size: 20
    });

    const userSettingsStore = useUserSettingsStore();
</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button quaternary @click="userSettingsStore.toggleFluidLayout" @mousedown.prevent>
                <IconLayoutNavbarExpand :size="size" v-if="userSettingsStore.hasFluidLayout" />
                <IconLayoutSidebarLeftExpand :size="size" v-else />
            </n-button>
        </template>
        {{
            t(userSettingsStore.hasFluidLayout ?
                "Switch to top navigation" :
                "Switch to side navigation")
        }}
    </n-tooltip>
</template>

<style lang="css" scoped></style>