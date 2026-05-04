<script setup lang="ts">
    import { useI18n } from "vue-i18n";
    import { NButton, NTooltip } from 'naive-ui';
    import { IconLayoutSidebarLeftExpand, IconLayoutNavbarExpand } from '@tabler/icons-vue';
    import { useUserSettingsStore } from '../../stores/userSettings';

    const { t } = useI18n();

    interface SwitchFluidLayoutButtonProps {
        iconSize?: number,
    };

    withDefaults(defineProps<SwitchFluidLayoutButtonProps>(), {
        iconSize: 20
    });

    const userSettingsStore = useUserSettingsStore();
</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button quaternary @click="userSettingsStore.toggleNavigationMode" @mousedown.prevent>
                <IconLayoutNavbarExpand :size="iconSize" v-if="userSettingsStore.sideNavigationMode" />
                <IconLayoutSidebarLeftExpand :size="iconSize" v-else />
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