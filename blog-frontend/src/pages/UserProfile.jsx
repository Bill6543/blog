import { useState, useEffect } from 'react'
import { Form, Input, Button, Card, Avatar, message } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import { useAuthStore } from '../store/authStore'
import { updateUserInfo } from '../api/user'

const UserProfile = () => {
  const { user, setUser } = useAuthStore()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (user) {
      form.setFieldsValue({
        nickname: user.nickname,
        bio: user.bio,
      })
    }
  }, [user, form])

  const onFinish = async (values) => {
    setLoading(true)
    try {
      const updatedUser = await updateUserInfo(user.id, values)
      setUser(updatedUser)
      message.success('更新成功')
    } catch (error) {
      console.error('Update failed:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ maxWidth: 600, margin: '0 auto' }}>
      <Card title="个人中心">
        <div style={{ textAlign: 'center', marginBottom: 32 }}>
          <Avatar src={user?.avatar} size={100} icon={<UserOutlined />} />
          <h2 style={{ marginTop: 16 }}>{user?.username}</h2>
        </div>

        <Form form={form} layout="vertical" onFinish={onFinish}>
          <Form.Item name="nickname" label="昵称">
            <Input placeholder="请输入昵称" />
          </Form.Item>

          <Form.Item name="bio" label="个人简介">
            <Input.TextArea 
              placeholder="请输入个人简介" 
              maxLength={500}
              showCount
              rows={4}
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              保存
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default UserProfile
