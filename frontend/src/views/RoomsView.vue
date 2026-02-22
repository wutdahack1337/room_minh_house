<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { getRooms, createRoom } from '../services/api'

interface Room {
  room_number: number
  price: number
  is_available: boolean
}

const rooms = ref<Room[]>([])
const loading = ref(false)
const error = ref('')

const showModal = ref(false)
const submitting = ref(false)
const formError = ref('')
const form = reactive({ room_number: '' as number | '', price: '' as number | '', is_available: false })

const resetForm = () => { form.room_number = ''; form.price = ''; form.is_available = false; formError.value = '' }

const fetchRooms = async () => {
  loading.value = true
  try {
    rooms.value = await getRooms()
    error.value = ''
  } catch (e: any) {
    error.value = e.message || 'Không thể tải danh sách phòng'
  } finally {
    loading.value = false
  }
}

onMounted(fetchRooms)

const openModal = () => { resetForm(); showModal.value = true }
const closeModal = () => { showModal.value = false }

const handleSubmit = async () => {
  if (!form.room_number) { formError.value = 'Vui lòng nhập số phòng'; return }
  submitting.value = true
  formError.value = ''
  try {
    await createRoom({ room_number: Number(form.room_number), price: Number(form.price) || 0, is_available: form.is_available })
    closeModal()
    await fetchRooms()
  } catch (e: any) {
    formError.value = e.response?.data?.message || e.message || 'Không thể tạo phòng'
  } finally {
    submitting.value = false
  }
}

const formatPrice = (price: number) =>
  new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(price)
</script>

<template>
  <div>
    <div class="mb-6 flex items-start justify-between">
      <button @click="openModal" class="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/></svg>
        Thêm phòng
      </button>
    </div>

    <!-- Modal thêm phòng -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/40" @click="closeModal"></div>
      <div class="relative bg-white rounded-xl shadow-xl w-full max-w-md mx-4 p-6">
        <h3 class="text-lg font-semibold text-gray-800 mb-5">Thêm phòng mới</h3>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Số phòng <span class="text-red-500">*</span></label>
            <input v-model="form.room_number" type="number" min="1" placeholder="1" required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Giá thuê (VNĐ)</label>
            <input v-model="form.price" type="number" min="0" placeholder="999999"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500" />
          </div>
          <div class="flex items-center gap-2">
            <input v-model="form.is_available" type="checkbox" id="is_available" class="w-4 h-4 rounded border-gray-300 text-blue-600" />
            <label for="is_available" class="text-sm font-medium text-gray-700">Phòng trống (chưa có người thuê)</label>
          </div>
          <p v-if="formError" class="text-xs text-red-500">{{ formError }}</p>
          <div class="flex gap-3 pt-2">
            <button type="button" @click="closeModal" class="flex-1 px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors">Hủy</button>
            <button type="submit" :disabled="submitting" class="flex-1 px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors">
              {{ submitting ? 'Đang lưu...' : 'Thêm phòng' }}
            </button>
          </div>
        </form>
      </div>
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
            <td class="px-6 py-4 text-left font-medium text-gray-800">{{ room.room_number }}</td>
            <td class="px-6 py-4 text-left text-gray-700">{{ formatPrice(room.price) }}</td>
            <td class="px-6 py-4 text-left">
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
