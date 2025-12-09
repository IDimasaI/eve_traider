<script lang="ts" setup>
import { inject } from "vue";
import type { Favorite  } from "../storage";

const { name_item } = defineProps(["name_item", "favoriteItems"]);

const { toggleFavorite, favoriteItems } = inject<Favorite>("Favourite")!;

</script>

<template>
    <button  @click="toggleFavorite(name_item);" class="cursor-pointer">
        <span :class="{
            'text-yellow-500 text-xl':
                favoriteItems.includes(name_item),
            'text-gray-400 text-xl':
                !favoriteItems.includes(name_item),
        }">
            {{ favoriteItems.includes(name_item) ? "★" : "☆" }}
        </span>
    </button>
</template>