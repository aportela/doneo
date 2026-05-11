<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, NTooltip } from 'naive-ui';
    import { IconLayoutSidebarLeftExpand, IconLayoutNavbarExpand } from '@tabler/icons-vue';

    import { useUserSettingsStore } from '../../stores/userSettings';


    interface SwitchFluidLayoutButtonProps {
        iconSize?: number,
    };

    withDefaults(defineProps<SwitchFluidLayoutButtonProps>(), {
        iconSize: 20
    });

    const { t } = useI18n();

    const userSettingsStore = useUserSettingsStore();
</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button quaternary @click="userSettingsStore.toggleNavigationMode" @mousedown.prevent>
                <n-icon :size="iconSize"
                    :component="userSettingsStore.sideNavigationMode ? IconLayoutNavbarExpand : IconLayoutSidebarLeftExpand" />
            </n-button>
        </template>
        {{
            t(userSettingsStore.sideNavigationMode ?
                "Switch to top navigation" :
                "Switch to side navigation")
        }}
    </n-tooltip>
</template>

<style lang="css" scoped></style>