<script setup lang="ts">

    import { computed, type Component } from 'vue';
    import { useRoute } from 'vue-router';

    import { NBreadcrumb, NBreadcrumbItem, NIcon } from 'naive-ui';

    import { Bookmark, CircleUser, FileCog, FolderCog, FolderKanban, Form, Goal, Home, LibraryBig, ListTodo, Notebook, Route, Settings, UserCog, UserKey, Users } from '@lucide/vue';

    import { useSessionStore } from '../../../stores/session';

    type AppRouteName =
        | 'manageProjects'
        | 'projectPage'
        | 'projectTab'
        | 'manageTasks'
        | 'manageWiki'
        | 'manageUsers'
        | 'manageRoles'
        | 'manageProjectTypes'
        | 'manageProjectPriorities'
        | 'manageProjectStatuses'
        | 'manageTaskPriorities'
        | 'manageTaskStatuses'
        | 'profile';

    interface BreadcrumbItem {
        id: string;
        label: () => string;
        icon?: Component;
        clickable?: boolean;
    }

    const route = useRoute();

    const sessionStore = useSessionStore();

    const breadcrumbConfig: Record<AppRouteName, BreadcrumbItem[]> = {
        manageProjects: [
            { id: "workspace", label: () => 'Workspace', icon: Notebook, clickable: false },
            { id: "manageProjects", label: () => 'Projects', icon: FolderKanban }
        ],

        projectPage: [
            { id: "workspace", label: () => 'Workspace', icon: Notebook, clickable: false },
            { id: "manageProjects", label: () => 'Projects', icon: FolderKanban },
            { id: "project", label: () => 'project slug', icon: Form }
        ],

        projectTab: [
            { id: "workspace", label: () => 'Workspace', icon: Notebook, clickable: false },
            { id: "manageProjects", label: () => 'Projects', icon: FolderKanban },
            { id: "project", label: () => 'project slug', icon: Form }
        ],

        manageTasks: [
            { id: "workspace", label: () => 'Workspace', icon: Notebook, clickable: false },
            { id: "manageTasks", label: () => 'Tasks', icon: ListTodo },
        ],

        manageWiki: [
            { id: "workspace", label: () => 'Workspace', icon: Notebook, clickable: false },
            { id: "manageWiki", label: () => 'Wiki', icon: LibraryBig }
        ],

        manageUsers: [
            { id: "settings", label: () => 'Settings', icon: Settings, clickable: false },
            { id: "manageUsers", label: () => 'Manage users', icon: Users }
        ],

        manageRoles: [
            { id: "settings", label: () => 'Settings', icon: Settings, clickable: false },
            { id: "ManageRoles", label: () => 'Manage roles', icon: UserKey }
        ],

        manageProjectTypes: [
            { id: "settings", label: () => 'Settings', icon: Settings, clickable: false },
            { id: "projectSettings", label: () => 'Project settings', icon: FolderCog },
            { id: "manageProjectTypes", label: () => 'Project types', icon: Bookmark }
        ],

        manageProjectPriorities: [
            { id: "settings", label: () => 'Settings', icon: Settings, clickable: false },
            { id: "projectSettings", label: () => 'Project settings', icon: FolderCog },
            { id: "manageProjectPriorities", label: () => 'Project priorities', icon: Goal }
        ],

        manageProjectStatuses: [
            { id: "settings", label: () => 'Settings', icon: Settings, clickable: false },
            { id: "projectSettings", label: () => 'Project settings', icon: FolderCog },
            { id: "manageProjectStatuses", label: () => 'Project statuses', icon: Route }
        ],

        manageTaskPriorities: [
            { id: "settings", label: () => 'Settings', icon: Settings, clickable: false },
            { id: "taskSettings", label: () => 'Task settings', icon: FileCog },
            { id: "manageTaskPriorities", label: () => 'Task priorities', icon: Goal }
        ],

        manageTaskStatuses: [
            { id: "settings", label: () => 'Settings', icon: Settings, clickable: false },
            { id: "taskSettings", label: () => 'Task settings', icon: FileCog },
            { id: "manageTaskStatuses", label: () => 'Task statuses', icon: Route }
        ],

        profile: [
            {
                id: "myAccount",
                label: () => sessionStore.sessionUserName ?? "anonymous",
                icon: CircleUser
            },
            {
                id: "profile",
                label: () => 'Profile',
                icon: UserCog
            }
        ]
    };

    const homeBreadcrumb: BreadcrumbItem = {
        id: 'home',
        label: () => 'Home',
        icon: Home,
        clickable: false
    };

    const breadcrumbs = computed(() => {
        if (!route.name || !(route.name in breadcrumbConfig)) {
            return [homeBreadcrumb];
        }

        return [
            homeBreadcrumb,
            ...(breadcrumbConfig[route.name as keyof typeof breadcrumbConfig] ?? [])
        ];
    });
</script>

<template>
    <n-breadcrumb>
        <n-breadcrumb-item v-for="item in breadcrumbs" :key="item.id" :clickable="false">
            <n-icon :component="item.icon" />
            {{ item.label() }}
        </n-breadcrumb-item>
    </n-breadcrumb>
</template>

<style lang="css" scoped></style>