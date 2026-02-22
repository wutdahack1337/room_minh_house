<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getRooms } from '../services/api'

interface Room {
  room_number: number
  price: number
  is_available: boolean
}

const rooms = ref<Room[]>([])
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  loading.value = true
  try {
    rooms.value = await getRooms()
  } catch (e: any) {
    error.value = e.message || 'Không thể tải danh sách phòng'
  } finally {
    loading.value = false
  }
})

const formatPrice = (price: number) =>
  new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(price)
</script>

<template>
  <div>
    <div class="mb-6">
      <h2 class="text-2xl font-bold text-gray-800">Phòng</h2>
      <p class="text-sm text-gray-500 mt-1">Danh sách tất cả các phòng trong hệ thống</p>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <div v-if="loading" class="p-8 text-center text-gray-500 text-sm">Đang tải...</div>
      <div v-else-if="error" class="p-8 text-center text-red-500 text-sm">{{ error }}</div>
      <table v-else class="w-full text-sm">
        <thead>
          <tr class="bg-gray-50 border-b border-gray-200">
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Số phòng</th>
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Giá thuê</th>
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Trạng thái</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="rooms.length === 0">
            <td colspan="3" class="px-6 py-8 text-center text-gray-400">Chưa có phòng nào</td>
          </tr>
          <tr
            v-for="room in rooms"
            :key="room.room_number"
            class="border-b border-gray-100 hover:bg-gray-50 transition-colors"
          >
            <td class="px-6 py-4 font-medium text-gray-800">P{{ room.room_number }}</td>
            <td class="px-6 py-4 text-gray-700">{{ formatPrice(room.price) }}</td>
            <td class="px-6 py-4">
              <span
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                :class="room.is_available
                  ? 'bg-green-100 text-green-700'
                  : 'bg-blue-100 text-blue-700'"
              >
                {{ room.is_available ? 'Trống' : 'Đang thuê' }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
