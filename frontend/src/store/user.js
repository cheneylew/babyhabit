import { defineStore } from 'pinia'
import api from '../api'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token') || '',
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
        const response = await api.post('/login', { username, password })
        this.token = response.data.token
        this.user = response.data.user
        localStorage.setItem('token', response.data.token)
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
        const response = await api.post('/register', data)
        this.token = response.data.token
        this.user = response.data.user
        localStorage.setItem('token', response.data.token)
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
        const response = await api.get('/user/info')
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
      localStorage.removeItem('token')
    }
  }
})