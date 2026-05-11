<script setup lang="ts">
    import { NAlert } from 'naive-ui'

    type AlertType = "default" | "info" | "warning" | "error" | "success";

    interface RemoteAPIAlertProps {
        title?: string
        type?: AlertType
        message?: string
    }

    const props = withDefaults(defineProps<RemoteAPIAlertProps>(), {
        title: "",
        type: "default",
        message: "",
    })

</script>

<template>
    <n-alert :title="title" :type="type">
        <pre class="alert-message">{{ message }}</pre>
    </n-alert>
</template>

<style lang="css" scoped>
    .alert-message {
        white-space: pre-wrap;
        margin: 0;
    }
</style>