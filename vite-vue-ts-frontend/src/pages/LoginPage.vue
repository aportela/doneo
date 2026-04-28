<script setup lang="ts">
    import { ref, reactive } from 'vue';
    import { useRouter } from "vue-router";
    import { api } from '../composables/api';
    import { useSessionStore } from "../stores/session";
    import { type AjaxState as AjaxStateInterface, defaultAjaxState } from "../types/ajaxState";


    const invalidEmailField = ref(false);
    const invalidEmailFeedbackMessage = ref<string>("");
    const invalidPasswordField = ref(false);
    const invalidPasswordFeedbackMessage = ref<string>("");

    const router = useRouter();

    const sessionStore = useSessionStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });


    const email = ref("admin@localhost");
    const password = ref("secret");


    const onSubmit = () => {
        invalidEmailField.value = false;
        invalidPasswordField.value = false;
        if (email.value && password.value) {
            state.ajaxRunning = true;
            api.auth.signIn(email.value, password.value)
                .then((successResponse: any) => {
                    if (successResponse.data.accessToken) {
                        sessionStore.setAccessToken(successResponse.data.accessToken.token, successResponse.data.accessToken.expiresAtTimestamp);
                        router.push(
                            { name: "home" }
                        ).catch((e) => {
                            console.error(e);
                        });
                    } else {
                        console.error("Invalid response");
                    }
                })
                .catch((errorResponse) => {
                    state.ajaxErrors = true;
                    switch (errorResponse.status) {
                        case 404:
                            invalidEmailField.value = true;
                            invalidEmailFeedbackMessage.value = "Email not found";
                            break;
                        case 401:
                            invalidPasswordField.value = true;
                            invalidPasswordFeedbackMessage.value = "Invalid password";
                            break;
                    }
                })
                .finally(() => {
                    state.ajaxRunning = false;
                });
        }
    }
</script>

<template>
    <div class="row g-0 flex-fill">
        <div class="col-12 col-lg-6 col-xl-4 d-flex flex-column justify-content-center">
            <div class="container container-tight my-5 px-lg-5">
                <div class="text-center mb-4">
                    <img src="../assets/hero.png" width="110" height="32" />
                    <strong style="margin-left: 1em; font-size: 2em;">Doneo</strong>
                </div>
                <h2 class="h3 text-center mb-3">Login to your account</h2>
                <form action="/api/auth/signin" method="post" autocomplete="off" @submit.prevent="onSubmit">
                    <div class="mb-3">
                        <label class="form-label">Email address</label>
                        <input type="email" v-model.trim="email" :disabled="state.ajaxRunning" class="form-control"
                            :class="{ 'is-invalid': invalidEmailField }" placeholder="your@email.com" autocomplete="off"
                            required>
                        <div class="invalid-feedback" v-if="invalidEmailField">{{ invalidEmailFeedbackMessage }}</div>
                    </div>
                    <div class="mb-2">
                        <label class="form-label">
                            Password
                            <!--
                            <span class="form-label-description">
                                <a href="./forgot-password.html">I forgot password</a>
                            </span>
                            -->
                        </label>
                        <div class="input-group input-group-flat">
                            <input type="password" class="form-control" :class="{ 'is-invalid': invalidPasswordField }"
                                placeholder="Your password" autocomplete="off" v-model="password"
                                :disabled="state.ajaxRunning" required minlength="5">
                            <div class="invalid-feedback" v-if="invalidPasswordField">{{ invalidPasswordFeedbackMessage
                            }}
                            </div>
                        </div>
                    </div>
                    <div class="form-footer">
                        <button type="submit" :disabled="state.ajaxRunning" class="btn btn-primary w-100">Sign
                            in</button>
                    </div>
                </form>
                <!--
                <div class="text-center text-secondary mt-3">Don't have account yet? <router-link
                        to="/auth/register">Sign up</router-link>
                </div>
                -->
            </div>
        </div>
        <div class="col-12 col-lg-6 col-xl-8 d-none d-lg-block">
            <!-- Photo -->
            <div class="bg-cover h-100 min-vh-100"
                style="background-image: url(https://picsum.photos/960/540); filter: grayscale(1) blur(1px);">
            </div>
        </div>
    </div>
</template>