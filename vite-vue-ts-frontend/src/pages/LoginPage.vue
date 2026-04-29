<script setup lang="ts">
    //import { useRouter } from "vue-router";
    import { default as LoginForm } from "../components/forms/LoginForm.vue"
    import { NCard, NSpace, NButton } from 'naive-ui'
    import { IconDatabaseStar, IconMoon } from '@tabler/icons-vue';

    import { default as GithubButton } from "../components/buttons/GithubButton.vue";
    import { default as SwitchLocaleButton } from "../components/buttons/SwitchLocaleButton.vue";

    /*
    const router = useRouter();

    const onSuccessLogin = () => {
        router.push(
            { name: "home" }
        ).catch((e) => {
            console.error(e);
        });
    };
    */

</script>

<template>
    <div class="login-page">
        <div class="left">
            <n-card class="login-card" bordered>
                <h2 class="title">Login to your account</h2>
                <LoginForm />
                <template #action>
                    <n-space justify="space-around">
                        <GithubButton :size="20" />
                        <n-button quaternary>
                            <IconMoon :size="20" />
                        </n-button>
                        <SwitchLocaleButton />
                    </n-space>
                </template>
            </n-card>
        </div>
        <div class="right">
            <div class="right-content">
                <h1>
                    <IconDatabaseStar :size="72" />Doneo
                </h1>
                <h2>Turn chaos into progress.</h2>
            </div>
        </div>

    </div>
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