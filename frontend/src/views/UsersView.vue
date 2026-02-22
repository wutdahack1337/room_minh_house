<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUsers } from '../services/api'

interface IdCard {
  name?: string
  sex?: string
  date_of_birth?: string
  place_of_birth?: string
}

interface User {
  pid: string
  phone_number?: string
  role?: string
  id_card?: IdCard
}

const users = ref<User[]>([])
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  loading.value = true
  try {
    users.value = await getUsers()
  } catch (e: any) {
    error.value = e.message || 'Không thể tải danh sách người thuê'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <div class="mb-6">
      <h2 class="text-2xl font-bold text-gray-800">Người thuê</h2>
      <p class="text-sm text-gray-500 mt-1">Danh sách tất cả người thuê trong hệ thống</p>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <div v-if="loading" class="p-8 text-center text-gray-500 text-sm">Đang tải...</div>
      <div v-else-if="error" class="p-8 text-center text-red-500 text-sm">{{ error }}</div>
      <table v-else class="w-full text-sm">
        <thead>
          <tr class="bg-gray-50 border-b border-gray-200">
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Họ tên</th>
            <th class="px-6 py-3 text-left font-semibold text-gray-600">CCCD</th>
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Số điện thoại</th>
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Vai trò</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="users.length === 0">
            <td colspan="4" class="px-6 py-8 text-center text-gray-400">Chưa có người thuê nào</td>
          </tr>
          <tr
            v-for="user in users"
            :key="user.pid"
            class="border-b border-gray-100 hover:bg-gray-50 transition-colors"
          >
            <td class="px-6 py-4 font-medium text-gray-800">{{ user.id_card?.name || '—' }}</td>
            <td class="px-6 py-4 text-gray-600 font-mono text-xs">{{ user.pid }}</td>
            <td class="px-6 py-4 text-gray-700">{{ user.phone_number || '—' }}</td>
            <td class="px-6 py-4">
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-700">
                {{ user.role || 'tenant' }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
