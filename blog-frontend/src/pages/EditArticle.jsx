import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Form, Input, Button, Card, Select, Upload, message, Space, Spin } from 'antd'
import { UploadOutlined, RobotOutlined } from '@ant-design/icons'
import { getArticle, updateArticle } from '../api/article'
import { getCategoryList } from '../api/category'
import { getTagList } from '../api/tag'
import { generateSummary, generateCover } from '../api/ai'

const { TextArea } = Input

const EditArticle = () => {
  const { id } = useParams()
  const navigate = useNavigate()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [initialLoading, setInitialLoading] = useState(true)
  const [categories, setCategories] = useState([])
  const [tags, setTags] = useState([])
  const [coverFile, setCoverFile] = useState(null)
  const [article, setArticle] = useState(null)
  const [aiLoading, setAiLoading] = useState(false)

  useEffect(() => {
    loadArticle()
    loadCategories()
    loadTags()
  }, [id])

  const loadArticle = async () => {
    try {
      const data = await getArticle(id)
      setArticle(data)
      
      // 填充表单数据
      form.setFieldsValue({
        title: data.title,
        summary: data.summary,
        content: data.content,
        category_id: data.category_id,
        tag_ids: data.tags ? data.tags.map(tag => tag.id) : [],
        status: data.status !== undefined ? data.status : 1
      })
    } catch (error) {
      console.error('Failed to load article:', error)
      message.error('加载文章失败')
      navigate('/')
    } finally {
      setInitialLoading(false)
    }
  }

  const loadCategories = async () => {
    try {
      const data = await getCategoryList()
      setCategories(data)
    } catch (error) {
      console.error('Failed to load categories:', error)
    }
  }

  const loadTags = async () => {
    try {
      const data = await getTagList()
      setTags(data)
    } catch (error) {
      console.error('Failed to load tags:', error)
    }
  }

  const handleGenerateSummary = async () => {
    const content = form.getFieldValue('content')
    if (!content || content.trim().length === 0) {
      message.warning('请先输入文章内容')
      return
    }

    setAiLoading(true)
    try {
      const result = await generateSummary({ content })
      form.setFieldValue('summary', result.summary)
      message.success('摘要生成成功')
    } catch (error) {
      console.error('Generate summary failed:', error)
      message.error('摘要生成失败')
    } finally {
      setAiLoading(false)
    }
  }

  const handleGenerateCover = async () => {
    const content = form.getFieldValue('content')
    if (!content || content.trim().length === 0) {
      message.warning('请先输入文章内容')
      return
    }

    setAiLoading(true)
    try {
      const result = await generateCover({ content })
      // 将返回的封面URL转换为文件对象
      const response = await fetch(result.cover_url)
      const blob = await response.blob()
      const file = new File([blob], 'ai-generated-cover.jpg', { type: 'image/jpeg' })
      setCoverFile(file)
      message.success('封面生成成功')
    } catch (error) {
      console.error('Generate cover failed:', error)
      message.error('封面生成失败')
    } finally {
      setAiLoading(false)
    }
  }

  const onFinish = async (values) => {
    setLoading(true)
    try {
      const formData = new FormData()
      formData.append('title', values.title)
      formData.append('content', values.content)
      if (values.summary) formData.append('summary', values.summary)
      if (values.category_id && values.category_id !== '') formData.append('category_id', String(values.category_id))
      if (values.status !== undefined && values.status !== null) formData.append('status', values.status)

      // Always send tag_ids as array (empty if no tags selected)
      const tagIds = values.tag_ids || []
      if (tagIds.length > 0) {
        tagIds.forEach(tagId => {
          formData.append('tag_ids', String(tagId))
        })
      }
      if (coverFile) {
        formData.append('cover_image_file', coverFile)
      }

      await updateArticle(id, formData)
      message.success('文章更新成功')
      navigate(`/article/${id}`)
    } catch (error) {
      console.error('Update article failed:', error)
      message.error('更新文章失败')
    } finally {
      setLoading(false)
    }
  }

  const uploadProps = {
    beforeUpload: (file) => {
      const isImage = file.type.startsWith('image/')
      if (!isImage) {
        message.error('只能上传图片文件')
        return false
      }
      const isLt5M = file.size / 1024 / 1024 < 5
      if (!isLt5M) {
        message.error('图片大小不能超过 5MB')
        return false
      }
      setCoverFile(file)
      return false
    },
    onRemove: () => {
      setCoverFile(null)
    },
    file: coverFile
  }

  if (initialLoading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', marginTop: '100px' }}>
        <Spin size="large" />
      </div>
    )
  }

  return (
    <div style={{ maxWidth: 900, margin: '0 auto' }}>
      <Card title="编辑文章">
        <Form form={form} layout="vertical" onFinish={onFinish}>
          <Form.Item
            name="title"
            label="标题"
            rules={[{ required: true, message: '请输入文章标题' }]}
          >
            <Input placeholder="请输入标题" maxLength={200} showCount />
          </Form.Item>

          <Form.Item
            name="summary"
            label={
              <Space>
                摘要
                <Button
                  type="link"
                  size="small"
                  icon={<RobotOutlined />}
                  onClick={handleGenerateSummary}
                  loading={aiLoading}
                >
                  AI生成
                </Button>
              </Space>
            }
          >
            <TextArea
              placeholder="请输入文章摘要（可选）"
              maxLength={500}
              showCount
              rows={3}
            />
          </Form.Item>

          <Form.Item
            name="content"
            label="内容"
            rules={[{ required: true, message: '请输入文章内容' }]}
          >
            <TextArea 
              placeholder="请输入文章内容" 
              rows={15}
            />
          </Form.Item>

          <Form.Item
            name="category_id"
            label="分类"
          >
            <Select placeholder="请选择分类" allowClear>
              {categories.map(cat => (
                <Select.Option key={cat.id || `cat-${Math.random()}`} value={cat.id}>{cat.name}</Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="tag_ids"
            label="标签"
          >
            <Select mode="multiple" placeholder="请选择标签" allowClear>
              {tags.map(tag => (
                <Select.Option key={tag.id || `tag-${Math.random()}`} value={tag.id}>{tag.name}</Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="status"
            label="状态"
          >
            <Select placeholder="请选择状态">
              <Select.Option value={1}>已发布</Select.Option>
              <Select.Option value={0}>草稿</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item label="封面图（可选）">
            <Space>
              <Upload {...uploadProps} showUploadList={false}>
                <Button icon={<UploadOutlined />}>更换封面</Button>
              </Upload>
              <Button
                icon={<RobotOutlined />}
                onClick={handleGenerateCover}
                loading={aiLoading}
              >
                AI生成封面
              </Button>
            </Space>
            {coverFile ? (
              <div style={{ marginTop: 8, display: 'flex', alignItems: 'center', gap: 8 }}>
                <span style={{ color: '#52c41a' }}>已选择：{coverFile.name}</span>
                <Button size="small" danger onClick={() => setCoverFile(null)}>
                  清除封面
                </Button>
              </div>
            ) : article?.cover_image ? (
              <div style={{ marginTop: 8 }}>
                <span style={{ color: '#999' }}>当前有封面图片，可重新选择更换</span>
              </div>
            ) : (
              <div style={{ marginTop: 8, color: '#999' }}>
                将使用默认封面
              </div>
            )}
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" loading={loading}>
                更新
              </Button>
              <Button onClick={() => navigate(`/article/${id}`)}>取消</Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default EditArticle
