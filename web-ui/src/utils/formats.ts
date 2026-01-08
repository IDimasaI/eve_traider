export const formatDeltaOrderVolume = (order: number, volume: number): string => {
    return `${(volume / order).toFixed(1)}`
}
 
export const formatPrice = (price: number): string => {
    if (price == null) return '0.00';

    const formattedPrice = price.toLocaleString('ru-RU', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
    });
    return formattedPrice.endsWith(',00') ?  formattedPrice.slice(0, -3) : formattedPrice

}

