<script setup lang="ts">
    import { computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NProgress, NFlex } from 'naive-ui';

    import { geti18nTimeParts } from '../../composables/datetime';

    interface Props {
        estimated: number;
        spent: number;
        height?: number;
        processing?: boolean;
    }

    const props = defineProps<Props>();

    const { t } = useI18n();

    const percent = computed(() => {
        if (props.estimated <= 0) return 0;
        return Math.min(
            Math.round((props.spent / props.estimated) * 100),
            100
        );
    });

    const color = computed(() => {
        if (props.spent > props.estimated) {
            return { stops: ['IndianRed', 'red'] };
        } else {
            return { stops: ['green', 'LightGreen'] }
        }
    });

    const spentLabel = computed(() => geti18nTimeParts(props.spent).map(({ key, count }) => `${count} ${t(key, count)}`).join(", "));
    const estimatedLabel = computed(() => geti18nTimeParts(props.estimated).map(({ key, count }) => `${count} ${t(key, count)}`).join(", "));
    const overrunLabel = computed(() => geti18nTimeParts(props.spent - props.estimated).map(({ key, count }) => `${count} ${t(key, count)}`).join(", "));
</script>

<template>
    <div>
        <n-flex justify="space-between">
            <div>{{ t("shared.components.progress.TaskSpentEstimatedPercent.labels.taskProgress") }}</div>
            <div class="spent-estimated-labels">
                <span v-if="spentLabel">
                    <strong>{{ t("shared.components.progress.TaskSpentEstimatedPercent.labels.spent") }}</strong>
                    {{ spentLabel }}
                </span>
                <span v-if="estimatedLabel">
                    <strong>{{ t("shared.components.progress.TaskSpentEstimatedPercent.labels.estimated") }}</strong>
                    {{ estimatedLabel }}
                </span>
                <span v-if="props.estimated > 0 && props.spent > props.estimated">
                    <strong>{{ t("shared.components.progress.TaskSpentEstimatedPercent.labels.overrun") }}</strong>
                    {{ overrunLabel }}
                </span>
            </div>
        </n-flex>
        <n-progress type="line" :percentage="percent" :color="color" :height="height" style="margin-top: 9px;"
            :processing="props.processing" />
    </div>
</template>

<style lang="css" scoped>

    div.spent-estimated-labels {
        padding-right: 48px;
    }

    strong {
        margin-left: 16px;
    }
</style>