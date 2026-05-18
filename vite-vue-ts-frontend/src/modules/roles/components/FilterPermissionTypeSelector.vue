<script setup lang="ts">
    import { computed } from "vue";
    import { useI18n } from "vue-i18n";

    import { NSelect } from 'naive-ui';
    import type { SelectMixedOption } from "naive-ui/es/select/src/interface";
    import type { SelectSize } from "naive-ui";

    import { PermissionTypeFilterValue, type PermissionTypeFilter } from "../types/filter-permission-type";

    interface Props {
        size?: SelectSize,
        placeholder?: string;
        clearable?: boolean;
    }

    const { t } = useI18n();

    const props = withDefaults(defineProps<Props>(), {
        size: "medium",
        clearable: false,
    });

    const options = computed<SelectMixedOption[]>(() => [
        { label: t("rolePermissionTypeSelectorValueAll"), value: PermissionTypeFilterValue.All },
        { label: t("rolePermissionTypeSelectorValueAllowed"), value: PermissionTypeFilterValue.Allowed },
        { label: t("rolePermissionTypeSelectorValueNotAllowed"), value: PermissionTypeFilterValue.NotAllowed }
    ]);


    const model = defineModel<PermissionTypeFilter>({
        default: PermissionTypeFilterValue.All,
    });
</script>

<template>
    <n-select :size="props.size" :options="options" v-model:value="model" :placeholder="props.placeholder"
        :clearable="props.clearable" />
</template>

<style lang="css" scoped></style>