<template>
  <div class="bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
    <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <!-- Header -->
      <div class="text-center mb-12">
        <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">
          統計分析
        </h1>
        <p class="text-xl text-gray-300">
          短網址的詳細使用數據
        </p>
      </div>

      <!-- URL Input Section -->
      <div class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20 mb-8">
        <h2 class="text-2xl font-bold text-white mb-6">查詢統計</h2>
        <div class="flex flex-col md:flex-row gap-4">
          <input
            v-model="shortCode"
            type="text"
            placeholder="輸入短碼，例如：123"
            class="flex-1 px-6 py-4 bg-white/10 border border-white/20 rounded-xl text-white placeholder-gray-400 focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all duration-200"
          />
          <button
            @click="fetchStats"
            :disabled="loading || !shortCode"
            class="bg-gradient-to-r from-purple-600 to-pink-600 text-white py-4 px-8 rounded-xl font-semibold hover:from-purple-700 hover:to-pink-700 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="loading" class="flex items-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              查詢中...
            </span>
            <span v-else>查詢統計</span>
          </button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading && !stats" class="text-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-500 mx-auto"></div>
        <p class="text-gray-300 mt-4">載入中...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-20">
        <div class="bg-red-500/20 border border-red-500/30 rounded-xl p-8 max-w-md mx-auto">
          <p class="text-red-300">{{ error }}</p>
        </div>
      </div>

      <!-- Stats Content -->
      <div v-else-if="stats" class="space-y-8">
        <!-- URL Info -->
        <div class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20">
          <h2 class="text-2xl font-bold text-white mb-6">網址資訊</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">短網址</label>
              <p class="text-white font-mono text-lg">http://localhost:8080/{{ stats.short_code }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">原始網址</label>
              <p class="text-gray-300 break-all">{{ stats.original_url }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">建立時間</label>
              <p class="text-gray-300">{{ new Date(stats.created_at).toLocaleString() }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">總點擊數</label>
              <p class="text-3xl font-bold text-purple-400">{{ stats.total_clicks }}</p>
            </div>
          </div>
        </div>

        <!-- Device Stats -->
        <div v-if="stats.device_stats && stats.device_stats.length > 0" class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20">
          <h2 class="text-2xl font-bold text-white mb-6">裝置分布</h2>
          <div class="space-y-4">
            <div v-for="device in stats.device_stats" :key="device.user_agent" class="bg-white/5 rounded-xl p-6">
              <div class="flex justify-between items-center mb-3">
                <span class="text-gray-300 font-medium">{{ getDeviceName(device.user_agent) }}</span>
                <span class="text-white font-semibold text-lg">{{ device.count }}</span>
              </div>
              <div class="bg-white/10 rounded-full h-3">
                <div 
                  class="bg-gradient-to-r from-purple-500 to-pink-500 h-3 rounded-full transition-all duration-500" 
                  :style="{ width: `${(device.count / stats.total_clicks) * 100}%` }"
                ></div>
              </div>
              <p class="text-gray-400 text-sm mt-2 truncate">{{ device.user_agent }}</p>
            </div>
          </div>
        </div>

        <!-- Referrer Stats -->
        <div v-if="stats.referrer_stats && stats.referrer_stats.length > 0" class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20">
          <h2 class="text-2xl font-bold text-white mb-6">來源統計</h2>
          <div class="space-y-4">
            <div v-for="referrer in stats.referrer_stats" :key="referrer.referrer" class="bg-white/5 rounded-xl p-6">
              <div class="flex justify-between items-center mb-3">
                <span class="text-gray-300 font-medium">{{ referrer.referrer || '直接訪問' }}</span>
                <span class="text-white font-semibold text-lg">{{ referrer.count }}</span>
              </div>
              <div class="bg-white/10 rounded-full h-3">
                <div 
                  class="bg-gradient-to-r from-blue-500 to-cyan-500 h-3 rounded-full transition-all duration-500" 
                  :style="{ width: `${(referrer.count / stats.total_clicks) * 100}%` }"
                ></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="!stats.device_stats?.length && !stats.referrer_stats?.length" class="text-center py-20">
          <div class="bg-white/5 rounded-2xl p-8 max-w-md mx-auto">
            <p class="text-gray-300">還沒有點擊數據</p>
            <p class="text-gray-400 text-sm mt-2">分享你的短網址來開始收集數據</p>
          </div>
        </div>
      </div>

      <!-- Initial State -->
      <div v-else class="text-center py-20">
        <div class="bg-white/5 rounded-2xl p-8 max-w-md mx-auto">
          <p class="text-gray-300">輸入短碼來查看統計數據</p>
          <p class="text-gray-400 text-sm mt-2">例如：123, gg, test</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '../config/api'

interface DeviceStat {
  user_agent: string
  count: number
}

interface ReferrerStat {
  referrer: string
  count: number
}

interface Stats {
  short_code: string
  original_url: string
  total_clicks: number
  created_at: string
  device_stats: DeviceStat[]
  referrer_stats: ReferrerStat[]
}

const route = useRoute()
const stats = ref<Stats | null>(null)
const loading = ref(false)
const error = ref('')
const shortCode = ref('')

// 從路由參數獲取短碼
onMounted(() => {
  if (route.params.short_code) {
    shortCode.value = route.params.short_code as string
    fetchStats()
  }
})

const fetchStats = async () => {
  if (!shortCode.value) {
    error.value = '請輸入短碼'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const response = await api.get(`/api/stats/${shortCode.value}`)
    stats.value = response.data
  } catch (err: any) {
    if (err.response?.data?.error) {
      error.value = err.response.data.error
    } else {
      error.value = '獲取統計資料時發生錯誤'
    }
    stats.value = null
  } finally {
    loading.value = false
  }
}

// 解析裝置名稱
const getDeviceName = (userAgent: string): string => {
  if (userAgent.includes('Chrome')) return 'Chrome'
  if (userAgent.includes('Firefox')) return 'Firefox'
  if (userAgent.includes('Safari')) return 'Safari'
  if (userAgent.includes('Edge')) return 'Edge'
  if (userAgent.includes('curl')) return 'API 調用'
  return '其他'
}
</script>
