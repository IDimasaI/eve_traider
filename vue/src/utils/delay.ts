import { watch, type Ref } from "vue";

/**
 * @name time_watch
 * @description Задержка для выполнения функции, выполняющейся при изменении значения `watched`. Значение `watched` должно быть реактивным.
 * @param fn Выполняемая функция в которую передается новое значение `watched`
 * @param watched Название переменной, за которой нужно следить
 * @param delay Время задержки в миллисекундах
 * @returns watch хук
 */
export function time_watch(fn: CallableFunction, watched: Ref, delay?: number) {
  let timeout: any;
  return watch(watched, (newWatched) => {
    clearTimeout(timeout);
    timeout = setTimeout(async () => {
      fn(newWatched);
    }, delay);
  });
}
