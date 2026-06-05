import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, List, Tag, Avatar, Space, Typography, Pagination, Empty, Spin } from 'antd'
import { EyeOutlined, LikeOutlined, MessageOutlined, UserOutlined } from '@ant-design/icons'
import { getArticleList } from '../api/article'
import { getTagList } from '../api/tag'
import { getCategoryList } from '../api/category'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

const { Title, Paragraph } = Typography

// 配置 dayjs
dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const Home = () => {
  const navigate = useNavigate()
  const [articles, setArticles] = useState([])
  const [tags, setTags] = useState([])
  const [categories, setCategories] = useState([])
  const [loading, setLoading] = useState(false)
  const [displayMode, setDisplayMode] = useState('all') // 'all', 'category', 'tag'
  const [pagination, setPagination] = useState({
    page: 1,
    pageSize: 10,
    total: 0,
  })

  // 加载文章列表
  const loadArticles = async (page = 1, pageSize = 10, categoryId = null, tagId = null) => {
    setLoading(true)
    try {
      const params = { page, page_size: pageSize }
      if (categoryId) {
        params.category_id = categoryId
      }
      if (tagId) {
        params.tag_id = tagId
      }
      const data = await getArticleList(params)
      setArticles(Array.isArray(data.list) ? data.list : [])
      setPagination(prev => ({
        ...prev,
        page,
        pageSize,
        total: data.total || 0,
        categoryId,
        tagId
      }))
    } catch (error) {
      console.error('Failed to load articles:', error)
    } finally {
      setLoading(false)
    }
  }

  // 加载标签列表
  const loadTags = async () => {
    try {
      const data = await getTagList()
      setTags(Array.isArray(data) ? data : [])
    } catch (error) {
      console.error('Failed to load tags:', error)
    }
  }

  // 加载分类列表
  const loadCategories = async () => {
    try {
      const data = await getCategoryList()
      setCategories(Array.isArray(data) ? data : [])
    } catch (error) {
      console.error('Failed to load categories:', error)
    }
  }

  useEffect(() => {
    loadArticles()
    loadTags()
    loadCategories()
  }, [])


  const handlePageChange = (page, pageSize) => {
    loadArticles(page, pageSize, pagination.categoryId, pagination.tagId)
  }

  return (
    <div style={{ maxWidth: 1200, margin: '0 auto' }}>
      <Space direction="vertical" style={{ width: '100%' }} size="large">
        {/* 显示模式选择 */}
        <Card>
          <Space direction="vertical" style={{ width: '100%' }}>
            <Title level={5}>显示方式</Title>
            <Space wrap>
              <Tag 
                color={displayMode === 'all' ? 'blue' : 'default'} 
                style={{ cursor: 'pointer' }}
                onClick={() => {
                  setDisplayMode('all')
                  loadArticles(1, pagination.pageSize, null, null)
                }}
              >
                全部文章
              </Tag>
              <Tag 
                color={displayMode === 'category' ? 'blue' : 'default'} 
                style={{ cursor: 'pointer' }}
                onClick={() => {
                  setDisplayMode('category')
                  loadArticles(1, pagination.pageSize, null, null)
                }}
              >
                按分类
              </Tag>
              <Tag 
                color={displayMode === 'tag' ? 'blue' : 'default'} 
                style={{ cursor: 'pointer' }}
                onClick={() => {
                  setDisplayMode('tag')
                  loadArticles(1, pagination.pageSize, null, null)
                }}
              >
                按标签
              </Tag>
            </Space>
            
            {/* 分类过滤 */}
            {displayMode === 'category' && categories.length > 0 && (
              <div style={{ marginTop: 16 }}>
                <Title level={5}>选择分类</Title>
                <Space wrap>
                  <Tag 
                    color={!pagination.categoryId ? 'green' : 'default'} 
                    style={{ cursor: 'pointer' }}
                    onClick={() => loadArticles(1, pagination.pageSize, null, null)}
                  >
                    全部分类
                  </Tag>
                  {categories.map(category => (
                    <Tag 
                      key={category.id} 
                      color={pagination.categoryId === category.id ? 'green' : 'default'} 
                      style={{ cursor: 'pointer' }}
                      onClick={() => loadArticles(1, pagination.pageSize, category.id, null)}
                    >
                      {category.name}
                    </Tag>
                  ))}
                </Space>
              </div>
            )}
            
            {/* 标签过滤 */}
            {displayMode === 'tag' && tags.length > 0 && (
              <div style={{ marginTop: 16 }}>
                <Title level={5}>选择标签</Title>
                <Space wrap>
                  <Tag 
                    color={!pagination.tagId ? 'green' : 'default'} 
                    style={{ cursor: 'pointer' }}
                    onClick={() => loadArticles(1, pagination.pageSize, null, null)}
                  >
                    全部标签
                  </Tag>
                  {tags.map(tag => (
                    <Tag 
                      key={tag.id} 
                      color={pagination.tagId === tag.id ? 'green' : 'default'} 
                      style={{ cursor: 'pointer' }}
                      onClick={() => loadArticles(1, pagination.pageSize, null, tag.id)}
                    >
                      {tag.name}
                    </Tag>
                  ))}
                </Space>
              </div>
            )}
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
                    onClick={() => navigate(`/article/${article.id}`)}
                    style={{ cursor: 'pointer' }}
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
                        avatar={<Avatar src={article.author?.avatar} icon={!article.author?.avatar && <UserOutlined />} />}
                        title={
                          <div>
                            <Title level={4} style={{ marginBottom: 4 }}>
                              {article.title}
                            </Title>
                            <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
                              <span style={{ color: '#666', fontSize: '14px' }}>
                                {article.author?.nickname || article.author?.username || '匿名用户'}
                              </span>
                              <span style={{ color: '#999', fontSize: '12px' }}>
                                {dayjs(article.created_at).fromNow()}
                              </span>
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

export default Home
