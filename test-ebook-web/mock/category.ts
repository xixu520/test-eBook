import { MockMethod } from 'vite-plugin-mock'

export default [
  {
    url: '/api/v1/categories',
    method: 'get',
    response: () => {
      return {
        code: 200,
        message: '获取成功',
        data: [
          { id: 1, parent_id: 0, name: '地基基础', sort_order: 1, doc_count: 128 },
          { id: 2, parent_id: 0, name: '建筑材料', sort_order: 2, doc_count: 85 },
          { id: 3, parent_id: 1, name: '桩基础', sort_order: 1, doc_count: 43 },
          { id: 4, parent_id: 1, name: '基坑工程', sort_order: 2, doc_count: 85 },
          { id: 5, parent_id: 2, name: '水泥混凝土', sort_order: 1, doc_count: 30 },
          { id: 6, parent_id: 2, name: '钢筋家电', sort_order: 2, doc_count: 55 },
        ],
      }
    },
  },
  {
    url: new RegExp('/api/v1/categories/.*'),
    method: 'delete',
    response: ({ url }: any) => {
      const id = url.split('/').pop()
      if (id === '1') {
        return { code: 400, message: '删除失败：该分类下存在子分类，请先移除或转移子分类' }
      }
      if (id === '2') {
        return { code: 400, message: '删除失败：该分类下仍有标准文件，请先转移或删除文件' }
      }
      return { code: 200, message: '分类删除成功' }
    }
  }
] as MockMethod[]
