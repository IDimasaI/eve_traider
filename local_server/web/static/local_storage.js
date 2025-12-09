export class LocalStorage {
  constructor() {
    this.storage = window.localStorage;
  }

  setItem(key, value) {
    this.storage.setItem(key, value);
  }

  getItem(key) {
    return this.storage.getItem(key);
  }

  removeItem(key) {
    this.storage.removeItem(key);
  }

  getAllKeys() {
    return Object.keys(this.storage);
  }
}
