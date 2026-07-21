<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NCard } from 'naive-ui';
    import RolesTable from '../components/RolesTable.vue';

    const { t } = useI18n();
</script>

<template>
    <n-card :title="t('modules.role.components.ManageRolesPage.header.title')">
        <RolesTable />
    </n-card>
</template>

<style lang="css"></style>