import { defineStore } from 'pinia'
import api, { setAccountId, getAccountId, setToken, getToken } from '@/api'

export const useAccountStore = defineStore('account', {
  state: () => ({
    accounts: [],
    currentId: getAccountId(),
    userLoggedIn: !!getToken(),
    userInfo: null,
  }),
  getters: {
    current(state) {
      return state.accounts.find((a) => String(a.id) === String(state.currentId)) || null
    },
  },
  actions: {
    // 设置 token（供 Login.vue 调用）
    setToken(token) {
      setToken(token)
      this.userLoggedIn = !!token
    },
    
    // 获取当前用户信息
    async loadUserInfo() {
      try {
        const { data } = await api.get('/api/users/me')
        if (data.ok) {
          this.userInfo = data.user
          return data.user
        }
        return null
      } catch (e) {
        return null
      }
    },
    
    // 用户登出
    logout() {
      setToken('')
      this.userLoggedIn = false
      this.userInfo = null
      this.accounts = []
      this.currentId = ''
      setAccountId('')
    },
    
    // 加载账号列表
    async loadAccounts() {
      const { data } = await api.get('/api/accounts')
      this.accounts = (data && data.data) || data.accounts || data.list || []
      const exists = this.accounts.some((a) => String(a.id) === String(this.currentId))
      if (this.accounts.length && !exists) {
        this.switchAccount(this.accounts[0].id)
      } else if (!this.accounts.length) {
        this.currentId = ''
        setAccountId('')
      }
      return this.accounts
    },
    
    // 切换账号
    switchAccount(id) {
      const next = id ? String(id) : ''
      this.currentId = next
      setAccountId(next)
    },
  },
})
