import { createWebHistory, createRouter } from 'vue-router'

import HomeView from './index.vue'
import AboutView from './new.vue'

const routes = [
  { path: '/new', component: HomeView },
  { path: '/', component: AboutView },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})