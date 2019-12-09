const path = require('path');

module.exports = {
    entry: {
        app: './src/*.js'
    },
    output: {
        filename: '[name].min.js',
        path: path.resolve(__dirname, 'dist')
    },
    mode: "production"
};
