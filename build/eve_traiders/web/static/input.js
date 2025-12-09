class Debouncer {
  constructor(delay) {
    this.delay = delay;
    this.timeout = null;
  }

  debounce(fn) {
    return (...args) => {
      clearTimeout(this.timeout);
      this.timeout = setTimeout(() => fn(...args), this.delay);
    };
  }

  cancel() {
    clearTimeout(this.timeout);
    this.timeout = null;
  }

  flush(fn, ...args) {
    this.cancel();
    fn(...args);
  }
}

export class SearchComponent {
  /**
   *
   * @param {string} id
   */
  constructor(id) {
    this.id = id;
    /** @private */
    this.eventListeners = {
      change: [],
      search: [],
    };

    /** @private */
    this.debouncer = new Debouncer(1100);

    this.init_input();
  }

  /**@private */
  init_input() {
    const input = document.getElementById(this.id);

    // Создаем debounced обработчик
    const debouncedChange = this.debouncer.debounce((value) => {
      this.triggerChange(value);
    });

    input.addEventListener("input", (event) => {
      debouncedChange(event.target.value);
    });

    input.addEventListener("keydown", (event) => {
      if (event.key === "Enter") {
        this.debouncer.cancel(); // Отменяем отложенный вызов
        this.triggerChange(event.target.value);
      }
    });
  }

  /**
   * Универсальная функция debounce
   * @param {number} delay - Задержка в миллисекундах
   */
  setDelay(delay) {
    this.debouncer.delay = delay;
  }

  /**
   * Подписка на события
   * @param {string} event
   * @param {function} callback
   * @returns {function(): void} Функция для отмены подписки
   */
  on(event, callback) {
    if (!this.eventListeners[event]) {
      this.eventListeners[event] = [];
    }
    this.eventListeners[event].push(callback);

    return () => {
      const index = this.eventListeners[event].indexOf(callback);
      if (index > -1) {
        this.eventListeners[event].splice(index, 1);
      }
    };
  }

  /**
   * Подписка на изменения (синтаксический сахар)
   * @param {function(string): void} callback
   * @returns {function(): void} Функция для отмены подписки
   */
  onChange(callback) {
    return this.on("change", callback);
  }

  /**@private */
  triggerChange(value) {
    const listeners = this.eventListeners.change || [];
    listeners.forEach((callback) => {
      try {
        callback(value);
      } catch (error) {
        console.error("Error in search callback:", error);
      }
    });
  }
}
