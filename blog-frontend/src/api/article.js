import request from '../utils/request'

// 获取文章列表
export const getArticleList = (params) => {
  return request({
    url: '/articles',
    method: 'get',
    params,
  })
}

// 获取文章详情
export const getArticle = (id) => {
  return request({
    url: `/articles/${id}`,
    method: 'get',
  })
}

// 创建文章
export const createArticle = (data) => {
  return request({
    url: '/articles',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

// 更新文章
export const updateArticle = (id, data) => {
  return request({
    url: `/articles/${id}`,
    method: 'put',
    data,
  })
}

// 删除文章
export const deleteArticle = (id) => {
  return request({
    url: `/articles/${id}`,
    method: 'delete',
  })
}

// 点赞文章
export const likeArticle = (id) => {
  return request({
    url: `/articles/${id}/like`,
    method: 'post',
  })
}

// 取消点赞
export const unlikeArticle = (id) => {
  return request({
    url: `/articles/${id}/like`,
    method: 'delete',
  })
}

// 获取点赞状态
export const getLikeStatus = (id) => {
  return request({
    url: `/articles/${id}/like`,
    method: 'get',
  })
}

// 获取用户文章列表
export const getUserArticles = (userId, params) => {
  return request({
    url: `/users/${userId}/articles`,
    method: 'get',
    params,
  })
}

// 恢复已删除的文章
export const restoreArticle = (id) => {
  return request({
    url: `/articles/${id}/restore`,
    method: 'post',
  })
}

// 增加文章浏览量
export const increaseViewCount = (id) => {
  return request({
    url: `/articles/${id}/view`,
    method: 'post',
  })
}
