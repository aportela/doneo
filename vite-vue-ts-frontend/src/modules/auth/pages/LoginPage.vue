<script setup lang="ts">
    import { useRouter } from "vue-router";
    import { useI18n } from "vue-i18n";

    import { NCard, NFlex, NIcon } from 'naive-ui'

    import LoginForm from "../components/LoginForm.vue";
    import GithubButton from "../../../shared/components/buttons/GithubButton.vue";
    import SwitchColorSchemeButton from "../../../shared/components/buttons/SwitchColorSchemeButton.vue";
    import SwitchLocaleDropdown from "../../../shared/components/dropdowns/SwitchLocaleDropdown.vue";
    import RemoteAPIAlertModal from "../../../shared/components/modals/RemoteAPIAlertModal.vue";

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
                    <n-icon :size="72">
                        <!-- credits IconDatabaseStar icon from tabler icons -->
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
                            class="icon icon-tabler icons-tabler-outline icon-tabler-database-star">
                            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                            <path d="M4 6c0 1.657 3.582 3 8 3s8 -1.343 8 -3s-3.582 -3 -8 -3s-8 1.343 -8 3" />
                            <path d="M4 6v6c0 1.43 2.67 2.627 6.243 2.927" />
                            <path d="M20 10.5v-4.5" />
                            <path d="M4 12v6c0 1.546 3.12 2.82 7.128 2.982" />
                            <path
                                d="M17.8 20.817l-2.172 1.138a.392 .392 0 0 1 -.568 -.41l.415 -2.411l-1.757 -1.707a.389 .389 0 0 1 .217 -.665l2.428 -.352l1.086 -2.193a.392 .392 0 0 1 .702 0l1.086 2.193l2.428 .352a.39 .39 0 0 1 .217 .665l-1.757 1.707l.414 2.41a.39 .39 0 0 1 -.567 .411l-2.172 -1.138" />
                        </svg>
                    </n-icon>
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