<script setup lang="ts">

    import { useRoute } from 'vue-router';
    import { NBreadcrumb, NBreadcrumbItem, NIcon } from 'naive-ui';

    import { useSessionStore } from '../stores/session';

    import { Bookmark, CircleUser, FileCog, FolderCog, FolderKanban, Goal, Home, ListTodo, Notebook, Route, Settings, UserCog, UserKey, Users } from '@lucide/vue';

    const route = useRoute();

    const sessionStore = useSessionStore();
</script>

<template>
    <n-breadcrumb>
        <n-breadcrumb-item><n-icon :component="Home" /> Home</n-breadcrumb-item>

        <n-breadcrumb-item v-if="route.name === 'manageProjects' || route.name === 'manageTasks'" :clickable="false">
            <n-icon :component="Notebook" /> Workspace
        </n-breadcrumb-item>
        <n-breadcrumb-item v-if="route.name === 'manageProjects'">
            <n-icon :component="FolderKanban" /> Projects
        </n-breadcrumb-item>
        <n-breadcrumb-item v-if="route.name === 'manageTasks'">
            <n-icon :component="ListTodo" /> Tasks
        </n-breadcrumb-item>

        <n-breadcrumb-item
            v-if="route.name === 'manageUsers' || route.name === 'manageRoles' || route.name === 'manageProjectTypes' || route.name === 'manageProjectPriorities' || route.name === 'manageProjectStatuses' || route.name === 'manageTaskPriorities' || route.name === 'manageTaskStatuses'"
            :clickable="false">
            <n-icon :component="Settings" /> Settings
        </n-breadcrumb-item>
        <n-breadcrumb-item v-if="route.name === 'manageUsers'">
            <n-icon :component="Users" /> Manage users
        </n-breadcrumb-item>
        <n-breadcrumb-item v-if="route.name === 'manageRoles'">
            <n-icon :component="UserKey" /> Manage roles
        </n-breadcrumb-item>

        <n-breadcrumb-item
            v-if="route.name === 'manageProjectTypes' || route.name === 'manageProjectPriorities' || route.name === 'manageProjectStatuses'">
            <n-icon :component="FolderCog" /> Project settings
        </n-breadcrumb-item>
        <n-breadcrumb-item v-if="route.name === 'manageProjectTypes'">
            <n-icon :component="Bookmark" /> Project types
        </n-breadcrumb-item>
        <n-breadcrumb-item v-if="route.name === 'manageProjectPriorities'">
            <n-icon :component="Goal" /> Project priorities
        </n-breadcrumb-item>
        <n-breadcrumb-item v-if="route.name === 'manageProjectStatuses'">
            <n-icon :component="Route" /> Project statuses
        </n-breadcrumb-item>

        <n-breadcrumb-item v-if="route.name === 'manageTaskPriorities' || route.name === 'manageTaskStatuses'">
            <n-icon :component="FileCog" /> Task settings
        </n-breadcrumb-item>
        <n-breadcrumb-item v-if="route.name === 'manageTaskPriorities'">
            <n-icon :component="Goal" /> Task priorities
        </n-breadcrumb-item>
        <n-breadcrumb-item v-if="route.name === 'manageTaskStatuses'">
            <n-icon :component="Route" /> Task statuses
        </n-breadcrumb-item>

        <n-breadcrumb-item v-if="route.name === 'profile'">
            <n-icon :component="CircleUser" /> {{ sessionStore.sessionUserName }}
        </n-breadcrumb-item>
        <n-breadcrumb-item v-if="route.name === 'profile'">
            <n-icon :component="UserCog" /> Profile
        </n-breadcrumb-item>
    </n-breadcrumb>
</template>

<style lang="css" scoped></style>