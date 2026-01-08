<template>
  <div class="price-chart-container" :class="{ 'dark-theme': isDarkTheme }">

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
          <span class="stat-value">{{ formatPrice(currentPrice) }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">Средняя цена:</span>
          <span class="stat-value average-price">{{ formatPrice(averagePrice) }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">Минимальная:</span>
          <span class="stat-value min-price">{{ formatPrice(minPrice) }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">Максимальная:</span>
          <span class="stat-value max-price">{{ formatPrice(maxPrice) }}</span>
        </div>
      </div>
      <div id="more" class="text-center">
        <button class="more-button" @click="toggleMore">Показать
          больше</button>
        <div v-if="moreInfo.open" class="more-info flex flex-col">
          <div class="w-full">
            <select v-model="selectedPeriod" class="period-select">
              <option v-for="period in moreInfo.periods" :key="period" :value="period">{{ period }}</option>
            </select>
          </div>
          <section class="w-full flex flex-col  gap-1 justify-evenly text-xs">
            <div class="grid grid-cols-4 w-full table-header">
              <li>Дата</li>
              <li>Объем</li>
              <li>Заказов</li>
              <li>~Обьем заказов</li>
            </div>
            <template v-for="(_, index) in filteredData.times" :key="index">
              <div class="grid grid-cols-4 w-full table-rows">
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
import { find_id, type PriceItem, get_price } from '../utils/API'
import { formatDeltaOrderVolume, formatPrice } from '../utils/formats'
import { CurrentTheme, Themes } from '../composables/Theme'
// Регистрируем компоненты Chart.js 
Chart.register(...registerables)

type MoreData = {
  open: boolean;
  periods: string[];
  times: string[];
  volume_per_month: number[];
  order_counts: number[]
}

const props = defineProps<{
  name_item: string;
}>()

const SKIP_EVERY_X_ROW = ref<number>(1)
const dailyPrices = ref<PriceItem[]>([])
const moreInfo = ref<MoreData>({
  open: false,
  periods: [],
  times: [],
  volume_per_month: [],
  order_counts: []
})

// Добавляем состояние для темы
const isDarkTheme = computed<boolean>(() => {
  return CurrentTheme.value === Themes.dark
})

let id_item = await find_id(props.name_item)

const toggleMore = async () => {
  moreInfo.value.open = !moreInfo.value.open
  if (moreInfo.value.open) {
    if (!id_item) return
    await getMoreData(id_item, 10000002)
  }
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

    const month = (date.getMonth() + 1).toString().padStart(2, '0')
    const year = date.getFullYear()

    periodsSet.add(`${year}-${month}`)
  }

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

 
const generateData = async (): Promise<void> => {
  if(!id_item) return
  let prices = await get_price(id_item)
  let iter = 4
  if (prices[0] && prices[0].error) {
    for (let i = 0; i < iter; i++) {
      prices = await get_price(id_item)
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

// Динамические цвета для графика в зависимости от темы
const chartLineColor = computed(() => isDarkTheme.value ? '#4CAF50' : '#2E7D32')
const chartBgColor = computed(() => isDarkTheme.value ? 'rgba(76, 175, 80, 0.15)' : 'rgba(76, 175, 80, 0.1)')

const chartData = computed<ChartData<'line'>>(() => ({
  labels: dailyPrices.value.map((item, i) => {
    return i % 6 == 0 ? item.day : ``
  }),
  datasets: [
    {
      label: 'Цена, Isk',
      data: dailyPrices.value.map(item => parseFloat(item.price)),
      borderColor: chartLineColor.value,
      backgroundColor: chartBgColor.value,
      borderWidth: 2,
      fill: true,
      tension: 0.4,
      pointBackgroundColor: chartLineColor.value,
      pointBorderColor: isDarkTheme.value ? '#1e293b' : '#ffffff',
      pointBorderWidth: 2,
      pointRadius: 4,
      pointHoverRadius: 6
    }
  ]
}))

const chartOptions = computed<ChartOptions<'line'>>(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: true,
      position: 'top',
      labels: {
        color: isDarkTheme.value ? '#e2e8f0' : '#333',
        font: {
          size: 14
        }
      }
    },
    tooltip: {
      mode: 'index',
      intersect: false,
      backgroundColor: isDarkTheme.value ? 'rgba(30, 41, 59, 0.95)' : 'rgba(0, 0, 0, 0.8)',
      titleColor: isDarkTheme.value ? '#e2e8f0' : '#ffffff',
      bodyColor: isDarkTheme.value ? '#e2e8f0' : '#ffffff',
      borderColor: isDarkTheme.value ? 'rgba(148, 163, 184, 0.3)' : 'rgba(255, 255, 255, 0.2)',
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
          return `💰 ${formatPrice(context?.raw)}`
        },
      }
    }
  },
  scales: {
    y: {
      beginAtZero: false,
      grid: {
        color: isDarkTheme.value ? 'rgba(148, 163, 184, 0.15)' : 'rgba(0, 0, 0, 0.05)'
      },
      ticks: {
        color: isDarkTheme.value ? '#94a3b8' : '#666',
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
        color: isDarkTheme.value ? 'rgba(148, 163, 184, 0.15)' : 'rgba(0, 0, 0, 0.05)'
      },
      ticks: {
        color: isDarkTheme.value ? '#94a3b8' : '#666'
      }
    }
  },
  interaction: {
    intersect: false,
    mode: 'nearest'
  }
}))

onMounted(async () => {
  watch(() => props.name_item, async () => {
    id_item = await find_id(props.name_item)
    await generateData()
    if (id_item) {
      await getMoreData(id_item, 10000002)
    }
  }, { immediate: false })

  await generateData()
})
</script>