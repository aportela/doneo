<script setup lang="ts">
    import { ref, computed, onMounted, onBeforeUnmount } from "vue";
    import { useI18n } from "vue-i18n";

    import { NButton, NIcon, NTooltip, NDropdown } from 'naive-ui';
    import { IconStopwatch } from '@tabler/icons-vue';

    import { renderIcon } from "../../composables/naive-ui-icon";

    import { useNotify } from '../../../shared/composables/notification';

    interface SwitchNotificationsButtonProps {
        iconSize?: number,
    };

    withDefaults(defineProps<SwitchNotificationsButtonProps>(), {
        iconSize: 20
    });

    const { t } = useI18n();
    const { notify } = useNotify();

    const timer = ref<boolean>(false);

    const color = computed<string>(() => {
        return timer.value ? "red" : "black";
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

    const timerDropdownOptions = computed(() => [
        {
            label: 'Pause ' + elapsedSeconds.value + " seconds",
            key: 'pause',
            icon: renderIcon(IconStopwatch)(commonIconSize)
        },
        {
            label: 'Resume',
            key: 'resume',
            icon: renderIcon(IconStopwatch)(commonIconSize)
        },
        {
            label: 'Cancel',
            key: 'cancel',
            icon: renderIcon(IconStopwatch)(commonIconSize)
        }
    ]);

    const elapsedSeconds = computed(() => {
        return Math.floor((now.value - (start.value ?? 0)) / 1000)
    })

    const onTimerDropDownSelect = (key: string | number) => {
        switch (key) {
            case "pause":
                break;
            case "resume":
                break;
            case "cancel":
                break;
        }
    };

    let interval: number | undefined;

    onMounted(() => {
        interval = setInterval(() => {
            now.value = Date.now()
        }, 1000)
    });

    onBeforeUnmount(() => {
        clearInterval(interval)
    })
</script>

<template>
    <n-tooltip trigger="hover" v-if="!timer">
        <template #trigger>
            <n-button quaternary @click.prevent="onToggleTimer" @mousedown.prevent>
                <n-icon :size="iconSize" :component="IconStopwatch" :color="color"
                    :class="{ 'doneo-timer-animated-icon': timer }" />
            </n-button>
        </template>
        {{
            t(!timer ?
                "shared.components.buttons.timer.enable.toolTip" :
                "shared.components.buttons.timer.disable.toolTip")
        }}
    </n-tooltip>
    <n-dropdown v-else :options="timerDropdownOptions" placement="bottom" trigger="hover"
        @select="onTimerDropDownSelect">
        <n-button quaternary @click.prevent="onToggleTimer" @mousedown.prevent>
            <n-icon :size="iconSize" :component="IconStopwatch" :color="color"
                :class="{ 'doneo-timer-animated-icon': timer }" />
        </n-button>
    </n-dropdown>
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