<script setup lang="ts">
    import { computed, ref, watch } from 'vue';

    import dayjs from 'dayjs'

    import { NTabs, NTabPane, NCard, NAvatar, NFlex, NInputGroup, NInput, NButton, NButtonGroup, NPopover, NIcon, NGrid, NGridItem } from 'naive-ui';
    import { IconFileUpload, IconTrash, IconSun, IconMoon, IconInfoCircle } from '@tabler/icons-vue';

    import { useColorSchemeStore } from '../../../stores/colorScheme';
    import { useUserSettingsStore } from '../../../stores/userSettings';

    const colorSchemeStore = useColorSchemeStore();
    const userSettingsStore = useUserSettingsStore();

    const currentTab = ref<string>("mySettings");

    const currentDatetimeMask = ref<string | null>(userSettingsStore.currentDatetimeMask);

    watch(() => currentDatetimeMask.value, (newValue) => {
        userSettingsStore.setDatetimeMask(newValue || "YYYY-MM-DD HH:MM:SS")
    });

    const currentDatetimeMaskPreview = computed<string | null>({
        get() {
            return currentDatetimeMask.value ? dayjs().format(currentDatetimeMask.value) : null;
        },
        set(_value) { }
    });

</script>

<template>
    <n-tabs placement="left" type="line" animated v-model:value="currentTab">
        <n-tab-pane name="myAccount" tab="My account">
            <n-card bordered>
                <h1>My account</h1>
                <h2>Profile details</h2>
                <n-flex style="align-items:center;">
                    <n-avatar :size="128" src="/api/wc/avatars/128/user/019dfcaa-aa06-755b-aca0-bb05e4ea44a7" />
                    <n-button>
                        <template #icon>
                            <n-icon>
                                <IconFileUpload />
                            </n-icon>
                        </template>
                        Change avatar
                    </n-button>
                    <n-button>
                        <template #icon>
                            <n-icon>
                                <IconTrash />
                            </n-icon>
                        </template>
                        Delete avatar
                    </n-button>
                </n-flex>
            </n-card>
        </n-tab-pane>
        <n-tab-pane name="mySettings" tab="My settings">
            <n-card bordered>
                <h1>My settings</h1>

                <h2>Locale</h2>

                <h3>Datetime format mask
                    <n-popover>
                        <template #trigger>
                            <n-icon class="doneo-cursor-help" :size="20" :component="IconInfoCircle" />
                        </template>
                        <template #header>
                            List of all available formats
                        </template>
                        <n-grid :cols="2" :x-gap="128">
                            <n-grid-item>
                                <table>
                                    <thead>
                                        <tr>
                                            <th>Format</th>
                                            <th>Output</th>
                                            <th>Description</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        <tr>
                                            <td><code>YY</code></td>
                                            <td>18</td>
                                            <td>Two-digit year</td>
                                        </tr>
                                        <tr>
                                            <td><code>YYYY</code></td>
                                            <td>2018</td>
                                            <td>Four-digit year</td>
                                        </tr>
                                        <tr>
                                            <td><code>M</code></td>
                                            <td>1-12</td>
                                            <td>The month, beginning at 1</td>
                                        </tr>
                                        <tr>
                                            <td><code>MM</code></td>
                                            <td>01-12</td>
                                            <td>The month, 2-digits</td>
                                        </tr>
                                        <tr>
                                            <td><code>MMM</code></td>
                                            <td>Jan-Dec</td>
                                            <td>The abbreviated month name</td>
                                        </tr>
                                        <tr>
                                            <td><code>MMMM</code></td>
                                            <td>January-December</td>
                                            <td>The full month name</td>
                                        </tr>
                                        <tr>
                                            <td><code>D</code></td>
                                            <td>1-31</td>
                                            <td>The day of the month</td>
                                        </tr>
                                        <tr>
                                            <td><code>DD</code></td>
                                            <td>01-31</td>
                                            <td>The day of the month, 2-digits</td>
                                        </tr>
                                        <tr>
                                            <td><code>d</code></td>
                                            <td>0-6</td>
                                            <td>The day of the week, with Sunday as 0</td>
                                        </tr>
                                        <tr>
                                            <td><code>dd</code></td>
                                            <td>Su-Sa</td>
                                            <td>The min name of the day of the week</td>
                                        </tr>
                                        <tr>
                                            <td><code>ddd</code></td>
                                            <td>Sun-Sat</td>
                                            <td>The short name of the day of the week</td>
                                        </tr>
                                        <tr>
                                            <td><code>dddd</code></td>
                                            <td>Sunday-Saturday</td>
                                            <td>The name of the day of the week</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </n-grid-item>
                            <n-grid-item>
                                <table>
                                    <thead>
                                        <tr>
                                            <th>Format</th>
                                            <th>Output</th>
                                            <th>Description</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        <tr>
                                            <td><code>H</code></td>
                                            <td>0-23</td>
                                            <td>The hour</td>
                                        </tr>
                                        <tr>
                                            <td><code>HH</code></td>
                                            <td>00-23</td>
                                            <td>The hour, 2-digits</td>
                                        </tr>
                                        <tr>
                                            <td><code>h</code></td>
                                            <td>1-12</td>
                                            <td>The hour, 12-hour clock</td>
                                        </tr>
                                        <tr>
                                            <td><code>hh</code></td>
                                            <td>01-12</td>
                                            <td>The hour, 12-hour clock, 2-digits</td>
                                        </tr>
                                        <tr>
                                            <td><code>m</code></td>
                                            <td>0-59</td>
                                            <td>The minute</td>
                                        </tr>
                                        <tr>
                                            <td><code>mm</code></td>
                                            <td>00-59</td>
                                            <td>The minute, 2-digits</td>
                                        </tr>
                                        <tr>
                                            <td><code>s</code></td>
                                            <td>0-59</td>
                                            <td>The second</td>
                                        </tr>
                                        <tr>
                                            <td><code>ss</code></td>
                                            <td>00-59</td>
                                            <td>The second, 2-digits</td>
                                        </tr>
                                        <tr>
                                            <td><code>SSS</code></td>
                                            <td>000-999</td>
                                            <td>The millisecond, 3-digits</td>
                                        </tr>
                                        <tr>
                                            <td><code>Z</code></td>
                                            <td>+05:00</td>
                                            <td>The offset from UTC, ±HH:mm</td>
                                        </tr>
                                        <tr>
                                            <td><code>ZZ</code></td>
                                            <td>+0500</td>
                                            <td>The offset from UTC, ±HHmm</td>
                                        </tr>
                                        <tr>
                                            <td><code>A</code></td>
                                            <td>AM PM</td>
                                            <td></td>
                                        </tr>
                                        <tr>
                                            <td><code>a</code></td>
                                            <td>am pm</td>
                                            <td></td>
                                        </tr>
                                    </tbody>
                                </table>
                            </n-grid-item>
                        </n-grid>
                    </n-popover>
                </h3>
                <n-input-group>
                    <n-input placeholder="Type datetime format current mask" v-model:value="currentDatetimeMask" />
                    <n-input placeholder="mask preview (current datetime)" v-model:value="currentDatetimeMaskPreview"
                        readonly />
                </n-input-group>

                <h2>Theme</h2>

                <n-button-group>
                    <n-button :disabled="colorSchemeStore.light" @click="colorSchemeStore.toggle">
                        <template #icon>
                            <n-icon>
                                <IconSun />
                            </n-icon>
                        </template>
                        light
                    </n-button>
                    <n-button :disabled="colorSchemeStore.dark" @click="colorSchemeStore.toggle">
                        <template #icon>
                            <n-icon>
                                <IconMoon />
                            </n-icon>
                        </template>
                        dark
                    </n-button>
                </n-button-group>
                <h2>Layout</h2>
            </n-card>
        </n-tab-pane>
    </n-tabs>
</template>

<style lang="css" scoped></style>