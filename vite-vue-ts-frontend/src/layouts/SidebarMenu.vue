<script setup lang="ts">
    import { NMenu } from 'naive-ui';
    import { IconDatabaseStar } from '@tabler/icons-vue';
    import { useRoute } from 'vue-router'
    import { menuOptions } from '../types/menu';

    const route = useRoute()
    const commonIconSize = 20;
</script>

<template>
    <div class="brand-container">
        <IconDatabaseStar :size="32" class="brand-icon" />
        <span class="brand-name">Doneo</span>
    </div>
    <n-menu :collapsed-width="64" :collapsed-icon-size="commonIconSize" :options="menuOptions"
        :value="route.name as string" accordion />
</template>

<style lang="css" scoped>
    .brand-container {
        padding-left: 32px;
        display: flex;
        align-items: center;
    }

    .brand-name {
        font-weight: bold;
        font-size: 30px;
    }

    .brand-icon {
        margin-right: 8px;
    }
</style>