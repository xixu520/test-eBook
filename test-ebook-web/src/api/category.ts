import request from '@/utils/request'

export function getCategories() {
  return request({
    url: '/categories',
    method: 'get',
  })
}

export function addCategory(data: { name: string, parent_id: number, order: number }) {
  return request({
    url: '/categories',
    method: 'post',
    data
  })
}

export function updateCategory(id: number, data: { name: string, parent_id: number, order: number }) {
  return request({
    url: `/categories/${id}`,
    method: 'put',
    data
  })
}

export function deleteCategory(id: number) {
  return request({
    url: `/categories/${id}`,
    method: 'delete',
  })
}
