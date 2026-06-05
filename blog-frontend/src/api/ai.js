import request from '../utils/request'

// 生成文章摘要
export const generateSummary = (data) => {
  return request({
    url: '/ai/summary',
    method: 'post',
    data,
  })
}

// 生成文章封面
export const generateCover = (data) => {
  return request({
    url: '/ai/cover',
    method: 'post',
    data,
  })
}
