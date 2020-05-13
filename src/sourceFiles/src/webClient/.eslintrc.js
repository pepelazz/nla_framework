module.exports = {
  root: true,

  parserOptions: {
    parser: 'babel-eslint',
    sourceType: 'module'
  },

  env: {
    browser: true
  },

  extends: [
    // https://github.com/vuejs/eslint-plugin-vue#priority-a-essential-error-prevention
    // consider switching to `plugin:vue/strongly-recommended` or `plugin:vue/recommended` for stricter rules.
    'plugin:vue/essential',
    '@vue/standard'
  ],

  // required to lint *.vue files
  plugins: [
    'vue'
  ],

  globals: {
    'ga': true, // Google Analytics
    'cordova': true,
    '__statics': true,
    'process': true
  },

  // add your custom rules here
  rules: {
    // allow async-await
    'generator-star-spacing': 'off',
    // allow paren-less arrow functions
    'arrow-parens': 'off',
    'one-var': 'off',
    'no-tabs': 0,
    'no-multiple-empty-lines': 0,
    'object-property-newline': 0,
    'key-spacing': 0,
    'comma-spacing': 0,
    'no-mixed-spaces-and-tabs': 0,
    'no-multi-spaces': 0,
    'no-return-assign': 'off',
    'padded-blocks': 'off',
    'quotes': 'off',

    'import/first': 'off',
    'import/named': 'error',
    'import/namespace': 'error',
    'import/default': 'error',
    'import/export': 'error',
    'import/extensions': 'off',
    'import/no-unresolved': 'off',
    'import/no-extraneous-dependencies': 'off',
    'prefer-promise-reject-errors': 'off',
    'object-curly-spacing': 'off',
    'standard/object-curly-even-spacing': 'off',

    // allow console.log during development only
    'no-console': process.env.NODE_ENV === 'production' ? 'off' : 'off',
    // allow debugger during development only
    'no-debugger': process.env.NODE_ENV === 'production' ? 'off' : 'off',
    'no-trailing-spaces': 'off',
    'brace-style': 'off',
    // 'object-curly-spacing': [2, "never"],
    'space-before-function-paren': 'off',
    'semi': 'off',
    "comma-dangle": 'off',
    "yoda": 'off',
    "indent": 'off',
    "camelcase": 'off',
    "vue/require-v-for-key": 'off',
    "vue/no-side-effects-in-computed-properties": 'off',
    "vue/valid-v-model": 'off',
    "template-curly-spacing" : "off"
  }
}
