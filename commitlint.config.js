module.exports = {
    extends: ['@commitlint/config-angular', '@commitlint/config-conventional'],
    rules: {
        'body-max-line-length': [2, 'always', 260],
        'footer-max-line-length': [2, 'always', 260],
        'scope-case': [2, 'always', ['lower-case', 'upper-case']]
    }
};
