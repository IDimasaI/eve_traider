<template>
  <div class="price-chart-container">
    <div class="chart-header">
      <h3>Динамика цен за месяц<small>(за все время, приложение еще не проработало и 20 дней).</small></h3>
    </div>

    <div class="chart-wrapper">
      <Line :data="chartData" :options="chartOptions" />
    </div>

    <div class="chart-info">
      <div class="stats">
        <div class="stat-item">
          <span class="stat-label">Текущая цена:</span>
          <span class="stat-value">{{ formatCurrency(currentPrice) }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">Средняя цена:</span>
          <span class="stat-value">{{ formatCurrency(averagePrice) }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">Минимальная:</span>
          <span class="stat-value min-price">{{ formatCurrency(minPrice) }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">Максимальная:</span>
          <span class="stat-value max-price">{{ formatCurrency(maxPrice) }}</span>
        </div>
      </div>
      <div id="more">
        <button class="more-button" @click="toggleMore">Показать больше</button>
        <div v-if="moreInfo.open" class="more-info flex flex-col">
          <div class="w-full">
            <select v-model="selectedPeriod" class="w-full">
              <option v-for="period in moreInfo.periods" :key="period" :value="period">{{ period }}</option>
            </select>
          </div>
          <section class="w-full flex flex-col  gap-1 justify-evenly text-xs">
            <div class="grid grid-cols-4 w-full">
              <li>Дата</li>
              <li>Объем</li>
              <li>Заказов</li>
              <li>~Обьем заказов</li>
            </div>
            <template v-for="(_, index) in filteredData.times" :key="index">
              <div class="grid grid-cols-4 w-full odd:bg-gray-100">
                <li>{{ filteredData.times[index] }}</li>
                <li>{{ filteredData.volume_per_month[index] }}</li>
                <li>{{ filteredData.order_counts[index] }}</li>
                <li>
                  {{ formatDeltaOrderVolume(filteredData.order_counts[index] || 0,
                    filteredData.volume_per_month[index] || 0) }}
                </li>
              </div>
            </template>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { Chart, registerables, type ChartData, type ChartOptions } from 'chart.js'
import { Line } from 'vue-chartjs'

// Регистрируем компоненты Chart.js 
Chart.register(...registerables)

type PriceItem = {
  item_id: number;
  timestamp: number;
  price: string;
  error: boolean;
  day: string;
}


type MoreData = {
  open: boolean;
  periods: string[];
  times: string[];
  volume_per_month: number[];
  order_counts: number[]
}

interface ItemIdMap {
  name: string;
  id: number;
}

interface EveIdsResponse {
  inventory_types: Array<{ id: number }>;
}



const props = defineProps<{
  name_item: string;
}>()

const SKIP_EVERY_X_ROW = ref<number>(1)
//const isDaily = ref<boolean>(true)
const dailyPrices = ref<PriceItem[]>([])
//const weeklyPrices = ref<PriceItem[]>([])
const moreInfo = ref<MoreData>({
  open: false,
  periods: [],
  times: [],
  volume_per_month: [],
  order_counts: []
})

const id_item = await find_id(props.name_item)

const toggleMore = async () => {
  moreInfo.value.open = !moreInfo.value.open
  if (moreInfo.value.open) {
    if (!id_item) return
    await getMoreData(id_item, 10000002)
  }
  console.log(moreInfo.value)
}


async function getMoreData(id: number, region_id: number): Promise<void> {
  const res = await fetch(`https://esi.evetech.net/markets/${region_id}/history?type_id=${id}`)
  if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
  const data = await res.json()

  moreInfo.value.periods = get_periods(data.map((item: any) => item.date))

  moreInfo.value.times = data.map((item: any) => item.date)
  moreInfo.value.volume_per_month = data.map((item: any) => item.volume)
  moreInfo.value.order_counts = data.map((item: any) => item.order_count)
}

function get_periods(times: string[]): string[] {
  const periodsSet = new Set<string>()

  for (const time of times) {
    if (!time) continue

    const date = new Date(time)

    // Пропускаем невалидные даты
    if (isNaN(date.getTime())) continue

    // Месяц: 0=январь, 11=декабрь → нужно +1
    const month = (date.getMonth() + 1).toString().padStart(2, '0')
    const year = date.getFullYear()

    periodsSet.add(`${year}-${month}`)
  }

  // Сортируем периоды (опционально)
  return Array.from(periodsSet).sort()
}

const selectedPeriod = ref<string>(`${new Date().getFullYear()}-${(new Date().getMonth() + 1).toString().padStart(2, '0')}`)

const filteredData = computed(() => {
  if (!selectedPeriod.value) {
    return moreInfo.value
  }

  const info: MoreData = {
    open: moreInfo.value.open,
    times: [],
    volume_per_month: [],
    order_counts: [],
    periods: []
  }

  for (let i = 0; i < moreInfo.value.times.length; i++) {
    const time = moreInfo.value.times[i]
    if (time) {
      // Извлекаем год и месяц из time
      const parts = time.split('-')
      const yearMonth = parts.slice(0, 2).join('-')
      if (yearMonth === selectedPeriod.value) {
        info.times.push(time)
        info.volume_per_month.push(moreInfo.value.volume_per_month[i] || 0)
        info.order_counts.push(moreInfo.value.order_counts[i] || 0)
      }
    }
  }
  return info
})

const formatDeltaOrderVolume = (order: number, volume: number): string => {
  return `${(volume / order).toFixed(1)}`
}


async function find_id(name: string): Promise<number | null> {
  const stored = localStorage.getItem("items_id_name")
  if (stored) {
    try {
      const items: ItemIdMap[] = JSON.parse(stored)
      const foundItem = items.find(item => item.name === name)
      if (foundItem) {

        return foundItem.id
      }
    } catch (e) {
      console.error("Ошибка парсинга localStorage:", e)
    }
  }

  // Запрос к API, если не найдено локально
  try {
    const res = await fetch("https://esi.evetech.net/universe/ids", {
      method: "POST",
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify([name]),
    })

    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)

    const data: EveIdsResponse = await res.json()
    return data.inventory_types[0]?.id || null
  } catch (error) {
    console.error("Ошибка при запросе ID:", error)
    return null
  }
}
const generateData = async (): Promise<void> => {


  let prices = await get_price(id_item)

  let iter = 4

  if (prices[0] && prices[0].error) {
    for (let i = 0; i < iter; i++) {
      // prices = await get_price(id_item)
      if (!prices[0]?.error) {
        break
      } else {
        await new Promise(resolve => setTimeout(resolve, 1000))
      }
      if (i === iter - 1) {
        console.error("Не удалось получить данные для:" + props.name_item)
        break
      }
    }
  }

  // Если prices - массив
  if (Array.isArray(prices)) {
    const sortedPrices = [...prices].sort((a, b) => {
      const dateA = new Date(a.timestamp).getTime()
      const dateB = new Date(b.timestamp).getTime()
      return dateA - dateB
    })

    let counter = 0
    const filteredPrices: PriceItem[] = []

    sortedPrices.forEach(item => {
      if (item.item_id == id_item) {
        if (counter % SKIP_EVERY_X_ROW.value === 0) {
          filteredPrices.push({
            ...item,
            day: `${new Date(item.timestamp).getMonth() + 1}/${new Date(item.timestamp).getDate()}`
          })
        }
        counter++
      }
    })

    dailyPrices.value = filteredPrices
  }
}



const get_price = async (id: number | null): Promise<PriceItem[]> => {
  if (!id) {
    return [{ item_id: 0, timestamp: 0, price: "0", error: true, day: "" }]
  }

  try {
    const res = await fetch(`/api/v2/get_prices?id=${id}`)
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`)
    return await res.json()
  } catch (error) {
    console.error("Ошибка при получении цен:", error)
    return [{ item_id: id, timestamp: Date.now(), price: "0", error: true, day: "" }]
  }
}

const formatCurrency = (value: number): string => {
  return `${value.toFixed(2)} Isk`
}

const averagePrice = computed<number>(() => {
  if (dailyPrices.value.length === 0) return 0
  const sum = dailyPrices.value.reduce((acc, item) => acc + parseFloat(item.price), 0)
  return sum / dailyPrices.value.length
})

const minPrice = computed<number>(() => {
  if (dailyPrices.value.length === 0) return 0
  return Math.min(...dailyPrices.value.map(item => parseFloat(item.price)))
})

const maxPrice = computed<number>(() => {
  if (dailyPrices.value.length === 0) return 0
  return Math.max(...dailyPrices.value.map(item => parseFloat(item.price)))
})

const currentPrice = computed<number>(() => {
  if (dailyPrices.value.length === 0) return 0
  const lastItem = dailyPrices.value[dailyPrices.value.length - 1]
  return lastItem ? parseFloat(lastItem.price) : 0
})

const chartData = computed<ChartData<'line'>>(() => ({
  labels: dailyPrices.value.map(item => `${item.day} день`),
  datasets: [
    {
      label: 'Цена, Isk',
      data: dailyPrices.value.map(item => parseFloat(item.price)),
      borderColor: '#4CAF50',
      backgroundColor: 'rgba(76, 175, 80, 0.1)',
      borderWidth: 2,
      fill: true,
      tension: 0.4,
      pointBackgroundColor: '#4CAF50',
      pointBorderColor: '#ffffff',
      pointBorderWidth: 2,
      pointRadius: 4,
      pointHoverRadius: 6
    }
  ]
}))



const chartOptions = ref<ChartOptions<'line'>>({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: true,
      position: 'top'
    },
    tooltip: {
      mode: 'index',
      intersect: false,
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      titleColor: '#ffffff',
      bodyColor: '#ffffff',
      borderColor: 'rgba(255, 255, 255, 0.2)',
      borderWidth: 1,
      padding: 12,
      boxPadding: 6,
      usePointStyle: true,
      callbacks: {
        title: (tooltipItems) => {
          if (tooltipItems.length > 0) {
            const index = tooltipItems[0]!.dataIndex
            const item = dailyPrices.value[index]
            if (item && item.timestamp) {
              const date = new Date(item.timestamp)
              return `📅 ${date.toLocaleDateString('ru-RU')} 🕒 ${date.toLocaleTimeString('ru-RU', {
                hour: '2-digit',
                minute: '2-digit'
              })}`
            }
          }
          return ''
        },
        label: (context: any) => {
          return `💰 Цена: ${formatCurrency(context?.raw)}`
        },
        afterLabel: (context: any) => {
          const index = context.dataIndex
          const item = dailyPrices.value[index]
          if (item && item.timestamp) {
            const date = new Date(item.timestamp)
            return `⏱️ Время обновления: ${date.toLocaleTimeString('ru-RU', {
              hour: '2-digit',
              minute: '2-digit',
              second: '2-digit'
            })}`
          }
          return ''
        }
      }
    }
  },
  scales: {
    y: {
      beginAtZero: false,
      grid: {
        color: 'rgba(0, 0, 0, 0.05)'
      },
      ticks: {
        callback: function (value: number | string) {
          if (typeof value === 'number') {
            return value.toLocaleString('ru-RU')
          }
          return value
        }
      }
    },
    x: {
      grid: {
        color: 'rgba(0, 0, 0, 0.05)'
      }
    }
  },
  interaction: {
    intersect: false,
    mode: 'nearest'
  }
})

onMounted(async () => {
  watch(() => props.name_item, async () => {
    await generateData()
  }, { immediate: false })

  await generateData()
})
</script>

<style scoped>
.price-chart-container {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  font-family: 'Segoe UI', Arial, sans-serif;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.chart-header h3 {
  margin: 0;
  color: #333;
  font-size: 1.5rem;
  font-weight: 600;
}

.controls {
  display: flex;
  gap: 12px;
  align-items: center;
}

select {
  padding: 8px 16px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: white;
  font-size: 14px;
  cursor: pointer;
  outline: none;
  transition: border-color 0.3s;
}

select:focus {
  border-color: #4CAF50;
}

.view-toggle {
  padding: 8px 16px;
  background: #4CAF50;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s;
}

.view-toggle:hover {
  background: #45a049;
}

.chart-wrapper {
  height: 400px;
  margin-bottom: 24px;
}

.chart-info {
  border-top: 1px solid #eee;
  padding-top: 20px;
}

.stats {
  display: flex;
  justify-content: space-around;
  flex-wrap: wrap;
  gap: 20px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 20px;
  background: #f8f9fa;
  border-radius: 8px;
  min-width: 150px;
}

.stat-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
  color: #333;
}

.min-price {
  color: #f44336;
}

.max-price {
  color: #4CAF50;
}

@media (max-width: 768px) {
  .chart-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }

  .controls {
    width: 100%;
    flex-direction: column;
  }

  select,
  .view-toggle {
    width: 100%;
  }

  .stats {
    flex-direction: column;
    align-items: center;
  }
}
</style>