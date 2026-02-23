<script setup lang="ts">
import { ref, onMounted, reactive, computed, watch } from 'vue'
import {
  listContracts,
  listInvoices,
  createInvoice,
  listInvoiceItems,
  updateInvoiceItem,
} from '../services/api'

// ─── Types ───────────────────────────────────────────────────────────────────
interface Contract {
  contract_id: string
  user_pid: string
  room_number: number
  start_date: string
  duration_months: number
  status: string
}

interface Invoice {
  contract_id: string
  period: string
  total_amount: number
  status: string
}

interface InvoiceItem {
  contract_id: string
  period: string
  item_type: string
  quantity: number
  unit_price: number
  amount: number
}

// ─── State ────────────────────────────────────────────────────────────────────
const contracts = ref<Contract[]>([])
const selectedContractId = ref('')
const invoices = ref<Invoice[]>([])
const loading = ref(false)
const contractsLoading = ref(false)
const error = ref('')

// Create modal
const showCreateModal = ref(false)
const submitting = ref(false)
const formError = ref('')
const form = reactive({
  period: '',
  electricity_unit_price: 3000,
  water_unit_price: 10000,
})

// Detail modal (items)
const showDetailModal = ref(false)
const detailInvoice = ref<Invoice | null>(null)
const invoiceItems = ref<InvoiceItem[]>([])
const itemsLoading = ref(false)
const itemsError = ref('')
const editingItem = reactive<Record<string, number>>({})
const savingItem = ref('')

// ─── Helpers ─────────────────────────────────────────────────────────────────
const formatCurrency = (value: number) =>
  new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(value)

const itemTypeLabel: Record<string, string> = {
  room_rent: 'Tiền phòng',
  electricity: 'Điện',
  water: 'Nước',
}

const itemTypeUnit: Record<string, string> = {
  room_rent: 'tháng',
  electricity: 'kWh',
  water: 'm³',
}

const selectedContract = computed(() =>
  contracts.value.find((c) => c.contract_id === selectedContractId.value)
)

// ─── Fetch ───────────────────────────────────────────────────────────────────
const fetchContracts = async () => {
  contractsLoading.value = true
  try {
    contracts.value = await listContracts()
  } catch {
    error.value = 'Không thể tải danh sách hợp đồng.'
  } finally {
    contractsLoading.value = false
  }
}

const fetchInvoices = async () => {
  if (!selectedContractId.value) {
    invoices.value = []
    return
  }
  loading.value = true
  error.value = ''
  try {
    invoices.value = await listInvoices(selectedContractId.value)
  } catch {
    error.value = 'Không thể tải danh sách hóa đơn.'
    invoices.value = []
  } finally {
    loading.value = false
  }
}

watch(selectedContractId, fetchInvoices)

onMounted(fetchContracts)

// ─── Create Invoice ───────────────────────────────────────────────────────────
const openCreateModal = () => {
  formError.value = ''
  form.period = ''
  form.electricity_unit_price = 3000
  form.water_unit_price = 10000
  showCreateModal.value = true
}

const closeCreateModal = () => {
  showCreateModal.value = false
}

const handleCreate = async () => {
  if (!selectedContractId.value) {
    formError.value = 'Vui lòng chọn hợp đồng trước.'
    return
  }
  if (!/^\d{2}-\d{4}$/.test(form.period)) {
    formError.value = 'Kỳ hóa đơn phải có định dạng MM-YYYY (ví dụ: 01-2026).'
    return
  }
  if (form.electricity_unit_price <= 0 || form.water_unit_price <= 0) {
    formError.value = 'Đơn giá điện và nước phải lớn hơn 0.'
    return
  }
  submitting.value = true
  formError.value = ''
  try {
    await createInvoice({
      contract_id: selectedContractId.value,
      period: form.period,
      electricity_unit_price: form.electricity_unit_price,
      water_unit_price: form.water_unit_price,
    })
    closeCreateModal()
    await fetchInvoices()
  } catch (e: any) {
    formError.value = e?.response?.data?.message ?? 'Tạo hóa đơn thất bại.'
  } finally {
    submitting.value = false
  }
}

// ─── Detail / Items ────────────────────────────────────────────────────────────
const openDetail = async (invoice: Invoice) => {
  detailInvoice.value = invoice
  invoiceItems.value = []
  itemsError.value = ''
  showDetailModal.value = true
  itemsLoading.value = true
  try {
    const items = await listInvoiceItems(invoice.contract_id, invoice.period)
    invoiceItems.value = items
    items.forEach((item: InvoiceItem) => {
      editingItem[item.item_type] = item.unit_price
    })
  } catch {
    itemsError.value = 'Không thể tải mục hóa đơn.'
  } finally {
    itemsLoading.value = false
  }
}

const closeDetail = () => {
  showDetailModal.value = false
  detailInvoice.value = null
}

const saveItem = async (item: InvoiceItem) => {
  savingItem.value = item.item_type
  try {
    const newUnitPrice = editingItem[item.item_type] ?? item.unit_price
    await updateInvoiceItem({
      ...item,
      unit_price: newUnitPrice,
      amount: item.quantity * newUnitPrice,
    })
    // Refresh items and re-fetch invoices (total_amount changes)
    const updated = await listInvoiceItems(item.contract_id, item.period)
    invoiceItems.value = updated
    updated.forEach((i: InvoiceItem) => {
      editingItem[i.item_type] = i.unit_price
    })
    await fetchInvoices()
    // Update detail invoice's total_amount display
    const refreshed = invoices.value.find(
      (inv) => inv.period === item.period && inv.contract_id === item.contract_id
    )
    if (refreshed) detailInvoice.value = refreshed
  } catch {
    itemsError.value = 'Lưu thất bại.'
  } finally {
    savingItem.value = ''
  }
}
</script>

<template>
  <!-- ── Page header ─────────────────────────────────────────────────── -->
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <button
        :disabled="!selectedContractId"
        @click="openCreateModal"
        class="inline-flex items-center gap-2 px-4 py-2 bg-orange-500 hover:bg-orange-600 disabled:bg-orange-300 disabled:cursor-not-allowed text-white text-sm font-medium rounded-lg transition-colors"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
        </svg>
        Tạo hóa đơn
      </button>
    </div>

    <!-- ── Contract selector ───────────────────────────────────────────── -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-4">
      <label class="block text-sm font-medium text-gray-700 mb-2">Chọn hợp đồng</label>
      <select
        v-model="selectedContractId"
        class="w-full sm:w-96 px-3 py-2 border border-gray-300 rounded-lg text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-orange-400"
      >
        <option value="">-- Chọn hợp đồng --</option>
        <option v-for="c in contracts" :key="c.contract_id" :value="c.contract_id">
          {{ c.contract_id }} — Phòng {{ c.room_number }} · {{ c.user_pid }}
        </option>
      </select>
      <p v-if="contractsLoading" class="text-xs text-gray-400 mt-1">Đang tải hợp đồng…</p>
    </div>

    <!-- ── Invoice table ───────────────────────────────────────────────── -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
        <h3 class="text-sm font-semibold text-gray-700">
          Danh sách hóa đơn
          <span v-if="selectedContract" class="font-normal text-gray-400">
            — Phòng {{ selectedContract.room_number }}
          </span>
        </h3>
        <span class="text-xs text-gray-400">{{ invoices.length }} hóa đơn</span>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="px-5 py-10 text-center text-sm text-gray-400">Đang tải…</div>

      <!-- Error -->
      <div
        v-else-if="error"
        class="px-5 py-4 text-sm text-red-600 bg-red-50 border-b border-red-100"
      >
        {{ error }}
      </div>

      <!-- No contract selected -->
      <div
        v-else-if="!selectedContractId"
        class="px-5 py-10 text-center text-sm text-gray-400"
      >
        Vui lòng chọn hợp đồng để xem hóa đơn.
      </div>

      <!-- Empty -->
      <div
        v-else-if="invoices.length === 0"
        class="px-5 py-10 text-center text-sm text-gray-400"
      >
        Chưa có hóa đơn nào cho hợp đồng này.
      </div>

      <!-- Table -->
      <table v-else class="min-w-full divide-y divide-gray-100">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-5 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wide">
              Kỳ
            </th>
            <th class="px-5 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wide">
              Tổng tiền
            </th>
            <th class="px-5 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wide">
              Trạng thái
            </th>
            <th class="px-5 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wide">
              Thao tác
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-100 bg-white">
          <tr v-for="inv in invoices" :key="inv.period" class="hover:bg-gray-50">
            <td class="px-5 py-3.5 text-left text-sm font-medium text-gray-900">{{ inv.period }}</td>
            <td class="px-5 py-3.5 text-left text-sm text-gray-700">{{ formatCurrency(inv.total_amount) }}</td>
            <td class="px-5 py-3.5 text-left">
              <span
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                :class="inv.status === 'paid'
                  ? 'bg-green-100 text-green-700'
                  : 'bg-red-100 text-red-700'"
              >
                {{ inv.status === 'paid' ? 'Đã thanh toán' : 'Chưa thanh toán' }}
              </span>
            </td>
            <td class="px-5 py-3.5 text-right">
              <div class="flex items-center justify-end gap-2">
                <button
                  @click="openDetail(inv)"
                  class="px-2.5 py-1 text-xs font-medium text-orange-600 bg-orange-50 hover:bg-orange-100 rounded-md transition-colors"
                >
                  Xem
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>

  <!-- ══ Create Invoice Modal ════════════════════════════════════════════ -->
  <div
    v-if="showCreateModal"
    class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 p-4"
    @click.self="closeCreateModal"
  >
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-md">
      <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
        <h3 class="text-base font-semibold text-gray-900">Tạo hóa đơn mới</h3>
        <button @click="closeCreateModal" class="text-gray-400 hover:text-gray-600">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form @submit.prevent="handleCreate" class="px-6 py-5 space-y-4">
        <!-- Contract (readonly) -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Hợp đồng</label>
          <input
            :value="selectedContractId"
            readonly
            class="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm text-gray-500 bg-gray-50"
          />
        </div>

        <!-- Period -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Kỳ hóa đơn <span class="text-gray-400 font-normal">(MM-YYYY)</span>
          </label>
          <input
            v-model="form.period"
            type="text"
            placeholder="01-2026"
            maxlength="7"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-400"
          />
        </div>

        <!-- Electricity unit price -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Đơn giá điện (VNĐ/kWh)</label>
          <input
            v-model.number="form.electricity_unit_price"
            type="number"
            min="0"
            step="100"
            placeholder="3500"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-400"
          />
        </div>

        <!-- Water unit price -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Đơn giá nước (VNĐ/m³)</label>
          <input
            v-model.number="form.water_unit_price"
            type="number"
            min="0"
            step="100"
            placeholder="15000"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-orange-400"
          />
        </div>

        <p v-if="formError" class="text-sm text-red-600">{{ formError }}</p>

        <div class="flex justify-end gap-2 pt-1">
          <button
            type="button"
            @click="closeCreateModal"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors"
          >
            Hủy
          </button>
          <button
            type="submit"
            :disabled="submitting"
            class="px-4 py-2 text-sm font-medium text-white bg-orange-500 hover:bg-orange-600 disabled:opacity-60 rounded-lg transition-colors"
          >
            {{ submitting ? 'Đang tạo…' : 'Tạo hóa đơn' }}
          </button>
        </div>
      </form>
    </div>
  </div>

  <!-- ══ Invoice Items Detail Modal ══════════════════════════════════════ -->
  <div
    v-if="showDetailModal"
    class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 p-4"
    @click.self="closeDetail"
  >
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-xl">
      <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
        <div>
          <h3 class="text-base font-semibold text-gray-900">
            Chi tiết hóa đơn — {{ detailInvoice?.period }}
          </h3>
          <p v-if="detailInvoice" class="text-xs text-gray-400 mt-0.5">
            Tổng: {{ formatCurrency(detailInvoice.total_amount) }} ·
            <span
              :class="detailInvoice.status === 'paid' ? 'text-green-600' : 'text-red-500'"
            >
              {{ detailInvoice.status === 'paid' ? 'Đã thanh toán' : 'Chưa thanh toán' }}
            </span>
          </p>
        </div>
        <button @click="closeDetail" class="text-gray-400 hover:text-gray-600">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="px-6 py-5">
        <div v-if="itemsLoading" class="text-center text-sm text-gray-400 py-6">Đang tải…</div>
        <div v-else-if="itemsError" class="text-sm text-red-600">{{ itemsError }}</div>
        <table v-else class="min-w-full text-sm">
          <thead>
            <tr class="text-xs text-gray-500 uppercase">
              <th class="pb-2 text-left font-medium">Mục</th>
              <th class="pb-2 text-right font-medium">Số lượng</th>
              <th class="pb-2 text-right font-medium">Đơn giá (VNĐ)</th>
              <th class="pb-2 text-right font-medium">Thành tiền</th>
              <th class="pb-2 text-right font-medium"></th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="item in invoiceItems" :key="item.item_type">
              <td class="py-2.5 text-gray-900 font-medium">
                {{ itemTypeLabel[item.item_type] ?? item.item_type }}
              </td>
              <td class="py-2.5 text-right text-gray-600">
                {{ item.quantity }} {{ itemTypeUnit[item.item_type] ?? '' }}
              </td>
              <td class="py-2.5 text-right">
                <input
                  v-model.number="editingItem[item.item_type]"
                  type="number"
                  min="0"
                  step="100"
                  class="w-28 px-2 py-1 border border-gray-300 rounded text-sm text-right focus:outline-none focus:ring-1 focus:ring-orange-400"
                />
              </td>
              <td class="py-2.5 text-right text-gray-700">
                {{ formatCurrency(item.amount) }}
              </td>
              <td class="py-2.5 text-right pl-3">
                <button
                  :disabled="savingItem === item.item_type || editingItem[item.item_type] === item.unit_price"
                  @click="saveItem(item)"
                  class="px-2.5 py-1 text-xs font-medium text-orange-600 bg-orange-50 hover:bg-orange-100 disabled:opacity-40 disabled:cursor-not-allowed rounded-md transition-colors"
                >
                  {{ savingItem === item.item_type ? 'Lưu…' : 'Lưu' }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
