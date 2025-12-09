import type { Ref } from "vue";

 
export interface Favorite {
    toggleFavorite: (item: string) => void,
    favoriteItems: Ref<string[]>
}