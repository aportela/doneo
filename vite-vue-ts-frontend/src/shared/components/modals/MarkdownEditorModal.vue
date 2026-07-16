<script setup lang="ts">
    import { NModal } from 'naive-ui';

    import { MdEditor } from 'md-editor-v3';
    import 'md-editor-v3/lib/style.css';

    import { useUserSettingsStore } from '../../../stores/userSettings';

    const userSettingsStore = useUserSettingsStore();

    const show = defineModel<boolean>('show', { default: false });
    const markDown = defineModel<string>('markDown');
</script>

<template>
    <n-modal v-model:show="show" closable close-on-esc mask-closable preset="card" size="huge" bordered>
        <MdEditor v-model="markDown" :theme="userSettingsStore.darkTheme ? 'dark' : 'light'" language="en-US" />
    </n-modal>
</template>

<style lang="css" scoped></style>