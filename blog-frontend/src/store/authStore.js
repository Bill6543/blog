import { create } from 'zustand'
import { getCurrentUser, logout as apiLogout } from '../api/auth'

export const useAuthStore = create((set, get) => ({
  user: JSON.parse(sessionStorage.getItem('user') || 'null'),
  token: sessionStorage.getItem('token'),
  
  // 设置用户信息
  setUser: (user) => set({ user }),
  
  // 设置 token
  setToken: (token) => {
    sessionStorage.setItem('token', token)
    set({ token })
  },
  
  // 获取当前用户信息
  fetchCurrentUser: async () => {
    try {
      const user = await getCurrentUser()
      set({ user })
      // 同步更新 sessionStorage
      if (user) {
        sessionStorage.setItem('user', JSON.stringify(user))
      }
      return user
    } catch (error) {
      console.error('Failed to fetch current user:', error)
      throw error
    }
  },
  
  // 检查是否已认证
  isAuthenticated: () => {
    const state = get()
    return !!state.token && !!state.user
  },
  
  // 登出
  logout: async () => {
    try {
      await apiLogout()
    } catch (error) {
      console.error('Logout failed:', error)
    } finally {
      sessionStorage.removeItem('token')
      sessionStorage.removeItem('user')
      set({ user: null, token: null })
    }
  },
}))
