import { useState, useEffect } from 'react'
import { Table, Button, Popconfirm, message, Modal, Form, Input } from 'antd'
import { getCategoryList, createCategory, updateCategory, deleteCategory } from '../../api/category'

const AdminCategories = () => {
  const [categories, setCategories] = useState([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [currentCategory, setCurrentCategory] = useState(null)
  const [form] = Form.useForm()

  const loadCategories = async () => {
    setLoading(true)
    try {
      const data = await getCategoryList()
      setCategories(data)
    } catch (error) {
      console.error('Failed to load categories:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadCategories()
  }, [])

  const handleCreate = () => {
    setCurrentCategory(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (category) => {
    setCurrentCategory(category)
    form.setFieldsValue({
      name: category.name,
      description: category.description,
    })
    setModalVisible(true)
  }

  const handleSubmit = async (values) => {
    try {
      if (currentCategory) {
        await updateCategory(currentCategory.id, values)
        message.success('更新成功')
      } else {
        await createCategory(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      loadCategories()
    } catch (error) {
      console.error('Operation failed:', error)
    }
  }

  const handleDelete = async (id) => {
    try {
      await deleteCategory(id)
      message.success('删除成功')
      loadCategories()
    } catch (error) {
      console.error('Delete failed:', error)
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '描述', dataIndex: 'description', key: 'description' },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <>
          <Button size="small" onClick={() => handleEdit(record)}>编辑</Button>
          <Popconfirm
            title="确定删除此分类？"
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
        <Button type="primary" onClick={handleCreate}>创建分类</Button>
      </div>
      
      <Table
        columns={columns}
        dataSource={categories}
        loading={loading}
        rowKey="id"
      />

      <Modal
        title={currentCategory ? '编辑分类' : '创建分类'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item 
            name="name" 
            label="名称"
            rules={[{ required: true, message: '请输入分类名称' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={4} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit">保存</Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default AdminCategories
