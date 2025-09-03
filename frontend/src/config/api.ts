import axios from 'axios'

// 設定 axios 基礎配置
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || (import.meta.env.PROD ? window.location.origin : 'http://localhost:8080'),
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 請求攔截器
api.interceptors.request.use(
  (config) => {
    // 可以在這裡添加認證 token
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 回應攔截器
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // 統一錯誤處理
    if (error.response?.status === 404) {
      console.error('API 端點不存在')
    } else if (error.response?.status >= 500) {
      console.error('伺服器錯誤')
    }
    return Promise.reject(error)
  }
)

export default api
