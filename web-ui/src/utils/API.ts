interface ItemIdMap {
  name: string;
  id: number;
}

interface EveIdsResponse {
  inventory_types: Array<{ id: number }>;
}
export async function find_id(name: string): Promise<number | null> {
  if (!name) return 0
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


export type PriceItem = {
  item_id: number;
  timestamp: number;
  price: string;
  error: boolean;
  day: string;
}


export async function get_price(id: number | null): Promise<PriceItem[]> {
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



export type MarketData = {
  loading: boolean;
  Region: string;
  data: Array<{
    duration: number;
    is_buy_order: boolean;
    issued: string;
    location_id: number;
    min_volume: number;
    order_id: number;
    price: number;
    range: string;
    system_id: number;
    type_id: number;
    volume_remain: number;
    volume_total: number;
  }>
}
export type Order = {
  duration: number;
  is_buy_order: boolean;
  issued: string;
  location_id: number;
  min_volume: number;
  order_id: number;
  price: number;
  range: string;
  system_id: number;
  type_id: number;
  volume_remain: number;
  volume_total: number;
}
export type Items = {
  name: string;
  category: string;
  id: number;
}[];

export type HashMapItems = {
  [key: string]: Items
}

export async function get_all_items() {
  const res = await fetch("/api/v2/get_all_items");
  return (await res.json()) as Items;
}

// export async function get_max_buy_price(id: number) {
//   const res = await fetch("");
//   return (await res.json());
// }