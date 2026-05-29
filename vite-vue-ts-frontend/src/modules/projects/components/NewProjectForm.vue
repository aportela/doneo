<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, watch, type CSSProperties, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NInput, NFlex, NButton, NForm, NFormItem, type FormItemRule, type FormInst, type FormRules, NIcon } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconPlus } from '@tabler/icons-vue';

    import { Project, MAX_KEY_LENGTH, MAX_SUMMARY_LENGTH } from '../models/project';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { projectService } from '../services/project';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { ProjectResponse, AddRequest } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';

    interface NewProjectFormProps {
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'cancel'])

    const props = defineProps<NewProjectFormProps>();

    const { t } = useI18n();

    const project = ref<Project>(new Project());

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const newProjectFormRef = ref<FormInst | null>(null)

    const newProjectFormRules: FormRules =
    {
        key: {
            required: true,
            validator: (_rule: FormItemRule, value: string) => {
                if (state.ajaxRunning) {
                    return true;
                }
                if (!value?.trim()) {
                    return new Error(t("shared.warningMessages.fieldIsRequired"));
                }
                else if (value.length > MAX_KEY_LENGTH) {
                    return new Error(t("shared.warningMessages.fieldExceedsMaxLength"));
                } else if (serverErrors.value.name) {
                    return new Error(t(serverErrors.value.name));
                } else {
                    return true;
                }
            },
            trigger: ['blur'],
        },
        summary: {
            required: true,
            validator: (_rule: FormItemRule, value: string) => {
                if (state.ajaxRunning) {
                    return true;
                }
                if (!value?.trim()) {
                    return new Error(t("shared.warningMessages.fieldIsRequired"));
                }
                else if (value.length > MAX_SUMMARY_LENGTH) {
                    return new Error(t("shared.warningMessages.fieldExceedsMaxLength"));
                } else if (serverErrors.value.name) {
                    return new Error(t(serverErrors.value.name));
                } else {
                    return true;
                }
            },
            trigger: ['blur'],
        },
    };

    watch(() => project.value.key, () => { delete serverErrors.value.name });

    const serverErrors = ref<Record<string, string>>({});

    const isSaveDisabled = computed<boolean>(() => {
        return !project.value.key || !project.value.summary;
    });

    const onSave = async () => {
        serverErrors.value = {};
        newProjectFormRef.value?.restoreValidation();
        try {
            await newProjectFormRef.value?.validate();
            await onAdd();
        }
        catch (error: any) {
            console.warn("Warning", { file: "NewProjectForm.vue", method: "onSave", details: "form validation error", error: error });
        }
    };

    const onCancel = () => {
        emit('cancel')
    }


    const onAdd = async () => {
        serverErrors.value = {};
        newProjectFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: AddRequest = {
                key: project.value.key ?? "",
                summary: project.value.summary ?? "",
                description: project.value.description,
                type: { id: project.value.type.id ?? "" },
                priority: { id: project.value.priority.id ?? "" },
                status: { id: project.value.status.id ?? "" },
            };
            const addedProject: ProjectResponse = await projectService.add(payload);
            emit('add', addedProject)
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "NewProjectForm.onAdd" } });
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                serverErrors.value.name = "modules.project.components.NewProjectForm.warnings.nameAlreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.project.components.NewProjectForm.errors.addError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.project.components.NewProjectForm.errors.addError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.project.components.NewProjectForm.errors.addError");
                    console.error("Unhandled API error", { file: "NewProjectForm.vue", method: "onAdd" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    newProjectFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };


    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("NewProjectForm.onAdd")) {
                onAdd();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-card :style="style" bordered>
        <template #header>
            <div class="doneo-flex-center-align">
                <n-icon :component="IconPlus" />
                {{ t("modules.project.components.NewProjectForm.headers.addProject") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="projectTypeFormRef" :model="project" :rules="newProjectFormRules" :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.project.components.NewProjectForm.inputs.key.label')" path="key"
                show-feedback>
                <n-input type="text"
                    :placeholder="t('modules.project.components.NewProjectForm.inputs.key.placeholder')"
                    v-model:value="project.key" :maxlength="MAX_KEY_LENGTH" :show-count="true" clearable required
                    autofocus>
                </n-input>
            </n-form-item>
            <n-form-item :label="t('modules.project.components.NewProjectForm.inputs.summary.label')" path="key"
                show-feedback>
                <n-input type="text"
                    :placeholder="t('modules.project.components.NewProjectForm.inputs.summary.placeholder')"
                    v-model:value="project.summary" :maxlength="MAX_KEY_LENGTH" :show-count="true" clearable required
                    autofocus>
                </n-input>
            </n-form-item>
        </n-form>
        <template #action>
            <n-flex>
                <n-button @click="onSave" :disabled="isSaveDisabled">
                    <template #icon>
                        <n-icon :component="IconDeviceFloppy" />
                    </template>
                    {{ t("shared.buttons.Save.label") }}
                </n-button>
                <n-button @click="onCancel" :disabled="state.ajaxRunning">
                    <template #icon>
                        <n-icon :component="IconCancel" />
                    </template>
                    {{ t("shared.buttons.Cancel.label") }}
                </n-button>
            </n-flex>
        </template>
    </n-card>

</template>

<style lang="css" scoped></style>