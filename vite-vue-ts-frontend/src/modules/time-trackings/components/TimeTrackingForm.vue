<script setup lang="ts">
    import { ref, reactive, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSpin, NCard, NFlex, NButton, NForm, type FormItemRule, type FormInst, type FormRules, NIcon, NFormItem, NInput, NSelect, type SelectOption } from 'naive-ui';

    import { IconCancel, IconDeviceFloppy } from '@tabler/icons-vue';

    import { TimeTracking } from '../models/time-tracking.ts';
    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { timeTrackingService } from '../services/time-tracking.ts';
    import { handleAPIError } from '../../../api/client/errorHandler';
    import type { TimeTrackingResponse, AddRequest, } from '../types/dto';
    import { appBus } from '../../../shared/composables/bus';
    import TimeFieldsInput from '../../../shared/components/form-blocks/TimeFieldsInput.vue';
    import { userTimerService } from '../../user-timer/services/user-timer.ts';
    import { type UserTimerResponse } from '../../user-timer/types/dto.ts';
    import { formatDuration } from '../../../shared/composables/datetime.ts';
    import { DONEO_ICON_ADD } from '../../../shared/types/icons.ts';

    interface Props {
        projectId: string;
        taskId: string;
    }

    const props = defineProps<Props>();

    const emit = defineEmits(['add', 'cancel'])

    const { t } = useI18n();

    const currentUserTimers = ref<UserTimerResponse[]>([]);

    const opts = computed<SelectOption[]>(() => {
        if (currentUserTimers.value.length > 0) {
            return [{ label: "Enter time manually", value: "0" }, { type: "group", label: "Current timers", key: "currentTimers", children: currentUserTimers.value.map((item) => { return ({ label: item.summary, value: item.id }); }) }];
        } else {
            return [{ label: "Enter time manually", value: "0" }];
        }
    });


    const selectedOpt = ref<string>("0");

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

    watch(selectedOpt, (newValue) => {
        if (newValue === "0") {
            timeTracking.value.spentTime = 0;
        } else {
            let selectedUserTimer = currentUserTimers.value.find((item) => item.id === newValue);
            timeTracking.value.spentTime = Math.round(((selectedUserTimer?.finishedAt ?? 0) - (selectedUserTimer?.startedAt ?? 0)) / 1000);
        }
    });

    const onGetTimers = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response = await userTimerService.search();
            currentUserTimers.value = response.userTimers.filter((item) => item.finishedAt !== null).map((item) => {
                item.summary = `${item.summary} (${formatDuration(Math.round(((item.finishedAt ?? Date.now()) - item.startedAt) / 1000))})`;
                return item;
            });
        } catch (error) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TimeTrackingForm.onGetTimers" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("shared.components.popOvers.TimerPopOver.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("shared.components.popOvers.TimeTrackingForm.errors.refreshError");
                    console.error("Unhandled API error", { file: "TimeTrackingForm.vue", method: "onGetTimers" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    // TODO: allow updates
    const onSave = async () => {
        serverErrors.value = {};
        timeTrackingFormRef.value?.restoreValidation();
        try {
            await timeTrackingFormRef.value?.validate();
            await onAdd();
        }
        catch (error: any) {
            console.warn("Warning", { file: "TimeTrackingForm.vue", method: "onSave", details: "form validation error", error: error });
        }
    };

    const onCancel = () => {
        emit('cancel');
    }

    const onAdd = async () => {
        if (props.projectId) {
            serverErrors.value = {};
            timeTrackingFormRef.value?.restoreValidation();
            Object.assign(state, defaultAjaxStateRunning);
            try {
                let selectedUserTimer = currentUserTimers.value.find((item) => item.id === selectedOpt.value);
                const payload: AddRequest = {
                    summary: timeTracking.value.summary,
                    spentTime: selectedOpt.value === "0" ? timeTracking.value.spentTime : Math.round(((selectedUserTimer?.finishedAt ?? 0) - (selectedUserTimer?.startedAt ?? 0)) / 1000),
                    userTimerId: selectedOpt.value !== "0" ? selectedOpt.value : null,
                };
                const addedTimeTracking: TimeTrackingResponse = await timeTrackingService.addTaskTimeTracking(props.projectId, props.taskId, payload);
                if (selectedOpt.value !== "0") {
                    appBus.emit({ type: "refreshUserTimers", payload: {} });
                }
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
        onGetTimers();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("TimeTrackingForm.onAdd")) {
                onAdd();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });
</script>

<template>
    <n-card bordered>
        <template #header>
            <div class="doneo-flex-center-align">
                <n-icon class="doneo-mr-4px" :component="DONEO_ICON_ADD" />
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
                    v-model:value="timeTracking.summary" :disabled="state.ajaxRunning" />
            </n-form-item>
            <n-form-item>
                <n-select :options="opts" v-model:value="selectedOpt"
                    :disabled="state.ajaxRunning || currentUserTimers.length < 1" />
            </n-form-item>
            <TimeFieldsInput input-type="spent" v-model:seconds="timeTracking.spentTime" :disabled="state.ajaxRunning"
                v-if="selectedOpt === '0'" />
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