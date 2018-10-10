module.exports = {
  devServer: {
    proxy: {
      "/@nv": {
        target: "http://red:7001",
        changeOrigin: true
      }
    }
  },
  baseUrl: process.env.NODE_ENV === "production" ? "/static/" : "/"
};
