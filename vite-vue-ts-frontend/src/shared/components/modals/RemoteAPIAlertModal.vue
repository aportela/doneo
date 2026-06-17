<script setup lang="ts">
    import { ref, onMounted, onBeforeUnmount } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal } from 'naive-ui';

    import { appBus } from '../../../shared/composables/bus';
    import RemoteAPIAlert from '../alerts/RemoteAPIAlert.vue';

    const { t } = useI18n();

    const show = defineModel('show', { default: false });

    const closable = ref<boolean>(true);

    const bodyStyle = {
        width: '600px'
    };


    let busListener: () => void;
    const errorMessage = ref<string>("");

    onMounted(async () => {
        busListener = appBus.on("remoteAPIError", (payload) => {
            errorMessage.value = payload.errorMessage
            if (payload.denyCloseDialog) {
                closable.value = false;
            }
            show.value = true;
        });
    });

    onBeforeUnmount(() => {
        busListener();
    });
</script>

<template>
    <n-modal v-model:show="show" :closable="closable" :close-on-esc="closable" :mask-closable="closable" preset="card"
        size="medium" :bordered="true" :style="bodyStyle">
        <template #header>
        </template>
        <RemoteAPIAlert type="error" :title="t('shared.errorMessages.Error')" :message="errorMessage" />
    </n-modal>
</template>

<style lang="css" scoped></style>