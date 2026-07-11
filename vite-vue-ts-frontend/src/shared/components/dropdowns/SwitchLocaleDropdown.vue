<script setup lang="ts">
    import { ref } from 'vue';
    import { useI18n } from "vue-i18n";

    import { NDropdown, NButton, NIcon, type ButtonSize } from 'naive-ui';

    import { Globe, ChevronsUpDown } from '@lucide/vue';
    import { availableLocaleSelectorOptionItems, getlocaleSelectorOptionItem } from '../../../i18n';
    import { useI18nStore } from '../../../stores/i18n';

    import { DEFAULT_SINGLE_ACTION_BUTTON_SIZE, DEFAULT_BUTTON_ICON_SIZE } from '../../../constants';

    interface IProps {
        size?: ButtonSize;
        iconSize?: number,
        disabled?: boolean,
    };

    const props = withDefaults(defineProps<IProps>(), {
        size: DEFAULT_SINGLE_ACTION_BUTTON_SIZE,
        iconSize: DEFAULT_BUTTON_ICON_SIZE,
        disabled: false,
    });

    const { locale } = useI18n();
    const i18NStore = useI18nStore();

    const selected = ref(availableLocaleSelectorOptionItems[0]);
    const selectedLocale = ref<string | null>(getlocaleSelectorOptionItem(i18NStore.currentLocale).label);

    const onChangeLocale = (key: string) => {
        selectedLocale.value = getlocaleSelectorOptionItem(key).label
        locale.value = key;
        i18NStore.setLocale(key);
    };
</script>

<template>
    <n-dropdown trigger="click" @select="onChangeLocale" :options="availableLocaleSelectorOptionItems"
        v-model="selected">
        <n-button :size="props.size" quaternary :disabled="props.disabled">
            <n-icon :size="props.iconSize" :component="Globe" />
            <span class="selected_locale">{{ selectedLocale }}</span>
            <n-icon :size="props.iconSize" :component="ChevronsUpDown" />
        </n-button>
    </n-dropdown>
</template>

<style lang="css" scoped>
    .selected_locale {
        margin: 0em 0.5em;
    }
</style>