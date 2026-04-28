<script setup lang="ts">
    import { ref, reactive, watch } from 'vue';
    import { useRouter } from "vue-router";
    import { NIcon, NSpin, NForm, NFormItem, NInput, NButton, type FormItemRule, type FormInst, type FormRules } from 'naive-ui'
    import { IconEye, IconEyeCancel } from '@tabler/icons-vue';
    import { api } from '../../composables/api';
    import { required, minLength, validEmail, runValidators } from '../../composables/form-validators';
    import { useSessionStore } from "../../stores/session";
    import { type AjaxState as AjaxStateInterface, defaultAjaxState } from "../../types/ajaxState";

    import { createStorageEntry } from '../../composables/localStorage';


    const router = useRouter();

    const sessionStore = useSessionStore();

    const lastUsedEmail = createStorageEntry<string | null>("lastUsedEmail", null);

    const savedEmail = lastUsedEmail.get();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const signInFormRef = ref<FormInst | null>(null)

    const formValues = ref({
        email: savedEmail,
        password: null
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
            state.ajaxRunning = true;
            api.auth.signIn(formValues.value.email, formValues.value.password)
                .then((successResponse: any) => {
                    if (successResponse.data.accessToken) {
                        sessionStore.setAccessToken(successResponse.data.accessToken.token, successResponse.data.accessToken.expiresAtTimestamp);
                        lastUsedEmail.set(formValues.value.email);

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
                            serverErrors.value.email = "Email not found";
                            break;
                        case 401:
                            serverErrors.value.password = "Invalid password";
                            break;
                    }
                    signInFormRef.value?.restoreValidation();
                    signInFormRef.value?.validate().catch(() => { });
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
                <n-button secondary @click="validateForm" block :disabled="state.ajaxRunning">Sign in</n-button>
            </n-form-item>
        </n-form>
    </n-spin>
</template>

<style lang="css" scoped></style>