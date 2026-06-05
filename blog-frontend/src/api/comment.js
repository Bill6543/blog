import request from '../utils/request'

// 获取评论列表
export const getCommentList = (articleId) => {
  return request({
    url: '/comments',
    method: 'get',
    params: { article_id: articleId },
  })
}

// 创建评论
export const createComment = (data) => {
  return request({
    url: '/comments',
    method: 'post',
    data,
  })
}

// 更新评论
export const updateComment = (id, data) => {
  return request({
    url: `/comments/${id}`,
    method: 'put',
    data,
  })
}

// 删除评论
export const deleteComment = (id) => {
  return request({
    url: `/comments/${id}`,
    method: 'delete',
  })
}

// 获取待审核评论列表（管理员专用）
export const getPendingComments = () => {
  return request({
    url: '/admin/comments/pending',
    method: 'get',
  })
}

// 获取待审核评论数量（管理员专用）
export const getPendingCommentCount = () => {
  return request({
    url: '/admin/comments/pending/count',
    method: 'get',
  })
}

// 审核通过评论（管理员专用）
export const approveComment = (id) => {
  return request({
    url: `/admin/comments/${id}/approve`,
    method: 'put',
  })
}

// 审核拒绝评论（管理员专用）
export const rejectComment = (id) => {
  return request({
    url: `/admin/comments/${id}/reject`,
    method: 'put',
  })
}
