module.exports = {
    env: {
        browser: true,
        es6: true,
    },
    extends: 'airbnb-base',
    globals: {
        Atomics: 'readonly',
        SharedArrayBuffer: 'readonly',
    },
    parserOptions: {
        ecmaVersion: 2018,
        sourceType: 'module',
    },
    rules: {
        "linebreak-style": 0,
        "no-console": "off",
        "spaced-comment": "off",
        "indent": "off",
        "import/newline-after-import": "off",

    },
};