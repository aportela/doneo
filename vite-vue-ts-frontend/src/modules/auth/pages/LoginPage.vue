<script setup lang="ts">
    import { useRouter } from "vue-router";
    import { useI18n } from "vue-i18n";

    import { NCard, NFlex, NIcon } from 'naive-ui'

    import LoginForm from "../components/LoginForm.vue";
    import GithubButton from "../../../shared/components/buttons/GithubButton.vue";
    import SwitchColorSchemeButton from "../../../shared/components/buttons/SwitchColorSchemeButton.vue";
    import SwitchLocaleDropdown from "../../../shared/components/dropdowns/SwitchLocaleDropdown.vue";
    import RemoteAPIAlertModal from "../../../shared/components/modals/RemoteAPIAlertModal.vue";
    import Doneo from "../../../shared/components/icons/Doneo.vue";

    const router = useRouter();

    const { t } = useI18n();

    const onSuccessLogin = () => {
        router.push(
            { name: "home" }
        ).catch((e) => {
            console.error(e);
        });
    };
</script>

<template>
    <div class="login-page">
        <div class="left">
            <n-card class="login-card">
                <h2 class="title">{{ t("modules.auth.components.LoginPage.headerMessage") }}</h2>
                <LoginForm @success="onSuccessLogin" />
                <template #action>
                    <n-flex justify="space-around">
                        <GithubButton />
                        <SwitchColorSchemeButton />
                        <SwitchLocaleDropdown />
                    </n-flex>
                </template>
            </n-card>
        </div>
        <div class="right">
            <div class="right-content">
                <h1>
                    <n-icon :size="72" :component="Doneo" />
                    Doneo
                </h1>
                <h2>{{ t("modules.auth.components.LoginPage.headerSlogan") }}</h2>
            </div>
        </div>
    </div>
    <RemoteAPIAlertModal />
</template>

<style lang="css" scoped>
    .login-page {
        height: 100vh;
        display: flex;
    }

    .left {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: center;
        background: #fff;
        padding: 40px;
    }

    .login-card {
        width: 100%;
        max-width: 380px;
    }

    .right {
        flex: 1;
        background: #ff4d8d;
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        padding: 40px;
    }

    .right-content {
        max-width: 420px;
    }

    .right h1 {
        font-size: 36px;
        margin-bottom: 12px;
    }

    .right h2 {
        opacity: 0.8;
        font-size: 16px;
    }

    @media (max-width: 768px) {
        .right {
            display: none;
        }
    }
</style>