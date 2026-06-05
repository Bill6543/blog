import request from '../utils/request'

// 获取用户列表
export const getUserList = (params) => {
  return request({
    url: '/users',
    method: 'get',
    params,
  })
}

// 获取用户详情
export const getUserInfo = (id) => {
  return request({
    url: `/users/${id}`,
    method: 'get',
  })
}

// 更新用户信息
export const updateUserInfo = (id, data) => {
  return request({
    url: `/users/${id}`,
    method: 'put',
    data,
  })
}

// 删除用户（管理员）
export const deleteUser = (id) => {
  return request({
    url: `/users/${id}`,
    method: 'delete',
  })
}
