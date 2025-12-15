<script lang="ts" setup>
import LeftPanel from '../components/new/LeftPanel.vue';
import Header from '../components/new/Header.vue'


function update() {
    //@ts-ignore
    window.ipc.postMessage(`-update`)
}
</script>
<template>
    <section class="bg-[#141414] min-h-screen">
        <Header></Header>
        <LeftPanel>
            <button @click="update" class=" hover:bg-slate-500">Клик</button>
        </LeftPanel>
    </section>
</template>