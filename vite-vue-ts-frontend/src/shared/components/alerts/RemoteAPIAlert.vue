<script setup lang="ts">
    import { NAlert } from 'naive-ui'

    interface IProps {
        bordered?: boolean;
        closable?: boolean;
        showIcon?: boolean;
        title: string
        type?: "default" | "info" | "warning" | "error" | "success"
        message: string,
        preformattedMessage?: boolean;
    }

    const props = withDefaults(defineProps<IProps>(), {
        bordered: true,
        closable: false,
        showIcon: true,
        type: "default",
        preformattedMessage: false,
    });
</script>

<template>
    <n-alert :bordered="props.bordered" :closable="props.closable" :show-icon="props.showIcon" :title="props.title"
        :type="props.type">
        <pre class="preformatted-message" v-if="props.preformattedMessage">{{ props.message }}</pre>
        <div v-else>{{ props.message }}</div>
    </n-alert>
</template>

<style lang="css" scoped>
    .preformatted-message {
        white-space: pre-wrap;
        margin: 0;
    }
</style>