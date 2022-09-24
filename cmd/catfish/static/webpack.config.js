const path = require('path');
const webpack = require('webpack');
const CopyPlugin = require('copy-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');

const appDirectory = path.resolve(__dirname, './');

var apiEndpoint = "";

switch (process.env.TARGET) {
  case 'local':
    apiEndpoint = "http://localhost:8081";
    break;
}

const imageLoaderConfiguration = {
  test: /\.(gif|jpe?g|png|svg)$/,
  use: {
    loader: 'file-loader',
    options: {
      name: '[name].[ext]',
    },
  },
};

const babelLoaderConfiguration = {
  test: /\.(js|jsx)$/,
  exclude: /node_modules/,
  use: {
    loader: 'babel-loader',
    options: {
      presets: ['@babel/preset-env', '@babel/preset-react']
    }
  }
};

module.exports = {
  cache: {
    type: 'memory',
  },
  entry: path.resolve(appDirectory, 'src/index.jsx'),
  devtool: 'source-map',
  mode: process.env.WEBPACK_SERVE ? 'development' : 'production',

  output: {
    filename: 'bundle.js',
    publicPath: '/assets/',
    path: path.resolve(appDirectory, './public/assets'),
  },
  module: {
    rules: [
      imageLoaderConfiguration,
      babelLoaderConfiguration,
    ],
  },
  devServer: {
    allowedHosts: 'all',
    historyApiFallback: {
      rewrites: [
        { from: /./, to: '/index.html' }
      ]
    },
    static: path.resolve(appDirectory, 'public'),
    host: '0.0.0.0',
  },

  plugins: [
    new CleanWebpackPlugin({
      cleanStaleWebpackAssets: !process.env.WEBPACK_SERVE,
    }),
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV || 'development'),
      __API_ENDPOINT__: JSON.stringify(apiEndpoint),
    }),
    new CopyPlugin({
      patterns: [
        {
          from: 'node_modules/bootstrap/dist/css/bootstrap.min.css',
          to: path.resolve(appDirectory, './public/assets'),
        },
      ],
    }),
  ],
};
