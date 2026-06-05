import { useState, useEffect, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Typography, Tag, Avatar, Space, Button, Divider, Spin, Alert, message, Input, Popconfirm } from 'antd'
import { LikeOutlined, EyeOutlined, MessageOutlined, UserOutlined, DeleteOutlined } from '@ant-design/icons'
import MarkdownIt from 'markdown-it'
import { getArticle, likeArticle, unlikeArticle, deleteArticle, getLikeStatus, increaseViewCount } from '../api/article'
import { getCommentList, createComment, deleteComment } from '../api/comment'
import { useAuthStore } from '../store/authStore'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'
import 'github-markdown-css/github-markdown.css'

const { Title, Text } = Typography
const { TextArea } = Input
const md = new MarkdownIt()

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const ArticleDetail = () => {
  const { id } = useParams()
  const navigate = useNavigate()
  const { user, isAuthenticated } = useAuthStore()
  const [article, setArticle] = useState(null)
  const [comments, setComments] = useState([])
  const [loading, setLoading] = useState(false)
  const [liked, setLiked] = useState(false)
  const [commentContent, setCommentContent] = useState('')
  const [replyingTo, setReplyingTo] = useState(null) // { id, userNickname }
  const [replyContent, setReplyContent] = useState('')
  const viewCountIncreasedRef = useRef(false) // 跟踪浏览量是否已增加
  const [commentDeleting, setCommentDeleting] = useState(null) // 跟踪正在删除的评论ID

  // 加载文章详情
  const loadArticle = async () => {
    setLoading(true)
    try {
      const data = await getArticle(id)
      setArticle(data)

      // 首次加载时增加浏览量
      if (!viewCountIncreasedRef.current) {
        try {
          await increaseViewCount(id)
          viewCountIncreasedRef.current = true
        } catch (error) {
          console.error('Failed to increase view count:', error)
          // 浏览量增加失败不影响文章加载
        }
      }

      // 如果已登录，获取点赞状态
      if (isAuthenticated()) {
        const { liked } = await getLikeStatus(id)
        setLiked(liked)
      }
    } catch (error) {
      console.error('Failed to load article:', error)
      message.error('加载文章失败')
    } finally {
      setLoading(false)
    }
  }

  // 加载评论列表
  const loadComments = async () => {
    try {
      const data = await getCommentList(id)
      setComments(data)
    } catch (error) {
      console.error('Failed to load comments:', error)
    }
  }

  useEffect(() => {
    loadArticle()
    loadComments()
  }, [id])

  // 处理点赞
  const handleLike = async () => {
    if (!isAuthenticated()) {
      message.warning('请先登录')
      return
    }

    try {
      if (liked) {
        await unlikeArticle(id)
        setLiked(false)
        message.success('已取消点赞')
      } else {
        await likeArticle(id)
        setLiked(true)
        message.success('点赞成功')
      }
      loadArticle()
    } catch (error) {
      console.error('Like operation failed:', error)
    }
  }

  // 处理删除文章
  const handleDelete = async () => {
    try {
      await deleteArticle(id)
      message.success('文章已删除')
      navigate('/')
    } catch (error) {
      console.error('Delete failed:', error)
    }
  }

  // 发表评论
  const submitComment = async () => {
    if (!commentContent.trim()) {
      message.warning('请输入评论内容')
      return
    }

    try {
      await createComment({
        article_id: parseInt(id),
        content: commentContent,
      })
      message.success('评论成功')
      setCommentContent('')
      await loadComments() // 确保等待评论列表重新加载
      await loadArticle()  // 确保等待文章详情（更新评论数）
    } catch (error) {
      console.error('Comment failed:', error)
    }
  }

  // 发表回复
  const submitReply = async () => {
    if (!replyContent.trim()) {
      message.warning('请输入回复内容')
      return
    }

    try {
      await createComment({
        article_id: parseInt(id),
        content: replyContent,
        parent_id: replyingTo.id,
      })
      message.success('回复成功')
      setReplyContent('')
      setReplyingTo(null)
      await loadComments()
      await loadArticle()
    } catch (error) {
      console.error('Reply failed:', error)
    }
  }

  // 删除评论
  const handleDeleteComment = async (commentId) => {
    try {
      setCommentDeleting(commentId)
      await deleteComment(commentId)
      message.success('评论已删除')
      // 重新加载评论和文章信息
      await loadComments()
      await loadArticle()
    } catch (error) {
      console.error('Delete comment failed:', error)
      message.error('删除评论失败')
    } finally {
      setCommentDeleting(null)
    }
  }

  // 取消回复
  const cancelReply = () => {
    setReplyingTo(null)
    setReplyContent('')
  }

  if (loading) {
    return <Spin size="large" style={{ display: 'flex', justifyContent: 'center', marginTop: '100px' }} />
  }

  if (!article && !loading) {
    return (
      <Card style={{ textAlign: 'center', marginTop: 50 }}>
        <Alert message="文章不存在或加载失败" type="warning" showIcon />
        <Button onClick={() => navigate('/')} style={{ marginTop: 20 }}>返回首页</Button>
      </Card>
    )
  }

  const isAuthor = user?.id === article.author_id
  const isAdmin = user?.role === 'admin'

  return (
    <div style={{ maxWidth: 900, margin: '0 auto' }}>
      {/* 文章详情 */}
      <Card>
        <Space direction="vertical" style={{ width: '100%' }} size="large">
          <Title level={2}>{article.title}</Title>

          <Space split={<span>|</span>} wrap>
            <Avatar src={article.author?.avatar} icon={<UserOutlined />} />
            <span>{article.author?.nickname || article.author?.username || '未知作者'}</span>
            <span><EyeOutlined /> {article.view_count}</span>
            <span><LikeOutlined /> {article.like_count}</span>
            <span><MessageOutlined /> {article.comment_count}</span>
            <span>{dayjs(article.created_at).format('YYYY-MM-DD HH:mm')}</span>
          </Space>

          <Space wrap>
            {article.category?.name && <Tag color="green">{article.category.name}</Tag>}
            {Array.isArray(article.tags) && article.tags.map(tag => (
              <Tag key={tag.id} color="blue">{tag.name}</Tag>
            ))}
          </Space>

          {article.cover_image && (
            <img 
              src={article.cover_image} 
              alt={article.title}
              style={{ width: '100%', maxHeight: 400, objectFit: 'cover', borderRadius: 8 }}
            />
          )}

          <Divider />

          {/* 文章内容（Markdown 渲染） */}
          <div 
            className="markdown-body"
            dangerouslySetInnerHTML={{ __html: md.render(article.content || '') }}
            style={{ lineHeight: 1.8, fontSize: 16, padding: '20px 0' }}
          />

          {/* 操作按钮 */}
          <Divider />
          <Space>
            <Button 
              type={liked ? 'primary' : 'default'} 
              icon={<LikeOutlined />}
              onClick={handleLike}
            >
              {liked ? '已点赞' : '点赞'}
            </Button>
            
            {(isAuthor || isAdmin) && (
              <>
                <Button onClick={() => navigate(`/edit/${id}`)}>编辑</Button>
                <Button danger onClick={handleDelete}>删除</Button>
              </>
            )}
          </Space>
        </Space>
      </Card>

      {/* 评论区 */}
      <Card title={`评论 (${article.comment_count})`} style={{ marginTop: 24 }}>
        <Space direction="vertical" style={{ width: '100%' }}>
          {/* 发表评论 */}
          {isAuthenticated() ? (
            <div>
              <TextArea
                value={commentContent}
                onChange={(e) => setCommentContent(e.target.value)}
                placeholder="写下你的评论..."
                rows={4}
                style={{ marginBottom: 8 }}
              />
              <Button type="primary" onClick={submitComment}>发表评论</Button>
            </div>
          ) : (
            <Alert message="请登录后发表评论" type="info" />
          )}

          {/* 评论列表 */}
          <Divider />
          {Array.isArray(comments) && comments.length > 0 ? (
            comments.map(comment => (
              <div key={comment.id} style={{ marginBottom: 16 }}>
                {/* 父评论 */}
                <div style={{ padding: '16px 0', borderBottom: '1px solid #f0f0f0' }}>
                  <Space align="start">
                    <Avatar src={comment.user?.avatar} />
                    <div style={{ flex: 1 }}>
                      <div style={{ fontWeight: 500 }}>{comment.user?.nickname}</div>
                      <div style={{ marginTop: 8 }}>{comment.content}</div>
                      <div style={{ marginTop: 8, fontSize: 12, color: '#999' }}>
                        {dayjs(comment.created_at).fromNow()}
                        {comment.replies && comment.replies.length > 0 && (
                          <span style={{ marginLeft: 8, color: '#666', padding: '2px 6px', background: '#f5f5f5', borderRadius: 3 }}>
                            {comment.replies.length} 条回复
                          </span>
                        )}
                        {isAuthenticated() && (
                          <Button 
                            type="link" 
                            size="small" 
                            style={{ marginLeft: 8, padding: 0 }}
                            onClick={() => setReplyingTo({ id: comment.id, userNickname: comment.user?.nickname })}
                          >
                            回复
                          </Button>
                        )}
                        {(user?.id === comment.user_id || user?.role === 'admin') && (
                          <Popconfirm
                            title="确定删除此评论？"
                            description="删除后无法恢复"
                            onConfirm={() => handleDeleteComment(comment.id)}
                            okText="确定"
                            cancelText="取消"
                          >
                            <Button 
                              type="link" 
                              size="small" 
                              danger
                              loading={commentDeleting === comment.id}
                              style={{ marginLeft: 8, padding: 0 }}
                              icon={<DeleteOutlined />}
                            >
                              删除
                            </Button>
                          </Popconfirm>
                        )}
                      </div>
                    </div>
                  </Space>
                </div>

                {/* 回复输入框（仅当正在回复此评论时显示） */}
                {replyingTo && replyingTo.id === comment.id && (
                  <div style={{ marginLeft: 40, marginTop: 12 }}>
                    <div style={{ fontSize: 12, color: '#666', marginBottom: 8 }}>
                      回复 @{replyingTo.userNickname}
                    </div>
                    <TextArea
                      value={replyContent}
                      onChange={(e) => setReplyContent(e.target.value)}
                      placeholder="写下你的回复..."
                      rows={3}
                      style={{ marginBottom: 8 }}
                    />
                    <Space>
                      <Button type="primary" size="small" onClick={submitReply}>发表回复</Button>
                      <Button size="small" onClick={cancelReply}>取消</Button>
                    </Space>
                  </div>
                )}

                {/* 子评论（回复） */}
                {comment.replies && comment.replies.length > 0 && (
                  <div style={{ marginLeft: 40, marginTop: 12 }}>
                    {/* 子评论列表 */}
                    <div style={{ 
                      background: '#fafafa', 
                      borderRadius: 8, 
                      padding: '12px 16px',
                      borderLeft: '3px solid #1890ff'
                    }}>
                      {comment.replies.map(reply => (
                        <div key={reply.id} style={{ padding: '10px 0', borderBottom: '1px solid #e8e8e8' }}>
                          {/* 父评论引用 */}
                          {reply.parent && (
                            <div style={{ 
                              marginBottom: 8, 
                              padding: '6px 10px', 
                              background: '#f0f0f0', 
                              borderRadius: 4,
                              fontSize: 12,
                              color: '#666',
                              borderLeft: '2px solid #d9d9d9'
                            }}>
                              <span style={{ color: '#999', marginRight: 4 }}>
                                @{reply.parent.user?.nickname}:
                              </span>
                              {reply.parent.content.length > 50 
                                ? reply.parent.content.substring(0, 50) + '...' 
                                : reply.parent.content}
                            </div>
                          )}
                          
                          <Space align="start">
                            <Avatar src={reply.user?.avatar} size="small" />
                            <div style={{ flex: 1 }}>
                              <div style={{ fontWeight: 500, fontSize: 13 }}>
                                <span style={{ color: '#1890ff' }}>{reply.user?.nickname}</span>
                                {reply.parent && (
                                  <span style={{ color: '#999', marginLeft: 4 }}>
                                    <MessageOutlined style={{ fontSize: 11, marginRight: 2 }} />
                                    回复
                                  </span>
                                )}
                              </div>
                              <div style={{ marginTop: 6, fontSize: 14, color: '#333' }}>{reply.content}</div>
                              <div style={{ marginTop: 6, fontSize: 12, color: '#999' }}>
                                {dayjs(reply.created_at).fromNow()}
                                {isAuthenticated() && (
                                  <Button 
                                    type="link" 
                                    size="small" 
                                    style={{ marginLeft: 8, padding: 0, fontSize: 12 }}
                                    onClick={() => setReplyingTo({ id: comment.id, userNickname: reply.user?.nickname })}
                                  >
                                    回复
                                  </Button>
                                )}
                                {(user?.id === reply.user_id || user?.role === 'admin') && (
                                  <Popconfirm
                                    title="确定删除此回复？"
                                    description="删除后无法恢复"
                                    onConfirm={() => handleDeleteComment(reply.id)}
                                    okText="确定"
                                    cancelText="取消"
                                  >
                                    <Button 
                                      type="link" 
                                      size="small" 
                                      danger
                                      loading={commentDeleting === reply.id}
                                      style={{ marginLeft: 8, padding: 0, fontSize: 12 }}
                                      icon={<DeleteOutlined />}
                                    >
                                      删除
                                    </Button>
                                  </Popconfirm>
                                )}
                              </div>
                            </div>
                          </Space>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ))
          ) : (
            <div style={{ textAlign: 'center', padding: 20, color: '#999' }}>暂无评论</div>
          )}
        </Space>
      </Card>
    </div>
  )
}

export default ArticleDetail
