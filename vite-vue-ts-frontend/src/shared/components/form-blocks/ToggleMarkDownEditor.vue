<script setup lang="ts">
    import { computed } from 'vue';

    import { MdEditor, MdPreview } from 'md-editor-v3';
    import 'md-editor-v3/lib/style.css';
    import 'md-editor-v3/lib/preview.css';

    import { useColorSchemeStore } from '../../../stores/colorScheme';

    interface IProps {
        disabled?: boolean;
        readOnly?: boolean;
    };

    const props = withDefaults(defineProps<IProps>(), {
        disabled: false,
        readOnly: false,
    });

    const colorSchemeStore = useColorSchemeStore();

    const value = defineModel<string | null>('value');

    const markDown = computed<string>({
        get: () => value.value ?? "",
        set(markDownStr: string) {
            value.value = markDownStr === "" ? null : markDownStr;
        }
    });

</script>

<template>
    <MdEditor v-model="markDown" :theme="colorSchemeStore.dark ? 'dark' : 'light'" language="en-US"
        :disabled="props.disabled" :read-only="props.readOnly" v-if="!props.readOnly" no-upload-img auto-focus />
    <div v-else class="doneo-md-preview">
        <MdPreview id="mdeditor" v-model:model-value="markDown" no-img-zoom-in
            :theme="colorSchemeStore.dark ? 'dark' : 'light'" language="en-US" />
    </div>
</template>

<style lang="css" scoped>
    .doneo-md-preview {
        border: 1px solid #e0e0e6;
        border-radius: var(--n-border-radius);
        padding: 4px 12px;
        width: 100%;
    }
</style>