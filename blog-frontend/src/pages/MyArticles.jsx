import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, List, Tag, Avatar, Space, Typography, Pagination, Empty, Spin, Button, message, Popconfirm } from 'antd'
import { EyeOutlined, LikeOutlined, MessageOutlined, UserOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { getUserArticles, restoreArticle, deleteArticle } from '../api/article'
import { useAuthStore } from '../store/authStore'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

const { Title, Paragraph } = Typography

// 配置 dayjs
dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const MyArticles = () => {
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const [articles, setArticles] = useState([])
  const [loading, setLoading] = useState(false)
  const [statusFilter, setStatusFilter] = useState('all') // 'all', '1', '0', '-1'
  const [pagination, setPagination] = useState({
    page: 1,
    pageSize: 10,
    total: 0,
  })

  // 加载文章列表
  const loadArticles = async (page = 1, pageSize = 10, status = 'all') => {
    if (!user?.id) return
    
    setLoading(true)
    try {
      const params = { page, page_size: pageSize }
      if (status !== 'all') {
        params.status = status
      }
      const data = await getUserArticles(user.id, params)
      setArticles(Array.isArray(data.list) ? data.list : [])
      setPagination(prev => ({
        ...prev,
        page,
        pageSize,
        total: data.total || 0,
      }))
    } catch (error) {
      console.error('Failed to load articles:', error)
      message.error('加载文章列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadArticles()
  }, [user])

  const handlePageChange = (page, pageSize) => {
    loadArticles(page, pageSize, statusFilter)
  }

  const handleStatusFilter = (status) => {
    setStatusFilter(status)
    loadArticles(1, pagination.pageSize, status)
  }

  const handleEdit = (articleId) => {
    navigate(`/edit/${articleId}`)
  }

  const handleDelete = async (articleId, isDeleted) => {
    try {
      if (isDeleted) {
        // 恢复文章
        await restoreArticle(articleId)
        message.success('恢复成功')
      } else {
        // 删除文章
        await deleteArticle(articleId)
        message.success('删除成功')
      }
      loadArticles(pagination.page, pagination.pageSize, statusFilter)
    } catch (error) {
      message.error(isDeleted ? '恢复失败' : '删除失败')
    }
  }

  const getStatusTag = (status) => {
    switch (status) {
      case 1:
        return <Tag color="green">已发布</Tag>
      case 0:
        return <Tag color="orange">草稿</Tag>
      case -1:
        return <Tag color="red">已删除</Tag>
      default:
        return <Tag color="default">未知</Tag>
    }
  }

  return (
    <div style={{ maxWidth: 1200, margin: '0 auto' }}>
      <Space direction="vertical" style={{ width: '100%' }} size="large">
        {/* 状态筛选 */}
        <Card>
          <Space direction="vertical" style={{ width: '100%' }}>
            <Title level={5}>文章状态</Title>
            <Space wrap>
              <Tag 
                color={statusFilter === 'all' ? 'blue' : 'default'} 
                style={{ cursor: 'pointer' }}
                onClick={() => handleStatusFilter('all')}
              >
                全部文章
              </Tag>
              <Tag 
                color={statusFilter === '1' ? 'blue' : 'default'} 
                style={{ cursor: 'pointer' }}
                onClick={() => handleStatusFilter('1')}
              >
                已发布
              </Tag>
              <Tag 
                color={statusFilter === '0' ? 'blue' : 'default'} 
                style={{ cursor: 'pointer' }}
                onClick={() => handleStatusFilter('0')}
              >
                草稿
              </Tag>
              <Tag 
                color={statusFilter === '-1' ? 'blue' : 'default'} 
                style={{ cursor: 'pointer' }}
                onClick={() => handleStatusFilter('-1')}
              >
                已删除
              </Tag>
            </Space>
          </Space>
        </Card>

        {/* 文章列表 */}
        <Spin spinning={loading}>
          {articles.length > 0 ? (
            <>
              <List
                itemLayout="vertical"
                dataSource={articles}
                renderItem={(article) => (
                  <List.Item
                    actions={[
                      article.status !== -1 && (
                        <Button 
                          key="edit" 
                          type="link" 
                          icon={<EditOutlined />}
                          onClick={() => handleEdit(article.id)}
                        >
                          编辑
                        </Button>
                      ),
                      article.status !== -1 ? (
                        <Popconfirm
                          key="delete"
                          title="确定删除此文章？"
                          description="删除后文章将移至回收站，可以恢复"
                          onConfirm={() => handleDelete(article.id, false)}
                          okText="确定"
                          cancelText="取消"
                        >
                          <Button
                            type="link"
                            danger
                            icon={<DeleteOutlined />}
                          >
                            删除
                          </Button>
                        </Popconfirm>
                      ) : (
                        <Button
                          key="restore"
                          type="link"
                          icon={<DeleteOutlined />}
                          onClick={() => handleDelete(article.id, true)}
                        >
                          恢复
                        </Button>
                      )
                    ].filter(Boolean)}
                  >
                    <Card 
                      hoverable
                      style={{ width: '100%' }}
                      cover={article.cover_image && (
                        <img 
                          alt={article.title} 
                          src={article.cover_image}
                          style={{ height: 200, objectFit: 'cover' }}
                        />
                      )}
                    >
                      <Card.Meta
                        title={
                          <div>
                            <Title level={4} style={{ marginBottom: 4 }}>
                              {article.title}
                            </Title>
                            <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
                              <span style={{ color: '#999', fontSize: '12px' }}>
                                {dayjs(article.created_at).fromNow()}
                              </span>
                              {getStatusTag(article.status)}
                            </div>
                          </div>
                        }
                        description={
                          <Space direction="vertical" style={{ width: '100%' }}>
                            <Space wrap>
                              {article.category && <Tag color="green">{article.category.name}</Tag>}
                              {article.tags && article.tags.length > 0 && article.tags.map(tag => (
                                <Tag key={tag.id} color="blue">{tag.name}</Tag>
                              ))}
                            </Space>
                            <Paragraph 
                              ellipsis={{ rows: 2 }}
                              style={{ marginBottom: 8, color: '#666' }}
                            >
                              {article.summary || (article.content ? article.content.substring(0, 200) : '暂无摘要')}
                            </Paragraph>
                            <Space split={<span>|</span>}>
                              <span>
                                <EyeOutlined /> {article.view_count}
                              </span>
                              <span>
                                <LikeOutlined /> {article.like_count}
                              </span>
                              <span>
                                <MessageOutlined /> {article.comment_count}
                              </span>
                            </Space>
                          </Space>
                        }
                      />
                    </Card>
                  </List.Item>
                )}
              />

              {/* 分页 */}
              <Pagination
                current={pagination.page}
                pageSize={pagination.pageSize}
                total={pagination.total}
                onChange={handlePageChange}
                showSizeChanger
                showTotal={(total) => `共 ${total} 篇文章`}
                style={{ textAlign: 'center', marginTop: 24 }}
              />
            </>
          ) : (
            !loading && <Empty description="暂无文章" />
          )}
        </Spin>
      </Space>
    </div>
  )
}

export default MyArticles
