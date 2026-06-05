import { useEffect } from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import Layout from './components/Layout'
import Login from './pages/Login'
import Register from './pages/Register'
import Home from './pages/Home'
import ArticleDetail from './pages/ArticleDetail'
import CreateArticle from './pages/CreateArticle'
import EditArticle from './pages/EditArticle'
import UserProfile from './pages/UserProfile'
import MyArticles from './pages/MyArticles'
import AdminUsers from './pages/admin/Users'
import AdminCategories from './pages/admin/Categories'
import AdminTags from './pages/admin/Tags'
import AdminComments from './pages/admin/Comments'
import { useAuthStore } from './store/authStore'

// 路由守卫组件
const ProtectedRoute = ({ children, requireAdmin = false }) => {
  const { user, isAuthenticated } = useAuthStore()

  if (!isAuthenticated()) {
    return <Navigate to="/login" replace />
  }

  if (requireAdmin && user?.role !== 'admin') {
    return <Navigate to="/" replace />
  }

  return children
}

function App() {
  const { fetchCurrentUser, token, user } = useAuthStore()

  // 应用启动时，如果本地有 token 但没有用户信息，尝试获取一次
  useEffect(() => {
    if (token && !user) {
      fetchCurrentUser().catch(() => {
        // 如果获取失败（token失效等），由 fetchCurrentUser 或拦截器清理状态
      })
    }
  }, [])

  return (
    <ConfigProvider locale={zhCN}>
      <BrowserRouter future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
        <Routes>
          {/* 公开路由 */}
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          
          {/* 主布局路由 */}
          <Route path="/" element={<Layout />}>
            <Route index element={<Home />} />
            <Route path="article/:id" element={<ArticleDetail />} />
            
            {/* 需要登录的路由 */}
            <Route 
              path="create" 
              element={
                <ProtectedRoute>
                  <CreateArticle />
                </ProtectedRoute>
              } 
            />
            <Route 
              path="edit/:id" 
              element={
                <ProtectedRoute>
                  <EditArticle />
                </ProtectedRoute>
              } 
            />
            <Route
              path="profile"
              element={
                <ProtectedRoute>
                  <UserProfile />
                </ProtectedRoute>
              }
            />
            <Route
              path="my-articles"
              element={
                <ProtectedRoute>
                  <MyArticles />
                </ProtectedRoute>
              }
            />
            
            {/* 管理员路由 */}
            <Route 
              path="admin/users" 
              element={
                <ProtectedRoute requireAdmin={true}>
                  <AdminUsers />
                </ProtectedRoute>
              } 
            />
            <Route 
              path="admin/categories" 
              element={
                <ProtectedRoute requireAdmin={true}>
                  <AdminCategories />
                </ProtectedRoute>
              } 
            />
            <Route 
              path="admin/tags" 
              element={
                <ProtectedRoute requireAdmin={true}>
                  <AdminTags />
                </ProtectedRoute>
              } 
            />
            <Route 
              path="admin/comments" 
              element={
                <ProtectedRoute requireAdmin={true}>
                  <AdminComments />
                </ProtectedRoute>
              } 
            />
          </Route>
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  )
}

export default App
