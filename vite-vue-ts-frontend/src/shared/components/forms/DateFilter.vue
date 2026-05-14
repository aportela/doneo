<script setup lang="ts">
    import { ref, watch, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSelect, NDatePicker } from 'naive-ui';
    import type { SelectMixedOption } from 'naive-ui/es/select/src/interface';
    import { type TimestampRange, getRange } from '../../composables/timestamps';


    const { t } = useI18n();

    const options = computed<SelectMixedOption[]>(() => [
        { label: t("Any date"), value: 0 },
        { label: t("Custom date"), value: 1 },
        { label: t("Yesterday"), value: 2 },
        { label: t("Today"), value: 3 },
        { label: t("Tomorrow"), value: 4 },
        { label: t("Last week"), value: 5 },
        { label: t("This week"), value: 6 },
        { label: t("Next week"), value: 7 },
        { label: t("Last month"), value: 8 },
        { label: t("This month"), value: 9 },
        { label: t("Next month"), value: 10 },
        { label: t("Last year"), value: 11 },
        { label: t("This year"), value: 12 },
        { label: t("Next year"), value: 13 },
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
</script>

<template>
    <n-select v-if="isSelectorVisible" v-model:value="selectorValue" :options="options" size="small" />
    <n-date-picker :placeholder="t('select date')" v-else v-model:value="datepickerValue" type="date" clearable
        size="small" v-model:show="isDatePickerVisible" @clear="onClearDate" :actions="['clear']" />
</template>

<style lang="css" scoped></style>