import { Renderer } from "./render.js";
import { LocalStorage } from "./local_storage.js";
import { SearchComponent } from "./input.js";

/** @type string[] */
let names = [];
/** @type string[] */
let all_ids = [];
/** @type string[] */
let all_prices = [];

document.addEventListener("DOMContentLoaded", async () => {
  const app = app_data();
  console.log(app);
  const data = await get_prices_from_bd();

  const localStorage = new LocalStorage();

  if (localStorage.getItem("names") === null) {
    names = await get_all_names_from_bd();
    localStorage.setItem("names", JSON.stringify(names).toLocaleLowerCase());
  }

  names = JSON.parse(localStorage.getItem("names"));

  const render = new Renderer();
  render.all_count_prices(data[data.length - 1].ids.split(",").length);

  for (const item of data) {
    /** @type string[] */
    const ids = item.ids.split(",");

    /** @type string[] */
    const prices = item.price.split(",");

    all_ids.push(...ids);
    all_prices.push(...prices);
  }

  init_search();
});

function app_data() {
  return JSON.parse(
    document.getElementById("app")?.getAttribute("data-app") || "{}",
  );
}

async function get_prices_from_bd() {
  return await (await fetch("/api/all_prices")).json();
}

async function get_all_names_from_bd() {
  return await (await fetch("/api/all_names")).json();
}

function init_search() {
  const search = new SearchComponent("search_input");

  // Можно изменить задержку
  search.setDelay(500);

  // Подписываемся на изменения с debounce

  search.onChange((value) => {
    const filteredNames = names.filter((name) => name.includes(value));
    document.getElementById("search_results").innerHTML =
      filteredNames.join(", ");
    document.getElementById("search_results").style.display = "block";
  });
}
