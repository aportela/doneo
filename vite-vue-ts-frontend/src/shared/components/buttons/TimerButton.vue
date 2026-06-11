<script setup lang="ts">
    import { ref, reactive, computed, onMounted, onBeforeUnmount } from "vue";
    //import { useI18n } from "vue-i18n";

    import { NButton, NIcon, NPopover, NCard, NList, NListItem } from 'naive-ui';
    import { IconAlarm, IconClock, IconClockCancel, IconClockPlay, IconClockStop } from '@tabler/icons-vue';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from "../../../shared/types/ajaxState";
    import { useNotify } from '../../../shared/composables/notification';

    import { timerService } from "../../../modules/timer/services/timer";
    import { type TimerResponse } from "../../../modules/timer/types/dto";
    import { IDate } from "../../types/idate";

    interface SwitchNotificationsButtonProps {
        iconSize?: number,
    };

    withDefaults(defineProps<SwitchNotificationsButtonProps>(), {
        iconSize: 20
    });

    //const { t } = useI18n();
    const { notify } = useNotify();

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const timer = ref<boolean>(false);

    const color = computed<string | undefined>(() => {
        return hasTimers.value ? "red" : undefined;
    });

    const start = ref<number | null>(null);
    const stop = ref<number | null>(null);
    const now = ref(Date.now())


    const onToggleTimer = () => {
        if (timer.value) {
            stop.value = new Date().getTime();
            const totalTime = stop.value - (start.value ?? 0);
            notify('success', "Total time: " + totalTime + " ms");
            console.log("Total time: " + totalTime / 1000 + "sec");
        } else {
            start.value = new Date().getTime();
        }
        timer.value = !timer.value;
    };

    const commonIconSize = 22;

    const hasTimers = computed(() => timers.value.length > 0);
    const currentActiveTimer = computed<TimerResponse | undefined>(() => timers.value.find((timer) => timer.finishedAt === null));
    const hasTimerRunning = computed(() => typeof currentActiveTimer.value !== "undefined");

    const currentTimerElapsedSeconds = computed(() => {
        return Math.floor((now.value - (start.value ?? 0)) / 1000)
    })

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

    const onTimerDropDownSelect = (key: string | number) => {
        switch (key) {
            case "pause":
                break;
            case "resume":
                break;
            case "start":
                onStartTimer();
                break;
            case "stop":
                console.log(currentActiveTimer.value);
                onStopTimer(currentActiveTimer.value?.id ?? "");
                break;
            case "clear":
                onClearTimers();
                break;
        }
    };

    let interval: number | undefined;

    const timers = ref<TimerResponse[]>([]);

    const finishedTimers = computed(() => timers.value.filter((timer) => timer.finishedAt !== null));

    const hasFinishedTimers = computed(() => finishedTimers.value.length > 0);

    const onGetTimers = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response = await timerService.getTimers();
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
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await timerService.start();
            await onGetTimers();
        } catch (e) {
            // TODO:
            console.error(e);
        } finally {
            state.ajaxRunning = false;
        }
    };

    const onStopTimer = async (id: string) => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await timerService.stop(id);
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
            await onGetTimers();
        } catch (e) {
            // TODO:
            console.error(e);
        } finally {
            state.ajaxRunning = false;
        }
    };

    onMounted(() => {
        interval = setInterval(() => {
            now.value = Date.now()
        }, 1000);
        onGetTimers();
    });

    onBeforeUnmount(() => {
        clearInterval(interval)
    })
</script>

<template>
    <n-popover placement="bottom" trigger="hover" @select="onTimerDropDownSelect">
        <template #trigger>
            <n-button quaternary @click.prevent="onToggleTimer" @mousedown.prevent :disabled="state.ajaxRunning">
                <n-icon :size="iconSize" :component="IconAlarm" :color="color"
                    :class="{ 'doneo-timer-animated-icon': hasTimerRunning }" />
            </n-button>
        </template>
        <n-card size="small" segmented>
            <template #header>
                <n-button block @click="onStopTimer(currentActiveTimer?.id ?? '')" v-if="hasTimerRunning"
                    :disabled="state.ajaxRunning">
                    <template #icon>
                        <n-icon :component="IconClockStop" :size="commonIconSize" />
                    </template>
                    Stop current timer: {{ formatDuration(currentTimerElapsedSeconds > 0 ? currentTimerElapsedSeconds :
                        0) }}
                </n-button>
                <n-button block @click="onStartTimer" v-else :disabled="state.ajaxRunning">
                    <template #icon>
                        <n-icon :component="IconClockPlay" />
                    </template>
                    Start new timer
                </n-button>
            </template>
            <template #footer v-if="hasFinishedTimers">
                <n-list>
                    <template #header>
                        <n-icon :component="IconClock" :size="commonIconSize" /> Previous timers
                    </template>
                    <n-list-item v-for="timer in finishedTimers" :key="timer.id">
                        <n-icon :component="IconClockPlay" :size="commonIconSize" /> {{ new
                            IDate(timer.startedAt).toLocaleString() }}
                        <br>
                        <n-icon :component="IconClockStop" :size="commonIconSize" /> {{ new
                            IDate(timer.finishedAt).toLocaleString() }}
                        <br>
                        <n-icon :component="IconClock" :size="commonIconSize" />Total timer:
                        {{
                            formatDuration(

                                ((timer.finishedAt ?? new Date().getTime()) - timer.startedAt) / 1000

                            )
                        }}
                    </n-list-item>
                </n-list>
            </template>
            <template #action v-if="hasTimers">
                <n-button block @click="onClearTimers" :disabled="state.ajaxRunning">
                    <template #icon>
                        <n-icon :component="IconClockCancel" :size="commonIconSize" />
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