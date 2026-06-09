import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { Form, Input, Button, Card, message, Alert, Space } from 'antd'
import { MailOutlined, CopyOutlined, ArrowRightOutlined } from '@ant-design/icons'
import { forgotPassword } from '../api/auth'

const ForgotPassword = () => {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [submitted, setSubmitted] = useState(false)
  const [resetUrl, setResetUrl] = useState('')

  const onFinish = async (values) => {
    setLoading(true)
    try {
      const data = await forgotPassword(values)
      // 处理邮箱不存在的情况（后端返回 data=null）
      if (data && data.reset_url) {
        setResetUrl(data.reset_url)
      }
      setSubmitted(true)
      message.success('如果该邮箱已注册，重置链接将发送到您的邮箱')
    } catch (error) {
      console.error('Forgot password failed:', error)
      message.error('操作失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  // 复制链接到剪贴板
  const handleCopy = () => {
    navigator.clipboard.writeText(resetUrl).then(() => {
      message.success('链接已复制到剪贴板')
    }).catch(() => {
      message.error('复制失败，请手动复制')
    })
  }

  // 跳转到重置页面
  const handleGoToReset = () => {
    window.location.href = resetUrl
  }

  return (
    <div style={{ 
      display: 'flex', 
      justifyContent: 'center', 
      alignItems: 'center', 
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
    }}>
      <Card style={{ width: 500, boxShadow: '0 20px 60px rgba(0,0,0,0.3)' }}>
        <h2 style={{ textAlign: 'center', marginBottom: 30 }}>忘记密码</h2>
        
        {submitted ? (
          <div>
            <Alert
              message="✅ 请求已提交"
              description="如果该邮箱已注册，重置链接将发送到您的邮箱。"
              type="success"
              showIcon
              style={{ marginBottom: 24 }}
            />
            
            {resetUrl && (
              <div style={{ marginBottom: 24 }}>
                <p style={{ marginBottom: 8, fontWeight: 'bold' }}>重置链接：</p>
                <div style={{ 
                  padding: '12px', 
                  background: '#f5f5f5', 
                  borderRadius: '4px',
                  wordBreak: 'break-all',
                  fontSize: '12px',
                  color: '#666',
                  marginBottom: 16
                }}>
                  {resetUrl}
                </div>
                
                <Space direction="vertical" style={{ width: '100%' }}>
                  <Button 
                    type="primary" 
                    size="large" 
                    block
                    icon={<ArrowRightOutlined />}
                    onClick={handleGoToReset}
                  >
                    立即重置密码
                  </Button>
                  <Button 
                    size="large" 
                    block
                    icon={<CopyOutlined />}
                    onClick={handleCopy}
                  >
                    复制链接
                  </Button>
                </Space>
              </div>
            )}

            <Alert
              message="安全提示"
              description={
                <ul style={{ margin: 0, paddingLeft: 20 }}>
                  <li>链接将在 1 小时后自动失效</li>
                  <li>重置成功后链接将自动失效</li>
                  <li>如未收到邮件，请检查垃圾邮件文件夹</li>
                </ul>
              }
              type="warning"
              showIcon
              style={{ marginBottom: 16 }}
            />
          </div>
        ) : (
          <Form
            name="forgot-password"
            onFinish={onFinish}
            autoComplete="off"
            size="large"
          >
            <Form.Item
              name="email"
              rules={[
                { required: true, message: '请输入邮箱' },
                { type: 'email', message: '请输入有效的邮箱地址' }
              ]}
            >
              <Input 
                prefix={<MailOutlined />} 
                placeholder="请输入注册邮箱"
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
                发送重置链接
              </Button>
            </Form.Item>
          </Form>
        )}

        <div style={{ textAlign: 'center', marginTop: 16 }}>
          <Link to="/login">返回登录</Link>
        </div>
      </Card>
    </div>
  )
}

export default ForgotPassword
