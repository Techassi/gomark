const path = require('path');

module.exports = {
    entry: {
        app: './src/app.js'
    },
    output: {
        filename: '[name].min.js',
        path: path.resolve(__dirname, 'dist')
    },
    mode: "production"
};
