<script lang="ts" setup>
import ET from '../../assets/svg/ET.vue';

</script>
<template>
    <div class="bg-[#292929] text-white grid grid-cols-3  border-b-[1.5px] border-[#5AD0D0] h-auto text-xl">
        <div class="flex flex-row justify-center gap-4">
            <a href="/new" class="break-spaces text-center">
                Магазин Traiders
            </a>
            <div class="break-spaces text-center">
                Магазин EVE
            </div>
        </div>
        <div class="flex flex-row justify-center text-center">
            <ET width="80" height="60" />
        </div>
        <div class="flex flex-row justify-center">
            <p class="break-spaces text-center">Сравнить цены</p>
        </div>
    </div>
</template>

<style scoped>
.break-spaces {
    white-space: normal;
    width: min-content;
}
</style>