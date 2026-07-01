<script setup lang="ts">
    import { computed } from "vue";
    import { useI18n } from "vue-i18n";

    import { NFlex, NFormItem, NInputNumber } from 'naive-ui';

    import { DAY_SECONDS, HOUR_SECONDS, MINUTE_SECONDS, getSecondsDatetimeParts } from "../../types/datetime";

    interface IProps {
        disabled?: boolean;
        readOnly?: boolean;
        inputType: "estimated" | "spent"
    };

    const props = withDefaults(defineProps<IProps>(), {
        disabled: false,
        readOnly: false,
    });

    const { t } = useI18n();

    const seconds = defineModel<number>("seconds", { default: 0 });

    const parts = computed(() => getSecondsDatetimeParts(seconds.value || 0));

    const updatePart = (part: "days" | "hours" | "minutes", value: number | null) => {
        const v = value ?? 0;

        const total = seconds.value || 0;

        const days = Math.floor(total / DAY_SECONDS);
        const hours = Math.floor((total % DAY_SECONDS) / HOUR_SECONDS);
        const minutes = Math.floor((total % HOUR_SECONDS) / MINUTE_SECONDS);

        switch (part) {
            case "days":
                seconds.value = v * DAY_SECONDS + hours * HOUR_SECONDS + minutes * MINUTE_SECONDS;
                break;

            case "hours":
                seconds.value = days * DAY_SECONDS + v * HOUR_SECONDS + minutes * MINUTE_SECONDS;
                break;

            case "minutes":
                seconds.value = days * DAY_SECONDS + hours * HOUR_SECONDS + v * MINUTE_SECONDS;
                break;
        }
    };
</script>

<template>
    <n-flex justify="space-between">
        <n-form-item
            :label="t(props.inputType === 'estimated' ? 'shared.components.formBlocks.TimeFieldsInput.inputs.estimatedDays.label' : 'shared.components.formBlocks.TimeFieldsInput.inputs.spentDays.label')">
            <n-input-number :min="0"
                :placeholder="t(props.inputType === 'estimated' ? 'shared.components.formBlocks.TimeFieldsInput.inputs.estimatedDays.placeholder' : 'shared.components.formBlocks.TimeFieldsInput.inputs.spentDays.placeholder')"
                :value="parts.days" @update:value="val => updatePart('days', val)" clearable :disabled="props.disabled"
                :readonly="props.readOnly">
            </n-input-number>
        </n-form-item>
        <n-form-item
            :label="t(props.inputType === 'estimated' ? 'shared.components.formBlocks.TimeFieldsInput.inputs.estimatedHours.label' : 'shared.components.formBlocks.TimeFieldsInput.inputs.spentHours.label')">
            <n-input-number :min="0"
                :placeholder="t(props.inputType === 'estimated' ? 'shared.components.formBlocks.TimeFieldsInput.inputs.estimatedHours.placeholder' : 'shared.components.formBlocks.TimeFieldsInput.inputs.spentHours.placeholder')"
                :value="parts.hours" @update:value="val => updatePart('hours', val)" clearable
                :disabled="props.disabled" :readonly="props.readOnly">
            </n-input-number>
        </n-form-item>
        <n-form-item
            :label="t(props.inputType === 'estimated' ? 'shared.components.formBlocks.TimeFieldsInput.inputs.estimatedMinutes.label' : 'shared.components.formBlocks.TimeFieldsInput.inputs.spentMinutes.label')">
            <n-input-number :min="0"
                :placeholder="t(props.inputType === 'estimated' ? 'shared.components.formBlocks.TimeFieldsInput.inputs.estimatedMinutes.placeholder' : 'shared.components.formBlocks.TimeFieldsInput.inputs.spentMinutes.placeholder')"
                :value="parts.minutes" @update:value="val => updatePart('minutes', val)" clearable
                :disabled="props.disabled" :readonly="props.readOnly">
            </n-input-number>
        </n-form-item>
    </n-flex>
</template>

<style lang="css" scoped></style>