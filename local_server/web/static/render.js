export class Renderer {
  constructor() {
    this.curent_count_prices = 0;
  }
  all_count_prices(count) {
    document.getElementById("all_count_prices").textContent = count;
    this.curent_count_prices = count;
  }

  init_components() {
    //  this.init_search_component();
  }
}
