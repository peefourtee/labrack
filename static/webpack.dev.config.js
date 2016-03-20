module.exports = {
  entry: "./assets/javascript/index.js",
  output: {
    path: __dirname,
    filename: "bundle.js"
  },
  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        loader: "babel-loader",
        query: {
          presets: ["react", "es2015"],
          cacheDirectory: "./build"
        }
      }
    ]
  },
  plugins: [],
  devtool: "source-map"
};
