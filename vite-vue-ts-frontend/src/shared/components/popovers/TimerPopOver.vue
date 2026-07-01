<script setup lang="ts">
    import { ref, reactive, computed, watch, onMounted, onBeforeUnmount, nextTick } from "vue";
    import { useI18n } from "vue-i18n";

    import { NInputGroup, NInput, NButton, NIcon, NPopover, NCard, type InputInst, NCollapse, NCollapseItem, type ButtonType } from 'naive-ui';
    import { IconAlarm, IconClockCancel, IconClockPlay, IconClockStop, IconTrash } from '@tabler/icons-vue';

    import { appBus } from '../../../shared/composables/bus';

    import { useUserSettingsStore } from '../../../stores/userSettings.ts';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from "../../../shared/types/ajaxState";
    import { userTimerService } from "../../../modules/user-timer/services/user-timer";
    import { handleAPIError } from "../../../api/client/errorHandler.ts";
    import { type UserTimerResponse } from "../../../modules/user-timer/types/dto";

    import { IDate } from "../../types/idate";

    import { BUTTON_DEFAULT_ICON_SIZE } from '../../../constants';

    interface IProps {
        iconSize?: number,
        disabled?: boolean,
    };

    const props = withDefaults(defineProps<IProps>(), {
        iconSize: BUTTON_DEFAULT_ICON_SIZE,
        disabled: false,
    });

    const { t } = useI18n();
    const userSettingsStore = useUserSettingsStore();

    const newTimerSummary = ref<string>("");
    const newTimerSummaryRef = ref<InputInst | null>(null);

    const showPopOver = ref<boolean>(false);

    const start = ref<number | null>(Date.now());
    const now = ref<number>(Date.now())
    const timers = ref<UserTimerResponse[]>([]);

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const hasTimers = computed<boolean>(() => timers.value.length > 0);
    const currentActiveTimer = computed<UserTimerResponse | undefined>(() => timers.value.find((timer) => timer.finishedAt === null));
    const hasTimerRunning = computed<boolean>(() => typeof currentActiveTimer.value !== "undefined");

    const currentButtonIconColor = computed<string | undefined>(() => {
        return hasTimers.value ? "primary" : undefined;
    });

    const currentButtonType = computed<ButtonType>(() => {
        if (hasTimers.value) {
            if (hasTimerRunning.value) {
                return ("primary");
            } else {
                return ("info");
            }
        } else {
            return ("default");
        }
    });

    const currentTimerElapsedSeconds = computed<number>(() => Math.floor((now.value - (start.value ?? 0)) / 1000));

    const hasFinishedTimers = computed<boolean>(() => finishedTimers.value.length > 0);

    const finishedTimers = computed<UserTimerResponse[]>(() => timers.value.filter((timer: UserTimerResponse) => timer.finishedAt !== null));

    watch(hasTimerRunning, (running: boolean) => {
        if (running) {
            onStartInterval();
        } else {
            onStopInterval();
        }
    });

    const formatDuration = (spentTime: number): string => {
        // TODO: i18N
        const days = Math.floor(spentTime / 86400);
        const hours = Math.floor((spentTime % 86400) / 3600);
        const minutes = Math.floor((spentTime % 3600) / 60);
        const seconds = spentTime % 60;
        const parts: string[] = [];
        if (days > 0) parts.push(`${days}d`);
        if (hours > 0) parts.push(`${hours}h`);
        if (minutes > 0) parts.push(`${minutes}m`);
        if (seconds > 0 || parts.length === 0) parts.push(`${seconds}s`);
        return parts.join(" ");
    }

    const onGetTimers = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response = await userTimerService.search();
            timers.value = response.userTimers;
            start.value = timers.value.find((timer) => timer.finishedAt === null)?.startedAt ?? null;
        } catch (error) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TimerPopOver.onGetTimers" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("shared.components.popOvers.TimerPopOver.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("shared.components.popOvers.TimerPopOver.errors.refreshError");
                    console.error("Unhandled API error", { file: "TimerPopOver.vue", method: "onGetTimers" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onStartTimer = async () => {
        if (newTimerSummary.value) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await userTimerService.start(newTimerSummary.value);
                showPopOver.value = false;
                await onGetTimers();
                newTimerSummary.value = "";
            } catch (error) {
                state.ajaxErrors = true;
                handleAPIError(error,
                    (apiError) => {
                        switch (apiError.response?.status) {
                            case 401:
                                state.ajaxErrors = false;
                                appBus.emit({ type: "reauthRequired", payload: { emitter: "TimerPopOver.onStartTimer" } });
                                break;
                            default:
                                state.ajaxErrorMessage = t("shared.components.popOvers.TimerPopOver.errors.startError");
                                break;
                        }
                    },
                    (fatalError) => {
                        state.ajaxErrorMessage = t("shared.components.popOvers.TimerPopOver.errors.startError");
                        console.error("Unhandled API error", { file: "TimerPopOver.vue", method: "onStartTimer" }, { err: fatalError });
                    });
            } finally {
                state.ajaxRunning = false;
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                }
            }
        }
    };

    let stopTimerId: string = "";

    const onStopTimer = async (id: string) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await userTimerService.stop(id);
            showPopOver.value = false;
            await onGetTimers();
        } catch (error) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TimerPopOver.onStopTimer" } });
                            stopTimerId = id;
                            break;
                        default:
                            state.ajaxErrorMessage = t("shared.components.popOvers.TimerPopOver.errors.stopError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("shared.components.popOvers.TimerPopOver.errors.stopError");
                    console.error("Unhandled API error", { file: "TimerPopOver.vue", method: "onStopTimer" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    let deleteTimerId: string = "";

    const onDeleteTimer = async (id: string) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await userTimerService.delete(id);
            await onGetTimers();
        } catch (error) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TimerPopOver.onDeleteTimer" } });
                            deleteTimerId = id;
                            break;
                        default:
                            state.ajaxErrorMessage = t("shared.components.popOvers.TimerPopOver.errors.deleteError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("shared.components.popOvers.TimerPopOver.errors.deleteError");
                    console.error("Unhandled API error", { file: "TimerPopOver.vue", method: "onDeleteTimer" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onClearTimers = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await userTimerService.clear();
            showPopOver.value = false;
            await onGetTimers();
        } catch (error) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "TimerPopOver.onClearTimers" } });
                            break;
                        default:
                            state.ajaxErrorMessage = t("shared.components.popOvers.TimerPopOver.errors.clearError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("shared.components.popOvers.TimerPopOver.errors.clearError");
                    console.error("Unhandled API error", { file: "TimerPopOver.vue", method: "onClearTimers" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onShowPopOver = () => {
        if (!state.ajaxRunning) {
            if (!hasTimerRunning.value) {
                nextTick().then(() => {
                    newTimerSummaryRef.value?.focus();
                }).catch((e) => { console.error(e); });
            }
        }
    };

    let interval: number | undefined;

    const onStartInterval = () => {
        interval = setInterval(() => {
            now.value = Date.now()
        }, 1000);
    };

    const onStopInterval = () => {
        clearInterval(interval)
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onGetTimers();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("TimerPopOver.onGetTimers")) {
                onGetTimers();
            } else if (payload.to.includes("TimerPopOver.onStartTimer")) {
                onStartTimer();
            } else if (payload.to.includes("TimerPopOver.onStopTimer")) {
                onStopTimer(stopTimerId);
            } else if (payload.to.includes("TimerPopOver.onDeleteTimer")) {
                onDeleteTimer(deleteTimerId);
            } else if (payload.to.includes("TimerPopOver.onClearTimers")) {
                onClearTimers();
            }
        });
    });

    onBeforeUnmount(() => {
        onStopInterval();
        stopBusReauthListener();
    });
</script>

<template>
    <n-popover placement="bottom" trigger="hover" v-model:show="showPopOver" @update-show="onShowPopOver"
        class="doneo-disable-user-select">
        <template #trigger>
            <n-button quaternary :disabled="props.disabled || state.ajaxRunning" :type="currentButtonType">
                <template #icon>
                    <n-icon :component="IconAlarm" :size="props.iconSize" :color="currentButtonIconColor"
                        :class="{ 'doneo-timer-animated-icon': hasTimerRunning }" />
                </template>
            </n-button>
        </template>
        <n-card size="small" segmented>
            <template #header>
                <n-button size="small" block @click="onStopTimer(currentActiveTimer?.id ?? '')" v-if="hasTimerRunning"
                    :disabled="props.disabled || state.ajaxRunning">
                    <template #icon>
                        <n-icon :component="IconClockStop" :size="props.iconSize" />
                    </template>
                    {{ t("shared.components.popOvers.TimerPopOver.buttons.stopCurrentTimer.label") }} {{
                        formatDuration(currentTimerElapsedSeconds > 0 ? currentTimerElapsedSeconds :
                            0) }}
                </n-button>
                <n-input-group v-else>
                    <n-input size="small" ref="newTimerSummaryRef"
                        :placeholder="t('shared.components.popOvers.TimerPopOver.inputs.newTimer.placeholder')"
                        :minlength="1" :maxlength="32" show-count v-model:value="newTimerSummary"
                        :disabled="props.disabled || state.ajaxRunning" @keydown.enter="onStartTimer" />
                    <n-button size="small" @click="onStartTimer"
                        :disabled="props.disabled || state.ajaxRunning || !newTimerSummary">
                        <template #icon>
                            <n-icon :component="IconClockPlay" :size="props.iconSize" />
                        </template>
                        {{ t("shared.components.popOvers.TimerPopOver.buttons.startNewTimer.label") }}
                    </n-button>
                </n-input-group>
            </template>
            <template #footer v-if="hasFinishedTimers">
                <p>{{ t("shared.components.popOvers.TimerPopOver.labels.previousTimers") }}</p>
                <n-collapse display-directive="if" accordion :trigger-areas="['main', 'arrow']">
                    <n-collapse-item v-for="timer in finishedTimers" :key="timer.id" :name="timer.id"
                        :title="timer.summary">
                        <template #header-extra>
                            <n-button block size="tiny" @click="onDeleteTimer(timer.id)"
                                :disabled="props.disabled || state.ajaxRunning">
                                <template #icon>
                                    <n-icon :component="IconTrash" :size="props.iconSize" />
                                </template>
                            </n-button>
                        </template>
                        <div>{{ t("shared.components.popOvers.TimerPopOver.labels.timerStarted") }} {{ new
                            IDate(timer.startedAt).toCustomMaskString(userSettingsStore.currentDatetimeMask) }}</div>
                        <div>{{ t("shared.components.popOvers.TimerPopOver.labels.timerFinished") }} {{ new
                            IDate(timer.finishedAt).toCustomMaskString(userSettingsStore.currentDatetimeMask) }}</div>
                        <div>{{ t("shared.components.popOvers.TimerPopOver.labels.totalTimer") }} {{
                            formatDuration(
                                Math.round(
                                    ((timer.finishedAt ?? new Date().getTime()) - timer.startedAt) / 1000
                                )
                            )
                        }}</div>
                    </n-collapse-item>
                </n-collapse>
            </template>
            <template #action v-if="hasTimers">
                <n-button size="small" block @click="onClearTimers" :disabled="props.disabled || state.ajaxRunning">
                    <template #icon>
                        <n-icon :component="IconClockCancel" :size="props.iconSize" />
                    </template>
                    {{ t("shared.components.popOvers.TimerPopOver.buttons.clearTimers.label") }}
                </n-button>
            </template>
        </n-card>
    </n-popover>
</template>

<style lang="css" scoped>
    .doneo-timer-animated-icon {
        animation: ring 1s ease-in-out infinite;
        transform-origin: top center;
    }

    @keyframes ring {
        0%,

        100% {
            transform: rotate(0deg);
        }

        25% {
            transform: rotate(-12deg);
        }

        75% {
            transform: rotate(12deg);
        }
    }
</style>