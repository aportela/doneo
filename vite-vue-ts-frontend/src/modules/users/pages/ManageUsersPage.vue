<script setup lang="ts">
    import { useI18n } from "vue-i18n";

    import { NCard } from "naive-ui";
    import UsersTable from '../components/UsersTable.vue';

    const { t } = useI18n();
</script>

<template>
    <n-card :title="t('modules.user.components.ManageUsersPage.header.title')">
        <UsersTable />
    </n-card>
</template>

<style lang="css" scoped></style>