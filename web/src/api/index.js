import axios from 'axios'
import router from '@/router'

// 用户 token 存储
const tokenKey = 'user_token'
const accountIdKey = 'current_account_id'

let currentAccountId = localStorage.getItem(accountIdKey) || ''

export function setAccountId(id) {
  currentAccountId = id || ''
  if (id) localStorage.setItem(accountIdKey, id)
  else localStorage.removeItem(accountIdKey)
}
export function getAccountId() {
  return currentAccountId
}
export function getToken() {
  return localStorage.getItem(tokenKey) || ''
}
export function setToken(t) {
  if (t) localStorage.setItem(tokenKey, t)
  else localStorage.removeItem(tokenKey)
}

const api = axios.create({
  baseURL: '/',
  timeout: 20000,
})

// 统一注入鉴权头 + 账号查询参数
// 使用 Authorization: Bearer <token> 进行用户鉴权
api.interceptors.request.use((config) => {
  const token = getToken()
  if (token) config.headers['Authorization'] = 'Bearer ' + token
  if (currentAccountId) {
    config.params = { ...(config.params || {}), accountId: currentAccountId }
  }
  return config
})

let lastNetToast = 0
api.interceptors.response.use(
  (res) => res,
  (error) => {
    const status = error.response?.status
    if (status === 401) {
      // 未登录或登录过期：清 token 跳登录（登录页自身除外）
      if (!location.pathname.includes('/login')) {
        setToken('')
        router.replace('/login')
      }
    } else if (!status) {
      // 网络层错误，限频提示
      const now = Date.now()
      if (now - lastNetToast > 5000) {
        lastNetToast = now
        console.warn('网络错误:', error.message)
      }
    }
    return Promise.reject(error)
  }
)

export const syncApi = {
  status: () => api.get('/api/sync/status'),
  checkAssets: () => api.post('/api/sync/assets/check'),
}

export default api
