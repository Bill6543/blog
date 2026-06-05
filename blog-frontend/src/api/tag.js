import request from '../utils/request'

// 获取标签列表
export const getTagList = () => {
  return request({
    url: '/tags',
    method: 'get',
  })
}

// 获取标签详情
export const getTag = (id) => {
  return request({
    url: `/tags/${id}`,
    method: 'get',
  })
}

// 创建标签（管理员）
export const createTag = (data) => {
  return request({
    url: '/tags',
    method: 'post',
    data,
  })
}

// 更新标签（管理员）
export const updateTag = (id, data) => {
  return request({
    url: `/tags/${id}`,
    method: 'put',
    data,
  })
}

// 删除标签（管理员）
export const deleteTag = (id) => {
  return request({
    url: `/tags/${id}`,
    method: 'delete',
  })
}
