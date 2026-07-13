<script setup lang="ts">

    import { computed } from 'vue';
    //import { useI18n } from "vue-i18n";

    import { NMenu } from 'naive-ui';
    import { buildMenuOptions } from '../../types/menu/mapper';
    import { sidebarMenuConfig } from './sidebar-menu-config';

    //const { t } = useI18n();

    const emit = defineEmits(['update'])
    const menuOptions = computed(() => buildMenuOptions(sidebarMenuConfig))

    function handleMenuSelect(key: string) {
        emit("update", key);
    }

</script>

<template>
    <n-menu :options="menuOptions" @update:value="handleMenuSelect" :indent="24" />
</template>

<style lang="css" scoped></style>