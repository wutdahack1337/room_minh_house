import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import RoomsView from '../views/RoomsView.vue'
import UsersView from '../views/UsersView.vue'
import ContractsView from '../views/ContractsView.vue'
import InvoicesView from '../views/InvoicesView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'dashboard', component: DashboardView },
    { path: '/rooms', name: 'rooms', component: RoomsView },
    { path: '/users', name: 'users', component: UsersView },
    { path: '/contracts', name: 'contracts', component: ContractsView },
    { path: '/invoices', name: 'invoices', component: InvoicesView },
  ],
})

export default router
