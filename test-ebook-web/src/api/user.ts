import request from '@/utils/request'

export function getUsers(params: { page: number; page_size: number }) {
  return request({
    url: '/admin/users',
    method: 'get',
    params
  })
}

export function updateUserStatus(id: number, isActive: boolean) {
  return request({
    url: `/admin/users/${id}/status`,
    method: 'put',
    data: { is_active: isActive },
  })
}

// Alias for compatibility
export const getUserList = () => getUsers({ page: 1, page_size: 100 })

export function deleteUser(id: number) {
  return request({
    url: `/admin/users/${id}`,
    method: 'delete',
  })
}

export function createUser(data: any) {
  return request({
    url: '/admin/users',
    method: 'post',
    data
  })
}

export function updateUser(id: number, data: any) {
  return request({
    url: `/admin/users/${id}`,
    method: 'put',
    data
  })
}

export function resetPassword(id: number, data: any) {
  return request({
    url: `/admin/users/${id}/password`,
    method: 'put',
    data
  })
}

export function updateTheme(theme: string) {

  return request({
    url: '/users/me/theme',
    method: 'put',
    data: { theme }
  })
}
