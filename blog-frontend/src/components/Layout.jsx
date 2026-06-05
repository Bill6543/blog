import { useState, useEffect } from 'react'
import { Outlet, Link, useNavigate, useLocation } from 'react-router-dom'
import { Layout as AntLayout, Menu, Avatar, Dropdown, Space, Badge } from 'antd'
import {
  HomeOutlined,
  PlusOutlined,
  UserOutlined,
  SettingOutlined,
  TeamOutlined,
  AppstoreOutlined,
  TagsOutlined,
  LogoutOutlined,
  MessageOutlined,
  CommentOutlined,
  FileTextOutlined
} from '@ant-design/icons'
import { useAuthStore } from '../store/authStore'
import { getPendingCommentCount } from '../api/comment'
import './Layout.css'

const { Header, Content, Footer } = AntLayout

const Layout = () => {
  const [collapsed, setCollapsed] = useState(false)
  const [pendingCommentCount, setPendingCommentCount] = useState(0)
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout, isAuthenticated } = useAuthStore()

  // 获取待审核评论数量
  useEffect(() => {
    if (isAuthenticated() && user?.role === 'admin') {
      const fetchPendingCount = async () => {
        try {
          const count = await getPendingCommentCount()
          setPendingCommentCount(count)
        } catch (error) {
          console.error('Failed to fetch pending comment count:', error)
        }
      }
      fetchPendingCount()
      // 每30秒刷新一次
      const interval = setInterval(fetchPendingCount, 30000)
      return () => clearInterval(interval)
    }
  }, [isAuthenticated, user])

  // 用户菜单项
  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心',
    },
    {
      key: 'my-articles',
      icon: <FileTextOutlined />,
      label: '我的文章',
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
    },
  ]

  // 处理菜单点击
  const handleUserMenuClick = async ({ key }) => {
    if (key === 'profile') {
      navigate('/profile')
    } else if (key === 'my-articles') {
      navigate('/my-articles')
    } else if (key === 'logout') {
      await logout()
      navigate('/login')
    }
  }

  // 导航菜单项
  const navMenuItems = [
    {
      key: '/',
      icon: <HomeOutlined />,
      label: <Link to="/">首页</Link>,
    },
  ]

  if (isAuthenticated()) {
    navMenuItems.push({
      key: '/create',
      icon: <PlusOutlined />,
      label: <Link to="/create">创作</Link>,
    })
  }

  // 管理员菜单
  const adminMenuItems = user?.role === 'admin' ? [
    {
      key: 'admin',
      icon: pendingCommentCount > 0 ? <Badge dot><SettingOutlined /></Badge> : <SettingOutlined />,
      label: '管理',
      children: [
        {
          key: '/admin/users',
          icon: <TeamOutlined />,
          label: <Link to="/admin/users">用户管理</Link>,
        },
        {
          key: '/admin/categories',
          icon: <AppstoreOutlined />,
          label: <Link to="/admin/categories">分类管理</Link>,
        },
        {
          key: '/admin/tags',
          icon: <TagsOutlined />,
          label: <Link to="/admin/tags">标签管理</Link>,
        },
        {
          key: '/admin/comments',
          icon: <CommentOutlined />,
          label: <Link to="/admin/comments">评论审核 {pendingCommentCount > 0 && <span style={{ 
            display: 'inline-block', 
            minWidth: '18px', 
            height: '18px', 
            lineHeight: '18px', 
            textAlign: 'center', 
            borderRadius: '9px', 
            backgroundColor: '#ff4d4f', 
            color: 'white', 
            fontSize: '12px', 
            marginLeft: '8px' 
          }}>{pendingCommentCount}</span>}</Link>,
        },
      ],
    },
  ] : []

  // 添加管理员菜单到导航菜单
  if (user?.role === 'admin') {
    navMenuItems.push(...adminMenuItems)
  }

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      <Header style={{ 
        display: 'flex', 
        alignItems: 'center', 
        justifyContent: 'space-between',
        padding: '0 24px'
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '20px' }}>
          <h1 style={{ color: 'white', margin: 0, fontSize: '20px' }}>博客系统</h1>
          <Menu
            theme="dark"
            mode="horizontal"
            selectedKeys={[location.pathname]}
            items={navMenuItems}
            style={{ minWidth: 0, flex: 'none' }}
          />
        </div>

        <Space size="large">
          {isAuthenticated() ? (
            <>
              <Badge count={0} size="small">
                <MessageOutlined style={{ fontSize: 18, color: 'white' }} />
              </Badge>
              
              <Dropdown
                menu={{
                  items: userMenuItems,
                  onClick: handleUserMenuClick,
                }}
                placement="bottomRight"
              >
                <Space style={{ cursor: 'pointer', color: 'white' }}>
                  <Avatar src={user?.avatar} icon={<UserOutlined />} />
                  <span>{user?.nickname || user?.username}</span>
                </Space>
              </Dropdown>
            </>
          ) : (
            <Space>
              <Link to="/login" style={{ color: 'white' }}>登录</Link>
              <Link to="/register" style={{ color: 'white' }}>注册</Link>
            </Space>
          )}
        </Space>
      </Header>

      <Content style={{ padding: '24px 50px', minHeight: 'calc(100vh - 128px)' }}>
        <Outlet />
      </Content>

      <Footer style={{ textAlign: 'center' }}>
        博客系统 ©{new Date().getFullYear()} Created with React & Vite
      </Footer>
    </AntLayout>
  )
}

export default Layout
