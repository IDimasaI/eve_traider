<template>
    <div class="flex flex-col text-center">
        <h2>Количество товаров</h2>
        <p>{{ props.count ?? 0 }}</p>
    </div>
</template>

<script lang="ts" setup>
const props = defineProps(["count"]);
</script>
