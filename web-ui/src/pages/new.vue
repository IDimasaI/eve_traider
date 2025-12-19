<script lang="ts" setup>
import { nextTick } from 'vue';
import Header from '../components/new/Header.vue'
import type { Favorite } from "../storage.ts"
import type { MarketData } from "../utils/API.ts"


import { FuzzySearcher } from "../utils/search.ts";
import { ref, onMounted, provide } from "vue";
import { time_watch } from "../utils/delay.ts";
import { FavoritesManager } from "../utils/fovarites.ts";
import { find_id } from "../utils/API.ts";


import LeftPanel from "../components/new/LeftPanel.vue";
import ItemsPanel from '../components/CategoryPanel.vue';
import Fovorite from "../components/Favourite.vue";
import BarChart from "../components/BarChart.vue";
import RegionMarketInfo from "../components/RegionMarketInfo.vue";
import BaseToolbar from '../components/BaseToolbar.vue';
import MainContainer from '../components/new/MainContainer.vue';

import { type Items, get_all_items } from "../utils/API.ts";
import { useItemsHashMap } from '../composables/ItemsHashMap.ts';

enum SearchType {
    Base,
    Market
}




// Поисковые утилиты и переменные//
const addItem = async (): Promise<void> => {
    MAX_COUNT_ITEMS.value += 1000;
    const old_query = searchQuery.value;
    searchQuery.value = "";
    await nextTick();
    searchQuery.value = old_query
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
const itemNames = ref<{ name: string, id: number }[]>([]);

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
const { initializeMap } = useItemsHashMap()
onMounted(async () => {
    favoriteItems.value = favoritesManager.getFavorites();

    if (localStorage.getItem("items_id_name")) {
        all_names.value = JSON.parse(localStorage.getItem("items_id_name")!)!;
    } else {
        all_names.value = await get_all_items();
        localStorage.setItem("items_id_name", JSON.stringify(all_names.value));
    }

    initializeMap(all_names.value!);
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

            itemNames.value = results.map((result) => ({ name: result.item, id: result.id }));
            console.log(results);
        },
        searchQuery,
        500,
    );
});

</script>
<template>
    <section class="bg-[#141414] min-h-screen">
        <Header>
            <template v-slot:left>
                <a class="break-spaces text-center cursor-pointer" @click="type_search = SearchType.Base">
                    База данных
                </a>
                <a class="break-spaces text-center cursor-pointer" @click="type_search = SearchType.Market">
                    Магазин EVE
                </a>
            </template>
        </Header>
        <section class="flex flex-row mt-4">
            <LeftPanel>
                <template v-if="type_search == SearchType.Base">
                    <div class="flex flex-col mb-8 text-center text-white">
                        <label for="search">Поиск в базе</label>
                        <input id="search" type="text" placeholder="Search..." v-model="searchQuery"
                            class="mt-2 border border-y-slate-500 focus:border-slate-500 focus:outline-none focus:ring-2 focus:ring-slate-500 focus:ring-offset-2 focus:ring-offset-white rounded-md px-2 py-1" />
                        <p class="text-sm  pointer-events-none">
                            Найдено: {{ itemNames.length }} предметов
                        </p>
                    </div>
                    <div id="search" class="w-full bg-white rounded-md px-4 py-2"
                        v-if="itemNames.length && type_search == SearchType.Base">
                        <ul class="text-sm max-h-96 overflow-y-scroll">
                            <li v-for="item in itemNames" :key="item.id" class="odd:bg-gray-200">
                                <div class="flex flex-row justify-between">
                                    <i>
                                        <img class="w-6 h-6" :src="`https://images.evetech.net/types/${item.id}/icon`"
                                            loading="lazy" />
                                    </i>
                                    <p class=" hover:text-blue-400 hover:cursor-pointer h-auto w-full"
                                        @click="selectItem(item.name)" :href="`#${item.name}`">{{ item.name }}
                                    </p>
                                    <Fovorite :name_item="item.name" />
                                </div>
                            </li>
                            <button @click="addItem()"
                                class="btn hover:bg-gray-200 py-2 px-4 rounded w-full mx-auto">Показать все</button>
                        </ul>
                    </div>
                    <ItemsPanel v-if="favoriteItems" :favoriteItems="favoriteItems" @selectItem="selectItem" />
                </template>
                <template v-else-if="type_search == SearchType.Market">
                    <div class="text-center text-white w-full mx-auto">
                        <label for="search">Поиск в 5 магазинах по названию</label>
                        <input id="search" type="text" placeholder="Search..." v-model="searchQuery"
                            class="mt-2 w-full border border-y-slate-500 focus:border-slate-500 focus:outline-none focus:ring-2 focus:ring-slate-500 focus:ring-offset-2 focus:ring-offset-white rounded-md px-2 py-1" />
                        <button @click="search_in_market"
                            class="btn hover:bg-gray-200 py-2 px-4 rounded w-1/2 mx-auto">Начать
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
                    </div>
                </template>
            </LeftPanel>
            <MainContainer>
                <section id="Selected" v-show="type_search == SearchType.Base && selectedItem">
                    <BaseToolbar>
                        <template #left>
                            <p>Выбранный предмет</p>
                        </template>
                        <template #center>
                            <p :id="selectedItem">{{
                                selectedItem }}</p>
                        </template>
                        <template #right>
                            <button class="text-red-500 hover:text-red-700 transition-colors duration-200 text-2xl"
                                title="Закрыть выбранный товар" @click="selectedItem = ''">X</button>
                        </template>
                    </BaseToolbar>
                    <BarChart :name_item="selectedItem"></BarChart>
                </section>

                <div v-if="type_search == SearchType.Market" class="grid grid-cols-2">
                    <template v-for="item in marketData">
                        <RegionMarketInfo :info="item" v-if="item.data.length > 0" />
                    </template>
                </div>

                <section id="Favorites" class="w-full bg-[#1e293b]/50 rounded-md px-4 py-2 "
                    v-show="favoriteItems.length && type_search == SearchType.Base">
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
            </MainContainer>
        </section>
    </section>
</template>