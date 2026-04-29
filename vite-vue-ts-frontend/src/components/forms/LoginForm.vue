<script setup lang="ts">
    import { ref, reactive, watch } from 'vue';
    import { useRouter } from "vue-router";
    import { NIcon, NSpin, NForm, NFormItem, NInput, NButton, type FormItemRule, type FormInst, type FormRules } from 'naive-ui'
    import { IconEye, IconEyeCancel } from '@tabler/icons-vue';
    import { api } from '../../composables/api';
    import { required, minLength, validEmail, runValidators } from '../../composables/form-validators';
    import { useSessionStore } from "../../stores/session";
    import { type AjaxState as AjaxStateInterface, defaultAjaxState } from "../../types/ajaxState";
    import { default as RemoteAPIAlert } from '../alerts/RemoteAPIAlert.vue';

    import { useI18n } from "vue-i18n";

    import { createStorageEntry } from '../../composables/localStorage';

    const localStorageLastUsedEmail = createStorageEntry<string | null>("lastUsedEmail", null);

    const lastUsedEmail = localStorageLastUsedEmail.get();

    const router = useRouter();

    const { t } = useI18n();

    const sessionStore = useSessionStore();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const signInFormRef = ref<FormInst | null>(null)

    const formValues = ref({
        email: lastUsedEmail || null,
        password: null,
    });

    const serverErrors = ref<{ [key: string]: string }>({})

    const rules: FormRules = {
        email: {
            validator: (_rule: FormItemRule, value) => {
                const localResult = runValidators(value, [required('Email'), validEmail])
                if (localResult !== true) return localResult
                if (serverErrors.value.email) return new Error(serverErrors.value.email)
                return true
            },
            trigger: ['blur']
        },
        password: {
            validator: (_rule: FormItemRule, value) => {
                const localResult = runValidators(value, [required('Password'), minLength(4)])
                if (localResult !== true) return localResult
                if (serverErrors.value.password) return new Error(serverErrors.value.password)
                return true
            },
            trigger: ['input', 'blur']
        }
    }

    watch(() => formValues.value.email, () => {
        delete serverErrors.value.email
    });

    watch(() => formValues.value.password, () => {
        delete serverErrors.value.password
    });

    const onSubmit = () => {
        serverErrors.value = {}
        if (formValues.value.email && formValues.value.password) {
            Object.assign(state, defaultAjaxState);
            state.ajaxRunning = true;
            api.auth.signIn(formValues.value.email, formValues.value.password)
                .then((successResponse: any) => {
                    if (successResponse.data.accessToken) {
                        sessionStore.setAccessToken(successResponse.data.accessToken.token, successResponse.data.accessToken.expiresAtTimestamp);
                        localStorageLastUsedEmail.set(formValues.value.email);
                        router.push(
                            { name: "home" }
                        ).catch((e) => {
                            console.error(e);
                        });
                    } else {
                        state.ajaxErrorMessage = "Invalid response";
                    }
                })
                .catch((errorResponse) => {
                    state.ajaxErrors = true;
                    if (errorResponse.isAPIError) {
                        state.ajaxAPIErrorDetails = errorResponse.customAPIErrorDetails;
                        switch (errorResponse.status) {
                            case 404:
                                serverErrors.value.email = "Email not found";
                                break;
                            case 401:
                                serverErrors.value.password = "Invalid password";
                                break;
                            default:
                                state.ajaxErrorMessage = "API Error: fatal error";
                                break;

                        }
                        signInFormRef.value?.restoreValidation();
                        signInFormRef.value?.validate().catch(() => { });
                    }
                    else {
                        state.ajaxErrorMessage = `Uncaught exception (${errorResponse.status})\n\n${errorResponse}`;
                        console.error(errorResponse);
                    }
                })
                .finally(() => {
                    state.ajaxRunning = false;
                });
        }
    }

    const validateForm = async (e: MouseEvent | KeyboardEvent) => {
        e.preventDefault()
        serverErrors.value = {};
        try {
            await signInFormRef.value?.validate();
            onSubmit();
        }
        catch (e) {
            console.log(e);
        }
    }
</script>

<template>
    <!--
    <h4 class="q-mt-sm q-mb-md text-h4 text-weight-bolder">{{
        t(!!lastUsedEmail ? "Glad to see you again!" : "Welcome aboard!")
        }}</h4>
    <div class="text-color-secondary">{{
        t(!!lastUsedEmail ? "Let's get back to organizing." : "Let's start organizing.")
        }}
    </div>
    -->

    <n-spin :show="state.ajaxRunning" stroke="pink">
        <n-form ref="signInFormRef" :model="formValues" label-width="100px" :rules="rules">
            <n-form-item label="Email" path="email" show-feedback>
                <n-input type="text" v-model:value="formValues.email" placeholder="Enter your email address"
                    :disabled="state.ajaxRunning" />
            </n-form-item>

            <n-form-item label="Password" path="password" show-feedback>
                <n-input v-model:value="formValues.password" type="password" placeholder="Enter your password"
                    show-password-on="click" :disabled="state.ajaxRunning" @keydown.enter="validateForm">
                    <template #password-visible-icon>
                        <n-icon :size="16" :component="IconEyeCancel" />
                    </template>
                    <template #password-invisible-icon>
                        <n-icon :size="16" :component="IconEye" />
                    </template>
                </n-input>
            </n-form-item>
            <n-form-item>
                <n-button secondary @click="validateForm" block :disabled="state.ajaxRunning">{{
                    t("Sign in")
                }}</n-button>
            </n-form-item>
        </n-form>
    </n-spin>
    <RemoteAPIAlert v-if="state.ajaxErrorMessage" type="error" title="Login error" :message="state.ajaxErrorMessage" />
</template>

<style lang="css" scoped></style>