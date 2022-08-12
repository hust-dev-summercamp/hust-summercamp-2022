import request from '@/utils/request'

/**
 * 获取分类列表
 */
export const getFiles = (params) => {
  return request({
    url: '/file/query',
    method: 'POST',
    params,
  })
}

export const downloadUrl = (params, callback) => {
  return request({
    url: '/file/downloadurl',
    method: 'POST',
    params
  }).then((response) => {
    callback && callback(response)
    return response
  })
}

export const uploadFile = (params, callback) => {
  return request({
    url: '/file/upload',
    method: 'POST',
    headers: { "Content-Type": "multipart/form-data" },
    data: params
  }).then(response => {
    callback && callback(response)
    return response
  })
}

export const deleteFile = (params, callback) => {
  return request({
    url: '/file/delete',
    method: 'DELETE',
    params
  }).then(response => {
    callback && callback(response)
    return response
  })
}

export const updateFile = (params, callback) => {
  return request({
    url: '/file/update',
    method: 'PUT',
    params
  }).then(response => {
    callback && callback(response)
    return response
  })
}