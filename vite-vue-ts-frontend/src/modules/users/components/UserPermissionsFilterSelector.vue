<script setup lang="ts">
    import { computed } from "vue";
    import { useI18n } from "vue-i18n";

    import { NSelect } from 'naive-ui';
    import type { SelectMixedOption } from "naive-ui/es/select/src/interface";
    import type { SelectSize } from "naive-ui";

    import { UserPermissionFilterValue, type UserPermissionFilter } from "../types/user-admin-permission-filter";

    interface Props {
        disabled?: boolean;
        size?: SelectSize,
        placeholder?: string;
        clearable?: boolean;
    }

    const { t } = useI18n();

    const props = withDefaults(defineProps<Props>(), {
        disabled: false,
        size: "medium",
        clearable: false,
    });

    const options = computed<SelectMixedOption[]>(() => [
        { label: t("modules.user.components.UserPermissionsFilterSelector.options.any"), value: UserPermissionFilterValue.Any },
        { label: t("modules.user.components.UserPermissionsFilterSelector.options.onlyAdministrators"), value: UserPermissionFilterValue.OnlyAdministrators },
        { label: t("modules.user.components.UserPermissionsFilterSelector.options.onlyUsers"), value: UserPermissionFilterValue.OnlyUsers }
    ]);


    const model = defineModel<UserPermissionFilter>({
        default: UserPermissionFilterValue.Any,
    });
</script>

<template>
    <n-select :size="props.size" :disabled="props.disabled" :options="options" v-model:value="model"
        :placeholder="props.placeholder" :clearable="props.clearable" />
</template>

<style lang="css" scoped></style>