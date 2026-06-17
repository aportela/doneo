<script setup lang="ts">
    import { computed } from 'vue';

    import { NModal } from 'naive-ui';

    import { ProjectAttachment } from '../../modules/attachments/models/project-attachment';
    import { formatBytes } from '../composables/format';
    import PDFWrapper from './PDFWrapper.vue';

    interface PDFPreviewProps {
        projectId: string;
        items: ProjectAttachment[];
    };

    const props = defineProps<PDFPreviewProps>();

    const show = defineModel('show', { default: false });

    const currentIndex = defineModel<number>("currentIndex", { default: 0 });

    const currentItem = computed<ProjectAttachment | undefined>(() =>
        props.items[currentIndex.value] ?? undefined
    );
</script>

<template>
    <n-modal v-model:show="show" :bordered="true" class="doneo-pdf-preview-modal">
        <div style="background-color: rgba(250, 250, 252, 1); padding: 16px">
            <div v-if="currentItem" :key="currentItem.id ?? 0">
                <p class="doneo-text-center"><strong>{{ currentItem.name }}</strong> ({{ formatBytes(currentItem.size)
                    }})</p>
                <PDFWrapper :url="currentItem.getPreviewURL(props.projectId)"
                    inner-content-class="doneo-pdf-wrapper-inner-class" />
            </div>
        </div>
    </n-modal>
</template>

<style lang="css" scoped>

    .doneo-pdf-preview-modal {
        width: 64%;
        min-height: 92vh;
    }

    .doneo-pdf-wrapper-inner-class {
        width: 100%;
        height: 88vh;
    }
</style>