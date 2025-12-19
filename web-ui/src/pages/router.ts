import { createWebHistory, createRouter } from 'vue-router'

import HomeView from './index.vue'
import NewIndex from './new.vue'

const routes = [
  { path: '/old', component: HomeView },
  { path: '/', component: NewIndex },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})