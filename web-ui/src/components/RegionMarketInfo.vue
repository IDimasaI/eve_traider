<script setup lang="ts">
import type { PropType } from 'vue';
import type { MarketData } from '../utils/API';
import { computed, ref } from 'vue';
import { CurrentTheme, Themes } from '../composables/Theme';
const { info } = defineProps({
    info: {
        type: Object as PropType<MarketData>,
        required: true
    }
})

type SortOption = 'default' | 'price_desc' | 'price_asc' | 'volume_desc' | 'volume_asc';

const sortOption = ref<SortOption>('default');

// Отсортированные заказы
const sortedOrders = computed(() => {
    const orders = info?.data ? [...info.data] : [];

    switch (sortOption.value) {
        case 'price_desc':
            return orders.sort((a, b) => (b.price || 0) - (a.price || 0));
        case 'price_asc':
            return orders.sort((a, b) => (a.price || 0) - (b.price || 0));
        case 'volume_desc':
            return orders.sort((a, b) => (b.volume_remain || 0) - (a.volume_remain || 0));
        case 'volume_asc':
            return orders.sort((a, b) => (a.volume_remain || 0) - (b.volume_remain || 0));
        default:
            return orders;
    }
});

// Функция для расчета даты окончания
const calculateEndDate = (order: any): string => {
    if (!order?.issued) return 'Не указана';

    const issuedDate = new Date(order.issued);
    if (isNaN(issuedDate.getTime())) return 'Неверная дата';

    const endDate = new Date(issuedDate);
    endDate.setDate(endDate.getDate() + (order.duration || 0));

    return endDate.toLocaleDateString('ru-RU', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
    });
}

// Форматирование цены
const formatPrice = (price: number): string => {
    if (price == null) return '0.00';
    return price.toLocaleString('ru-RU', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
    });
}
const classes = computed(() => {
    return {
        color: CurrentTheme.value === Themes.dark ? 'text-[#f1f5f9]' : 'text-[#333]',
        bg: CurrentTheme.value === Themes.dark ? 'bg-[#1e293b]' : 'bg-white',
        border_color: CurrentTheme.value === Themes.dark ? 'border-gray-500' : 'border-gray-300'
    }
})
</script>

<template>
    <section :class="`p-4 ${classes.color} bg-[#1e293b]/50 ${CurrentTheme===Themes.dark ? 'dark-theme' : ''}`">
        <div class="">
            <h1 :class="`text-xl font-bold mb-4 text-[#f1f5f9]`">
                Информация о рынке {{ info?.market || 'Не указан' }}
            </h1>
            <div class="flex flex-col sm:flex-row items-start sm:items-center gap-2">
                <span class="font-medium text-[#f1f5f9]">Сортировка:</span>
                <select name="sort" v-model="sortOption"
                    class="period-select max-w-32 max-h-32">
                    <option value="default">По умолчанию</option>
                    <option value="price_desc">Цена ↓</option>
                    <option value="price_asc">Цена ↑</option>
                    <option value="volume_desc">Объем ↓</option>
                    <option value="volume_asc">Объем ↑</option>
                </select>
            </div>
        </div>
        <div v-if="sortedOrders.length" class="space-y-4">
            <div v-for="(order, index) in sortedOrders" :key="order.order_id || index"
                :class="`grid grid-cols-1 md:grid-cols-2 gap-4 p-4 ${classes.bg} rounded-lg`">

                <div class="space-y-2">
                    <h3 class="font-semibold text-lg">Ордер #{{ index + 1 }}</h3>
                    <p><span class="font-medium">Товар ID:</span> {{ order.type_id }}</p>
                    <p><span class="font-medium">Система ID:</span> {{ order.location_id }}</p>
                    <p><span class="font-medium">Создан:</span>
                        {{ order.issued ? new Date(order.issued).toLocaleString('ru-RU') : 'Не указан' }}
                    </p>
                    <p><span class="font-medium">Истекает:</span> {{ calculateEndDate(order) }}</p>
                </div>

                <div class="space-y-2">
                    <p><span class="font-medium">Цена:</span> {{ formatPrice(order.price) }} ISK</p>
                    <p><span class="font-medium">Объем:</span>
                        {{ order.volume_remain?.toLocaleString() || 0 }} /
                        {{ order.volume_total?.toLocaleString() || 0 }}
                    </p>
                    <div class="w-full bg-gray-200 rounded-full h-2.5">
                        <div class="bg-blue-600 h-2.5 rounded-full"
                            :style="{ width: `${(order.volume_remain / order.volume_total) * 100 || 0}%` }">
                        </div>
                    </div>
                    <p :class="`text-sm stat-label`">
                        Осталось {{ order.volume_remain || 0 }}
                    </p>
                </div>
            </div>
        </div>

        <div v-else class="text-center py-8 text-gray-500">
            Нет доступных ордеров на рынке
        </div>
    </section>
</template>