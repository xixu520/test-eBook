import request from '@/utils/request'

export function getCategories() {
  return request({
    url: '/standards/categories',
    method: 'get',
  })
}

export function addCategory(data: { name: string, parent_id: number, order: number }) {
  return request({
    url: '/standards/categories',
    method: 'post',
    data
  })
}

export function deleteCategory(id: number) {
  return request({
    url: `/standards/categories/${id}`,
    method: 'delete',
  })
}
