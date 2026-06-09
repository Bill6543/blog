import { useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { Form, Input, Button, Card, message } from 'antd'
import { LockOutlined } from '@ant-design/icons'
import { resetPassword } from '../api/auth'

const ResetPassword = () => {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [loading, setLoading] = useState(false)
  
  const token = searchParams.get('token')

  if (!token) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
      }}>
        <Card style={{ width: 400, boxShadow: '0 20px 60px rgba(0,0,0,0.3)' }}>
          <h2 style={{ textAlign: 'center', marginBottom: 30 }}>无效的链接</h2>
          <p style={{ textAlign: 'center', marginBottom: 24 }}>重置链接无效或已过期</p>
          <Button type="primary" block onClick={() => navigate('/forgot-password')}>
            重新获取重置链接
          </Button>
        </Card>
      </div>
    )
  }

  const onFinish = async (values) => {
    setLoading(true)
    try {
      await resetPassword({
        token: token,
        new_password: values.newPassword
      })
      message.success('密码重置成功，请使用新密码登录')
      navigate('/login')
    } catch (error) {
      console.error('Reset password failed:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ 
      display: 'flex', 
      justifyContent: 'center', 
      alignItems: 'center', 
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
    }}>
      <Card style={{ width: 400, boxShadow: '0 20px 60px rgba(0,0,0,0.3)' }}>
        <h2 style={{ textAlign: 'center', marginBottom: 30 }}>重置密码</h2>
        
        <Form
          name="reset-password"
          onFinish={onFinish}
          autoComplete="off"
          size="large"
        >
          <Form.Item
            name="newPassword"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 6, message: '密码至少 6 个字符' }
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="新密码"
            />
          </Form.Item>

          <Form.Item
            name="confirmPassword"
            dependencies={['newPassword']}
            rules={[
              { required: true, message: '请确认新密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('newPassword') === value) {
                    return Promise.resolve()
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'))
                }
              })
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="确认新密码"
            />
          </Form.Item>

          <Form.Item>
            <Button 
              type="primary" 
              htmlType="submit" 
              loading={loading}
              block
              size="large"
            >
              重置密码
            </Button>
          </Form.Item>
        </Form>

        <div style={{ textAlign: 'center', marginTop: 16 }}>
          <a onClick={() => navigate('/login')}>返回登录</a>
        </div>
      </Card>
    </div>
  )
}

export default ResetPassword
