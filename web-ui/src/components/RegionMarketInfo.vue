<script setup lang="ts">
import { type PropType, watch, computed, onMounted, ref } from 'vue';
import { type Order, type MarketData } from '../utils/API';
import { CurrentTheme, Themes } from '../composables/Theme';
import { formatPrice } from '../utils/formats';

const { info_all } = defineProps({
    info_all: {
        type: Array as PropType<MarketData[]>,
        required: true
    }
})

const ui_variants = ref<"eve_m" | "single">("single");

type SortOption = 'default' | 'price_desc' | 'price_asc' | 'volume_desc' | 'volume_asc' |
    'profitability_desc' | 'location_id_asc' | 'location_id_desc' |
    'system_name_asc' | 'system_name_desc';

type SystemData = {
    name: string;
    security_status?: number;
}

const sortOption_buy = ref<SortOption>('profitability_desc');
const sortOption_sell = ref<SortOption>('profitability_desc');
const location_names = ref(new Map<number, SystemData>());

// Создаем enum для уровней безопасности
enum SecurityLevel {
    HIGH = 'high',      // 1.0 - 0.5
    LOW = 'low',        // 0.5 - 0.0
    NULL = 'null',      // < 0.0
    ANY = 'any'         // все уровни
}

// Реф для выбранных уровней безопасности (может быть несколько)
const selectedSecurityLevels = ref<SecurityLevel[]>([SecurityLevel.HIGH, SecurityLevel.LOW, SecurityLevel.NULL]);

// Разделяем ордера на покупки и продажи сразу
const preparedOrders = computed(() => {
    const allOrders = info_all.flatMap(item => item.data || []);
    return {
        buy: allOrders.filter(order => order.is_buy_order === true),
        sell: allOrders.filter(order => order.is_buy_order === false)
    };
});

// Функция для расчета выгодности
const calculateProfitability = (order: any): number => {
    if (!order?.price || !order?.volume_remain) return 0;
    return order.price * order.volume_remain;
}

// Карта названий систем
const systemDataMap = computed<Map<number, SystemData>>(() => {
    const map = new Map<number, SystemData>();
    info_all.forEach(marketData => {
        if (marketData.data) {
            marketData.data.forEach(order => {
                if (order.system_id && marketData.Region) {
                    map.set(order.system_id, { name: marketData.Region });
                }
            });
        }
    });
    return map;
});


const getRegionName = (systemId: number): string => {
    return systemDataMap.value.get(systemId)?.name || `Регион ${systemId}`;
}

const getLocationName = (locationId: number): string => {
    return location_names.value.get(locationId)?.name || `Локация ${locationId}`;
}


// Нч функций безопасности
const getSecurityLevel = (securityStatus: number): SecurityLevel => {
    if (securityStatus >= 0.5) return SecurityLevel.HIGH;
    if (securityStatus >= 0) return SecurityLevel.LOW;
    return SecurityLevel.NULL;
}

//Сделать через мап
const getSystemSecurityStatus = (systemId: number): number => {
    //@ts-ignore
    for (const [locationId, systemData] of location_names.value.entries()) {
        // Здесь нужно иметь маппинг location_id -> system_id
        // Но так как его нет, мы будем использовать другой подход
        // Сначала найдем все ордера в этой системе
        const orderInSystem = info_all.flatMap(item => item.data || [])
            .find(order => order.system_id === systemId);

        if (orderInSystem) {
            const locationData = location_names.value.get(orderInSystem.location_id);
            if (locationData?.security_status !== undefined) {
                return locationData.security_status;
            }
        }
    }
    return 0;
}

const getSystemSecurityLevel = (systemId: number): SecurityLevel => {
    const securityStatus = getSystemSecurityStatus(systemId);
    return getSecurityLevel(securityStatus);
}

const passesSecurityFilter = (order: Order): boolean => {
    if (selectedSecurityLevels.value.length === 0 ||
        selectedSecurityLevels.value.includes(SecurityLevel.ANY)) {
        return true;
    }

    const securityLevel = getSystemSecurityLevel(order.system_id);
    return selectedSecurityLevels.value.includes(securityLevel);
}

const initSecurityStatus = async (system_id: number) => {
    try {
        const res = await fetch(`https://esi.evetech.net/universe/systems/${system_id}/?datasource=tranquility`);
        if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
        const data = await res.json();
        return data.security_status;
    } catch (error) {
        console.error('Error fetching security status:', error);
        return 0;
    }
}
//Кн функций безопасности



let need_update_local_storage = false;

const get_from_API_LocationName = async (locationId: number, systemId: number): Promise<void> => {
    if (location_names.value.has(locationId)) return;
    let locations: Record<string, SystemData> = JSON.parse(localStorage.getItem('locations') || "{}");

    if (!locations[locationId]) {
        try {
            const res = await fetch('https://esi.evetech.net/universe/names', {
                body: `[${locationId}]`,
                method: 'POST'
            });
            if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
            const data = await res.json();
            const security_status = await initSecurityStatus(systemId);

            locations[locationId] = {
                name: data[0].name,
                security_status: security_status
            };
        } catch (error) {
            const security_status = await initSecurityStatus(systemId);
            locations[locationId] = {
                name: `Локация ${locationId}`,
                security_status: security_status
            };
            console.error('Error fetching location name:', error);
        }
        need_update_local_storage = true;
    }

    location_names.value.set(locationId, locations[locationId]!);
}

// Вспомогательная функция для сортировки
const sortOrders = (orders: any[], sortOption: SortOption) => {
    const sorted = [...orders];

    switch (sortOption) {
        case 'price_desc':
            return sorted.sort((a, b) => (b?.price || 0) - (a?.price || 0));
        case 'price_asc':
            return sorted.sort((a, b) => (a?.price || 0) - (b?.price || 0));
        case 'volume_desc':
            return sorted.sort((a, b) => (b?.volume_remain || 0) - (a?.volume_remain || 0));
        case 'volume_asc':
            return sorted.sort((a, b) => (a?.volume_remain || 0) - (b?.volume_remain || 0));
        case 'profitability_desc':
            return sorted.sort((a, b) => calculateProfitability(b) - calculateProfitability(a));
        case 'location_id_asc':
            return sorted.sort((a, b) => (a?.location_id || 0) - (b?.location_id || 0));
        case 'location_id_desc':
            return sorted.sort((a, b) => (b?.location_id || 0) - (a?.location_id || 0));
        case 'system_name_asc':
            return sorted.sort((a, b) => {
                const nameA = getRegionName(a?.system_id || 0);
                const nameB = getRegionName(b?.system_id || 0);
                return nameA.localeCompare(nameB, 'ru');
            });
        case 'system_name_desc':
            return sorted.sort((a, b) => {
                const nameA = getRegionName(a?.system_id || 0);
                const nameB = getRegionName(b?.system_id || 0);
                return nameB.localeCompare(nameA, 'ru');
            });
        default:
            return sorted;
    }
};

 
const getFilteredSortedOrders = (orders: Order[], sortOption: SortOption) => {
    const filtered = orders.filter(passesSecurityFilter);
    return sortOrders(filtered, sortOption);
};

const filteredSortedOrders_sell = computed<Order[]>(() => {
    return getFilteredSortedOrders(preparedOrders.value.sell, sortOption_sell.value);
});

const filteredSortedOrders_buy = computed<Order[]>(() => {
    return getFilteredSortedOrders(preparedOrders.value.buy, sortOption_buy.value);
});

function sortOrders_sell(option: SortOption, changeOption: SortOption) {
    sortOption_sell.value = sortOption_sell.value === option ? changeOption : option;
}

function sortOrders_buy(option: SortOption, changeOption: SortOption) {
    sortOption_buy.value = sortOption_buy.value === option ? changeOption : option;
}

// Функция для расчета даты окончания ордера
const calculateEndDate = (order: any): string => {
    if (!order?.issued || !order?.duration) return 'Не указана';

    const issuedDate = new Date(order.issued);
    if (isNaN(issuedDate.getTime())) return 'Неверная дата';

    const endDate = new Date(issuedDate);
    endDate.setDate(endDate.getDate() + order.duration);

    return endDate.toLocaleDateString('ru-RU', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
    });
}


const getProfitabilityFormatted = (order: any): string => {
    const profitability = calculateProfitability(order);
    return formatPrice(profitability);
}


// Дизайн
const classes = computed(() => {
    return {
        color: CurrentTheme.value === Themes.dark ? 'text-[#f1f5f9]' : 'text-[#333]',
        bg: CurrentTheme.value === Themes.dark ? 'bg-[#1e293b]' : 'bg-white',
        border_color: CurrentTheme.value === Themes.dark ? 'border-gray-500' : 'border-gray-300'
    }
});

// Функция для переключения уровней безопасности
const toggleSecurityLevel = (level: SecurityLevel) => {
    const index = selectedSecurityLevels.value.indexOf(level);

    if (level === SecurityLevel.ANY) {
        // Если выбрали "Все", снимаем все остальные
        selectedSecurityLevels.value = [SecurityLevel.ANY];
        return;
    }

    // Если выбрали конкретный уровень, убираем "Все" если он был выбран
    const anyIndex = selectedSecurityLevels.value.indexOf(SecurityLevel.ANY);
    if (anyIndex > -1) {
        selectedSecurityLevels.value.splice(anyIndex, 1);
    }

    if (index > -1) {
        // Убираем уровень, если он уже выбран
        selectedSecurityLevels.value.splice(index, 1);
    } else {
        // Добавляем уровень
        selectedSecurityLevels.value.push(level);
    }

    // Если ничего не выбрано, выбираем "Все"
    if (selectedSecurityLevels.value.length === 0) {
        selectedSecurityLevels.value.push(SecurityLevel.ANY);
    }
};

// Проверка, выбран ли уровень безопасности
const isSecurityLevelSelected = (level: SecurityLevel): boolean => {
    return selectedSecurityLevels.value.includes(level);
};

//Ватчеры
watch(() => info_all, async (newInfo: any[]) => {
    if (newInfo) {
        const allOrders = newInfo.flatMap((item: { data: any; }) => item.data || []);

        // Создаем массив уникальных пар location_id + system_id
        const uniquePairs = new Map<number, number>();

        allOrders.forEach((order: any) => {
            if (order.location_id && order.system_id) {
                uniquePairs.set(order.location_id, order.system_id);
            }
        });

        // Загружаем названия для всех уникальных location_id с их system_id
        const promises = Array.from(uniquePairs.entries()).map(
            ([locationId, systemId]) => get_from_API_LocationName(locationId, systemId)
        );

        await Promise.allSettled(promises);

        if (need_update_local_storage) {
            console.log("Обновляем локальное хранилище");
            localStorage.setItem('locations', JSON.stringify(Object.fromEntries(location_names.value)));
            need_update_local_storage = false;
        }
    }
}, { immediate: true, deep: true });

onMounted(async () => {
    // Ничего не делаем специального при монтировании
})
</script>

<template>
    <section :class="`p-4 ${classes.color} bg-[#1e293b]/50 ${CurrentTheme === Themes.dark ? 'dark-theme' : ''}`">

        <button class="mb-6" @click="ui_variants = ui_variants === 'eve_m' ? 'single' : 'eve_m'">
            {{ ui_variants === 'eve_m' ? 'Еве М' : 'Еве С' }}
        </button>

        <!-- Компонент фильтра по безопасности -->
        <div class="mb-4 p-3 bg-gray-800/50 rounded-lg">
            <h3 class="font-medium text-[#f1f5f9] mb-2">Фильтр по безопасности системы:</h3>
            <div class="flex flex-wrap gap-2">
                <button @click="toggleSecurityLevel(SecurityLevel.ANY)"
                    :class="`px-3 py-1 rounded ${isSecurityLevelSelected(SecurityLevel.ANY) ? 'bg-blue-600 text-white' : 'bg-gray-700 text-gray-300 hover:bg-gray-600'}`"
                    title="Показать все ордера независимо от уровня безопасности">
                    Все
                </button>
                <button @click="toggleSecurityLevel(SecurityLevel.HIGH)"
                    :class="`px-3 py-1 rounded ${isSecurityLevelSelected(SecurityLevel.HIGH) ? 'bg-green-600 text-white' : 'bg-gray-700 text-gray-300 hover:bg-gray-600'}`"
                    title="Высокая безопасность (1.0 - 0.5)">
                    Высокая
                    <span class="text-xs ml-1 opacity-75">1.0-0.5</span>
                </button>
                <button @click="toggleSecurityLevel(SecurityLevel.LOW)"
                    :class="`px-3 py-1 rounded ${isSecurityLevelSelected(SecurityLevel.LOW) ? 'bg-yellow-600 text-white' : 'bg-gray-700 text-gray-300 hover:bg-gray-600'}`"
                    title="Низкая безопасность (0.5 - 0.0)">
                    Низкая
                    <span class="text-xs ml-1 opacity-75">0.5-0.0</span>
                </button>
                <button @click="toggleSecurityLevel(SecurityLevel.NULL)"
                    :class="`px-3 py-1 rounded ${isSecurityLevelSelected(SecurityLevel.NULL) ? 'bg-red-600 text-white' : 'bg-gray-700 text-gray-300 hover:bg-gray-600'}`"
                    title="Нулевая безопасность (< 0.0)">
                    Нулевая
                    <span class="text-xs ml-1 opacity-75">
                        < 0.0</span>
                </button>
            </div>
            <p class="text-sm text-gray-400 mt-2">
                Показаны только ордера в системах с выбранным уровнем безопасности.
                <span class="block">Выбрано:
                    <span v-if="isSecurityLevelSelected(SecurityLevel.ANY)">все уровни</span>
                    <template v-else>
                        <span v-if="isSecurityLevelSelected(SecurityLevel.HIGH)">Высокая</span>
                        <span v-if="isSecurityLevelSelected(SecurityLevel.LOW)">, Низкая</span>
                        <span v-if="isSecurityLevelSelected(SecurityLevel.NULL)">, Нулевая</span>
                    </template>
                </span>
            </p>
        </div>

        <template v-if="ui_variants == 'single'">
            <div class="mb-6">
                <h1 :class="`text-xl font-bold mb-4 text-[#f1f5f9]`">
                    Сравнение магазинов ({{ filteredSortedOrders_sell.length }} ордеров)
                </h1>
                <div class="flex flex-col sm:flex-row items-start sm:items-center gap-4">
                    <span class="font-medium text-[#f1f5f9]">Сортировка:</span>
                    <select name="sort" v-model="sortOption_sell"
                        class="period-select max-w-48 max-h-32 px-3 py-2 rounded border">
                        <option value="default">По умолчанию</option>
                        <option value="profitability_desc">Выгодность ↓</option>
                        <option value="price_desc">Цена ↓</option>
                        <option value="price_asc">Цена ↑</option>
                        <option value="volume_desc">Объем ↓</option>
                        <option value="volume_asc">Объем ↑</option>
                    </select>
                    <div class="text-sm text-[#94a3b8]">
                        <span v-if="sortOption_sell === 'profitability_desc'">(Сортировка по цене × доступный
                            объем)</span>
                        <span v-if="sortOption_sell == 'price_asc'">(Сортировка по убыванию цены)</span>
                        <span v-if="sortOption_sell == 'price_desc'">(Сортировка по возрастанию цены)</span>
                        <span v-if="sortOption_sell == 'volume_desc'">(Сортировка по убыванию объема)</span>
                        <span v-if="sortOption_sell == 'volume_asc'">(Сортировка по возрастанию объема)</span>
                    </div>
                </div>
            </div>

            <div v-show="filteredSortedOrders_sell.length" class="space-y-4">
                <div v-for="(order, index) in filteredSortedOrders_sell" :key="order?.order_id || index"
                    :class="`grid grid-cols-1 md:grid-cols-3 gap-4 p-4 ${classes.bg} rounded-lg border ${classes.border_color}`">
                    <div class="space-y-2">
                        <h3 class="font-semibold text-lg">Ордер #{{ index + 1 }}</h3>
                        <p><span class="font-medium">Товар ID:</span> {{ order?.type_id }}</p>
                        <p><span class="font-medium">Регион:</span> {{ getRegionName(order?.system_id) }}</p>
                        <p><span class="font-medium">Магазин:</span> {{ getLocationName(order?.location_id) }}</p>
                        <p>
                            <span class="font-medium">Безопасность:</span>
                            <span :class="{
                                'text-green-500': getSystemSecurityLevel(order.system_id) === SecurityLevel.HIGH,
                                'text-yellow-500': getSystemSecurityLevel(order.system_id) === SecurityLevel.LOW,
                                'text-red-500': getSystemSecurityLevel(order.system_id) === SecurityLevel.NULL
                            }">
                                {{ getSystemSecurityStatus(order.system_id).toFixed(2) }}
                                ({{ getSystemSecurityLevel(order.system_id) === SecurityLevel.HIGH ? 'Высокая' :
                                    getSystemSecurityLevel(order.system_id) === SecurityLevel.LOW ? 'Низкая' : 'Нулевая' }})
                            </span>
                        </p>
                        <p><span class="font-medium">Создан:</span>
                            {{ order?.issued ? new Date(order.issued).toLocaleString('ru-RU') : 'Не указан' }}
                        </p>
                        <p><span class="font-medium">Истекает:</span> {{ calculateEndDate(order) }}</p>
                    </div>

                    <div class="space-y-2">
                        <p><span class="font-medium">Цена:</span> {{ formatPrice(order?.price) }} ISK</p>
                        <p><span class="font-medium">Объем:</span>
                            {{ order?.volume_remain?.toLocaleString() || 0 }} /
                            {{ order?.volume_total?.toLocaleString() || 0 }}
                        </p>
                        <div class="w-full bg-gray-200 rounded-full h-2.5">
                            <div class="bg-blue-600 h-2.5 rounded-full" :style="{
                                width: `${(order?.volume_remain / order?.volume_total) * 100 || 0}%`
                            }">
                            </div>
                        </div>
                        <p :class="`text-sm ${classes.color}`">
                            Осталось {{ order?.volume_remain || 0 }} ({{
                                Math.round((order?.volume_remain / order?.volume_total) * 100) || 0
                            }}%)
                        </p>
                    </div>

                    <div class="space-y-2">
                        <div class="p-3 bg-gray-800/50 rounded">
                            <p class="font-medium text-[#f1f5f9]">Выгодность:</p>
                            <p class="text-2xl font-bold text-green-400">
                                {{ getProfitabilityFormatted(order) }} ISK
                            </p>
                            <p class="text-sm text-[#94a3b8]">
                                (Цена × Доступный объем)
                            </p>
                        </div>
                        <div v-if="sortOption_sell === 'profitability_desc'" class="mt-2 flex items-center">
                            <span class="text-lg font-bold" :class="index < 3 ? 'text-yellow-400' : 'text-gray-400'">
                                #{{ index + 1 }}
                            </span>
                            <span v-if="index === 0" class="ml-2 text-xs bg-yellow-500 text-black px-2 py-1 rounded">
                                Самый выгодный
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        </template>
        <template v-else>
            <div class="flex flex-row justify-around mb-4">
                <a href="#sell" class="text-2xl font-bold text-green-400">Продажи</a>
                <a href="#buy" class="text-2xl font-bold text-green-400">Покупки</a>
            </div>
            <div class="grid grid-cols-2 max-lg:grid-cols-1 gap-4">
                <!-- ТАБЛИЦА ПРОДАЖ -->
                <div id="sell" class="flex flex-col">
                    <table>
                        <thead class="text-left">
                            <tr>
                                <th @click="sortOrders_sell('volume_desc', 'volume_asc')">Объем
                                    <span v-if="sortOption_sell == 'volume_desc'">↓</span>
                                    <span v-if="sortOption_sell == 'volume_asc'">↑</span>
                                </th>
                                <th @click="sortOrders_sell('price_desc', 'price_asc')">Цена
                                    <span v-if="sortOption_sell == 'price_desc'">↓</span>
                                    <span v-if="sortOption_sell == 'price_asc'">↑</span>
                                </th>
                                <th @click="sortOrders_sell('system_name_desc', 'system_name_asc')">Регион
                                    <span v-if="sortOption_sell == 'system_name_desc'">↓</span>
                                    <span v-if="sortOption_sell == 'system_name_asc'">↑</span>
                                </th>
                                <th>sec</th>
                                <th>Локация</th>
                            </tr>
                        </thead>
                        <tbody class="text-sm">
                            <tr v-for="(order, index) in filteredSortedOrders_sell"
                                :key="`sell-${index}-${order.volume_remain}`">
                                <td>{{ order.volume_remain }}</td>
                                <td>{{ formatPrice(order.price) }}</td>
                                <td>{{ getRegionName(order.system_id) }}</td>
                                <td :class="{
                                    'text-green-500': getSystemSecurityLevel(order.system_id) === SecurityLevel.HIGH,
                                    'text-yellow-500': getSystemSecurityLevel(order.system_id) === SecurityLevel.LOW,
                                    'text-red-500': getSystemSecurityLevel(order.system_id) === SecurityLevel.NULL
                                }">
                                    {{ getSystemSecurityStatus(order.system_id).toFixed(2) }}
                                </td>
                                <td class="max-lg:w-2xs">{{ getLocationName(order.location_id) }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <!-- ТАБЛИЦА ПОКУПОК -->
                <div id="buy" class="flex flex-col">
                    <table>
                        <thead class="text-left">
                            <tr>
                                <th @click="sortOrders_buy('volume_desc', 'volume_asc')">Объем
                                    <span v-if="sortOption_buy == 'volume_desc'">↓</span>
                                    <span v-if="sortOption_buy == 'volume_asc'">↑</span>
                                </th>
                                <th @click="sortOrders_buy('price_desc', 'price_asc')">Цена
                                    <span v-if="sortOption_buy == 'price_desc'">↓</span>
                                    <span v-if="sortOption_buy == 'price_asc'">↑</span>
                                </th>
                                <th @click="sortOrders_buy('system_name_desc', 'system_name_asc')">Регион
                                    <span v-if="sortOption_buy == 'system_name_desc'">↓</span>
                                    <span v-if="sortOption_buy == 'system_name_asc'">↑</span>
                                </th>
                                <th>sec</th>
                                <th>Локация</th>
                            </tr>
                        </thead>
                        <tbody class="text-sm">
                            <tr v-for="(order, index) in filteredSortedOrders_buy"
                                :key="`buy-${index}-${order.volume_remain}`">
                                <td>{{ order.volume_remain }}</td>
                                <td>{{ formatPrice(order.price) }}</td>
                                <td>{{ getRegionName(order.system_id) }}</td>
                                <td :class="{
                                    'text-green-500': getSystemSecurityLevel(order.system_id) === SecurityLevel.HIGH,
                                    'text-yellow-500': getSystemSecurityLevel(order.system_id) === SecurityLevel.LOW,
                                    'text-red-500': getSystemSecurityLevel(order.system_id) === SecurityLevel.NULL
                                }">
                                    {{ getSystemSecurityStatus(order.system_id).toFixed(2) }}
                                </td>
                                <td class="max-lg:w-2xs">{{ getLocationName(order.location_id) }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </template>
    </section>
</template>