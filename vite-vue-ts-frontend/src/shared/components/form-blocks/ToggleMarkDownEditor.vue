<script setup lang="ts">
    import { computed } from 'vue';

    import { useMarkdown } from '../../composables/useMarkdown';
    import { MdEditor, MdPreview, type ToolbarNames } from 'md-editor-v3';
    import 'md-editor-v3/lib/style.css';
    import 'md-editor-v3/lib/preview.css';

    import { useColorSchemeStore } from '../../../stores/colorScheme';

    interface IProps {
        disabled?: boolean;
        readOnly?: boolean;
        maxLength?: number;
    };

    const props = withDefaults(defineProps<IProps>(), {
        disabled: false,
        readOnly: false,
    });

    const colorSchemeStore = useColorSchemeStore();
    const { toMarkdown } = useMarkdown();

    const value = defineModel<string | null>('value');

    const markDown = computed<string>({
        get: () => value.value ?? "",
        set(markDownStr: string) {
            value.value = markDownStr === "" ? null : markDownStr;
        }
    });

    const insertAtCursor = (value: string) => {
        const el = document.activeElement as HTMLTextAreaElement
        if (!el) {
            markDown.value += value
            return
        }

        const start = el.selectionStart ?? markDown.value?.length
        const end = el.selectionEnd ?? markDown.value?.length

        markDown.value =
            markDown.value?.slice(0, start) +
            value +
            markDown.value?.slice(end)

        // restore cursor
        requestAnimationFrame(() => {
            el.selectionStart = el.selectionEnd = start + value.length
        })
    }

    const onPaste = (e: ClipboardEvent) => {
        const clipboard = e.clipboardData
        if (!clipboard) return

        const html = clipboard.getData('text/html')
        const plain = clipboard.getData('text/plain')

        let markdown = plain

        if (html) {
            markdown = toMarkdown(html)
        }

        e.preventDefault()

        insertAtCursor(markdown)
    };

    const excludedMDEditorToolBars: ToolbarNames[] = ["previewOnly", "htmlPreview", "catalog", "github"];

</script>

<template>
    <MdEditor v-model="markDown" :max-length="props.maxLength" :theme="colorSchemeStore.dark ? 'dark' : 'light'"
        language="en-US" :disabled="props.disabled" :read-only="props.readOnly" v-if="!props.readOnly" no-upload-img
        auto-focus @paste="onPaste" :toolbars-exclude="excludedMDEditorToolBars" :footers="[]" />
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