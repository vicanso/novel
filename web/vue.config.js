module.exports = {
  devServer: {
    proxy: {
      "/@nv": {
        target: "http://jenny.f3322.net:3015",
        changeOrigin: true
      }
    }
  },
  baseUrl: process.env.NODE_ENV === "production" ? "./static/" : "/"
};
