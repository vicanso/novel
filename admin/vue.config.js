module.exports = {
  devServer: {
    proxy: {
      "/@nv": {
        target: "http://red:3015",
        changeOrigin: true
      }
    }
  },
  // base url指定加载静态文件的前缀，需要需要配置
  baseUrl: process.env.NODE_ENV === "production" ? "./static/" : "/"
};
