<script setup lang="ts">
    import { ref, onMounted } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NModal, NButton, NSelect, NFlex, NFormItem, NIcon } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconImageGeneration } from '@tabler/icons-vue';

    import { generateAvatar, generateParams, getThemeNames } from 'avatarka';

    const { t } = useI18n();

    const emit = defineEmits(['confirm', 'cancel']);

    const show = defineModel<boolean>("show");

    const bodyStyle = {
        width: '600px'
    }
    const segmented = {
        content: 'soft',
        footer: 'soft'
    } as const

    const svg = ref<string>("");

    const themes = getThemeNames();
    const options = themes.map((t) => { return { label: t, value: t } })

    const selectedTheme = ref(null);

    const generateRandomAvatar = async () => {
        const theme = selectedTheme.value ? selectedTheme.value : themes[Math.floor(Math.random() * themes.length)];
        const params = {
            ...generateParams(theme),
            backgroundShape: "square" as const,

        }
        svg.value = generateAvatar(theme, params);
    }

    const onConfirm = () => {
        emit("confirm", svg.value)
    };
    const onCancel = () => {
        emit("cancel")
    };

    onMounted(() => {
        generateRandomAvatar();
    });

</script>

<template>
    <n-modal v-model:show="show" :title="t('shared.components.dialogs.GenerateAvatarModal.title')" :closable="true"
        preset="card" size="medium" :bordered="true" :segmented="segmented" :style="bodyStyle">
        <n-flex justify="space-between">
            <div class="avatar" v-html="svg" />
            <div>
                <n-form-item :label="t('shared.components.dialogs.GenerateAvatarModal.selectors.themeSelector.label')">
                    <n-select :options="options" v-model:value="selectedTheme"
                        :placeholder="t('shared.components.dialogs.GenerateAvatarModal.selectors.themeSelector.placeholder')"
                        clearable />
                </n-form-item>
                <n-button @click="generateRandomAvatar">
                    <template #icon>
                        <n-icon>
                            <IconImageGeneration />
                        </n-icon>
                    </template>
                    {{ t("shared.components.dialogs.GenerateAvatarModal.buttons.generate.label") }}
                </n-button>
            </div>
        </n-flex>
        <template #action>
            <n-flex justify="end">
                <n-button @click="onConfirm">
                    <template #icon>
                        <n-icon>
                            <IconDeviceFloppy />
                        </n-icon>
                    </template>
                    {{ t("shared.components.dialogs.GenerateAvatarModal.buttons.confirm.label") }}
                </n-button>
                <n-button @click="onCancel">
                    <template #icon>
                        <n-icon>
                            <IconCancel />
                        </n-icon>
                    </template>
                    {{ t("shared.components.dialogs.GenerateAvatarModal.buttons.cancel.label") }}
                </n-button>
            </n-flex>
        </template>
    </n-modal>
</template>

<style lang="css" scoped>
    .avatar {
        width: 256px;
        height: 256px;
    }

    .avatar :deep(svg) {
        width: 100% !important;
        height: 100% !important;
        display: block;
    }
</style>