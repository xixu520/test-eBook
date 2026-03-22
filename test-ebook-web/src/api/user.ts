import request from '@/utils/request'

export function getUserList() {
  return request({
    url: '/admin/users',
    method: 'get',
  })
}

export function updateUserStatus(id: number, status: number) {
  return request({
    url: `/admin/users/${id}/status`,
    method: 'put',
    data: { status },
  })
}

export function updateUserRole(id: number, role: string) {
  return request({
    url: `/admin/users/${id}/role`,
    method: 'put',
    data: { role },
  })
}

export function deleteUser(id: number) {
  return request({
    url: `/admin/users/${id}`,
    method: 'delete',
  })
}
