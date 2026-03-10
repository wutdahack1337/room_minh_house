<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { listRooms, createRoom, updateRoom, deleteRoom, listMeterReadings, createMeterReading, updateMeterReading } from '../services/api'
import ConfirmDialog from '../components/ConfirmDialog.vue'

interface Room {
  room_number: number
  price: number
  is_available: boolean
}

interface MeterReading {
  room_number: number
  period: string
  electricity_reading: number
  water_reading: number
}

// ─── Rooms ───────────────────────────────────────────────────────────────────
const rooms = ref<Room[]>([])
const loading = ref(false)
const error = ref('')

const fetchRooms = async () => {
  loading.value = true
  try {
    rooms.value = await listRooms()
    error.value = ''
  } catch (e: any) {
    error.value = e.message || 'Không thể tải danh sách phòng'
  } finally {
    loading.value = false
  }
}

onMounted(fetchRooms)

// ─── Modal thêm phòng ─────────────────────────────────────────────────────────
const showModal = ref(false)
const submitting = ref(false)
const formError = ref('')
const form = reactive({ room_number: '' as number | '', price: '' as number | '', is_available: false })

const resetForm = () => { form.room_number = ''; form.price = ''; form.is_available = false; formError.value = '' }
const openModal = () => { resetForm(); showModal.value = true }
const closeModal = () => { showModal.value = false }

const handleSubmit = async () => {
  if (!form.room_number) { formError.value = 'Vui lòng nhập số phòng'; return }
  submitting.value = true; formError.value = ''
  try {
    await createRoom({ room_number: Number(form.room_number), price: Number(form.price) || 0, is_available: form.is_available })
    closeModal(); await fetchRooms()
  } catch (e: any) {
    formError.value = e.response?.data?.message || e.message || 'Không thể tạo phòng'
  } finally {
    submitting.value = false
  }
}

// ─── Modal sửa phòng ──────────────────────────────────────────────────────────
const showEditModal = ref(false)
const editSubmitting = ref(false)
const editError = ref('')
const editForm = reactive({ room_number: 0, price: 0, is_available: false })

const openEditModal = (room: Room) => {
  editForm.room_number = room.room_number
  editForm.price = room.price
  editForm.is_available = room.is_available
  editError.value = ''
  showEditModal.value = true
}
const closeEditModal = () => { showEditModal.value = false }

const handleEditSubmit = async () => {
  editSubmitting.value = true; editError.value = ''
  try {
    await updateRoom({ room_number: editForm.room_number, price: Number(editForm.price), is_available: editForm.is_available })
    closeEditModal(); await fetchRooms()
  } catch (e: any) {
    editError.value = e.response?.data?.message || e.message || 'Không thể cập nhật phòng'
  } finally {
    editSubmitting.value = false
  }
}

// ─── Xóa phòng ────────────────────────────────────────────────────────────────
const showDeleteConfirm = ref(false)
const deletingRoom = ref<Room | null>(null)
const deleteSubmitting = ref(false)

const openDeleteConfirm = (room: Room) => { deletingRoom.value = room; showDeleteConfirm.value = true }
const closeDeleteConfirm = () => { showDeleteConfirm.value = false; deletingRoom.value = null }

const handleDelete = async () => {
  if (!deletingRoom.value) return
  deleteSubmitting.value = true
  try {
    await deleteRoom(deletingRoom.value.room_number)
    closeDeleteConfirm(); await fetchRooms()
  } catch (e: any) {
    error.value = e.response?.data?.message || e.message || 'Không thể xóa phòng'
    closeDeleteConfirm()
  } finally {
    deleteSubmitting.value = false
  }
}

// ─── Modal chỉ số điện nước ───────────────────────────────────────────────────
const showMeterModal = ref(false)
const selectedRoom = ref<number | null>(null)
const meterReadings = ref<MeterReading[]>([])
const meterLoading = ref(false)
const meterError = ref('')

const fetchMeterReadings = async (roomNumber: number) => {
  meterLoading.value = true; meterError.value = ''
  try {
    meterReadings.value = await listMeterReadings(roomNumber)
  } catch (e: any) {
    meterError.value = e.response?.data?.message || e.message || 'Không thể tải chỉ số'
  } finally {
    meterLoading.value = false
  }
}

const openMeterModal = async (room: Room) => {
  selectedRoom.value = room.room_number
  showMeterModal.value = true
  showAddForm.value = false
  editingMeter.value = null
  await fetchMeterReadings(room.room_number)
}
const closeMeterModal = () => { showMeterModal.value = false; selectedRoom.value = null }

// ── Sub-form: thêm chỉ số mới ────────────────────────────────────────────────
const showAddForm = ref(false)
const addSubmitting = ref(false)
const addError = ref('')
const addForm = reactive({ period: '', electricity_reading: '' as number | '', water_reading: '' as number | '' })
const PERIOD_REGEX = /^\d{2}-\d{4}$/

const resetAddForm = () => { addForm.period = ''; addForm.electricity_reading = ''; addForm.water_reading = ''; addError.value = '' }
const openAddForm = () => { resetAddForm(); showAddForm.value = true; editingMeter.value = null }
const cancelAddForm = () => { showAddForm.value = false }

const handleAddMeterReading = async () => {
  if (!PERIOD_REGEX.test(addForm.period)) { addError.value = 'Kỳ không hợp lệ, định dạng MM-YYYY (vd: 01-2026)'; return }
  if (addForm.electricity_reading === '' || addForm.water_reading === '') { addError.value = 'Vui lòng nhập đầy đủ chỉ số'; return }
  addSubmitting.value = true; addError.value = ''
  try {
    await createMeterReading({
      room_number: selectedRoom.value!,
      period: addForm.period,
      electricity_reading: Number(addForm.electricity_reading),
      water_reading: Number(addForm.water_reading),
    })
    showAddForm.value = false
    await fetchMeterReadings(selectedRoom.value!)
  } catch (e: any) {
    addError.value = e.response?.data?.message || e.message || 'Không thể thêm chỉ số'
  } finally {
    addSubmitting.value = false
  }
}

// ── Edit meter reading ────────────────────────────────────────────────────────
const editingMeter = ref<MeterReading | null>(null)
const editMeterForm = reactive({ electricity_reading: 0, water_reading: 0 })
const editMeterSubmitting = ref(false)
const editMeterError = ref('')

const openEditMeter = (r: MeterReading) => {
  editingMeter.value = r
  editMeterForm.electricity_reading = r.electricity_reading
  editMeterForm.water_reading = r.water_reading
  editMeterError.value = ''
  showAddForm.value = false
}
const cancelEditMeter = () => { editingMeter.value = null }

const handleEditMeterReading = async () => {
  if (!editingMeter.value) return
  editMeterSubmitting.value = true; editMeterError.value = ''
  try {
    await updateMeterReading({
      room_number: selectedRoom.value!,
      period: editingMeter.value.period,
      electricity_reading: Number(editMeterForm.electricity_reading),
      water_reading: Number(editMeterForm.water_reading),
    })
    editingMeter.value = null
    await fetchMeterReadings(selectedRoom.value!)
  } catch (e: any) {
    editMeterError.value = e.response?.data?.message || e.message || 'Không thể cập nhật chỉ số'
  } finally {
    editMeterSubmitting.value = false
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
            <button type="button" @click="closeModal" class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-lg hover:bg-gray-50 transition-colors">Hủy</button>
            <button type="submit" :disabled="submitting" class="flex-1 px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors">
              {{ submitting ? 'Đang lưu...' : 'Thêm phòng' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal sửa phòng -->
    <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/40" @click="closeEditModal"></div>
      <div class="relative bg-white rounded-xl shadow-xl w-full max-w-md mx-4 p-6">
        <h3 class="text-lg font-semibold text-gray-800 mb-5">Sửa phòng {{ editForm.room_number }}</h3>
        <form @submit.prevent="handleEditSubmit" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Giá thuê (VNĐ)</label>
            <input v-model="editForm.price" type="number" min="0"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500" />
          </div>
          <div class="flex items-center gap-2">
            <input v-model="editForm.is_available" type="checkbox" id="edit_is_available" class="w-4 h-4 rounded border-gray-300 text-blue-600" />
            <label for="edit_is_available" class="text-sm font-medium text-gray-700">Phòng trống</label>
          </div>
          <p v-if="editError" class="text-xs text-red-500">{{ editError }}</p>
          <div class="flex gap-3 pt-2">
            <button type="button" @click="closeEditModal" class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-lg hover:bg-gray-50 transition-colors">Hủy</button>
            <button type="submit" :disabled="editSubmitting" class="flex-1 px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors">
              {{ editSubmitting ? 'Đang lưu...' : 'Lưu thay đổi' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal chỉ số điện nước -->
    <div v-if="showMeterModal" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/40" @click="closeMeterModal"></div>
      <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-2xl mx-4 overflow-hidden">
        <!-- Header -->
        <div class="flex items-center justify-between px-6 py-4 border-b border-gray-200 bg-gray-50">
          <h3 class="text-base font-semibold text-gray-800">
            Chỉ số điện nước — Phòng {{ selectedRoom }}
          </h3>
          <button @click="closeMeterModal" class="p-1 rounded-lg hover:bg-gray-200 text-gray-500 transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <!-- Body -->
        <div class="p-6 space-y-4 max-h-[70vh] overflow-y-auto">

          <!-- Sub-form thêm chỉ số -->
          <div v-if="showAddForm" class="bg-green-50 border border-green-200 rounded-xl p-4">
            <p class="text-sm font-semibold text-green-800 mb-3">Thêm chỉ số mới</p>
            <div class="grid grid-cols-3 gap-3">
              <div>
                <label class="block text-xs font-medium text-gray-600 mb-1">Kỳ <span class="text-red-500">*</span></label>
                <input v-model="addForm.period" type="text" placeholder="MM-YYYY"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-green-500" />
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-600 mb-1">Điện (kWh) <span class="text-red-500">*</span></label>
                <input v-model="addForm.electricity_reading" type="number" min="0" step="0.01" placeholder="0"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-green-500" />
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-600 mb-1">Nước (m³) <span class="text-red-500">*</span></label>
                <input v-model="addForm.water_reading" type="number" min="0" step="0.01" placeholder="0"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-green-500" />
              </div>
            </div>
            <p v-if="addError" class="mt-2 text-xs text-red-500">{{ addError }}</p>
            <div class="flex gap-2 mt-3">
              <button @click="handleAddMeterReading" :disabled="addSubmitting"
                class="px-4 py-1.5 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 disabled:opacity-50 transition-colors">
                {{ addSubmitting ? 'Đang lưu...' : 'Lưu' }}
              </button>
              <button @click="cancelAddForm" class="px-4 py-1.5 border border-gray-300 text-gray-700 text-sm rounded-lg hover:bg-gray-50 transition-colors">Hủy</button>
            </div>
          </div>

          <!-- Sub-form sửa chỉ số -->
          <div v-if="editingMeter" class="bg-blue-50 border border-blue-200 rounded-xl p-4">
            <p class="text-sm font-semibold text-blue-800 mb-3">Sửa chỉ số kỳ {{ editingMeter.period }}</p>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-medium text-gray-600 mb-1">Điện (kWh) <span class="text-red-500">*</span></label>
                <input v-model="editMeterForm.electricity_reading" type="number" min="0" step="0.01"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500" />
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-600 mb-1">Nước (m³) <span class="text-red-500">*</span></label>
                <input v-model="editMeterForm.water_reading" type="number" min="0" step="0.01"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500" />
              </div>
            </div>
            <p v-if="editMeterError" class="mt-2 text-xs text-red-500">{{ editMeterError }}</p>
            <div class="flex gap-2 mt-3">
              <button @click="handleEditMeterReading" :disabled="editMeterSubmitting"
                class="px-4 py-1.5 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors">
                {{ editMeterSubmitting ? 'Đang lưu...' : 'Lưu thay đổi' }}
              </button>
              <button @click="cancelEditMeter" class="px-4 py-1.5 border border-gray-300 text-gray-700 text-sm rounded-lg hover:bg-gray-50 transition-colors">Hủy</button>
            </div>
          </div>

          <!-- Nút thêm -->
          <div v-if="!showAddForm && !editingMeter" class="flex justify-end">
            <button @click="openAddForm" class="flex items-center gap-1.5 px-3 py-1.5 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 transition-colors">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/></svg>
              Thêm chỉ số
            </button>
          </div>

          <!-- Bảng chỉ số -->
          <div class="rounded-xl border border-gray-200 overflow-hidden">
            <div v-if="meterLoading" class="p-6 text-center text-gray-500 text-sm">Đang tải...</div>
            <div v-else-if="meterError" class="p-6 text-center text-red-500 text-sm">{{ meterError }}</div>
            <table v-else class="w-full text-sm">
              <thead>
                <tr class="bg-gray-50 border-b border-gray-200">
                  <th class="px-4 py-3 text-left font-semibold text-gray-600">Kỳ</th>
                  <th class="px-4 py-3 text-left font-semibold text-gray-600">Điện (kWh)</th>
                  <th class="px-4 py-3 text-left font-semibold text-gray-600">Nước (m³)</th>
                  <th class="px-4 py-3 text-right font-semibold text-gray-600">Thao tác</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="meterReadings.length === 0">
                  <td colspan="4" class="px-4 py-6 text-center text-gray-400">Chưa có chỉ số nào</td>
                </tr>
                <tr v-for="r in meterReadings" :key="r.period" class="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                  <td class="px-4 py-3 text-left font-medium text-gray-800">{{ r.period }}</td>
                  <td class="px-4 py-3 text-left text-gray-700">{{ r.electricity_reading }}</td>
                  <td class="px-4 py-3 text-left text-gray-700">{{ r.water_reading }}</td>
                  <td class="px-4 py-3 text-right">
                    <button @click="openEditMeter(r)"
                      class="px-2.5 py-1 text-xs font-medium text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-md transition-colors">
                      Sửa
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

        </div>
      </div>
    </div>

    <!-- Confirm xóa phòng -->
    <ConfirmDialog
      :show="showDeleteConfirm"
      title="Xóa phòng"
      :message="`Xóa phòng ${deletingRoom?.room_number} sẽ xóa toàn bộ chỉ số điện nước và tất cả hợp đồng, hóa đơn liên quan. Thao tác này không thể hoàn tác.`"
      confirm-text="Xóa"
      :danger="true"
      :loading="deleteSubmitting"
      @confirm="handleDelete"
      @cancel="closeDeleteConfirm"
    />

    <!-- Bảng phòng -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <div v-if="loading" class="p-8 text-center text-gray-500 text-sm">Đang tải...</div>
      <div v-else-if="error" class="p-8 text-center text-red-500 text-sm">{{ error }}</div>
      <table v-else class="w-full text-sm">
        <thead>
          <tr class="bg-gray-50 border-b border-gray-200">
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Số phòng</th>
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Giá thuê</th>
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Trạng thái</th>
            <th class="px-6 py-3 text-left font-semibold text-gray-600">Thao tác</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="rooms.length === 0">
            <td colspan="4" class="px-6 py-8 text-center text-gray-400">Chưa có phòng nào</td>
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
                :class="room.is_available ? 'bg-green-100 text-green-700' : 'bg-blue-100 text-blue-700'"
              >
                {{ room.is_available ? 'Trống' : 'Đang thuê' }}
              </span>
            </td>
            <td class="px-6 py-4 text-left">
              <div class="flex items-center gap-2">
                <button @click="openMeterModal(room)"
                  class="px-2.5 py-1 text-xs font-medium text-green-600 bg-green-50 hover:bg-green-100 rounded-md transition-colors">
                  Chỉ số
                </button>
                <button @click="openEditModal(room)"
                  class="px-2.5 py-1 text-xs font-medium text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-md transition-colors">
                  Sửa
                </button>
                <button @click="openDeleteConfirm(room)"
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
