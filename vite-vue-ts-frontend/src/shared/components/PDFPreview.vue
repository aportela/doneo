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

    const onPrevious = () => {
        if (currentIndex.value > 0) {
            currentIndex.value--;
        } else {
            currentIndex.value = props.items.length - 1
        }
    };

    const onNext = () => {
        if (currentIndex.value < props.items.length - 1) {
            currentIndex.value++;
        } else {
            currentIndex.value = 0;
        }
    };

    const onClose = () => {
        show.value = false;
    }
</script>

<template>
    <n-modal v-model:show="show" :bordered="true" class="doneo-pdf-preview-modal">
        <div style="background-color: rgba(250, 250, 252, 1); padding: 16px">
            <div v-if="currentItem" :key="currentItem.id ?? 0">
                <p class="doneo-text-center"><strong>{{ currentItem.name }}</strong> ({{ formatBytes(currentItem.size)
                    }})</p>
                <PDFWrapper :url="currentItem.getPreviewURL(props.projectId)"
                    inner-content-class="doneo-pdf-wrapper-inner-class" />
                <div class="doneo-pdf-preview-toolbar">
                    <svg @click="onPrevious" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path
                            d="M6 5C5.75454 5 5.55039 5.17688 5.50806 5.41012L5.5 5.5V14.5C5.5 14.7761 5.72386 15 6 15C6.24546 15 6.44961 14.8231 6.49194 14.5899L6.5 14.5V5.5C6.5 5.22386 6.27614 5 6 5ZM13.8536 5.14645C13.68 4.97288 13.4106 4.9536 13.2157 5.08859L13.1464 5.14645L8.64645 9.64645C8.47288 9.82001 8.4536 10.0894 8.58859 10.2843L8.64645 10.3536L13.1464 14.8536C13.3417 15.0488 13.6583 15.0488 13.8536 14.8536C14.0271 14.68 14.0464 14.4106 13.9114 14.2157L13.8536 14.1464L9.70711 10L13.8536 5.85355C14.0488 5.65829 14.0488 5.34171 13.8536 5.14645Z"
                            fill="currentColor"></path>
                    </svg>
                    <svg @click="onNext" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path
                            d="M13.5 5C13.7455 5 13.9496 5.17688 13.9919 5.41012L14 5.5V14.5C14 14.7761 13.7761 15 13.5 15C13.2545 15 13.0504 14.8231 13.0081 14.5899L13 14.5V5.5C13 5.22386 13.2239 5 13.5 5ZM5.64645 5.14645C5.82001 4.97288 6.08944 4.9536 6.28431 5.08859L6.35355 5.14645L10.8536 9.64645C11.0271 9.82001 11.0464 10.0894 10.9114 10.2843L10.8536 10.3536L6.35355 14.8536C6.15829 15.0488 5.84171 15.0488 5.64645 14.8536C5.47288 14.68 5.4536 14.4106 5.58859 14.2157L5.64645 14.1464L9.79289 10L5.64645 5.85355C5.45118 5.65829 5.45118 5.34171 5.64645 5.14645Z"
                            fill="currentColor"></path>
                    </svg>
                    <svg viewBox="0 0 16 16" version="1.1" xmlns="http://www.w3.org/2000/svg">
                        <g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
                            <g fill="currentColor" fill-rule="nonzero">
                                <path
                                    d="M3.5,13 L12.5,13 C12.7761424,13 13,13.2238576 13,13.5 C13,13.7454599 12.8231248,13.9496084 12.5898756,13.9919443 L12.5,14 L3.5,14 C3.22385763,14 3,13.7761424 3,13.5 C3,13.2545401 3.17687516,13.0503916 3.41012437,13.0080557 L3.5,13 L12.5,13 L3.5,13 Z M7.91012437,1.00805567 L8,1 C8.24545989,1 8.44960837,1.17687516 8.49194433,1.41012437 L8.5,1.5 L8.5,10.292 L11.1819805,7.6109127 C11.3555469,7.43734635 11.6249713,7.4180612 11.8198394,7.55305725 L11.8890873,7.6109127 C12.0626536,7.78447906 12.0819388,8.05390346 11.9469427,8.2487716 L11.8890873,8.31801948 L8.35355339,11.8535534 C8.17998704,12.0271197 7.91056264,12.0464049 7.7156945,11.9114088 L7.64644661,11.8535534 L4.1109127,8.31801948 C3.91565056,8.12275734 3.91565056,7.80617485 4.1109127,7.6109127 C4.28447906,7.43734635 4.55390346,7.4180612 4.7487716,7.55305725 L4.81801948,7.6109127 L7.5,10.292 L7.5,1.5 C7.5,1.25454011 7.67687516,1.05039163 7.91012437,1.00805567 L8,1 L7.91012437,1.00805567 Z">
                                </path>
                            </g>
                        </g>
                    </svg>
                    <svg @click="onClose" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path
                            d="M4.089 4.216l.057-.07a.5.5 0 0 1 .638-.057l.07.057L10 9.293l5.146-5.147a.5.5 0 0 1 .638-.057l.07.057a.5.5 0 0 1 .057.638l-.057.07L10.707 10l5.147 5.146a.5.5 0 0 1 .057.638l-.057.07a.5.5 0 0 1-.638.057l-.07-.057L10 10.707l-5.146 5.147a.5.5 0 0 1-.638.057l-.07-.057a.5.5 0 0 1-.057-.638l.057-.07L9.293 10L4.146 4.854a.5.5 0 0 1-.057-.638l.057-.07l-.057.07z"
                            fill="currentColor"></path>
                    </svg>
                </div>
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

    .doneo-pdf-preview-toolbar {
        z-index: 100;
        position: absolute;
        bottom: 40px;
        left: 50%;
        transform: translateX(-50%);
        border-radius: 24px;
        padding: 12px;
        background: rgba(0, 0, 0, .35);
        color: rgba(255, 255, 255, .9);
        transition: color .3s var(--n-bezier);
        display: flex;
        flex-flow: wrap;
        justify-content: space-around;
        gap: 12px 16px;
        align-items: center;
    }

    .doneo-pdf-preview-toolbar svg {
        cursor: pointer;
        width: 28px;
        height: 28px;
    }
</style>