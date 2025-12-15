// utils/sound.ts - дополняем класс SimpleTimer
export class SimpleTimer {
    private timer: number | null = null
    private interval: number = 1000
    private events: (() => Promise<void> | void)[] = []
    
    constructor(interval?: number) {
        if (interval) {
            this.interval = interval
        }
    }
    
    // Добавить событие
    addEvent(event: () => Promise<void> | void): void {
        this.events.push(event)
    }
    
    // Удалить событие
    removeEvent(event: () => Promise<void> | void): void {
        const index = this.events.indexOf(event)
        if (index > -1) {
            this.events.splice(index, 1)
        }
    }
    
    // Очистить все события
    clearEvents(): void {
        this.events = []
    }
    
    // Получить количество событий
    getEventCount(): number {
        return this.events.length
    }
    
    // Выполнить все события
    async executeEvents(): Promise<void> {
        for (const event of this.events) {
            try {
                const result = event()
                if (result instanceof Promise) {
                    await result
                }
            } catch (error) {
                console.error('Ошибка в событии:', error)
            }
        }
    }
    
    // Запустить таймер с выполнением событий
    start(): void {
        this.stop()
        
        // Сразу выполняем события
        this.executeEvents()
        
        // Запускаем интервал
        this.timer = setInterval(() => {
            this.executeEvents()
        }, this.interval)
    }
    
    stop(): void {
        if (this.timer) {
            clearInterval(this.timer)
            this.timer = null
        }
    }
}
export class Audio_player {
    private audio: HTMLAudioElement
    
    constructor(src: string) {
        this.audio = new Audio(src)
        // Устанавливаем параметры для лучшей совместимости
        this.audio.preload = 'auto'
        this.audio.volume = 0.5
    }
    
    public async play(): Promise<void> {
        try {
            // Сбрасываем время для возможности повторного воспроизведения
            this.audio.currentTime = 0
            await this.audio.play()
        } catch (error) {
            console.error('Ошибка воспроизведения:', error)
            // Перебрасываем ошибку для обработки в вызывающем коде
            throw error
        }
    }
    
    // Дополнительные полезные методы
    public pause(): void {
        this.audio.pause()
    }
    
    public setVolume(volume: number): void {
        this.audio.volume = Math.max(0, Math.min(1, volume))
    }
}