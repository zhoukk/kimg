
// ref: https://umijs.org/config/
export default {
  treeShaking: true,
  hash: true,
  targets: {
    ie: 9,
  },
  theme: {
    '@primary-color': 'darkslateblue'
  },
  outputPath: '../www',
  proxy: {
    '/image': {
      target: "http://localhost/",
      changeOrigin: true,
    },
    '/info': {
      target: "http://localhost/",
      changeOrigin: true,
    }
  },
  plugins: [
    // ref: https://umijs.org/plugin/umi-plugin-react.html
    ['umi-plugin-react', {
      antd: true,
      dva: false,
      dynamicImport: { webpackChunkName: true },
      title: 'kimg',
      dll: true,
      locale: {},

      routes: {
        exclude: [
          /components\//,
        ],
      },
    }],
  ],
}
