<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount, type CSSProperties, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NFlex, NButton, NForm, type FormItemRule, type FormInst, type FormRules, NIcon, NFormItem, NInput } from 'naive-ui';
    import { IconCancel, IconDeviceFloppy, IconEdit, IconPlus } from '@tabler/icons-vue';

    import { TimeTracking } from '../models/time-tracking.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { timeTrackingService } from '../services/time-tracking.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { TimeTrackingResponse, AddRequest, } from '../types/dto';
    import type { FormMode } from '../../../shared/types/form-mode';
    import { appBus } from '../../../shared/composables/bus';
    import TimeFieldsInput from '../../../shared/components/form-blocks/TimeFieldsInput.vue';

    interface TimeTrackingFormProps {
        mode: FormMode;
        projectId: string;
        taskId: string;
        style?: string | CSSProperties;
    }

    const emit = defineEmits(['add', 'update', 'cancel'])

    const props = defineProps<TimeTrackingFormProps>();

    const { t } = useI18n();

    const timeTracking = ref<TimeTracking>(new TimeTracking());

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const timeTrackingFormRef = ref<FormInst | null>(null)

    const timeTrackingFormRules: FormRules =
    {
        summary: {
            required: true,
            validator: (_rule: FormItemRule, value: string) => {
                if (state.ajaxRunning) {
                    return true;
                }
                if (!value?.trim()) {
                    return new Error(t("shared.warningMessages.fieldIsRequired"));
                } else {
                    return true;
                }
            },
            trigger: ['blur'],
        },
    };

    const serverErrors = ref<Record<string, string>>({});

    const isSaveDisabled = computed<boolean>(() => {
        return !timeTracking.value.summary || timeTracking.value.spentTime <= 0;
    });

    // TODO: allow updates
    const onSave = async () => {
        serverErrors.value = {};
        timeTrackingFormRef.value?.restoreValidation();
        try {
            await timeTrackingFormRef.value?.validate();
            if (props.mode === "add") {
                await onAdd();
            } else {
                console.error("invalid form mode", { file: "TimeTrackingForm.vue", method: "onSave" });
            }
        }
        catch (error: any) {
            console.warn("Warning", { file: "TimeTrackingForm.vue", method: "onSave", details: "form validation error", error: error });
        }
    };

    const onCancel = () => {
        emit('cancel')
    }

    const onAdd = async () => {
        if (props.projectId) {
            serverErrors.value = {};
            timeTrackingFormRef.value?.restoreValidation();
            Object.assign(state, defaultAjaxStateRunning);
            try {
                const payload: AddRequest = {
                    summary: timeTracking.value.summary,
                    spentTime: timeTracking.value.spentTime,
                };
                const addedTimeTracking: TimeTrackingResponse = await timeTrackingService.addTaskTimeTracking(props.projectId, props.taskId, payload);
                emit('add', new TimeTracking(addedTimeTracking));
            } catch (error: unknown) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "TimeTrackingForm.onAdd" } });
                                break;
                            default:
                                state.ajaxErrorMessage = t("modules.projectPermission.components.TimeTrackingForm.errors.addError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("modules.projectPermission.components.TimeTrackingForm.errors.addError");
                        console.error("Unhandled API error", { file: "TimeTrackingForm.vue", method: "onAdd" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrors) {
                    if (state.ajaxErrorMessage) {
                        appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                    } else {
                        await nextTick();
                        timeTrackingFormRef.value?.validate().then(() => { }).catch(() => { });
                    }
                }
            }
        } else {
            console.error("project id not set", { file: "TimeTrackingForm.vue", method: "onAdd" });
        }
    };
    let stopBusReauthListener: () => void;

    onMounted(() => {
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("TimeTrackingForm.onAdd")) {
                onAdd();
            }
        });
        if (props.mode !== "add") {
            console.error("invalid form mode", { file: "TimeTrackingForm.vue", method: "onSave" });
        }
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-card :style="style" bordered>
        <template #header>
            <div class="doneo-flex-center-align">
                <!-- TOOD: icon alignment ??? -->
                <n-icon :component="props.mode == 'add' ? IconPlus : IconEdit" />
                {{ t("modules.timeTracking.components.TimeTrackingForm.headers.addTimeTracking") }}
            </div>
        </template>
        <template #header-extra>
            <n-spin v-if="state.ajaxRunning" size="small" />
        </template>
        <n-form ref="timeTrackingFormRef" :model="timeTracking" :rules="timeTrackingFormRules"
            :disabled="state.ajaxRunning">
            <n-form-item :label="t('modules.timeTracking.components.TimeTrackingForm.inputs.summary.label')">
                <n-input type="textarea"
                    :placeholder="t('modules.timeTracking.components.TimeTrackingForm.inputs.summary.placeholder')"
                    v-model:value="timeTracking.summary" />
            </n-form-item>
            <TimeFieldsInput input-type="spent" v-model:seconds="timeTracking.spentTime" />
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