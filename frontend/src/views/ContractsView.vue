<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getContracts } from '../services/api'

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

onMounted(async () => {
  loading.value = true
  try {
    contracts.value = await getContracts()
  } catch (e: any) {
    error.value = e.message || 'Không thể tải danh sách hóa đơn'
  } finally {
    loading.value = false
  }
})

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
    <div class="mb-6">
      <h2 class="text-2xl font-bold text-gray-800">Hợp đồng</h2>
      <p class="text-sm text-gray-500 mt-1">Danh sách hợp đồng</p>
    </div>

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
          </tr>
        </thead>
        <tbody>
          <tr v-if="contracts.length === 0">
            <td colspan="5" class="px-6 py-8 text-center text-gray-400">Chưa có hợp đồng nào</td>
          </tr>
          <tr
            v-for="c in contracts"
            :key="c.contract_id"
            class="border-b border-gray-100 hover:bg-gray-50 transition-colors"
          >
            <td class="px-6 py-4 font-mono text-xs text-gray-600">{{ c.contract_id }}</td>
            <td class="px-6 py-4 font-medium text-gray-800">{{ c.room_number }}</td>
            <td class="px-6 py-4 text-gray-600 font-mono text-xs">{{ c.user_pid }}</td>
            <td class="px-6 py-4 text-gray-700">{{ c.duration_months ? `${c.duration_months} tháng` : '—' }}</td>
            <td class="px-6 py-4">
              <span
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                :class="statusClass(c.status)"
              >
                {{ statusLabel(c.status) }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
