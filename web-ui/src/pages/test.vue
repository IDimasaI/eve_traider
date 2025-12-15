<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { SimpleTimer,Audio_player } from '../utils/sound'
const timer = new SimpleTimer(1000)
const isRunning = ref(false)

//const audio= new Audio_player("/src/assets/audio/update_price.mp3")

const audio = new Audio_player("https://assets.mixkit.co/active_storage/sfx/212/212-preview.mp3")
// Добавляем события
timer.addEvent(async () => {
    console.log('Событие 1: API запрос')
    const res = await fetch('https://jsonplaceholder.typicode.com/posts/1')
    const data = await res.json()
    console.log('Данные:', data.title)
})

timer.addEvent(async () => {
    console.log('Событие 2: Воспроизведение звука')
    audio.play().catch(() => console.log('Звук заблокирован'))
})

function toggle() {
    if (isRunning.value) {
        timer.stop()
    } else {
        timer.start()
    }
    isRunning.value = !isRunning.value
}
onMounted(() => {
    //toggle()
})  
</script>

<template>
    <button @click="toggle">
        {{ isRunning ? '⏹ Остановить' : '▶ Запустить' }}
    </button>
    <p>Статус: {{ isRunning ? 'Запрос каждые 5 секунд' : 'Остановлено' }}</p>
</template>