<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { DocumentDuplicateIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import QRCode from 'qrcode'
import api from '../config/api'

interface FormData {
  url: string
  customCode: string
}

interface Errors {
  url: string
  customCode: string
}

const formData = reactive<FormData>({
  url: '',
  customCode: ''
})

const errors = reactive<Errors>({
  url: '',
  customCode: ''
})

const loading = ref(false)
const shortUrl = ref('')
const error = ref('')
const qrCodeDataUrl = ref('')
const showCopySuccess = ref(false)
const copySuccessPosition = ref({ x: 0, y: 0 })
const errorPosition = ref({ x: 0, y: 0 })

const validateForm = (): boolean => {
  errors.url = ''
  errors.customCode = ''

  if (!formData.url) {
    errors.url = '請輸入網址'
    return false
  }

  // 支援沒有協議的網址
  let urlToValidate = formData.url
  if (!urlToValidate.startsWith('http://') && !urlToValidate.startsWith('https://')) {
    urlToValidate = 'https://' + urlToValidate
  }

  try {
    new URL(urlToValidate)
  } catch {
    errors.url = '請輸入有效的網址'
    return false
  }

  if (formData.customCode && !/^[a-zA-Z0-9-]+$/.test(formData.customCode)) {
    errors.customCode = '短碼只能包含字母、數字和連字符'
    return false
  }

  return true
}

const createShortUrl = async (event: MouseEvent) => {
  if (!validateForm()) return

  // 獲取按鈕位置用於錯誤提示
  const button = event.currentTarget as HTMLElement
  const rect = button.getBoundingClientRect()
  errorPosition.value = {
    x: rect.left + rect.width / 2,
    y: rect.bottom + 10
  }

  loading.value = true
  error.value = ''

  try {
    const response = await api.post('/api/shorten', {
      url: formData.url,
      custom_code: formData.customCode || undefined
    })

    shortUrl.value = response.data.short_url
    await generateQRCode()
  } catch (err: any) {
    if (err.response?.data?.error) {
      error.value = err.response.data.error
    } else {
      error.value = '建立短網址時發生錯誤'
    }
  } finally {
    loading.value = false
  }
}

const generateQRCode = async () => {
  if (!shortUrl.value) return

  try {
    const qrCodeUrl = await QRCode.toDataURL(shortUrl.value, {
      width: 300,
      margin: 2,
      color: {
        dark: '#000000',
        light: '#FFFFFF'
      }
    })
    qrCodeDataUrl.value = qrCodeUrl
  } catch (err) {
    console.error('QR Code 產生失敗:', err)
  }
}

const downloadQRCode = () => {
  if (!qrCodeDataUrl.value) return

  const link = document.createElement('a')
  link.download = 'qrcode.png'
  link.href = qrCodeDataUrl.value
  link.click()
}

const copyToClipboard = async (event: MouseEvent) => {
  // 獲取按鈕位置
  const button = event.currentTarget as HTMLElement
  const rect = button.getBoundingClientRect()
  
  // 設置提示位置在按鈕下方
  copySuccessPosition.value = {
    x: rect.left + rect.width / 2,
    y: rect.bottom + 10
  }

  try {
    await navigator.clipboard.writeText(shortUrl.value)
    showCopySuccess.value = true
    setTimeout(() => {
      showCopySuccess.value = false
    }, 2000)
  } catch {
    // 降級方案
    const textArea = document.createElement('textarea')
    textArea.value = shortUrl.value
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    showCopySuccess.value = true
    setTimeout(() => {
      showCopySuccess.value = false
    }, 2000)
  }
}

const testUrl = () => {
  window.open(shortUrl.value, '_blank')
}

const resetForm = () => {
  formData.url = ''
  formData.customCode = ''
  shortUrl.value = ''
  error.value = ''
  qrCodeDataUrl.value = ''
  showCopySuccess.value = false
  errors.url = ''
  errors.customCode = ''
}

// 監聽短網址變化，自動產生 QR Code
watch(shortUrl, (newUrl) => {
  if (newUrl) {
    generateQRCode()
  }
})
</script>

<template>
  <div class="bg-gradient-to-br from-slate-900 via-slate-800 to-purple-900">
    <!-- Hero Section -->
    <div class="relative overflow-hidden">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <div class="text-center mb-16">
          <h1 class="text-5xl md:text-7xl font-bold text-white mb-8">
            xsong
            <span class="block text-2xl md:text-3xl font-light text-gray-300 mt-4">
              短網址服務
            </span>
          </h1>
          <p class="text-xl text-gray-300 max-w-2xl mx-auto">
            簡潔、快速、個人化的短網址工具
          </p>
        </div>

        <!-- Main URL Input -->
        <div class="max-w-3xl mx-auto">
          <div class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20 shadow-2xl">
            <div class="space-y-6">
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">輸入網址</label>
                <input
                  v-model="formData.url"
                  type="url"
                  placeholder="https://example.com"
                  class="w-full text-lg px-6 py-4 bg-white/10 border border-white/20 rounded-xl text-white placeholder-gray-400 focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all duration-200"
                  :class="{ 'border-red-400': errors.url }"
                />
                <p v-if="errors.url" class="mt-2 text-sm text-red-400">{{ errors.url }}</p>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">自訂短碼 (可選)</label>
                <input
                  v-model="formData.customCode"
                  type="text"
                  placeholder="my-custom-code"
                  class="w-full px-6 py-4 bg-white/10 border border-white/20 rounded-xl text-white placeholder-gray-400 focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all duration-200"
                  :class="{ 'border-red-400': errors.customCode }"
                />
                <p v-if="errors.customCode" class="mt-2 text-sm text-red-400">{{ errors.customCode }}</p>
              </div>
              
              <button
                @click="createShortUrl"
                :disabled="loading"
                class="w-full bg-gradient-to-r from-purple-600 to-pink-600 text-white py-4 px-8 rounded-xl font-semibold text-lg hover:from-purple-700 hover:to-pink-700 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg"
              >
                <span v-if="loading" class="flex items-center justify-center">
                  <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  處理中...
                </span>
                <span v-else>建立短網址</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Results Section -->
    <div v-if="shortUrl" class="bg-black/20 backdrop-blur-sm border-t border-white/10">
      <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <!-- Shortened URL Card -->
          <div class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20">
            <h2 class="text-2xl font-bold text-white mb-6">您的短網址</h2>
            
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">短網址</label>
                <div class="flex">
                  <input
                    :value="shortUrl"
                    readonly
                    class="flex-1 px-4 py-3 bg-white/10 border border-white/20 rounded-l-xl text-white"
                  />
                  <button
                    @click="copyToClipboard"
                    class="px-4 py-3 bg-purple-600 text-white rounded-r-xl hover:bg-purple-700 transition-colors"
                  >
                    <DocumentDuplicateIcon class="h-5 w-5" />
                  </button>
                </div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">原始網址</label>
                <p class="px-4 py-3 bg-white/10 border border-white/20 rounded-xl text-gray-300 break-all">
                  {{ formData.url }}
                </p>
              </div>

              <div class="flex space-x-3">
                <button
                  @click="testUrl"
                  class="flex-1 bg-green-600 text-white py-3 px-4 rounded-xl hover:bg-green-700 transition-colors"
                >
                  測試連結
                </button>
                <button
                  @click="resetForm"
                  class="flex-1 bg-gray-600 text-white py-3 px-4 rounded-xl hover:bg-gray-700 transition-colors"
                >
                  建立新的
                </button>
              </div>
            </div>
          </div>

          <!-- QR Code Card -->
          <div class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20">
            <h2 class="text-2xl font-bold text-white mb-6">QR Code</h2>
            
            <div class="text-center">
              <div v-if="qrCodeDataUrl" class="mb-4">
                <img :src="qrCodeDataUrl" alt="QR Code" class="mx-auto w-48 h-48 bg-white rounded-xl p-2" />
              </div>
              <div v-else class="w-48 h-48 bg-white/10 rounded-xl mx-auto flex items-center justify-center border border-white/20">
                <span class="text-gray-400">QR Code 將在這裡顯示</span>
              </div>
              
              <button
                @click="downloadQRCode"
                :disabled="!qrCodeDataUrl"
                class="mt-4 bg-blue-600 text-white py-3 px-6 rounded-xl hover:bg-blue-700 transition-colors disabled:opacity-50"
              >
                下載 QR Code
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="fixed bg-red-500/90 backdrop-blur-sm border border-red-400 rounded-xl p-4 max-w-md z-50 transform -translate-x-1/2" :style="{ left: `${errorPosition.x}px`, top: `${errorPosition.y}px` }">
      <div class="flex">
        <ExclamationTriangleIcon class="h-5 w-5 text-red-200" />
        <div class="ml-3">
          <h3 class="text-sm font-medium text-white">錯誤</h3>
          <p class="mt-1 text-sm text-red-100">{{ error }}</p>
        </div>
              </div>
      </div>

      <!-- Copy Success Message -->
      <div v-if="showCopySuccess" class="fixed bg-green-500/90 backdrop-blur-sm border border-green-400 rounded-xl p-4 max-w-md z-50 transform -translate-x-1/2" :style="{ left: `${copySuccessPosition.x}px`, top: `${copySuccessPosition.y}px` }">
        <div class="flex">
          <svg class="h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 6 9 17l-5-5"></path>
          </svg>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-white">複製成功</h3>
            <p class="mt-1 text-sm text-green-100">短網址已複製到剪貼簿！</p>
          </div>
        </div>
      </div>
    </div>
  </template>
