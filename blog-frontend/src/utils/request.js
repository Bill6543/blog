import axios from 'axios'
import { message } from 'antd'

// 创建 axios 实例
const request = axios.create({
  baseURL: '/api',
  timeout: 180000, // 增加到180秒以匹配 Coze API 的封面生成超时时间
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const token = sessionStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const { code, message: msg, data } = response.data
    
    if (code === 0) {
      return data
    } else {
      message.error(msg || '请求失败')
      return Promise.reject(response.data)
    }
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 401:
          message.error(data?.message || '用户名或密码错误')
          sessionStorage.removeItem('token')
          // 如果已经在登录页，不要重复跳转
          if (window.location.pathname !== '/login') {
            window.location.href = '/login'
          }
          break
        case 403:
          message.error(data?.message || '权限不足')
          break
        case 404:
          message.error(data?.message || '资源不存在')
          break
        case 500:
          message.error(data?.message || '服务器错误')
          break
        default:
          message.error(data?.message || '请求失败')
      }
    } else {
      message.error('网络错误，请检查网络连接')
    }
    
    return Promise.reject(error)
  }
)

export default request
