import { useState, useEffect } from 'react'
import { Table, Button, Tag, Popconfirm, message, Modal, Form, Input, Select } from 'antd'
import { getUserList, deleteUser, updateUserInfo } from '../../api/user'

const { TextArea } = Input

const AdminUsers = () => {
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(false)
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 })
  const [modalVisible, setModalVisible] = useState(false)
  const [currentUser, setCurrentUser] = useState(null)
  const [form] = Form.useForm()

  const loadUsers = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const data = await getUserList({ page, page_size: pageSize })
      setUsers(data.list)
      setPagination(prev => ({ ...prev, current: page, pageSize, total: data.total }))
    } catch (error) {
      console.error('Failed to load users:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadUsers()
  }, [])

  const handleEdit = (user) => {
    setCurrentUser(user)
    form.setFieldsValue({
      nickname: user.nickname,
      role: user.role,
      status: user.status,
    })
    setModalVisible(true)
  }

  const handleUpdate = async (values) => {
    try {
      await updateUserInfo(currentUser.id, values)
      message.success('更新成功')
      setModalVisible(false)
      loadUsers()
    } catch (error) {
      console.error('Update failed:', error)
    }
  }

  const handleDelete = async (id) => {
    try {
      await deleteUser(id)
      message.success('删除成功')
      loadUsers()
    } catch (error) {
      console.error('Delete failed:', error)
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '用户名', dataIndex: 'username', key: 'username' },
    { title: '昵称', dataIndex: 'nickname', key: 'nickname' },
    { title: '邮箱', dataIndex: 'email', key: 'email' },
    { 
      title: '角色', 
      dataIndex: 'role', 
      key: 'role',
      render: (role) => <Tag color={role === 'admin' ? 'red' : 'blue'}>{role}</Tag>
    },
    { 
      title: '状态', 
      dataIndex: 'status', 
      key: 'status',
      render: (status) => <Tag color={status === 1 ? 'green' : 'default'}>{status === 1 ? '正常' : '禁用'}</Tag>
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <>
          <Button size="small" onClick={() => handleEdit(record)}>编辑</Button>
          <Popconfirm
            title="确定删除此用户？"
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
      <h2>用户管理</h2>
      <Table
        columns={columns}
        dataSource={users}
        loading={loading}
        pagination={pagination}
        onChange={(pag) => loadUsers(pag.current, pag.pageSize)}
        rowKey="id"
      />

      <Modal
        title="编辑用户"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
      >
        <Form form={form} layout="vertical" onFinish={handleUpdate}>
          <Form.Item name="nickname" label="昵称">
            <Input />
          </Form.Item>
          <Form.Item name="role" label="角色">
            <Select>
              <Select.Option value="user">普通用户</Select.Option>
              <Select.Option value="admin">管理员</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select>
              <Select.Option value={1}>正常</Select.Option>
              <Select.Option value={0}>禁用</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit">保存</Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default AdminUsers
