<script lang="ts" setup>
import { CurrentTheme, Themes } from '../composables/Theme';
import { computed } from 'vue';
const props = defineProps({
    color: String,
    bg: String,
    border_color: String
});

const classes = computed(() => {

    if (props.bg || props.border_color) {
        return {
            color: props.color,
            bg: props.bg,
            border_color: props.border_color
        }
    }


    return {
        color: CurrentTheme.value === Themes.dark ? 'text-[#f1f5f9]' : 'text-[#333]',
        bg: CurrentTheme.value === Themes.dark ? 'bg-[#1e293b]' : 'bg-white',
        border_color: CurrentTheme.value === Themes.dark ? 'border-gray-500' : 'border-gray-300'
    }
})
</script>
<template>
    <div
        :class="`${classes.color} flex justify-between w-full text-left ${classes.bg} shadow-xl rounded-t-md p-4 border-[1.5px] ${classes.border_color} mx-auto`">
        <div class="float-left">
            <slot name="left"></slot>
        </div>
        <div>
            <slot name="center"></slot>
        </div>
        <div class="float-right">
            <slot name="right"></slot>
        </div>
    </div>
</template>
