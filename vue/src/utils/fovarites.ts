
import { LocalStorage } from "./localStorage";

// Определяем интерфейс для фаворитов
export interface FavoritesStorage {
  getFavorites(): string[];
  addFavorite(item: string): void;
  removeFavorite(item: string): void;
  isFavorite(item: string): boolean;
  clearFavorites(): void;
}

export class FavoritesManager implements FavoritesStorage {
  private storage: LocalStorage;
  private readonly STORAGE_KEY = "favorites";

  constructor() {
    this.storage = new LocalStorage();
  }

  getFavorites(): string[] {
    try {
      const favorites = this.storage.getItem(this.STORAGE_KEY);
      return favorites ? JSON.parse(favorites) : [];
    } catch (error) {
      console.error('Error reading favorites:', error);
      this.storage.setItem(this.STORAGE_KEY, JSON.stringify([]));
      return [];
    }
  }

  addFavorite(item: string): void {
    const favorites = this.getFavorites();
    if (!favorites.includes(item)) {
      favorites.push(item);
      this.saveFavorites(favorites);
    }
  }

  removeFavorite(item: string): void {
    const favorites = this.getFavorites();
    const updatedFavorites = favorites.filter(i => i !== item);
    if (updatedFavorites.length !== favorites.length) {
      this.saveFavorites(updatedFavorites);
    }
  }

  isFavorite(item: string): boolean {
    const favorites = this.getFavorites();
    return favorites.includes(item);
  }

  clearFavorites(): void {
    this.saveFavorites([]);
  }

  toggleFavorite(item: string): boolean {
    if (this.isFavorite(item)) {
      this.removeFavorite(item);
      return false;
    } else {
      this.addFavorite(item);
      return true;
    }
  }

  private saveFavorites(favorites: string[]): void {
    this.storage.setItem(this.STORAGE_KEY, JSON.stringify(favorites));
  }
}
export const favoritesManager = new FavoritesManager();
