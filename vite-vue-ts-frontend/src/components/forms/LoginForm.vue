<script setup lang="ts">
    import { ref, reactive } from 'vue';
    import { useRouter } from "vue-router";
    import { api } from '../../composables/api';
    import { useSessionStore } from "../../stores/session";
    import { type AjaxState as AjaxStateInterface, defaultAjaxState } from "../../types/ajaxState";


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
    <form action="/api/auth/signin" method="post" autocomplete="off" @submit.prevent="onSubmit">
        <div class="mb-3">
            <label class="form-label">Email address</label>
            <input type="email" v-model.trim="email" :disabled="state.ajaxRunning" class="form-control"
                :class="{ 'is-invalid': invalidEmailField }" placeholder="your@email.com" autocomplete="off" required>
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
                    placeholder="Your password" autocomplete="off" v-model="password" :disabled="state.ajaxRunning"
                    required minlength="5">
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
</template>

<style lang="css" scoped></style>