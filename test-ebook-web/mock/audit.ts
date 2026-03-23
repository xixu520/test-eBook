import { MockMethod } from 'vite-plugin-mock'

const auditLogs = [
  {
    id: 1,
    timestamp: '2026-03-23 09:00:00',
    username: 'admin',
    action: 'LOGIN',
    ip: '192.168.1.1',
    details: '用户登录成功'
  },
  {
    id: 2,
    timestamp: '2026-03-23 09:15:00',
    username: 'admin',
    action: 'UPLOAD',
    ip: '192.168.1.1',
    details: '上传文件: GB 50007-2011'
  },
  {
    id: 3,
    timestamp: '2026-03-23 09:20:00',
    username: 'editor',
    action: 'EDIT',
    ip: '192.168.1.10',
    details: '修改分类: 地基基础'
  },
  {
    id: 4,
    timestamp: '2026-03-23 09:25:00',
    username: 'admin',
    action: 'DELETE',
    ip: '192.168.1.1',
    details: '删除文件: GBJ 202-83'
  }
]

export default [
  {
    url: '/api/v1/audit-logs',
    method: 'get',
    response: ({ query }: { query: any }) => {
      const { page = 1, size = 10, action } = query
      let filtered = [...auditLogs]
      if (action) {
        filtered = filtered.filter(i => i.action === action)
      }
      return {
        code: 200,
        message: '获取成功',
        data: {
          total: filtered.length,
          page: Number(page),
          size: Number(size),
          list: filtered.slice((Number(page)-1)*Number(size), Number(page)*Number(size))
        }
      }
    }
  }
] as MockMethod[]
