<script setup lang="ts">
    import { ref, watch, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSelect, NDatePicker, NButton } from 'naive-ui';
    import type { SelectMixedOption } from 'naive-ui/es/select/src/interface';
    import { type TimestampRange, getRange } from '../../composables/timestamps';

    interface Props {
        disabled?: boolean;
        clearable?: boolean;
    }

    const { t } = useI18n();

    const props = withDefaults(defineProps<Props>(), {
        disabled: false,
        clearable: false,
    });

    const options = computed<SelectMixedOption[]>(() => [
        { label: t("shared.components.selectors.dateFilter.options.anyDate"), value: 0 },
        { label: t("shared.components.selectors.dateFilter.options.customDate"), value: 1 },
        { label: t("shared.components.selectors.dateFilter.options.yesterday"), value: 2 },
        { label: t("shared.components.selectors.dateFilter.options.today"), value: 3 },
        { label: t("shared.components.selectors.dateFilter.options.tomorrow"), value: 4 },
        { label: t("shared.components.selectors.dateFilter.options.lastWeek"), value: 5 },
        { label: t("shared.components.selectors.dateFilter.options.thisWeek"), value: 6 },
        { label: t("shared.components.selectors.dateFilter.options.nextWeek"), value: 7 },
        { label: t("shared.components.selectors.dateFilter.options.lastMonth"), value: 8 },
        { label: t("shared.components.selectors.dateFilter.options.thisMonth"), value: 9 },
        { label: t("shared.components.selectors.dateFilter.options.nextMonth"), value: 10 },
        { label: t("shared.components.selectors.dateFilter.options.lastYear"), value: 11 },
        { label: t("shared.components.selectors.dateFilter.options.thisYear"), value: 12 },
        { label: t("shared.components.selectors.dateFilter.options.nextYear"), value: 13 },
    ]);

    const range = defineModel<TimestampRange>("range", {
        default: {
            from: null,
            to: null
        }
    });

    const selectorValue = ref<number>(0)
    const datepickerValue = ref<number | null>(null)
    const isDatePickerVisible = ref<boolean>(false);

    const recalcTimestamps = () => {
        switch (selectorValue.value) {
            case 0: // any date
                range.value.from = null;
                range.value.to = null;
                break;
            case 1: // custom  date
                if (datepickerValue.value) {
                    const from = new Date(datepickerValue.value);
                    from.setHours(0, 0, 0, 0);
                    const to = new Date(datepickerValue.value);
                    to.setHours(23, 59, 59, 999);
                    range.value.from = from.getTime();
                    range.value.to = to.getTime();
                } else {
                    range.value.from = null;
                    range.value.to = null;
                }
                break;
            case 2: { // yesterday
                const { from, to } = getRange('yesterday');
                range.value.from = from;
                range.value.to = to;
                break;
            }

            case 3: { // today
                const { from, to } = getRange('today');
                range.value.from = from;
                range.value.to = to;
                break;
            }

            case 4: { // tomorrow
                const { from, to } = getRange('tomorrow');
                range.value.from = from;
                range.value.to = to;
                break;
            }

            case 5: { // last week
                const { from, to } = getRange('last_week', { weekStartsOn: 1 });
                range.value.from = from;
                range.value.to = to;
                break;
            }

            case 6: { // this week
                const { from, to } = getRange('this_week', { weekStartsOn: 1 });
                range.value.from = from;
                range.value.to = to;
                break;
            }

            case 7: { // next week
                const { from, to } = getRange('next_week', { weekStartsOn: 1 });
                range.value.from = from;
                range.value.to = to;
                break;
            }

            case 8: { // last month
                const { from, to } = getRange('last_month');
                range.value.from = from;
                range.value.to = to;
                break;
            }

            case 9: { // this month
                const { from, to } = getRange('this_month');
                range.value.from = from;
                range.value.to = to;
                break;
            }

            case 10: { // next month
                const { from, to } = getRange('next_month');
                range.value.from = from;
                range.value.to = to;
                break;
            }

            case 11: { // last year
                const { from, to } = getRange('last_year');
                range.value.from = from;
                range.value.to = to;
                break;
            }

            case 12: { // this year
                const { from, to } = getRange('this_year');
                range.value.from = from;
                range.value.to = to;
                break;
            }

            case 13: { // next year
                const { from, to } = getRange('next_year');
                range.value.from = from;
                range.value.to = to;
                break;
            }
        }
    };

    watch(selectorValue, async (val) => {
        if (val !== 1) { // NOT custom date
            datepickerValue.value = null;
            isDatePickerVisible.value = false;
            recalcTimestamps();
            return;
        }
        isDatePickerVisible.value = true
        recalcTimestamps();
    });

    watch(isDatePickerVisible, (visible: boolean) => {
        if (!visible) {
            if (!datepickerValue.value) {
                // datepicker closed with empty value
                selectorValue.value = 0;
                recalcTimestamps();
            } else {
                // datepicker closed with non empty value
                recalcTimestamps();
            }
        } else {
            // waiting for datepicker date
        }
    });

    const onClearDate = () => {
        datepickerValue.value = null;
        selectorValue.value = 0;
        isDatePickerVisible.value = false;
        recalcTimestamps();
    };

    const isSelectorVisible = computed(() => selectorValue.value !== 1);

    const reset = () => {
        onClearDate();
    }

    defineExpose({ reset });
</script>

<template>
    <n-select v-if="isSelectorVisible" :disabled="props.disabled" v-model:value="selectorValue" :options="options"
        size="small" :clearable="props.clearable" />
    <n-date-picker :disabled="props.disabled" :placeholder="t('shared.components.selectors.dateFilter.placeholder')"
        v-else v-model:value="datepickerValue" type="date" clearable size="small" v-model:show="isDatePickerVisible"
        @clear="onClearDate" :actions="['now', 'clear']">
        <template #now="{ onNow }">
            <n-button size="tiny" @click="onNow">
                {{ t('shared.components.pickers.datepicker.buttons.now.label') }}
            </n-button>
        </template>
        <template #clear="{ onClear }">
            <n-button size="tiny" @click="onClear">
                {{ t('shared.components.pickers.datepicker.buttons.clear.label') }}
            </n-button>
        </template>
    </n-date-picker>

</template>

<style lang="css" scoped></style>