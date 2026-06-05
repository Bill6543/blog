import { useState, useEffect } from 'react'
import { Table, Button, Tag, Popconfirm, message, Card, Typography, Space, Tooltip } from 'antd'
import { getPendingComments, approveComment, rejectComment } from '../../api/comment'
import { CheckOutlined, CloseOutlined, EyeOutlined } from '@ant-design/icons'

const { Text, Paragraph } = Typography

const AdminComments = () => {
  const [comments, setComments] = useState([])
  const [loading, setLoading] = useState(false)

  const loadPendingComments = async () => {
    setLoading(true)
    try {
      const data = await getPendingComments()
      setComments(data)
    } catch (error) {
      console.error('Failed to load pending comments:', error)
      message.error('获取待审核评论失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadPendingComments()
  }, [])

  const handleApprove = async (id) => {
    try {
      await approveComment(id)
      message.success('评论审核通过')
      loadPendingComments()
    } catch (error) {
      console.error('Approve failed:', error)
      message.error('审核通过失败')
    }
  }

  const handleReject = async (id) => {
    try {
      await rejectComment(id)
      message.success('评论审核拒绝')
      loadPendingComments()
    } catch (error) {
      console.error('Reject failed:', error)
      message.error('审核拒绝失败')
    }
  }

  const getStatusTag = (status) => {
    switch (status) {
      case 1:
        return <Tag color="orange">待审核</Tag>
      case 2:
        return <Tag color="green">已通过</Tag>
      case -1:
        return <Tag color="red">已拒绝</Tag>
      default:
        return <Tag color="default">未知</Tag>
    }
  }

  const columns = [
    { 
      title: 'ID', 
      dataIndex: 'id', 
      key: 'id', 
      width: 60 
    },
    { 
      title: '评论内容', 
      dataIndex: 'content', 
      key: 'content',
      width: 300,
      render: (content) => (
        <Tooltip title={content}>
          <Paragraph 
            ellipsis={{ rows: 2, expandable: false }} 
            style={{ margin: 0, maxWidth: 280 }}
          >
            {content}
          </Paragraph>
        </Tooltip>
      )
    },
    { 
      title: '评论用户', 
      dataIndex: 'user', 
      key: 'user',
      width: 120,
      render: (user) => user?.nickname || user?.username || '未知用户'
    },
    { 
      title: '文章标题', 
      dataIndex: 'article', 
      key: 'article',
      width: 200,
      render: (article) => (
        <Tooltip title={article?.title}>
          <Text ellipsis style={{ maxWidth: 180 }}>
            {article?.title || '未知文章'}
          </Text>
        </Tooltip>
      )
    },
    { 
      title: '状态', 
      dataIndex: 'status', 
      key: 'status',
      width: 100,
      render: (status) => getStatusTag(status)
    },
    { 
      title: '创建时间', 
      dataIndex: 'created_at', 
      key: 'created_at',
      width: 150,
      render: (time) => new Date(time).toLocaleString()
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_, record) => (
        <Space>
          <Tooltip title="审核通过">
            <Button 
              size="small" 
              type="primary" 
              icon={<CheckOutlined />}
              onClick={() => handleApprove(record.id)}
            />
          </Tooltip>
          <Tooltip title="审核拒绝">
            <Popconfirm
              title="确定拒绝此评论？"
              onConfirm={() => handleReject(record.id)}
              okText="确定"
              cancelText="取消"
            >
              <Button 
                size="small" 
                danger 
                icon={<CloseOutlined />}
              />
            </Popconfirm>
          </Tooltip>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Card title="评论审核" extra={
        <Button 
          icon={<EyeOutlined />} 
          onClick={loadPendingComments}
          loading={loading}
        >
          刷新
        </Button>
      }>
        <Table
          columns={columns}
          dataSource={comments}
          loading={loading}
          rowKey="id"
          pagination={{
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条待审核评论`,
          }}
          scroll={{ x: 1000 }}
        />
      </Card>
    </div>
  )
}

export default AdminComments
