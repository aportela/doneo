<script setup lang="ts">
    import { useRouter } from 'vue-router'
    import { useI18n } from "vue-i18n";

    import { NResult, NSpace, NButton } from 'naive-ui'

    const { t } = useI18n();
    const router = useRouter()

    const goHome = () => {
        router.push('/')
    };

    const goBack = () => {
        router.back()
    };
</script>

<template>
    <div class="not-found">
        <n-result status="404" :title="t('layouts.errorNotFound.title')"
            :description="t('layouts.errorNotFound.description')" size="huge">
            <template #footer>
                <n-space justify="space-around">
                    <n-button type="primary" @click="goHome">
                        {{ t("layouts.errorNotFound.buttons.goHome") }}
                    </n-button>
                    <n-button tertiary @click="goBack">
                        {{ t("layouts.errorNotFound.buttons.goBack") }}
                    </n-button>
                </n-space>
            </template>
        </n-result>
    </div>
</template>

<style lang="css" scoped>
    .not-found {
        height: 100vh;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 24px;
    }
</style>