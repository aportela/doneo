<script setup lang="ts">
    import { ref } from 'vue';

    import { NIcon, NForm, NFormItem, NInput, NButton, type FormInst } from 'naive-ui'
    import { IconEye, IconEyeCancel } from '@tabler/icons-vue';


    const signInFormRef = ref<FormInst | null>(null)

    const rules = {

        formData: {
            email: {
                required: true,
                trigger: ['input', 'blur'],
                message: 'Email is required'
            },
            password: {
                required: true,
                trigger: ['input', 'blur'],
                message: 'password is required'
            }
        }
    };

    const formData = ref({
        email: '',
        password: ''
    });

    const validateForm = () => {
        signInFormRef.value?.validate((valid: any) => {
            if (valid) {
                alert('Formulario válido')
            } else {
                alert('Hay errores en el formulario')
            }
        })
    }

</script>

<template>

    <n-form ref="signInFormRef" :model="formData" label-width="100px" @submit.prevent="validateForm" :rules="rules">
        <n-form-item label="Email" prop="email" show-feedback>
            <n-input v-model:value="formData.email" placeholder="Enter your email address" />
        </n-form-item>

        <n-form-item label="Password" prop="password" show-feedback>
            <n-input v-model="formData.password" type="password" placeholder="Enter your password"
                show-password-on="click">
                <template #password-visible-icon>
                    <n-icon :size="16" :component="IconEyeCancel" />
                </template>
                <template #password-invisible-icon>
                    <n-icon :size="16" :component="IconEye" />
                </template>
            </n-input>
        </n-form-item>

        <n-form-item>
            <n-button secondary @click="validateForm" block>Sign in</n-button>
        </n-form-item>
    </n-form>

</template>

<style lang="css" scoped></style>