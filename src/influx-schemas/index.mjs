export default {
  performance: {
    fields: {
      heapTotalHeapSize: 'i',
      memoryUsageRss: 'i',
      memoryUsageHeapTotal: 'i',
      cpuUsageUsedPercent: 'i',
      cpuUsageUserUsedPercent: 'i',
      cpuUsageSystemUsedPercent: 'i',
      connectingCount: 'i',
      heapSpaceNewSpaceUsedSize: 'i',
      heapSpaceOldSpaceUsedSize: 'i',
      heapSpaceCodeSpaceUsedSize: 'i',
      heapSpaceMapSpaceUsedSize: 'i',
      heapSpaceLargeObjectSpaceUsedSize: 'i',
    },
    tags: {
      memory: ['low', 'mid', 'high', 'higher'],
      cpu: ['free', 'normal', 'busy'],
      connecting: ['fewer', 'few', 'medium', 'many'],
    },
    options: {
      stripUnknown: true,
    },
  },
  deprecate: {
    fields: {
      path: 's',
    },
    options: {
      stripUnknown: true,
    },
  },
  httpRoute: {
    fields: {
      use: 'i',
      url: 's',
    },
    tags: {
      method: '*',
      path: '*',
      spdy: '012345'.split(''),
    },
    options: {
      stripUnknown: true,
    },
  },
  exception: {
    fields: {
      code: 's',
      path: 's',
      message: 's',
    },
    tags: {
      type: 'EU'.split(''),
    },
    options: {
      stripUnknown: true,
    },
  },
  response: {
    fields: {
      connecting: 'i',
      total: 'i',
      use: 'i',
      bytes: 'i',
      code: 'i',
      ip: 's',
      url: 's',
      request: 'i',
    },
    // 根据koa-http-stats配置的指定
    tags: {
      type: '12345'.split(''),
      spdy: '012345'.split(''),
      size: '012345'.split(''),
      busy: '01234'.split(''),
      method: '*',
    },
    options: {
      stripUnknown: true,
    },
  },
  session: {
    fields: {
      account: 's',
      use: 'i',
    },
    tags: {
      spdy: '012345'.split(''),
    },
    options: {
      stripUnknown: true,
    },
  },
  userTracker: {
    fields: {
      ip: 's',
      account: 's',
      params: 's',
      use: 'i',
      message: 's',
    },
    tags: {
      category: '*',
      result: ['success', 'fail'],
    },
    options: {
      stripUnknown: true,
    },
  },
  mongo: {
    fields: {
      use: 'i',
      conditions: 's',
      options: 's',
      fields: 's',
      size: 'i',
      update: 's',
    },
    tags: {
      spdy: '01234'.split(''),
      model: '*',
      op: [
        'find',
        'findOne',
        'findOneAndUpdate',
        'findOneAndRemove',
        'save',
        'count',
      ],
    },
  },
  http: {
    tags: {
      category: ['request', 'response'],
      type: '12345'.split(''),
      method: '*',
      host: '*',
      spdy: '01234'.split(''),
    },
    fields: {
      requesting: 'i',
      url: 's',
      status: 'i',
      bytes: 'i',
      socket: 'i',
      dns: 'i',
      tcp: 'i',
      tls: 'i',
      processing: 'i',
      transfer: 'i',
      all: 'i',
      ip: 's',
    },
  },
  userLogin: {
    fields: {
      account: 's',
      ip: 's',
    },
    tags: {
      country: '*',
      region: '*',
      city: '*',
      isp: '*',
    },
    options: {
      stripUnknown: true,
    },
  },
};
