<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { listContracts, createContract, updateContract, deleteContract } from '../services/api'
import ConfirmDialog from '../components/ConfirmDialog.vue'

interface Contract {
  contract_id: string
  user_pid: string
  room_number: number
  start_date?: string
  end_date?: string
  duration_months?: number
  status?: string
}

const contracts = ref<Contract[]>([])
const loading = ref(false)
const error = ref('')

const showModal = ref(false)
const submitting = ref(false)
const formError = ref('')
const form = reactive({ user_pid: '', room_number: '' as number | '', start_date: '', duration_months: '' as number | '' })

const resetForm = () => { form.user_pid = ''; form.room_number = ''; form.start_date = ''; form.duration_months = ''; formError.value = '' }

const fetchContracts = async () => {
  loading.value = true
  try {
    contracts.value = await listContracts()
    error.value = ''
  } catch (e: any) {
    error.value = e.message || 'Không thể tải danh sách hợp đồng'
  } finally {
    loading.value = false
  }
}

onMounted(fetchContracts)

const openModal = () => { resetForm(); showModal.value = true }
const closeModal = () => { showModal.value = false }

const contractIsActive = (status?: string) => {
  const normalized = (status || '').trim().toLowerCase()
  return normalized !== 'expired' && normalized !== 'terminated'
}

const hasAnotherActiveContractInRoom = (roomNumber: number, excludeContractId?: string) => {
  return contracts.value.some((c) => {
    if (excludeContractId && c.contract_id === excludeContractId) return false
    return c.room_number === roomNumber && contractIsActive(c.status)
  })
}

const handleSubmit = async () => {
  if (!form.user_pid || !form.room_number || !form.start_date || !form.duration_months) {
    formError.value = 'Vui lòng điền đầy đủ thông tin'; return
  }
  const roomNumber = Number(form.room_number)
  if (hasAnotherActiveContractInRoom(roomNumber)) {
    formError.value = `Phòng ${roomNumber} đã có người thuê`;
    return
  }
  submitting.value = true; formError.value = ''
  try {
    await createContract({
      user_pid: form.user_pid,
      room_number: roomNumber,
      start_date: `${form.start_date}T00:00:00Z`,
      duration_months: Number(form.duration_months),
    })
    closeModal()
    await fetchContracts()
  } catch (e: any) {
    formError.value = e.response?.data?.message || e.message || 'Không thể tạo hợp đồng'
  } finally {
    submitting.value = false
  }
}

// ─── Edit contract ────────────────────────────────────────────────────────────
const showEditModal = ref(false)
const editSubmitting = ref(false)
const editError = ref('')
const editForm = reactive({
  contract_id: '',
  user_pid: '',
  room_number: 0,
  start_date: '',
  duration_months: 0,
  status: 'active',
})

const openEditModal = (c: Contract) => {
  editForm.contract_id = c.contract_id
  editForm.user_pid = c.user_pid
  editForm.room_number = c.room_number
  // strip time part for date input
  editForm.start_date = c.start_date ? c.start_date.substring(0, 10) : ''
  editForm.duration_months = c.duration_months || 1
  editForm.status = c.status || 'active'
  editError.value = ''
  showEditModal.value = true
}
const closeEditModal = () => { showEditModal.value = false }

const handleEditSubmit = async () => {
  if (!editForm.user_pid || !editForm.room_number || !editForm.start_date || !editForm.duration_months) {
    editError.value = 'Vui lòng điền đầy đủ thông tin'; return
  }
  const roomNumber = Number(editForm.room_number)
  if (contractIsActive(editForm.status) && hasAnotherActiveContractInRoom(roomNumber, editForm.contract_id)) {
    editError.value = `Phòng ${roomNumber} đã có người thuê`;
    return
  }
  editSubmitting.value = true; editError.value = ''
  try {
    await updateContract({
      contract_id: editForm.contract_id,
      user_pid: editForm.user_pid,
      room_number: roomNumber,
      start_date: `${editForm.start_date}T00:00:00Z`,
      duration_months: Number(editForm.duration_months),
      status: editForm.status,
    })
    closeEditModal()
    await fetchContracts()
  } catch (e: any) {
    editError.value = e.response?.data?.message || e.message || 'Không thể cập nhật hợp đồng'
  } finally {
    editSubmitting.value = false
  }
}

// ─── Delete contract ──────────────────────────────────────────────────────────
const showDeleteConfirm = ref(false)
const deletingContract = ref<Contract | null>(null)
const deleteSubmitting = ref(false)

const openDeleteConfirm = (c: Contract) => { deletingContract.value = c; showDeleteConfirm.value = true }
const closeDeleteConfirm = () => { showDeleteConfirm.value = false; deletingContract.value = null }

const handleDelete = async () => {
  if (!deletingContract.value) return
  deleteSubmitting.value = true
  try {
    await deleteContract(deletingContract.value.contract_id)
    closeDeleteConfirm()
    await fetchContracts()
  } catch (e: any) {
    error.value = e.response?.data?.message || e.message || 'Không thể xóa hợp đồng'
    closeDeleteConfirm()
  } finally {
    deleteSubmitting.value = false
  }
}

const statusLabel = (status?: string) => {
  switch (status) {
    case 'active': return 'Đang hiệu lực'
    case 'expired': return 'Hết hạn'
    case 'terminated': return 'Đã chấm dứt'
    default: return status || '—'
  }
}

const statusClass = (status?: string) => {
  switch (status) {
    case 'active': return 'bg-green-100 text-green-700'
    case 'expired': return 'bg-yellow-100 text-yellow-700'
    case 'terminated': return 'bg-red-100 text-red-700'
    default: return 'bg-gray-100 text-gray-600'
  }
}
</script>

<template>
  <div>
    <div class="mb-6 flex items-start justify-between">
      <button @click="openModal" class="flex items-center gap-2 px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 transition-colors">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/></svg>
        Thêm hợp đồng
      </button>
    </div>

    <!-- Modal thêm hợp đồng -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/40" @click="closeModal"></div>
      <div class="relative bg-white rounded-xl shadow-xl w-full max-w-md mx-4 p-6">
        <h3 class="text-lg font-semibold text-gray-800 mb-5">Thêm hợp đồng mới</h3>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">CCCD người thuê <span class="text-red-500">*</span></label>
            <input v-model="form.user_pid" type="text" placeholder="012345678910" required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-green-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Số phòng <span class="text-red-500">*</span></label>
            <input v-model="form.room_number" type="number" min="1" placeholder="1" required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-green-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Ngày bắt đầu <span class="text-red-500">*</span></label>
            <input v-model="form.start_date" type="date" required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-green-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Thời hạn (tháng) <span class="text-red-500">*</span></label>
            <input v-model="form.duration_months" type="number" min="1" placeholder="3" required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-green-500" />
          </div>
          <p v-if="formError" class="text-xs text-red-500">{{ formError }}</p>
          <div class="flex gap-3 pt-2">
            <button type="button" @click="closeModal" class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-lg hover:bg-gray-50 transition-colors">Hủy</button>
            <button type="submit" :disabled="submitting" class="flex-1 px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 disabled:opacity-50 transition-colors">
              {{ submitting ? 'Đang lưu...' : 'Thêm hợp đồng' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal sửa hợp đồng -->
    <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/40" @click="closeEditModal"></div>
      <div class="relative bg-white rounded-xl shadow-xl w-full max-w-md mx-4 p-6">
        <h3 class="text-lg font-semibold text-gray-800 mb-5">Sửa hợp đồng</h3>
        <p class="text-xs text-gray-400 mb-4 font-mono">{{ editForm.contract_id }}</p>
        <form @submit.prevent="handleEditSubmit" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">CCCD người thuê <span class="text-red-500">*</span></label>
            <input v-model="editForm.user_pid" type="text" required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-green-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Số phòng <span class="text-red-500">*</span></label>
            <input v-model="editForm.room_number" type="number" min="1" required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-green-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Ngày bắt đầu <span class="text-red-500">*</span></label>
            <input v-model="editForm.start_date" type="date" required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-green-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Thời hạn (tháng) <span class="text-red-500">*</span></label>
            <input v-model="editForm.duration_months" type="number" min="1" required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-green-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Trạng thái</label>
            <select v-model="editForm.status" class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-green-500">
              <option value="active">Đang hiệu lực</option>
              <option value="expired">Hết hạn</option>
              <option value="terminated">Đã chấm dứt</option>
            </select>
          </div>
          <p v-if="editError" class="text-xs text-red-500">{{ editError }}</p>
          <div class="flex gap-3 pt-2">
            <button type="button" @click="closeEditModal" class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-lg hover:bg-gray-50 transition-colors">Hủy</button>
            <button type="submit" :disabled="editSubmitting" class="flex-1 px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 disabled:opacity-50 transition-colors">
              {{ editSubmitting ? 'Đang lưu...' : 'Lưu thay đổi' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Confirm xóa hợp đồng -->
    <ConfirmDialog
      :show="showDeleteConfirm"
      title="Xóa hợp đồng"
      :message="`Xóa hợp đồng này sẽ xóa toàn bộ hóa đơn liên quan. Thao tác này không thể hoàn tác.`"
      confirm-text="Xóa"
      :danger="true"
      :loading="deleteSubmitting"
      @confirm="handleDelete"
      @cancel="closeDeleteConfirm"
    />

    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <div v-if="loading" class="p-8 text-center text-gray-500 text-sm">Đang tải...</div>
      <div v-else-if="error" class="p-8 text-center text-red-500 text-sm">{{ error }}</div>
      <table v-else class="w-full text-sm">
        <thead>
          <tr class="bg-gray-50 border-b border-gray-200">
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Mã hợp đồng</th>
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Phòng</th>
            <th class="px-6 py-3 text-left font-semibold text-gray-600">CCCD</th>
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Thời hạn</th>
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Trạng thái</th>
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Thao tác</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="contracts.length === 0">
            <td colspan="6" class="px-6 py-8 text-center text-gray-400">Chưa có hợp đồng nào</td>
          </tr>
          <tr
            v-for="c in contracts"
            :key="c.contract_id"
            class="border-b border-gray-100 hover:bg-gray-50 transition-colors"
          >
            <td class="px-6 py-4 text-left font-mono text-xs text-gray-600">{{ c.contract_id }}</td>
            <td class="px-6 py-4 text-left font-medium text-gray-800">{{ c.room_number }}</td>
            <td class="px-6 py-4 text-left text-gray-600 font-mono text-xs">{{ c.user_pid }}</td>
            <td class="px-6 py-4 text-left text-gray-700">{{ c.duration_months ? `${c.duration_months} tháng` : '—' }}</td>
            <td class="px-6 py-4 text-left">
              <span
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                :class="statusClass(c.status)"
              >
                {{ statusLabel(c.status) }}
              </span>
            </td>
            <td class="px-6 py-4 text-left">
              <div class="flex items-center gap-2">
                <button @click="openEditModal(c)"
                  class="px-2.5 py-1 text-xs font-medium text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-md transition-colors">
                  Sửa
                </button>
                <button @click="openDeleteConfirm(c)"
                  class="px-2.5 py-1 text-xs font-medium text-red-600 bg-red-50 hover:bg-red-100 rounded-md transition-colors">
                  Xóa
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
