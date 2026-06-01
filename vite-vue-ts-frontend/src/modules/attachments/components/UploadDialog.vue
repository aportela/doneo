<script setup lang="ts">
    import { NModal, NUpload, NUploadDragger, NIcon, NText, type UploadCustomRequestOptions } from 'naive-ui';
    import { IconUpload } from '@tabler/icons-vue';

    import { useSessionStore } from '../../../stores/session';

    interface UploadDialogProps {
        projectId: string;
    };

    const props = defineProps<UploadDialogProps>();

    const sessionStore = useSessionStore();

    const show = defineModel<boolean>("show");

    const uploadFile = async ({
        file,
        headers,
        onProgress,
        onFinish,
        onError
    }: UploadCustomRequestOptions) => {
        const formData = new FormData()
        formData.append('file', file.file as Blob)

        const xhr = new XMLHttpRequest()
        xhr.open('POST', `/api/projects/${props.projectId}/attachments/`)

        const headerObj = typeof headers === 'function' ? headers({ file }) : headers
        if (headerObj) {
            for (const key in headerObj) {
                xhr.setRequestHeader(key, headerObj[key])
            }
        }

        xhr.setRequestHeader('Authorization', `Bearer ${sessionStore.accessToken}`)

        xhr.upload.onprogress = (event) => {
            if (event.lengthComputable) {
                const percent = (event.loaded / event.total) * 100
                onProgress({ percent })
            }
        }

        xhr.onload = () => {
            if (xhr.status >= 200 && xhr.status < 300) {
                onFinish()
            } else {
                onError()
            }
        }

        xhr.onerror = () => {
            onError()
        }

        xhr.send(formData)
    }
</script>

<template>
    <n-modal v-model:show="show" preset="card" style="width: 98vw;">
        <template #default>
            <div style="height: 90vh;">
                <n-upload multiple directory-dnd :max="5" list-type="image" :custom-request="uploadFile">
                    <n-upload-dragger style="height: 40vh;">
                        <div style=" margin-bottom: 12px">
                            <n-icon size="48" :depth="3" :component="IconUpload" />
                        </div>
                        <n-text style="font-size: 16px">
                            Click or drag a file to this area to upload
                        </n-text>
                    </n-upload-dragger>
                </n-upload>
            </div>
        </template>
    </n-modal>
</template>

<style lang="css" scoped></style>