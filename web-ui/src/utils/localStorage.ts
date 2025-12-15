export class LocalStorage {
  private storage: Storage;

  constructor() {
    this.storage = window.localStorage;
  }

  setItem(key: string, value: string) {
    this.storage.setItem(key, value);
  }

  getItem(key: string) {
    return this.storage.getItem(key);
  }

  removeItem(key: string) {
    this.storage.removeItem(key);
  }

  getAllKeys() {
    return Object.keys(this.storage);
  }
}
