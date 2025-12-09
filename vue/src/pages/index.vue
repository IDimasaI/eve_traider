<script lang="ts" setup>
import type { Favorite } from "../storage.ts"

import CountItems from "../components/CountItems.vue";
import BaseToolbar from "../components/BaseToolbar.vue";
import { FuzzySearcher } from "../utils/search.ts";
import { ref, onMounted, provide } from "vue";
import { time_watch } from "../utils/delay.ts";
import { FavoritesManager } from "../utils/fovarites.ts";
import LeftPanel from "../components/LeftPanel.vue";
import Fovorite from "../components/Favourite.vue";

import BarChart from "../components/BarChart.vue";


type Items = {
    name: string;
    category: string;
    id: number;
}[];


async function get_all_items() {
    const res = await fetch("/api/v2/get_all_items");
    return (await res.json()) as Items;
}

// Поисковые утилиты и переменные//
const addItem = (): void => {
    MAX_COUNT_ITEMS.value += 10;
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




const favoritesManager = new FavoritesManager();

const searchQuery = ref("");
const itemNames = ref<string[]>([]);

const selectedItem = ref("");

let MAX_COUNT_ITEMS = ref(30);


// КН Поисковые утилиты и переменные//





const favoriteItems = ref<string[]>([]);

const all_names = ref<Items>();
provide<Favorite>('Favourite', {
    toggleFavorite,
    favoriteItems
})

onMounted(async () => {
    favoriteItems.value = favoritesManager.getFavorites();

    if (localStorage.getItem("items_id_name")) {
        all_names.value = JSON.parse(localStorage.getItem("items_id_name")!)!;
    } else {
        all_names.value = await get_all_items();
        localStorage.setItem("items_id_name", JSON.stringify(all_names.value));
    }

    const missileSearcher = new FuzzySearcher(all_names.value!.map((item) => item.name));
    time_watch(
        (query: string) => {
            if (query.length < 3) {
                itemNames.value = [];
                return;
            }
            const results = missileSearcher.search(query, {
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
            <div class="flex flex-col  mb-8">
                <label for="search">Поиск </label>
                <input id="search" type="text" placeholder="Search..." v-model="searchQuery"
                    class="border border-y-slate-500 focus:border-slate-500 focus:outline-none focus:ring-2 focus:ring-slate-500 focus:ring-offset-2 focus:ring-offset-white rounded-md px-2 py-1" />
                <p class="text-sm text-gray-500 pointer-events-none">
                    {{ itemNames.length }}/{{ MAX_COUNT_ITEMS }}
                    <span class="pointer-events-auto cursor-pointer" @click="addItem()"
                        title="добавить 10 максимальных предметов">+</span>
                </p>
            </div>
            <LeftPanel v-if="favoriteItems" :favoriteItems="favoriteItems" @selectItem="selectItem" />
        </section>
        <section class="flex flex-col items-center w-full">
            <div class="w-1/4 bg-gray-300 rounded-md px-4 py-2">
                <CountItems :count="all_names?.length" />
            </div>
            <div id="search" class="w-3/4 bg-white rounded-md px-4 py-2" v-if="itemNames.length">
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
            <section id="selected-item" class="w-full bg-white rounded-md px-4 py-2" v-if="selectedItem">
                <div>
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
