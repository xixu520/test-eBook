import { MockMethod } from 'vite-plugin-mock'

export default [
  {
    url: '/api/v1/stats/dashboard',
    method: 'get',
    response: () => {
      return {
        code: 200,
        message: '获取成功',
        data: {
          total_documents: 1256,
          today_uploaded: 12,
          pending_verify: 45,
          pending_ocr: 8,
          storage_used: '4.2 GB',
          recent_activities: [
            { id: 1, type: 'upload', content: '管理员上传了《建筑给水排水设计标准》', time: '10分钟前' },
            { id: 2, type: 'ocr', content: '《地基基础设计规范》OCR 识别完成', time: '25分钟前' },
            { id: 3, type: 'verify', content: '用户审核通过了 5 份文件', time: '1小时前' },
          ]
        },
      }
    },
  },
] as MockMethod[]
