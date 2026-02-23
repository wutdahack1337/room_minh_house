<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { listUsers, createUser } from '../services/api'

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

const showModal = ref(false)
const submitting = ref(false)
const formError = ref('')
const form = reactive({
  pid: '', phone_number: '', role: 'tenant',
  id_card: { name: '', sex: 'Nam', date_of_birth: '', place_of_birth: '' },
  account: { username: '', password: '' },
})

const resetForm = () => {
  form.pid = ''; form.phone_number = ''; form.role = 'tenant'
  form.id_card = { name: '', sex: 'Nam', date_of_birth: '', place_of_birth: '' }
  form.account = { username: '', password: '' }
  formError.value = ''
}

const fetchUsers = async () => {
  loading.value = true
  try {
    users.value = await listUsers()
  } catch (e: any) {
    error.value = e.message || 'Không thể tải danh sách người thuê'
  } finally {
    loading.value = false
  }
}

onMounted(fetchUsers)

const openModal = () => { resetForm(); showModal.value = true }
const closeModal = () => { showModal.value = false }

const handleSubmit = async () => {
  submitting.value = true; formError.value = ''
  try {
    await createUser({
      pid: form.pid, phone_number: form.phone_number, role: form.role,
      id_card: { ...form.id_card },
      account: { ...form.account },
    })
    closeModal()
    await fetchUsers()
  } catch (e: any) {
    formError.value = e.response?.data?.message || e.message || 'Không thể tạo người thuê'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div>
    <div class="mb-6 flex items-start justify-between">
      <button @click="openModal" class="flex items-center gap-2 px-4 py-2 bg-purple-600 text-white text-sm font-medium rounded-lg hover:bg-purple-700 transition-colors">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/></svg>
        Thêm người thuê
      </button>
    </div>

    <!-- Modal thêm người thuê -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/40" @click="closeModal"></div>
      <div class="relative bg-white rounded-xl shadow-xl w-full max-w-lg mx-4 p-6 max-h-[90vh] overflow-y-auto">
        <h3 class="text-lg font-semibold text-gray-800 mb-5">Thêm người thuê mới</h3>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">CCCD <span class="text-red-500">*</span></label>
              <input v-model="form.pid" type="text" placeholder="012345678910" required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Số điện thoại</label>
              <input v-model="form.phone_number" type="text" placeholder="0901234567"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Vai trò <span class="text-red-500">*</span></label>
            <select v-model="form.role" required class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500">
              <option value="tenant">Người thuê</option>
              <option value="admin">Quản lý</option>
            </select>
          </div>
          <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide pt-1">Thông tin CCCD</p>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Họ và tên <span class="text-red-500">*</span></label>
            <input v-model="form.id_card.name" type="text" placeholder="Nguyễn Văn A" required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Giới tính <span class="text-red-500">*</span></label>
              <select v-model="form.id_card.sex" required class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500">
                <option value="Nam">Nam</option>
                <option value="Nữ">Nữ</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Ngày sinh <span class="text-red-500">*</span></label>
              <input v-model="form.id_card.date_of_birth" type="date" required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Quê quán <span class="text-red-500">*</span></label>
            <input v-model="form.id_card.place_of_birth" type="text" placeholder="Cần Thơ" required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
          </div>
          <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide pt-1">Tài khoản đăng nhập</p>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Tên đăng nhập <span class="text-red-500">*</span></label>
              <input v-model="form.account.username" type="text" placeholder="username" required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Mật khẩu <span class="text-red-500">*</span></label>
              <input v-model="form.account.password" placeholder="password" required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500" />
            </div>
          </div>
          <p v-if="formError" class="text-xs text-red-500">{{ formError }}</p>
          <div class="flex gap-3 pt-2">
            <button type="button" @click="closeModal" class="flex-1 px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors">Hủy</button>
            <button type="submit" :disabled="submitting" class="flex-1 px-4 py-2 bg-purple-600 text-white text-sm font-medium rounded-lg hover:bg-purple-700 disabled:opacity-50 transition-colors">
              {{ submitting ? 'Đang lưu...' : 'Thêm người thuê' }}
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
            <td class="px-6 py-4 text-left font-medium text-gray-800">{{ user.id_card?.name || '—' }}</td>
            <td class="px-6 py-4 text-left text-gray-600 font-mono text-xs">{{ user.pid }}</td>
            <td class="px-6 py-4 text-left text-gray-700">{{ user.phone_number || '—' }}</td>
            <td class="px-6 py-4 text-left">
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
