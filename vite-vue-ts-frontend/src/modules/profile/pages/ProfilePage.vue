<script setup lang="ts">
    import { computed, ref, reactive, watch, onMounted, onBeforeUnmount, nextTick } from 'vue';
    import { useI18n } from "vue-i18n";

    import dayjs from 'dayjs'

    import { NTabs, NTabPane, NCard, NAvatar, NFlex, NFormItem, NInputGroup, NInput, NButton, NButtonGroup, NPopover, NIcon, NGrid, NGridItem, NDivider } from 'naive-ui';
    import { IconTrash, IconSun, IconMoon, IconInfoCircle, IconLayoutSidebarLeftExpand, IconLayoutNavbarExpand, IconDeviceFloppy, IconImageGeneration } from '@tabler/icons-vue';

    import { useNotify } from '../../../shared/composables/notification';
    import { appBus } from '../../../shared/composables/bus';

    import { type AjaxStateInterface, defaultAjaxState, defaultAjaxStateRunning } from '../../../shared/types/ajaxState';
    import { profileService } from '../services/profile';
    import { handleAPIError } from '../../../api/client/errorHandler';

    import { useLoadingStore } from '../../../stores/loading';
    import { useColorSchemeStore } from '../../../stores/colorScheme';
    import { useUserSettingsStore } from '../../../stores/userSettings';
    import { Profile } from '../models/profile';
    import type { ProfileResponse, UpdateRequest } from '../types/dto';
    import { MIN_PASSWORD_LENGTH } from '../../users/models/user';
    import { defaultDateTimeMask } from '../../../shared/composables/datetime.ts';

    import GenerateAvatarModal from '../../../shared/components/modals/GenerateAvatarModal.vue';

    const { t } = useI18n();
    const { notify } = useNotify();

    const loadingStore = useLoadingStore();
    const colorSchemeStore = useColorSchemeStore();
    const userSettingsStore = useUserSettingsStore();

    const profile = ref<Profile>(new Profile());

    const state: AjaxStateInterface = reactive({ ...defaultAjaxState });

    const newPassword = ref<string | null>(null);
    const confirmedPassword = ref<string | null>(null);

    const currentTab = ref<string>("myAccount");

    const lastAvatarTimestamp = ref<number>(new Date().getTime());

    const currentAvatarURL = computed(() => profile.value.id ? `/api/wc/avatars/user/${profile.value.id}?t=${lastAvatarTimestamp.value}` : undefined);

    const showAvatarGeneratorModal = ref<boolean>(false);

    const currentDatetimeMask = ref<string | null>(userSettingsStore.currentDatetimeMask);

    watch(state, (newValue: AjaxStateInterface) => {
        loadingStore.set(newValue.ajaxRunning);
    });

    watch(() => currentDatetimeMask.value, (newValue) => {
        userSettingsStore.setDatetimeMask(newValue || defaultDateTimeMask)
    });

    const matchedPasswords = computed<boolean>(() => newPassword.value === confirmedPassword.value);

    const currentDatetimeMaskPreview = computed<string | null>({
        get() {
            return currentDatetimeMask.value ? dayjs().format(currentDatetimeMask.value) : null;
        },
        set(_value) { }
    });

    const allowSubmit = computed<boolean>(() => !!profile.value.name && !!profile.value.email && matchedPasswords.value)

    const onGet = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const response: ProfileResponse = await profileService.get();
            profile.value = new Profile(response);
            lastAvatarTimestamp.value = new Date().getTime();
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProfilePage.onGet" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.projectType.components.ProfilePage.errors.refreshError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.projectType.components.ProfilePage.errors.refreshError");
                    console.error("Unhandled API error", { file: "ProfilePage.vue", method: "onGet" }, { err: fatalError });
                });
        }
        finally {
            state.ajaxRunning = false;
            if (state.ajaxErrorMessage) {
                appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
            }
        }
    };

    const onUpdate = async () => {
        //serverErrors.value = {};
        //userFormRef.value?.restoreValidation();
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload: UpdateRequest = {
                name: profile.value.name ?? "",
                password: newPassword.value || undefined,
                email: profile.value.email ?? "",
            };
            const response: ProfileResponse = await profileService.update(payload);
            profile.value = new Profile(response);
            notify('success', t("modules.profile.components.ProfilePage.notifications.profileUpdated"));
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProfilePage.onUpdate" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.profile.components.ProfilePage.errors.notFoundError");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                //serverErrors.value.name = "modules.user.components.UserForm.inputs.name.errors.alreadyExists";
                            } else if (apiError.details?.field === "email") {
                                //serverErrors.value.email = "modules.user.components.UserForm.inputs.email.errors.alreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.profile.components.ProfilePage.errors.updateError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.profile.components.ProfilePage.errors.updateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.profile.components.ProfilePage.errors.updateError");
                    console.error("Unhandled API error", { file: "ProfilePage.vue", method: "onUpdate" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    //userFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    const onSaveAvatar = async (svg: string) => {
        showAvatarGeneratorModal.value = false;
        Object.assign(state, defaultAjaxStateRunning);
        try {
            const payload = { svg: svg };
            await profileService.saveAvatar(payload);
            notify('success', t("modules.profile.components.ProfilePage.notifications.avatarUpdated"));
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProfilePage.onSaveAvatar" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.profile.components.ProfilePage.errors.notFoundError");
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.profile.components.ProfilePage.errors.updateAvatarError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.profile.components.ProfilePage.errors.updateAvatarError");
                    console.error("Unhandled API error", { file: "ProfilePage.vue", method: "onSaveAvatar" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            lastAvatarTimestamp.value = new Date().getTime();
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    //userFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    const onDeleteAvatar = async () => {
        Object.assign(state, defaultAjaxStateRunning);
        try {
            await profileService.deleteAvatar();
            notify('success', t("modules.profile.components.ProfilePage.notifications.avatarDeleted"));
        } catch (error: unknown) {
            state.ajaxErrors = true;
            handleAPIError(error,
                (apiError) => {
                    switch (apiError.response?.status) {
                        case 401:
                            state.ajaxErrors = false;
                            appBus.emit({ type: "reauthRequired", payload: { emitter: "ProfilePage.onDeleteAvatar" } });
                            break;
                        case 403:
                            state.ajaxErrorMessage = t("shared.errorMessages.unauthorizedOperation");
                            break;
                        case 404:
                            state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.notFoundError");
                            break;
                        case 409:
                            if (apiError.details?.field === "name") {
                                //serverErrors.value.name = "modules.user.components.UserForm.inputs.name.errors.alreadyExists";
                            } else if (apiError.details?.field === "email") {
                                //serverErrors.value.email = "modules.user.components.UserForm.inputs.email.errors.alreadyExists";
                            } else {
                                state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.updateError");
                            }
                            break;
                        default:
                            state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.updateError");
                            break;
                    }
                },
                (fatalError) => {
                    state.ajaxErrorMessage = t("modules.user.components.UserForm.errors.updateError");
                    console.error("Unhandled API error", { file: "ProfilePage.vue", method: "onDeleteAvatar" }, { err: fatalError });
                });
        } finally {
            state.ajaxRunning = false;
            lastAvatarTimestamp.value = new Date().getTime();
            if (state.ajaxErrors) {
                if (state.ajaxErrorMessage) {
                    appBus.emit({ type: "remoteAPIError", payload: { errorMessage: state.ajaxErrorMessage } });
                } else {
                    await nextTick();
                    //userFormRef.value?.validate().then(() => { }).catch(() => { });
                }
            }
        }
    };

    let stopBusReauthListener: () => void;

    onMounted(() => {
        onGet();
        stopBusReauthListener = appBus.on("reauthValidNotify", async (payload) => {
            if (payload.to.includes("ProfilePage.onGet")) {
                onGet();
            } else if (payload.to.includes("ProfilePage.onUpdate")) {
                onUpdate();
            } else if (payload.to.includes("ProfilePage.onDeleteAvatar")) {
                onDeleteAvatar();
            }
        });
    });

    onBeforeUnmount(() => {
        stopBusReauthListener();
    });

</script>

<template>
    <GenerateAvatarModal v-model:show="showAvatarGeneratorModal" @confirm="(svg: string) => onSaveAvatar(svg)"
        @cancel="showAvatarGeneratorModal = false;" />
    <n-tabs placement="top" type="line" animated v-model:value="currentTab">
        <n-tab-pane name="myAccount" tab="My account">
            <n-card bordered>
                <h1>My account</h1>
                <h2>Profile details</h2>
                <p>Account created on {{
                    profile.createdAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) }}</p>
                <p v-if="profile.updatedAt?.hasValue()">Account last update on {{
                    profile.updatedAt?.toCustomMaskString(userSettingsStore.currentDatetimeMask) }}</p>
                <n-flex style="align-items:center;">
                    <n-avatar :size="128" :src="currentAvatarURL" :key="lastAvatarTimestamp" color="transparent" />
                    <div>
                        <n-button @click="showAvatarGeneratorModal = true" block style="margin-bottom: 8px;">
                            <template #icon>
                                <n-icon :component="IconImageGeneration" />
                            </template>
                            Change avatar
                        </n-button>
                        <n-button tertiary type="error" @click="onDeleteAvatar" block>
                            <template #icon>
                                <n-icon :component="IconTrash" />
                            </template>
                            Delete avatar
                        </n-button>
                    </div>
                </n-flex>
                <n-divider />
                <n-form-item label="Name">
                    <n-input v-model:value="profile.name" />
                </n-form-item>
                <n-form-item label="Email">
                    <n-input v-model:value="profile.email" />
                </n-form-item>
                <p>Password</p>
                <n-flex align="center">
                    <n-form-item label="New password">
                        <n-input type="password" :min="MIN_PASSWORD_LENGTH" v-model:value="newPassword"
                            placeholder="type new password" clearable />
                    </n-form-item>
                    <n-form-item label="Confirm new password"
                        :feedback="!matchedPasswords ? 'passwords do not match' : undefined"
                        :validation-status="!matchedPasswords ? 'error' : undefined">
                        <n-input type="password" :min="MIN_PASSWORD_LENGTH" v-model:value="confirmedPassword"
                            placeholder="confirm new password" clearable />
                    </n-form-item>
                </n-flex>
                <n-button @click="onUpdate" :disabled="!allowSubmit">
                    <template #icon>
                        <n-icon :component="IconDeviceFloppy" />
                    </template>
                    {{ t("shared.buttons.Save.label") }}
                </n-button>
            </n-card>
        </n-tab-pane>
        <n-tab-pane name="mySettings" tab="My settings">
            <n-card bordered>
                <h1>My settings</h1>

                <h2>Locale</h2>

                <h3 class="doneo-flex-center-align">Datetime format mask
                    <n-popover>
                        <template #trigger>
                            <n-icon style="margin-left: 8px; margin-top: 2px;" class="doneo-cursor-help" :size="20"
                                :component="IconInfoCircle" />
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
                    <n-button disabled>Mask preview (current datetime)</n-button>
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

                <n-button-group>
                    <n-button :disabled="userSettingsStore.sideNavigationMode"
                        @click="userSettingsStore.toggleNavigationMode">
                        <template #icon>
                            <n-icon>
                                <IconLayoutSidebarLeftExpand />
                            </n-icon>
                        </template>
                        side
                    </n-button>
                    <n-button :disabled="userSettingsStore.topNavigationMode"
                        @click="userSettingsStore.toggleNavigationMode">
                        <template #icon>
                            <n-icon>
                                <IconLayoutNavbarExpand />
                            </n-icon>
                        </template>
                        top
                    </n-button>
                </n-button-group>
            </n-card>
        </n-tab-pane>
    </n-tabs>
</template>

<style lang="css" scoped></style>