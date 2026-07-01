<script setup lang="ts">
    import { ref, onMounted, onBeforeUnmount } from 'vue';
    import { useI18n } from "vue-i18n";
    import { useRouter } from "vue-router";

    import { NFlex, NButton, NIcon, NModal } from 'naive-ui';
    import { IconHome } from '@tabler/icons-vue';

    import { appBus } from '../../../shared/composables/bus';
    import RemoteAPIAlert from '../alerts/RemoteAPIAlert.vue';

    const router = useRouter();
    const { t } = useI18n();

    const show = defineModel('show', { default: false });

    const closable = ref<boolean>(true);

    const bodyStyle = {
        width: '600px'
    };

    const onGoHome = () => {
        router.push(
            { name: "home" }
        ).catch((e) => {
            console.error(e);
        });
        show.value = false;
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
        size="medium" bordered :style="bodyStyle">
        <RemoteAPIAlert :title="t('shared.errorMessages.Error')" type="error" :message="errorMessage"
            preformatted-message />
        <template #action v-if="!closable">
            <n-flex justify="end">
                <n-button strong secondary @click="onGoHome">
                    <template #icon>
                        <n-icon :component="IconHome" />
                    </template>
                    {{ t("shared.components.modals.RemoteAPIAlertModal.buttons.backToHome.label") }}
                </n-button>
            </n-flex>
        </template>
    </n-modal>
</template>

<style lang="css" scoped></style>