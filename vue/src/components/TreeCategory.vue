<template>
    <div v-for="(items, categoryName) in json" :key="getFullPath(categoryName)">
        <!-- Всегда показываем как категорию, если это объект или массив -->
        <div class="category">
            <!--Толькое если есть элементы в категории-->
            <button @click="toggleCategory(getFullPath(categoryName))" class="category-button">
                <Arrow_open v-if="isExpanded(getFullPath(categoryName))" />
                <Arrow_close v-else />
                <span class="ml-2"> {{ categoryName }}</span>
            </button>

            <!-- Содержимое категории -->
            <div v-if="isExpanded(getFullPath(categoryName))" class="category-content pl-3 border-l">
                <!-- Если это объект (не массив), рекурсивно вызываем для подкатегорий -->
                <div v-if="isObject(items) && !Array.isArray(items)">
                    <TreeCategory :json="items" :parent-path="getFullPath(categoryName)" :expanded-paths="expandedPaths"  
                        @toggle="handleToggle" @selectItem="selectItem" />
                </div>
                <!-- Если это массив, показываем элементы -->
                <div v-else-if="Array.isArray(items)" class="items">
                    <div v-for="(item, index) in items" :key="index" class="item">
                        <!-- Проверяем, является ли элемент объектом или примитивом -->
                        <div v-if="typeof item === 'object' && item !== null" class="nested-object">
                            <TreeCategory :json="{ [`обьект в массиве. ОШИБКА СИНТАКСИСА`]: item }"  
                                :parent-path="getFullPath(categoryName)" :expanded-paths="expandedPaths"
                                @toggle="handleToggle" @selectItem="selectItem"/>
                        </div>
                        <div v-else>
                            <button class="item cursor-pointer hover:bg-gray-100" @click="selectItem(item)">{{ item
                            }}</button>
                            <Fovorite :name_item="item"  />
                        </div>
                    </div>
                </div>
                <!-- Если это примитив (строка, число и т.д.) -->
                <div v-else class="items">
                    <button class="item" @click="selectItem(items)">{{ items }}</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import Arrow_open from '../assets/svg/Arrow_open.vue';
import Arrow_close from '../assets/svg/Arrow_close.vue';
import Fovorite from './Favourite.vue';
const props = defineProps({
    json: {
        type: Object,
        required: true
    },
    parentPath: {
        type: String,
        default: ''
    },
    expandedPaths: {
        type: Set<string>,
        required: true
    }, 
});
 
const emit = defineEmits(['toggle', 'selectItem']);

// Упрощенная проверка: объект ли это
const isObject = (items: any): boolean => {
    return items && typeof items === 'object';
};

const getFullPath = (categoryName: string): string => {
    return props.parentPath ? `${props.parentPath}.${categoryName}` : categoryName;
};

const isExpanded = (path: string): boolean => {
    return props.expandedPaths.has(path);
};

const toggleCategory = (path: string) => {
    emit('toggle', path);
};

const handleToggle = (path: string) => {
    emit('toggle', path);
};

const selectItem = (item: string) => {
    emit('selectItem', item);
}
 

</script>

<style scoped>
.category {
    margin-bottom: 4px;
}

.category-button {
    display: flex;
    align-items: center;
    padding: 6px 8px;
    background: #f3f4f6;
    border: 1px solid #e5e7eb;
    border-radius: 4px;
    cursor: pointer;
    width: 100%;
    text-align: left;
    font-weight: 500;
}

.category-button:hover {
    background: #e5e7eb;
}

.category-content {
    margin-top: 4px;
}

.items {
    
    background: #f9fafb;
    border-radius: 4px;
    
}

.item {
 
    border-bottom: 1px solid #e5e7eb;
}

.item:last-child {
    border-bottom: none;
}

.nested-object {
    margin-left: 10px;
}
</style>