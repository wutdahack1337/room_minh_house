<script setup lang="ts">
import { RouterLink } from 'vue-router'

const quickLinks = [
  { to: '/rooms', label: 'Phòng', desc: 'Xem và quản lý danh sách phòng trọ', color: 'bg-blue-50 border-blue-200 text-blue-700', icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6' },
  { to: '/users', label: 'Người thuê', desc: 'Quản lý thông tin người thuê phòng', color: 'bg-purple-50 border-purple-200 text-purple-700', icon: 'M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z' },
  { to: '/contracts', label: 'Hợp đồng', desc: 'Xem danh sách hợp đồng', color: 'bg-green-50 border-green-200 text-green-700', icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z' },
]
</script>

<template>
  <div>
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-gray-800">Dashboard</h2>
      <p class="text-sm text-gray-500 mt-1">Chào mừng đến hệ thống quản lý phòng trọ</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <RouterLink
        v-for="link in quickLinks"
        :key="link.to"
        :to="link.to"
        class="flex items-start gap-4 p-5 rounded-xl border transition-shadow hover:shadow-md"
        :class="link.color"
      >
        <div class="mt-0.5">
          <svg class="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" :d="link.icon" />
          </svg>
        </div>
        <div>
          <p class="font-semibold text-sm">{{ link.label }}</p>
          <p class="text-xs mt-0.5 opacity-75">{{ link.desc }}</p>
        </div>
      </RouterLink>
    </div>
  </div>
</template>
