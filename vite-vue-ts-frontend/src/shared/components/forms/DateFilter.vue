<script setup lang="ts">
    import { ref, watch, computed } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NSelect, NDatePicker } from 'naive-ui';
    import type { SelectMixedOption } from 'naive-ui/es/select/src/interface';

    const { t } = useI18n();

    const dateFilterOptions = computed<SelectMixedOption[]>(() => [
        { label: t("Any date"), value: 0 },
        { label: t("Custom date"), value: 1 },
        { label: t("Yesterday"), value: 2 },
        { label: t("Today"), value: 3 },
        { label: t("Tomorrow"), value: 4 },
        { label: t("This week"), value: 5 },
        { label: t("This month"), value: 6 },
        { label: t("This year"), value: 7 },
    ]);

    const selectorValue = ref<number>(0)
    const datepickerValue = ref<number | null>(null)
    const showPicker = ref<boolean>(false);

    watch(selectorValue, async (val) => {
        if (val !== 4) {
            datepickerValue.value = null
            showPicker.value = false
            return
        }
        showPicker.value = true
    });

    watch(showPicker, (newValue: boolean) => {
        if (!newValue && !datepickerValue.value) {
            selectorValue.value = 0;
        }
    });

    const onClearDate = () => {
        datepickerValue.value = null
        selectorValue.value = 0
        showPicker.value = false
    };

</script>

<template>
    <n-select v-if="selectorValue !== 1" v-model:value="selectorValue" :options="dateFilterOptions" size="small" />
    <n-date-picker :placeholder="t('select date')" v-else v-model:value="datepickerValue" type="date" clearable
        size="small" v-model:show="showPicker" @clear="onClearDate" :actions="['clear']" />
</template>

<style lang="css" scoped></style>