<script setup lang="ts">
    import { ref } from 'vue';
    import { NDropdown, NButton } from 'naive-ui';
    import { IconWorld, IconSelector } from '@tabler/icons-vue';
    import { availableLocaleSelectorOptionItems, getlocaleSelectorOptionItem } from '../../i18n';
    import { useI18n } from "vue-i18n";
    import { useI18nStore } from "../../stores/i18n";

    const { locale } = useI18n();

    const i18NStore = useI18nStore();

    interface SwitchLocaleButtonProps {
        iconSize?: number,
    };

    withDefaults(defineProps<SwitchLocaleButtonProps>(), {
        iconSize: 20
    });


    const selected = ref(availableLocaleSelectorOptionItems[0]);
    const selectedLocale = ref<string | null>(getlocaleSelectorOptionItem(i18NStore.currentLocale).label);

    const onChangeLocale = (key: string) => {
        selectedLocale.value = getlocaleSelectorOptionItem(key).label
        locale.value = key;
        i18NStore.setLocale(key);
    }
</script>

<template>
    <n-dropdown trigger="click" @select="onChangeLocale" :options="availableLocaleSelectorOptionItems"
        v-model="selected">
        <n-button quaternary>
            <IconWorld :size="iconSize" />
            <span class="selected_locale">{{ selectedLocale }}</span>
            <IconSelector />
        </n-button>
    </n-dropdown>
</template>

<style lang="css" scoped>
    .selected_locale {
        margin: 0em 0.5em;
    }
</style>