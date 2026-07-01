<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, NTooltip } from 'naive-ui';
    import { IconLayoutSidebarLeftExpand, IconLayoutNavbarExpand } from '@tabler/icons-vue';

    import { useUserSettingsStore } from "../../../stores/userSettings";

    import { BUTTON_DEFAULT_ICON_SIZE } from '../../../constants';

    interface IProps {
        iconSize?: number,
    };

    const props = withDefaults(defineProps<IProps>(), {
        iconSize: BUTTON_DEFAULT_ICON_SIZE
    });

    const { t } = useI18n();
    const userSettingsStore = useUserSettingsStore();
</script>

<template>
    <n-tooltip trigger="hover">
        <template #trigger>
            <n-button quaternary @click="userSettingsStore.toggleNavigationMode" @mousedown.prevent>
                <n-icon :size="props.iconSize"
                    :component="userSettingsStore.sideNavigationMode ? IconLayoutNavbarExpand : IconLayoutSidebarLeftExpand" />
            </n-button>
        </template>
        {{
            t(userSettingsStore.sideNavigationMode ?
                "shared.components.buttons.navigationMode.topNavigation.toolTip" :
                "shared.components.buttons.navigationMode.sideNavigation.toolTip")
        }}
    </n-tooltip>
</template>

<style lang="css" scoped></style>