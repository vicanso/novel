module.exports = {
  devServer: {
    proxy: {
      "/@nv": {
        target: "http://red:3015",
        changeOrigin: true
      }
    }
  },
  baseUrl: process.env.NODE_ENV === "production" ? "/@nv/static/" : "/"
};
