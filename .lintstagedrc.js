module.exports = {
  "*.{css,less,scss,md,json}": (filenames) => {
    return [
      `prettier --write ${filenames.join(" ")}`,
    ];
  },
};
