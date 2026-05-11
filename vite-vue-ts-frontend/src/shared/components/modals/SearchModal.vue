<script setup lang="ts">

    import { NModal, NInput, NEmpty } from 'naive-ui';

    const show = defineModel<boolean>();
    const bodyStyle = {
        width: '600px'
    }
    const segmented = {
        content: 'soft',
        footer: 'soft'
    } as const
</script>

<template>
    <n-modal v-model:show="show" title="Search" :closable="true" preset="card" size="medium" :bordered="true"
        :segmented="segmented" :style="bodyStyle">
        <n-input placeholder="search..." autofocus clearable></n-input>

        <n-empty description="No data results" style="margin-top: 16px;">
        </n-empty>
        <template #footer>
            <div class="shortcut">
                <kbd>↵</kbd><span>Select</span>
                <kbd>↑</kbd><kbd>↓</kbd><span>Toggle</span>
                <kbd>Del</kbd><span>Delete</span>
                <kbd>Esc</kbd><span>Close</span>
            </div>

            <!--
            <div class="shortcut">
                <kbd>Ctrl</kbd>
                <span>+</span>
                <kbd>K</kbd>
            </div>
            -->
        </template>
    </n-modal>
</template>

<style lang="css" scoped>
    .shortcut {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 12px;
    }

    kbd {
        padding: 2px 8px;
        border-radius: 4px;
        border: 1px solid #ccc;
        background: #f5f5f5;
        font-family: monospace;
        font-size: 14px;
    }
</style>