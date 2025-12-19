<script lang="ts" setup>
import { Updater } from '../../utils/updater';
import { ref, onMounted } from 'vue';

import ET from '../../assets/svg/ET.vue';
import { ToggleTheme, Themes, CurrentTheme } from '../../composables/Theme';


const updater = new Updater();
const needUpdate = ref(false);

onMounted(async () => {
    needUpdate.value = await updater.check_update();
});

</script>
<template>
    <div class="bg-[#292929] text-white grid grid-cols-3  border-b-[1.5px] border-[#5AD0D0] h-auto text-xl">
        <div class="flex flex-row justify-center gap-4">
            <slot name="left"></slot>
        </div>
        <div class="flex flex-row justify-center text-center">
            <ET width="80" height="60" />
        </div>
        <div class="flex flex-row justify-center gap-4">
            <div class="content-center">
                <button class="theme-toggle" @click="ToggleTheme">
                    {{ CurrentTheme === Themes.dark ? '☀️ Светлая тема' : '🌙 Темная тема' }}
                </button>
            </div>
            <div class="break-spaces text-center" v-if="needUpdate">
                <p class="hover:opacity-50 cursor-pointer" @click="updater.update_app">Новое обновление<span
                        class="size-1 rounded-full bg-red-500 text-red-500">О</span>
                </p>
            </div>
        </div>
    </div>
</template>