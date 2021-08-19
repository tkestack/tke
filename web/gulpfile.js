const { src, dest } = require("gulp");
const header = require("gulp-header");
const fs = require("fs");

function start() {
  return src(["./**/*.{ts,tsx,js,jsx}", "!./**/node_modules", "!./**/public"])
    .pipe(header(fs.readFileSync("./license")))
    .pipe(dest("."));
}

exports.default = start;
