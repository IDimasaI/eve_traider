// composable/useHashMap.ts
import { ref, readonly } from 'vue'
import type { Items } from '../utils/API'

const namesMap = ref<Record<string, number>>({})

export const useItemsHashMap = () => {
  const initializeMap = (items: Items) => {
    const stored = localStorage.getItem("HashMap")
    if (stored) {
      namesMap.value = JSON.parse(stored)
    } else if (items) {
      const map: Record<string, number> = {}
      items.forEach((item: { name: string | number; id: number }) => {
        if (item?.name) map[item.name] = item.id
      })
      namesMap.value = map
      localStorage.setItem("HashMap", JSON.stringify(map))
    }
  }

  const getById = (name: string): number | undefined => {
    return namesMap.value[name]
  }

  // Опционально: сброс кеша при изменении данных
  const updateMap = (newMap: Record<string, number>) => {
    namesMap.value = newMap
    localStorage.setItem("HashMap", JSON.stringify(newMap))
  }

  return {
    map: readonly(namesMap), // только для чтения
    getById,
    initializeMap,
    updateMap
  }
}