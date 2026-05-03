<script setup lang="ts">
    import { ref, reactive, watch } from 'vue';
    import { useI18n } from "vue-i18n";
    import { NIcon, NSpin, NForm, NFormItem, NInput, NButton, type FormItemRule, type FormInst, type FormRules } from 'naive-ui'
    import { IconEye, IconEyeCancel } from '@tabler/icons-vue';
    import { type AjaxStateInterface, defaultAjaxState } from "../../types/ajaxState";
    import { type SignInSucessResponse } from '../../types/apiResponses';
    import { type AxiosAPIError } from '../../composables/axios';
    import { api } from '../../composables/api';
    import { required, minLength, validEmail, runValidators } from '../../composables/form-validators';
    import { createStorageEntry } from '../../composables/localStorage';
    import { useSessionStore } from "../../stores/session";
    import { default as RemoteAPIAlert } from '../alerts/RemoteAPIAlert.vue';

    const emit = defineEmits(['success'])

    const { t } = useI18n();

    const sessionStore = useSessionStore();

    const localStorageLastUsedEmail = createStorageEntry<string | null>("lastUsedEmail", null);

    const lastUsedEmail = localStorageLastUsedEmail.get();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const signInFormRef = ref<FormInst | null>(null)

    type signInFormValuesInterface = {
        email: string | null,
        password: string | null,
    };

    const signinFormValues = ref<signInFormValuesInterface>({
        email: lastUsedEmail || null,
        password: null,
    });

    const serverErrors = ref<Record<string, string>>({});

    const signInFormRules: FormRules = {
        email: {
            validator: (_rule: FormItemRule, value) => {
                const localResult = runValidators(value, [required('Email'), validEmail])
                if (localResult !== true) return localResult
                if (serverErrors.value.email) return new Error(serverErrors.value.email)
                return true
            },
            trigger: ['blur'],
        },
        password: {
            validator: (_rule: FormItemRule, value) => {
                const localResult = runValidators(value, [required('Password'), minLength(4)])
                if (localResult !== true) return localResult
                if (serverErrors.value.password) return new Error(serverErrors.value.password)
                return true
            },
            trigger: ['blur']
        }
    }

    watch(() => signinFormValues.value.email, () => {
        delete serverErrors.value.email
    });

    watch(() => signinFormValues.value.password, () => {
        delete serverErrors.value.password
    });

    const onSubmit = () => {
        serverErrors.value = {}
        if (signinFormValues.value.email && signinFormValues.value.password) {
            Object.assign(state, defaultAjaxState);
            state.ajaxRunning = true;
            api.auth.signIn(signinFormValues.value.email, signinFormValues.value.password)
                .then((successResponse: SignInSucessResponse) => {
                    if (successResponse.data.accessToken) {
                        sessionStore.setAccessToken(successResponse.data.accessToken.token, successResponse.data.accessToken.expiresAtTimestamp);
                        localStorageLastUsedEmail.set(signinFormValues.value.email);
                        emit('success');
                    } else {
                        state.ajaxErrorMessage = t("Invalid API response");
                    }
                })
                .catch((errorResponse: AxiosAPIError) => {
                    state.ajaxErrors = true;
                    if (errorResponse.isAPIError) {
                        state.ajaxAPIErrorDetails = errorResponse.customAPIErrorDetails;
                        switch (errorResponse.status) {
                            case 404:
                                serverErrors.value.email = t("Email not found");
                                break;
                            case 401:
                                serverErrors.value.password = t("Invalid password");
                                break;
                            default:
                                state.ajaxErrorMessage = t("API Error: fatal error");
                                break;

                        }
                        signInFormRef.value?.restoreValidation();
                        signInFormRef.value?.validate().catch((e: any) => {
                            console.error("TODO: SignIn form validation error", e);
                        });
                    }
                    else {
                        // TODO: i18n
                        state.ajaxErrorMessage = `Uncaught exception (${errorResponse.status})\n\n${errorResponse}`;
                        console.error("Unhandled API error", errorResponse);
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
        catch (e: any) {
            //console.debug("SignIn form validation error", e)
        }
    }
</script>

<template>
    <h2 class="title">{{ t("Login to your account") }}</h2>
    <n-spin :show="state.ajaxRunning" stroke="pink">
        <n-form ref="signInFormRef" :model="signinFormValues" label-width="100px" :rules="signInFormRules">
            <n-form-item :label="t('Email')" path="email" show-feedback>
                <n-input type="text" :autofocus="true" v-model:value="signinFormValues.email"
                    placeholder="Enter your email address" :disabled="state.ajaxRunning" />
            </n-form-item>
            <n-form-item :label="t('Password')" path="password" show-feedback>
                <n-input v-model:value="signinFormValues.password" type="password" placeholder="Enter your password"
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
    <RemoteAPIAlert v-if="state.ajaxErrorMessage" type="error" :title="t('Sign in error')"
        :message="state.ajaxErrorMessage" />
</template>

<style lang="css" scoped></style>