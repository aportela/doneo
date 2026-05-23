<script setup lang="ts">
    import { NGrid, NGridItem, NCard } from 'naive-ui';

</script>

<template>
    <h1>Overview</h1>
    <h2>Dashboard</h2>

    <n-grid cols="1 1600:2" :x-gap="32" :y-gap="32" responsive="self" item-responsive style="margin-top: 1em;">
        <n-grid-item>
            <n-card>
                <h2>
                    WARNING
                </h2>


                <p>5 pending post-install settings</p>
                <ul>
                    <li>No project types found, <router-link :to="{ name: 'manageProjectTypes' }">create
                            one</router-link></li>
                    <li>No project priorities found, <router-link :to="{ name: 'manageProjectPriorities' }">create
                            one</router-link></li>
                    <li>No project statuses found, <router-link :to="{ name: 'manageProjectStatuses' }">create
                            one</router-link></li>
                    <li>No task priorities found, <router-link :to="{ name: 'manageTaskPriorities' }">create
                            one</router-link></li>
                    <li>No task statuses found, <router-link :to="{ name: 'manageTaskStatuses' }">create
                            one</router-link></li>
                </ul>
            </n-card>
        </n-grid-item>
        <n-grid-item>
            <n-card>
                <h2>Welcome back, Administrator</h2>

                <p>You have 5 new messages and 2 pending tasks.</p>
            </n-card>
        </n-grid-item>
        <n-grid-item>
            <n-card>
                <h2>Welcome back, Administrator</h2>

                <p>You have 5 new messages and 2 pending tasks.</p>
            </n-card>
        </n-grid-item>
        <n-grid-item>
            <n-card>
                <h2>Welcome back, Administrator</h2>

                <p>You have 5 new messages and 2 pending tasks.</p>
            </n-card>
        </n-grid-item>
        <n-grid-item>
            <n-card>
                <h2>Welcome back, Administrator</h2>

                <p>You have 5 new messages and 2 pending tasks.</p>
            </n-card>
        </n-grid-item>
        <n-grid-item>
            <n-card>
                <h2>Welcome back, Administrator</h2>

                <p>You have 5 new messages and 2 pending tasks.</p>
            </n-card>
        </n-grid-item>
    </n-grid>
</template>

<style lang="css" scoped></style>