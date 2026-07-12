<script setup lang="ts">

    import { computed } from 'vue';
    //import { useI18n } from "vue-i18n";

    import { NMenu } from 'naive-ui';
    import { buildMenuOptions } from '../../types/menu/mapper';
    import { sidebarMenuConfig } from './sidebar-menu-config';

    //const { t } = useI18n();

    const menuOptions = computed(() => buildMenuOptions(sidebarMenuConfig))

    function handleMenuSelect(key: string) {
        console.log(key);
    }

</script>

<template>
    <n-menu :options="menuOptions" @update:value="handleMenuSelect" />
</template>

<style lang="css" scoped></style>