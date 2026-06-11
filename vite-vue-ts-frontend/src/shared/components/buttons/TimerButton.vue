<script setup lang="ts">
    import { ref, reactive, computed, watch, onMounted, onBeforeUnmount, nextTick } from "vue";
    //import { useI18n } from "vue-i18n";

    import { NInputGroup, NInput, NButton, NIcon, NPopover, NCard, type InputInst, NCollapse, NCollapseItem, type ButtonType } from 'naive-ui';
    import { IconAlarm, IconClockCancel, IconClockPlay, IconClockStop, IconTrash } from '@tabler/icons-vue';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from "../../../shared/types/ajaxState";

    import { timerService } from "../../../modules/timer/services/timer";
    import { type TimerResponse } from "../../../modules/timer/types/dto";
    import { IDate } from "../../types/idate";

    interface TimerButtonProps {
        iconSize?: number,
    };

    withDefaults(defineProps<TimerButtonProps>(), {
        iconSize: 20
    });

    const newTimerSummary = ref<string>("");
    const newTimerSummaryRef = ref<InputInst | null>(null);

    const showPopOver = ref<boolean>(false);

    const start = ref<number | null>(Date.now());
    const now = ref<number>(Date.now())

    //const { t } = useI18n();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const timers = ref<TimerResponse[]>([]);

    const hasTimers = computed(() => timers.value.length > 0);
    const currentActiveTimer = computed<TimerResponse | undefined>(() => timers.value.find((timer) => timer.finishedAt === null));
    const hasTimerRunning = computed(() => typeof currentActiveTimer.value !== "undefined");

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

    const currentTimerElapsedSeconds = computed(() => {
        return Math.floor((now.value - (start.value ?? 0)) / 1000)
    });

    const hasFinishedTimers = computed(() => finishedTimers.value.length > 0);

    const finishedTimers = computed(() => timers.value.filter((timer) => timer.finishedAt !== null));

    watch(hasTimerRunning, (running) => {
        if (running) {
            onStartInterval();
        } else {
            onStopInterval();
        }
    });

    const formatDuration = (totalSeconds: number): string => {
        const days = Math.floor(totalSeconds / 86400);
        const hours = Math.floor((totalSeconds % 86400) / 3600);
        const minutes = Math.floor((totalSeconds % 3600) / 60);
        const seconds = totalSeconds % 60;
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
            const response = await timerService.search();
            timers.value = response.timers;
            start.value = timers.value.find((timer) => timer.finishedAt === null)?.startedAt ?? null;
        } catch (e) {
            // TODO:
            console.error(e);
        } finally {
            state.ajaxRunning = false;
        }
    };

    const onStartTimer = async () => {
        if (newTimerSummary.value) {
            Object.assign(state, defaultAjaxStateRunning);
            try {
                await timerService.start(newTimerSummary.value);
                showPopOver.value = false;
                await onGetTimers();
                newTimerSummary.value = "";

            } catch (e) {
                // TODO:
                console.error(e);
            } finally {
                state.ajaxRunning = false;
            }
        }
    };

    const onStopTimer = async (id: string) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await timerService.stop(id);
            showPopOver.value = false;
            await onGetTimers();
        } catch (e) {
            // TODO:
            console.error(e);
        } finally {
            state.ajaxRunning = false;
        }
    };

    const onDeleteTimer = async (id: string) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await timerService.delete(id);
            await onGetTimers();
        } catch (e) {
            // TODO:
            console.error(e);
        } finally {
            state.ajaxRunning = false;
        }
    };

    const onClearTimers = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await timerService.clear();
            showPopOver.value = false;
            await onGetTimers();
        } catch (e) {
            // TODO:
            console.error(e);
        } finally {
            state.ajaxRunning = false;
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
    }

    const onStopInterval = () => {
        clearInterval(interval)
    };

    onMounted(() => {
        onGetTimers();
    });

    onBeforeUnmount(() => {
        onStopInterval();
    })
</script>

<template>
    <n-popover placement="bottom" trigger="hover" v-model:show="showPopOver" @update-show="onShowPopOver">
        <template #trigger>
            <n-button quaternary :disabled="state.ajaxRunning" :type="currentButtonType">
                <template #icon>
                    <n-icon :component="IconAlarm" :size="iconSize" :color="currentButtonIconColor"
                        :class="{ 'doneo-timer-animated-icon': hasTimerRunning }" />
                </template>
            </n-button>
        </template>
        <n-card size="small" segmented>
            <template #header>
                <n-button size="small" block @click="onStopTimer(currentActiveTimer?.id ?? '')" v-if="hasTimerRunning"
                    :disabled="state.ajaxRunning">
                    <template #icon>
                        <n-icon :component="IconClockStop" />
                    </template>
                    Stop current timer: {{ formatDuration(currentTimerElapsedSeconds > 0 ? currentTimerElapsedSeconds :
                        0) }}
                </n-button>
                <n-input-group v-else>
                    <n-input size="small" ref="newTimerSummaryRef" placeholder="type new timer summary" :minlength="1"
                        :maxlength="32" show-count v-model:value="newTimerSummary" :disabled="state.ajaxRunning"
                        @keydown.enter="onStartTimer" />
                    <n-button size="small" @click="onStartTimer" :disabled="state.ajaxRunning || !newTimerSummary">
                        <template #icon>
                            <n-icon :component="IconClockPlay" />
                        </template>
                        Start new timer
                    </n-button>
                </n-input-group>
            </template>
            <template #footer v-if="hasFinishedTimers">
                <p>Previous timers</p>
                <n-collapse display-directive="if" accordion :trigger-areas="['main', 'arrow']">
                    <n-collapse-item v-for="timer in finishedTimers" :key="timer.id" :name="timer.id"
                        :title="timer.summary">
                        <template #header-extra>
                            <n-button block size="tiny" @click="onDeleteTimer(timer.id)" :disabled="state.ajaxRunning">
                                <template #icon>
                                    <n-icon :component="IconTrash" />
                                </template>
                            </n-button>
                        </template>
                        <dt>
                        <dd>Started: {{ new IDate(timer.startedAt).toLocaleString() }}</dd>
                        <dd>Finished: {{ new IDate(timer.finishedAt).toLocaleString() }}</dd>
                        <dd>Total elapsed time: {{
                            formatDuration(
                                Math.round(
                                    ((timer.finishedAt ?? new Date().getTime()) - timer.startedAt) / 1000
                                )
                            )
                        }}</dd>
                        </dt>
                    </n-collapse-item>
                </n-collapse>
            </template>
            <template #action v-if="hasTimers">
                <n-button size="small" block @click="onClearTimers" :disabled="state.ajaxRunning">
                    <template #icon>
                        <n-icon :component="IconClockCancel" />
                    </template>
                    Clear timers
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