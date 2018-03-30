const env = process.env.NODE_ENV || 'development';

const output = {};

if (env !== 'development') {
  output.publicPath = './';
}

module.exports = {
  devServer: {
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:5018',
      },
    },
  },
  configureWebpack: {
    output,
    // plugins: [
    //   new webpack.DefinePlugin({
    //     ENV: JSON.stringify(env),
    //   }),
    // ],
  },
};
