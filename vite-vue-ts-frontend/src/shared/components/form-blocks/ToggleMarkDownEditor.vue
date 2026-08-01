<script setup lang="ts">
    import { ref, computed } from 'vue';

    import { useMarkdown } from '../../composables/useMarkdown';
    import { MdEditor, MdPreview, type ToolbarNames, type ExposeParam } from 'md-editor-v3';
    import 'md-editor-v3/lib/style.css';
    import 'md-editor-v3/lib/preview.css';

    import { useUserSettingsStore } from '../../../stores/userSettings';

    interface Props {
        disabled?: boolean;
        readOnly?: boolean;
        maxLength?: number;
        hidePreview?: boolean;
        placeholder?: string;
        autoFocus?: boolean;
        noUploadImg?: boolean;
        simpleToolbars?: boolean;
    };

    const props = withDefaults(defineProps<Props>(), {
        disabled: false,
        readOnly: false,
        hidePreview: false,
        noUploadImg: true,
    });

    const userSettingsStore = useUserSettingsStore();
    const { toMarkdown } = useMarkdown();

    const value = defineModel<string | null>('value');

    const editorRef = ref<ExposeParam>();

    const markDown = computed<string>({
        get: () => value.value ?? "",
        set(markDownStr: string) {
            value.value = markDownStr === "" ? null : markDownStr;
        }
    });

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

        editorRef.value?.insert((_selectedText) => ({
            targetValue: markdown,
            select: false
        }));
    };

    //const excludedMDEditorToolBars: ToolbarNames[] = ["save", "previewOnly", "htmlPreview", "catalog", "github", "mermaid", "katex", "revoke", "sub", "sup", "prettier", "pageFullscreen", "fullscreen"];
    const simpleToolbarItems: ToolbarNames[] = [
        'bold',
        'underline',
        'italic',
        '-',
        'title',
        'strikeThrough',
        'sub',
        'sup',
        'quote',
        'unorderedList',
        'orderedList',
        'task',
        '-',
        'codeRow',
        'code',
        'link',
        'image',
        'table',
        '-',
        'preview',
    ];
</script>

<template>
    <MdEditor ref="editorRef" v-model="markDown" :max-length="props.maxLength"
        :theme="userSettingsStore.darkTheme ? 'dark' : 'light'" language="en-US" :disabled="props.disabled"
        :read-only="props.readOnly" v-if="!props.readOnly" :no-upload-img="props.noUploadImg"
        :auto-focus="props.autoFocus" @paste="onPaste" :toolbars="props.simpleToolbars ? simpleToolbarItems : undefined"
        :footers="[]" :preview="!props.hidePreview" :placeholder="props.placeholder" />
    <MdPreview v-else id="mdeditor" v-model:model-value="markDown" no-img-zoom-in
        :theme="userSettingsStore.darkTheme ? 'dark' : 'light'" language="en-US" :code-foldable="false"
        class="doneo-markdown-preview-container" />
</template>

<style lang="css" scoped>
    .doneo-markdown-preview-container {
        border: 1px solid #e0e0e6;
        border-radius: 3px;
        padding: 0em 1em;
    }
</style>