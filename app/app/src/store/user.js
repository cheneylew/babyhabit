import { defineStore } from 'pinia'
import api from '../api'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    token: uni.getStorageSync('token') || '',
    isLoading: false,
    error: null
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.user?.user_type === 1
  },

  actions: {
    async login(username, password) {
      this.isLoading = true
      this.error = null
      try {
        const response = await api.post('/api/login', { username, password })
        this.token = response.data.token
        this.user = response.data.user
        uni.setStorageSync('token', response.data.token)
        return response.data
      } catch (error) {
        this.error = error.response?.data?.error || '登录失败'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    async register(data) {
      this.isLoading = true
      this.error = null
      try {
        const response = await api.post('/api/register', data)
        this.token = response.data.token
        this.user = response.data.user
        uni.setStorageSync('token', response.data.token)
        return response.data
      } catch (error) {
        this.error = error.response?.data?.error || '注册失败'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    async getUserInfo() {
      if (!this.token) return
      this.isLoading = true
      try {
        const response = await api.get('/api/user/info')
        this.user = response.data.user
      } catch (error) {
        this.logout()
      } finally {
        this.isLoading = false
      }
    },

    logout() {
      this.user = null
      this.token = ''
      uni.removeStorageSync('token')
    }
  }
})