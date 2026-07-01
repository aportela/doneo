<script setup lang="ts">
    import { ref, onMounted, watch } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NButton, NSelect, NFlex, NFormItem, NIcon } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconImageGeneration } from '@tabler/icons-vue';

    import { generateAvatar, generateParams, getThemeNames, type ThemeName } from 'avatarka';

    import { BUTTON_DEFAULT_ICON_SIZE } from '../../../constants';

    const { t } = useI18n();

    const emit = defineEmits(['confirm', 'cancel']);

    const show = defineModel<boolean>("show");

    const bodyStyle = {
        width: '512px'
    };

    const segmented = {
        content: 'soft',
        footer: 'soft'
    } as const;

    const svg = ref<string>("");

    const themes: ThemeName[] = getThemeNames();
    const themeOptions = themes.map((item: ThemeName) => { return { label: item, value: item } })
    const selectedTheme = ref<ThemeName>(themes[Math.floor(Math.random() * themes.length)]);

    type BackgroundShape = "square" | "circle" | "rounded";

    const backgroundShapes: BackgroundShape[] = ["square", "circle", "rounded"];
    const shapeOptions = backgroundShapes.map((item) => { return { label: item, value: item } })
    const selectedShape = ref<BackgroundShape>(backgroundShapes[Math.floor(Math.random() * backgroundShapes.length)]);

    const themeParams = ref();

    const generateRandomAvatar = () => {
        themeParams.value = {
            ...generateParams(selectedTheme.value),
            backgroundShape: selectedShape.value
        }
        svg.value = generateAvatar(selectedTheme.value, themeParams.value);
    }

    const onConfirm = () => {
        emit("confirm", svg.value)
    };

    const onCancel = () => {
        emit("cancel")
    };

    watch(() => [selectedTheme.value, selectedShape.value], () => {
        generateRandomAvatar();
    });

    onMounted(() => {
        generateRandomAvatar();
    });
</script>

<template>
    <n-modal v-model:show="show" :title="t('shared.components.modals.GenerateAvatarModal.title')" :closable="true"
        preset="card" size="small" :bordered="true" :segmented="segmented" :style="bodyStyle">
        <n-flex justify="space-between" align="center">
            <div class="avatar" v-html="svg" />
            <div>
                <n-form-item :label="t('shared.components.modals.GenerateAvatarModal.selectors.themeSelector.label')">
                    <n-select :options="themeOptions" v-model:value="selectedTheme"
                        :placeholder="t('shared.components.modals.GenerateAvatarModal.selectors.themeSelector.placeholder')" />
                </n-form-item>
                <n-form-item :label="t('shared.components.modals.GenerateAvatarModal.selectors.shapeSelector.label')">
                    <n-select :options="shapeOptions" v-model:value="selectedShape"
                        :placeholder="t('shared.components.modals.GenerateAvatarModal.selectors.shapeSelector.placeholder')" />
                </n-form-item>
                <n-button @click="generateRandomAvatar" block>
                    <template #icon>
                        <n-icon :component="IconImageGeneration" :size="BUTTON_DEFAULT_ICON_SIZE" />
                    </template>
                    {{ t("shared.components.modals.GenerateAvatarModal.buttons.generate.label") }}
                </n-button>
            </div>
        </n-flex>
        <template #action>
            <n-flex justify="end">
                <n-button @click="onConfirm">
                    <template #icon>
                        <n-icon :component="IconDeviceFloppy" :size="BUTTON_DEFAULT_ICON_SIZE" />
                    </template>
                    {{ t("shared.components.modals.GenerateAvatarModal.buttons.confirm.label") }}
                </n-button>
                <n-button @click="onCancel">
                    <template #icon>
                        <n-icon :component="IconCancel" :size="BUTTON_DEFAULT_ICON_SIZE" />
                    </template>
                    {{ t("shared.components.modals.GenerateAvatarModal.buttons.cancel.label") }}
                </n-button>
            </n-flex>
        </template>
    </n-modal>
</template>

<style lang="css" scoped>
    .avatar {
        width: 192px;
        height: 192px;
    }

    .avatar :deep(svg) {
        width: 100% !important;
        height: 100% !important;
        display: block;
    }
</style>