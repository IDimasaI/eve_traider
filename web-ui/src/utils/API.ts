interface ItemIdMap {
  name: string;
  id: number;
}

interface EveIdsResponse {
  inventory_types: Array<{ id: number }>;
}
export async function find_id(name: string): Promise<number | null> {
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

export type MarketData = {
    loading: boolean;
    market: string;
    data: any
}

export type Items = {
    name: string;
    category: string;
    id: number;
}[];

export type HashMapItems= {
    [key: string]: Items
}

export async function get_all_items() {
    const res = await fetch("/api/v2/get_all_items");
    return (await res.json()) as Items;
}