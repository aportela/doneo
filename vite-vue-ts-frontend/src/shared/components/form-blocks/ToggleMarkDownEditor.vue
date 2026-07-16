<script setup lang="ts">
    import { ref, computed } from 'vue';

    import { useMarkdown } from '../../composables/useMarkdown';
    import { MdEditor, MdPreview, type ToolbarNames, type ExposeParam } from 'md-editor-v3';
    import 'md-editor-v3/lib/style.css';
    import 'md-editor-v3/lib/preview.css';

    import { useUserSettingsStore } from '../../../stores/userSettings';

    interface IProps {
        disabled?: boolean;
        readOnly?: boolean;
        maxLength?: number;
        hidePreview?: boolean;
        placeholder?: string;
        autoFocus?: boolean;
        noUploadImg?: boolean;
    };

    const props = withDefaults(defineProps<IProps>(), {
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

    const excludedMDEditorToolBars: ToolbarNames[] = ["save", "previewOnly", "htmlPreview", "catalog", "github"];

</script>

<template>
    <MdEditor ref="editorRef" v-model="markDown" :max-length="props.maxLength"
        :theme="userSettingsStore.darkTheme ? 'dark' : 'light'" language="en-US" :disabled="props.disabled"
        :read-only="props.readOnly" v-if="!props.readOnly" :no-upload-img="props.noUploadImg"
        :auto-focus="props.autoFocus" @paste="onPaste" :toolbars-exclude="excludedMDEditorToolBars" :footers="[]"
        :preview="!props.hidePreview" :placeholder="props.placeholder" />
    <MdPreview v-else id="mdeditor" v-model:model-value="markDown" no-img-zoom-in
        :theme="userSettingsStore.darkTheme ? 'dark' : 'light'" language="en-US" :code-foldable="false" />
</template>

<style lang="css" scoped></style>