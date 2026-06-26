<script setup lang="ts">
    import { computed } from "vue";
    import { useI18n } from "vue-i18n";

    import { NFlex, NFormItem, NInputNumber } from 'naive-ui';

    interface EstimatedTimeInputProps {
        disabled?: boolean;
        readOnly?: boolean;
    };

    const props = withDefaults(defineProps<EstimatedTimeInputProps>(), {
        disabled: false,
        readOnly: false,
    });

    const { t } = useI18n();

    const seconds = defineModel<number>("seconds", { default: 0 });

    const DAY = 86400;
    const HOUR = 3600;
    const MINUTE = 60;

    const parts = computed(() => {
        const total = seconds.value || 0;

        return {
            days: Math.floor(total / DAY),
            hours: Math.floor((total % DAY) / HOUR),
            minutes: Math.floor((total % HOUR) / MINUTE),
        };
    });

    function updatePart(part: "days" | "hours" | "minutes", value: number | null) {
        const v = value ?? 0;

        const total = seconds.value || 0;

        const days = Math.floor(total / DAY);
        const hours = Math.floor((total % DAY) / HOUR);
        const minutes = Math.floor((total % HOUR) / MINUTE);

        switch (part) {
            case "days":
                seconds.value = v * DAY + hours * HOUR + minutes * MINUTE;
                break;

            case "hours":
                seconds.value = days * DAY + v * HOUR + minutes * MINUTE;
                break;

            case "minutes":
                seconds.value = days * DAY + hours * HOUR + v * MINUTE;
                break;
        }
    }
</script>

<template>
    <n-flex justify="space-between">
        <n-form-item :label="t('shared.components.inputs.EstimatedTimeInput.inputs.estimatedDays.label')">
            <n-input-number :min="0"
                :placeholder="t('shared.components.inputs.EstimatedTimeInput.inputs.estimatedDays.placeholder')"
                :value="parts.days" @update:value="val => updatePart('days', val)" clearable :disabled="props.disabled"
                :readonly="props.readOnly">
            </n-input-number>
        </n-form-item>
        <n-form-item :label="t('shared.components.inputs.EstimatedTimeInput.inputs.estimatedHours.label')">
            <n-input-number :min="0"
                :placeholder="t('shared.components.inputs.EstimatedTimeInput.inputs.estimatedHours.placeholder')"
                :value="parts.hours" @update:value="val => updatePart('hours', val)" clearable
                :disabled="props.disabled" :readonly="props.readOnly">
            </n-input-number>
        </n-form-item>
        <n-form-item :label="t('shared.components.inputs.EstimatedTimeInput.inputs.estimatedMinutes.label')">
            <n-input-number :min="0"
                :placeholder="t('shared.components.inputs.EstimatedTimeInput.inputs.estimatedMinutes.placeholder')"
                :value="parts.minutes" @update:value="val => updatePart('minutes', val)" clearable
                :disabled="props.disabled" :readonly="props.readOnly">
            </n-input-number>
        </n-form-item>
    </n-flex>
</template>

<style lang="css" scoped></style>