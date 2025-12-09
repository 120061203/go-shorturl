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
        <!-- 視圖切換按鈕 -->
        <div class="flex justify-end gap-4 mb-6">
          <button
            @click="switchView('chart')"
            :class="[
              'px-6 py-2 rounded-lg font-medium transition-all duration-300',
              viewMode === 'chart'
                ? 'bg-gradient-to-r from-blue-500 to-cyan-500 text-white shadow-lg'
                : 'bg-white/10 text-gray-300 hover:bg-white/20'
            ]"
          >
            圖表視圖
          </button>
          <button
            @click="switchView('table')"
            :class="[
              'px-6 py-2 rounded-lg font-medium transition-all duration-300',
              viewMode === 'table'
                ? 'bg-gradient-to-r from-blue-500 to-cyan-500 text-white shadow-lg'
                : 'bg-white/10 text-gray-300 hover:bg-white/20'
            ]"
          >
            表格視圖
          </button>
        </div>

        <!-- 表格視圖 -->
        <div v-if="viewMode === 'table'" class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20">
          <h2 class="text-2xl font-bold text-white mb-6">點擊記錄列表</h2>
          <div v-if="loading" class="text-center py-10">
            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-500 mx-auto"></div>
            <p class="text-gray-300 mt-4">載入中...</p>
          </div>
          <div v-else-if="error && viewMode === 'table'" class="text-center py-10">
            <div class="bg-red-500/20 border border-red-500/30 rounded-xl p-6 max-w-md mx-auto">
              <p class="text-red-300">{{ error }}</p>
            </div>
          </div>
          <div v-else-if="clickList && clickList.clicks && clickList.clicks.length > 0" class="overflow-x-auto">
            <table class="w-full text-left">
              <thead>
                <tr class="border-b border-white/20">
                  <th class="pb-4 text-gray-300 font-semibold">時間</th>
                  <th class="pb-4 text-gray-300 font-semibold">IP地址</th>
                  <th class="pb-4 text-gray-300 font-semibold">地點</th>
                  <th class="pb-4 text-gray-300 font-semibold">裝置</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(click, index) in clickList.clicks"
                  :key="index"
                  class="border-b border-white/10 hover:bg-white/5 transition-colors"
                >
                  <td class="py-4 text-white font-mono text-sm">{{ formatDateTime(click.clicked_at) }}</td>
                  <td class="py-4 text-gray-300 font-mono text-sm">{{ click.ip_address || '未知' }}</td>
                  <td class="py-4 text-gray-300">{{ click.location || '未知' }}</td>
                  <td class="py-4 text-gray-300">{{ click.device_type || '未知' }}</td>
                </tr>
              </tbody>
            </table>
            <div class="mt-6 text-gray-400 text-sm text-center">
              共 {{ clickList.total }} 筆記錄
            </div>
          </div>
          <div v-else class="text-center py-10">
            <p class="text-gray-300">沒有點擊記錄</p>
          </div>
        </div>

        <!-- 圖表視圖 -->
        <div v-if="viewMode === 'chart'" class="space-y-8">
        <!-- URL Info -->
        <div class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20">
          <h2 class="text-2xl font-bold text-white mb-6">網址資訊</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">短網址</label>
              <p class="text-white font-mono text-lg">{{ getShortUrl(stats.short_code) }}</p>
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

        <!-- IP Stats -->
        <div v-if="stats.ip_stats && stats.ip_stats.length > 0" class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20">
          <h2 class="text-2xl font-bold text-white mb-6">IP 地址統計</h2>
          <div class="space-y-4">
            <div v-for="ip in stats.ip_stats" :key="ip.ip_address" class="bg-white/5 rounded-xl p-6">
              <div class="flex justify-between items-center mb-3">
                <span class="text-gray-300 font-mono font-medium">{{ ip.ip_address }}</span>
                <span class="text-white font-semibold text-lg">{{ ip.count }}</span>
              </div>
              <div class="bg-white/10 rounded-full h-3">
                <div 
                  class="bg-gradient-to-r from-green-500 to-emerald-500 h-3 rounded-full transition-all duration-500" 
                  :style="{ width: `${(ip.count / stats.total_clicks) * 100}%` }"
                ></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Time Distribution - 柱狀圖 -->
        <div v-if="stats.time_distribution && stats.time_distribution.length > 0" class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20">
          <h2 class="text-2xl font-bold text-white mb-6">點擊時間分布</h2>
          <div class="overflow-x-auto">
            <div class="relative" style="height: 320px;">
              <!-- 柱狀圖容器 - 固定在底部 -->
              <div class="absolute bottom-0 left-0 right-0 flex items-end justify-start gap-2" style="height: 250px; padding-bottom: 60px;">
                <div 
                  v-for="timeStat in stats.time_distribution" 
                  :key="timeStat.time" 
                  class="flex flex-col items-center justify-end group min-w-[40px] h-full"
                >
                  <div class="w-full flex flex-col items-center justify-end">
                    <span class="text-white font-semibold text-xs mb-1">{{ timeStat.count }}</span>
                    <div 
                      class="w-full bg-gradient-to-t from-yellow-500 to-orange-500 rounded-t-lg transition-all duration-500 hover:from-yellow-400 hover:to-orange-400 group-hover:opacity-90 cursor-pointer"
                      :style="{ height: `${Math.max((timeStat.count / getMaxTimeCount()) * 240, 8)}px` }"
                      :title="`${formatTime(timeStat.time)}: ${timeStat.count}次點擊`"
                    ></div>
                  </div>
                </div>
              </div>
              <!-- 時間標籤 - 固定在底部 -->
              <div class="absolute bottom-0 left-0 right-0 flex items-center justify-start gap-2" style="height: 60px;">
                <div 
                  v-for="timeStat in stats.time_distribution" 
                  :key="timeStat.time" 
                  class="flex items-center justify-center min-w-[40px] h-full"
                >
                  <span class="text-gray-400 text-xs text-center leading-tight whitespace-nowrap" style="transform: rotate(-45deg); transform-origin: center;">
                    {{ formatTimeLabel(timeStat.time) }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- OS Stats -->
        <div v-if="stats.os_stats && stats.os_stats.length > 0" class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20">
          <h2 class="text-2xl font-bold text-white mb-6">操作系統分布</h2>
          <div class="space-y-4">
            <div v-for="os in stats.os_stats" :key="os.os" class="bg-white/5 rounded-xl p-6">
              <div class="flex justify-between items-center mb-3">
                <span class="text-gray-300 font-medium text-lg">{{ os.os }}</span>
                <span class="text-white font-semibold text-lg">{{ os.count }}</span>
              </div>
              <div class="bg-white/10 rounded-full h-3">
                <div 
                  class="bg-gradient-to-r from-cyan-500 to-blue-500 h-3 rounded-full transition-all duration-500" 
                  :style="{ width: `${(os.count / stats.total_clicks) * 100}%` }"
                ></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Device Type Stats -->
        <div v-if="stats.device_type_stats && stats.device_type_stats.length > 0" class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20">
          <h2 class="text-2xl font-bold text-white mb-6">設備類型分布</h2>
          <div class="space-y-4">
            <div v-for="deviceType in stats.device_type_stats" :key="deviceType.device_type" class="bg-white/5 rounded-xl p-6">
              <div class="flex justify-between items-center mb-3">
                <span class="text-gray-300 font-medium text-lg">{{ deviceType.device_type }}</span>
                <span class="text-white font-semibold text-lg">{{ deviceType.count }}</span>
              </div>
              <div class="bg-white/10 rounded-full h-3">
                <div 
                  class="bg-gradient-to-r from-indigo-500 to-purple-500 h-3 rounded-full transition-all duration-500" 
                  :style="{ width: `${(deviceType.count / stats.total_clicks) * 100}%` }"
                ></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Location Stats -->
        <div v-if="stats.location_stats && stats.location_stats.length > 0" class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20">
          <h2 class="text-2xl font-bold text-white mb-6">地理位置分布</h2>
          <div class="space-y-4">
            <div v-for="location in stats.location_stats" :key="location.location" class="bg-white/5 rounded-xl p-6">
              <div class="flex justify-between items-center mb-3">
                <span class="text-gray-300 font-medium">{{ location.location }}</span>
                <span class="text-white font-semibold text-lg">{{ location.count }}</span>
              </div>
              <div class="bg-white/10 rounded-full h-3">
                <div 
                  class="bg-gradient-to-r from-teal-500 to-emerald-500 h-3 rounded-full transition-all duration-500" 
                  :style="{ width: `${(location.count / stats.total_clicks) * 100}%` }"
                ></div>
              </div>
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
        <div v-if="!stats.device_stats?.length && !stats.referrer_stats?.length && !stats.ip_stats?.length && !stats.time_distribution?.length && !stats.device_type_stats?.length && !stats.location_stats?.length && !stats.os_stats?.length" class="text-center py-20">
          <div class="bg-white/5 rounded-2xl p-8 max-w-md mx-auto">
            <p class="text-gray-300">還沒有點擊數據</p>
            <p class="text-gray-400 text-sm mt-2">分享你的短網址來開始收集數據</p>
          </div>
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

interface IPStat {
  ip_address: string
  count: number
}

interface TimeDistributionStat {
  time: string
  count: number
}

interface DeviceTypeStat {
  device_type: string
  count: number
}

interface LocationStat {
  location: string
  count: number
}

interface OSStat {
  os: string
  count: number
}

interface Stats {
  short_code: string
  original_url: string
  total_clicks: number
  created_at: string
  device_stats: DeviceStat[]
  referrer_stats: ReferrerStat[]
  ip_stats?: IPStat[]
  time_distribution?: TimeDistributionStat[]
  device_type_stats?: DeviceTypeStat[]
  location_stats?: LocationStat[]
  os_stats?: OSStat[]
}

interface ClickDetail {
  clicked_at: string
  ip_address: string
  location: string
  device_type: string
}

interface ClickList {
  short_code: string
  clicks: ClickDetail[]
  total: number
}

const route = useRoute()
const stats = ref<Stats | null>(null)
const clickList = ref<ClickList | null>(null)
const loading = ref(false)
const error = ref('')
const shortCode = ref('')
const viewMode = ref<'chart' | 'table'>('chart')

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

const fetchClickList = async () => {
  if (!shortCode.value) {
    error.value = '請輸入短碼'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const response = await api.get(`/api/clicks/${shortCode.value}`)
    clickList.value = response.data
  } catch (err: any) {
    if (err.response?.data?.error) {
      error.value = err.response.data.error
    } else {
      error.value = '獲取點擊列表時發生錯誤'
    }
    clickList.value = null
  } finally {
    loading.value = false
  }
}

const switchView = async (mode: 'chart' | 'table') => {
  viewMode.value = mode
  if (mode === 'table' && !clickList.value) {
    await fetchClickList()
  }
}

// 格式化時間
const formatDateTime = (dateStr: string): string => {
  try {
    const date = new Date(dateStr)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    const seconds = String(date.getSeconds()).padStart(2, '0')
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
  } catch {
    return dateStr
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

// 生成短網址
const getShortUrl = (shortCode: string): string => {
  return `https://xsong.us/url/${shortCode}`
}

// 格式化時間顯示
const formatTime = (timeStr: string): string => {
  try {
    // 格式: "2024-01-01 14:00"
    const [date, hour] = timeStr.split(' ')
    const [year, month, day] = date.split('-')
    return `${year}-${month}-${day} ${hour}:00`
  } catch {
    return timeStr
  }
}

// 格式化時間標籤（用於柱狀圖）
const formatTimeLabel = (timeStr: string): string => {
  try {
    // 格式: "2024-01-01 14:00"
    const [date, hour] = timeStr.split(' ')
    const [year, month, day] = date.split('-')
    // 只顯示日期和小時，簡化顯示
    return `${month}/${day} ${hour}:00`
  } catch {
    return timeStr
  }
}

// 獲取時間分布中的最大點擊數（用於計算百分比）
const getMaxTimeCount = (): number => {
  if (!stats.value?.time_distribution || stats.value.time_distribution.length === 0) {
    return 1
  }
  return Math.max(...stats.value.time_distribution.map(t => t.count))
}
</script>
