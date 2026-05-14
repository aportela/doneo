<script setup lang="ts">
    import { ref, watch, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSelect, NDatePicker } from 'naive-ui';
    import type { SelectMixedOption } from 'naive-ui/es/select/src/interface';
    import { getRange } from '../../composables/timestamps';

    const emit = defineEmits(['timestampRangeChange']);

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

    const selectorValue = ref<number>(0)
    const datepickerValue = ref<number | null>(null)
    const isDatePickerVisible = ref<boolean>(false);

    const recalcTimestamps = () => {
        switch (selectorValue.value) {
            case 0: // any date
                fromTimestamp.value = null;
                toTimestamp.value = null;
                break;
            case 1: // custom  date
                if (datepickerValue.value) {
                    // TODO
                }
                break;
            case 2: { // yesterday
                const { from, to } = getRange('yesterday');
                fromTimestamp.value = from;
                toTimestamp.value = to;
                break;
            }

            case 3: { // today
                const { from, to } = getRange('today');
                fromTimestamp.value = from;
                toTimestamp.value = to;
                break;
            }

            case 4: { // tomorrow
                const { from, to } = getRange('tomorrow');
                fromTimestamp.value = from;
                toTimestamp.value = to;
                break;
            }

            case 5: { // last week
                const { from, to } = getRange('last_week', { weekStartsOn: 1 });
                fromTimestamp.value = from;
                toTimestamp.value = to;
                break;
            }

            case 6: { // this week
                const { from, to } = getRange('this_week', { weekStartsOn: 1 });
                fromTimestamp.value = from;
                toTimestamp.value = to;
                break;
            }

            case 7: { // next week
                const { from, to } = getRange('next_week', { weekStartsOn: 1 });
                fromTimestamp.value = from;
                toTimestamp.value = to;
                break;
            }

            case 8: { // last month
                const { from, to } = getRange('last_month');
                fromTimestamp.value = from;
                toTimestamp.value = to;
                break;
            }

            case 9: { // this month
                const { from, to } = getRange('this_month');
                fromTimestamp.value = from;
                toTimestamp.value = to;
                break;
            }

            case 10: { // next month
                const { from, to } = getRange('next_month');
                fromTimestamp.value = from;
                toTimestamp.value = to;
                break;
            }

            case 11: { // last year
                const { from, to } = getRange('last_year');
                fromTimestamp.value = from;
                toTimestamp.value = to;
                break;
            }

            case 12: { // this year
                const { from, to } = getRange('this_year');
                fromTimestamp.value = from;
                toTimestamp.value = to;
                break;
            }

            case 13: { // next year
                const { from, to } = getRange('next_year');
                fromTimestamp.value = from;
                toTimestamp.value = to;
                break;
            }
        }
        emit("timestampRangeChange", { from: fromTimestamp.value, to: toTimestamp.value });
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

    const fromTimestamp = ref<number | null>(null);
    const toTimestamp = ref<number | null>(null);

    const isSelectorVisible = computed(() => selectorValue.value !== 1);
</script>

<template>
    <n-select v-if="isSelectorVisible" v-model:value="selectorValue" :options="options" size="small" />
    <n-date-picker :placeholder="t('select date')" v-else v-model:value="datepickerValue" type="date" clearable
        size="small" v-model:show="isDatePickerVisible" @clear="onClearDate" :actions="['clear']" />
</template>

<style lang="css" scoped></style>