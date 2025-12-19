import { ref } from 'vue'
export enum Themes {
    light = "light",
    dark = "dark"
}
export const CurrentTheme = ref(localStorage.getItem('theme') || Themes.light)

export const ToggleTheme = () => {
    CurrentTheme.value = CurrentTheme.value === Themes.light ? Themes.dark : Themes.light
    localStorage.setItem('theme', CurrentTheme.value)
}