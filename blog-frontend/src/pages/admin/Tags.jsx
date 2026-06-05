import { useState, useEffect } from 'react'
import { Table, Button, Popconfirm, message, Modal, Form, Input } from 'antd'
import { getTagList, createTag, updateTag, deleteTag } from '../../api/tag'

const AdminTags = () => {
  const [tags, setTags] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [currentTag, setCurrentTag] = useState(null)
  const [form] = Form.useForm()

  const loadTags = async () => {
    setLoading(true)
    try {
      const response = await getTagList()
      console.log('API Response:', response)
      
      // Response interceptor returns data directly (array), not wrapped in response object
      const tagData = Array.isArray(response) ? response : []
      console.log('Extracted tag data:', tagData)
      setTags(tagData)
      console.log('Tags set:', tagData)
    } catch (error) {
      console.error('Failed to load tags:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadTags()
  }, [])

  const handleCreate = () => {
    setCurrentTag(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (tag) => {
    setCurrentTag(tag)
    form.setFieldsValue({
      name: tag.name,
      color: tag.color,
    })
    setModalVisible(true)
  }

  const handleSubmit = async (values) => {
    try {
      if (currentTag) {
        await updateTag(currentTag.id, values)
        message.success('更新成功')
      } else {
        await createTag(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      loadTags()
    } catch (error) {
      console.error('Operation failed:', error)
    }
  }

  const handleDelete = async (id) => {
    try {
      await deleteTag(id)
      message.success('删除成功')
      loadTags()
    } catch (error) {
      console.error('Delete failed:', error)
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { 
      title: '名称', 
      dataIndex: 'name', 
      key: 'name',
      render: (text, record) => (
        <span style={{ 
          padding: '4px 12px', 
          borderRadius: 4, 
          backgroundColor: record.color || '#1890ff',
          color: '#fff'
        }}>
          {text}
        </span>
      )
    },
    { title: '颜色', dataIndex: 'color', key: 'color', render: (c) => c || '默认' },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <>
          <Button size="small" onClick={() => handleEdit(record)}>编辑</Button>
          <Popconfirm
            title="确定删除此标签？"
            onConfirm={() => handleDelete(record.id)}
          >
            <Button size="small" danger style={{ marginLeft: 8 }}>删除</Button>
          </Popconfirm>
        </>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" onClick={handleCreate}>创建标签</Button>
      </div>
      
      <Table
        columns={columns}
        dataSource={tags}
        loading={loading}
        rowKey="id"
      />

      <Modal
        title={currentTag ? '编辑标签' : '创建标签'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item 
            name="name" 
            label="名称"
            rules={[{ required: true, message: '请输入标签名称' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item name="color" label="颜色">
            <Input placeholder="如：#1890ff" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit">保存</Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default AdminTags
