<script lang="ts" setup>
import { inject } from "vue";
import type { Favorite  } from "../storage";

const { name_item } = defineProps(["name_item", "favoriteItems"]);

const { toggleFavorite, favoriteItems } = inject<Favorite>("Favourite")!;

</script>

<template>
    <button  @click="toggleFavorite(name_item);" class="cursor-pointer w-5 hover:bg-amber-100">
        <span :class="{
            'text-yellow-500 text-lg':
                favoriteItems.includes(name_item),
            'text-gray-400 text-lg':
                !favoriteItems.includes(name_item),
        }">
            {{ favoriteItems.includes(name_item) ? "★" : "☆" }}
        </span>
    </button>
</template>