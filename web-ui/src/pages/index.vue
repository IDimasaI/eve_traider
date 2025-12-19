<script lang="ts" setup>
import type { Favorite } from "../storage.ts"
import type { MarketData } from "../utils/API.ts"


import BaseToolbar from "../components/BaseToolbar.vue";
import { FuzzySearcher } from "../utils/search.ts";
import { ref, onMounted, provide } from "vue";
import { time_watch } from "../utils/delay.ts";
import { FavoritesManager } from "../utils/fovarites.ts";
import { find_id } from "../utils/API.ts";


import ItemsPanel from "../components/CategoryPanel.vue";
import Fovorite from "../components/Favourite.vue";
import BarChart from "../components/BarChart.vue";
import RegionMarketInfo from "../components/RegionMarketInfo.vue";

type Items = {
    name: string;
    category: string;
    id: number;
}[];

enum SearchType {
    Base,
    Market
}


async function get_all_items() {
    const res = await fetch("/api/v2/get_all_items");
    return (await res.json()) as Items;
}

// Поисковые утилиты и переменные//
const addItem = (): void => {
    MAX_COUNT_ITEMS.value += 1000;
};

const toggleFavorite = (item: string) => {
    favoritesManager.toggleFavorite(item);
    favoriteItems.value = favoritesManager.getFavorites();
};


const selectItem = (item: string) => {

    if (selectedItem.value === item) {
        return;
    }

    selectedItem.value = item;

}
function set_market_loading() {
    marketData.value.forEach(market => {
        market.loading = true
    })
}
const search_in_market = async () => {
    const id = await find_id(searchQuery.value).catch(() => { return })

    set_market_loading()
    const markets = [
        { id: 10000002, name: "Jita" },
        { id: 10000030, name: "Rens" },
        { id: 10000032, name: "Dodixie" },
        { id: 10000042, name: "Hek" },
        { id: 10000043, name: "Amarr" },
    ]

    const promises = markets.map(async market => {
        const res = await fetch(`https://esi.evetech.net/latest/markets/${market.id}/orders/?type_id=${id}&order_type=sell&language=en-us`)
        if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
        const data = await res.json()
        return { loading: false, market: market.name, data: data }
    })

    const results = await Promise.all(promises)
    marketData.value = results
    console.log(marketData.value)
    return
}


const favoritesManager = new FavoritesManager();

const searchQuery = ref("");
const itemNames = ref<string[]>([]);

const selectedItem = ref("");

let MAX_COUNT_ITEMS = ref(30);

const marketData = ref<MarketData[]>([]);
// КН Поисковые утилиты и переменные//





const favoriteItems = ref<string[]>([]);

const all_names = ref<Items>();
provide<Favorite>('Favourite', {
    toggleFavorite,
    favoriteItems
})


const type_search = ref<SearchType>(SearchType.Base)

onMounted(async () => {
    favoriteItems.value = favoritesManager.getFavorites();

    if (localStorage.getItem("items_id_name")) {
        all_names.value = JSON.parse(localStorage.getItem("items_id_name")!)!;
    } else {
        all_names.value = await get_all_items();
        localStorage.setItem("items_id_name", JSON.stringify(all_names.value));
    }

    const fuzzySearcher = new FuzzySearcher(all_names.value!);
    time_watch(
        (query: string) => {
            if (type_search.value == SearchType.Market) return
            if (query.length < 3) {
                itemNames.value = [];
                return;
            }
            const results = fuzzySearcher.search(query, {
                minScore: 0.1,
                exactMatchBonus: 0.4,
                maxResults: MAX_COUNT_ITEMS.value,
                exactMatchPriority: true,
            });

            itemNames.value = results.map((result) => result.item);
            console.log(itemNames.value);
        },
        searchQuery,
        500,
    );
});
</script>

<template>
    <section class="bg-white min-h-screen flex md:flex-row flex-col">
        <section
            class="text-center flex flex-col items-left md:min-w-[310px] md:max-w-2/12 max-md:w-full text-sm bg-gray-50">
            <div class="flex flex-row mt-8 justify-around p-1">
                <button class="btn hover:bg-gray-200 py-2 px-4 rounuded"
                    @click="type_search = SearchType.Base">База</button>
                <button class="btn hover:bg-gray-200 py-2 px-4 rounded"
                    @click="type_search = SearchType.Market">Магазин</button>
            </div>
            <template v-if="type_search == SearchType.Base">
                <div class="flex flex-col  mb-8">
                    <label for="search">Поиск в базе</label>
                    <input id="search" type="text" placeholder="Search..." v-model="searchQuery"
                        class="mt-2 border border-y-slate-500 focus:border-slate-500 focus:outline-none focus:ring-2 focus:ring-slate-500 focus:ring-offset-2 focus:ring-offset-white rounded-md px-2 py-1" />
                    <p class="text-sm text-gray-500 pointer-events-none">
                        {{ itemNames.length }}/{{ MAX_COUNT_ITEMS }}
                        <span class="pointer-events-auto cursor-pointer" @click="addItem()"
                            title="добавить 10 максимальных предметов">+</span>
                    </p>
                </div>
                <ItemsPanel v-if="favoriteItems" :favoriteItems="favoriteItems" @selectItem="selectItem" />
            </template>
            <template v-else-if="type_search == SearchType.Market">
                <label for="search">Поиск в 5 магазинах по названию</label>
                <input id="search" type="text" placeholder="Search..." v-model="searchQuery"
                    class="mt-2 border border-y-slate-500 focus:border-slate-500 focus:outline-none focus:ring-2 focus:ring-slate-500 focus:ring-offset-2 focus:ring-offset-white rounded-md px-2 py-1" />
                <button @click="search_in_market" class="btn hover:bg-gray-200 py-2 px-4 rounded w-1/2 mx-auto">Начать
                    поиск</button>
                <div class="flex flex-col  mb-8">
                    <div v-for="item in marketData">
                        <p class="text-sm text-gray-500 pointer-events-none" v-if="!item.loading">
                            {{ item.market }}: {{ item.data.length }}
                        </p>
                        <p class="text-sm text-gray-500 pointer-events-none" v-else>
                            {{ item.market }}: Загрузка
                        </p>
                    </div>
                </div>
            </template>
        </section>
        <section class="flex flex-col items-center w-full">
            <div id="search" class="w-3/4 bg-white rounded-md px-4 py-2" v-if="itemNames.length && type_search == SearchType.Base">
                <ul>
                    <li v-for="name in itemNames" :key="name">
                        <BaseToolbar>
                            <template #center>
                                <a class=" hover:text-blue-400 hover:cursor-pointer" @click="selectItem(name)"
                                    :href="`#${name}`">{{ name }}
                                </a>
                            </template>
                            <template #right>
                                <Fovorite :name_item="name" />
                            </template>
                        </BaseToolbar>
                    </li>
                </ul>
            </div>
            <section id="selected-item" class="w-full bg-white rounded-md px-4 py-2">
                <div v-if="type_search == SearchType.Base &&selectedItem">
                    <BaseToolbar>
                        <template #left>
                            <p>Выбранный предмет</p>
                        </template>
                        <template #center>
                            <p :id="selectedItem">{{ selectedItem }}</p>
                        </template>
                        <template #right>
                            <button class="text-red-500 hover:text-red-700 transition-colors duration-200 text-2xl"
                                title="Закрыть выбранный товар" @click="selectedItem = ''">X</button>
                        </template>
                    </BaseToolbar>
                    <BarChart :name_item="selectedItem"></BarChart>
                </div>
                <div v-else class="grid grid-cols-2">
                    <template v-for="item in marketData">
                        <RegionMarketInfo :info="item" v-if="item.data.length>0"/>
                    </template>
                </div>
            </section>
            <section id="Favorites" class="w-full bg-white rounded-md px-4 py-2 " v-if="favoriteItems.length">
                <ul class="grid grid-cols-2 gap-2">
                    <li v-for="name in favoriteItems">
                        <div>
                            <BaseToolbar>
                                <template #left>
                                    <p>Избранное</p>
                                </template>
                                <template #center>
                                    <p>{{ name }}</p>
                                </template>
                                <template #right>
                                    <Fovorite :name_item="name" :favoriteItems="favoriteItems" />
                                </template>
                            </BaseToolbar>
                            <BarChart :name_item="name"></BarChart>
                        </div>
                    </li>
                </ul>
            </section>
        </section>
    </section>
</template>
