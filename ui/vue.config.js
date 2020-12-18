module.exports = {
  devServer: {
    public : 'dev.meshify.app',
    clientLogLevel: 'debug',
    disableHostCheck: true,
    port: 8081,
  },
  "transpileDependencies": [
    "vuetify"
  ],
  configureWebpack: {
    devtool: 'source-map'
  }  
};
