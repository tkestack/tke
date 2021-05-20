const { src, dest, series, parallel } = require('gulp');
const uglifyjs = require('gulp-uglify');
const { pipeline } = require('readable-stream');
const cleanCSS = require('gulp-clean-css');
const htmlmin = require('gulp-html-minifier-terser');
const ejs = require('gulp-ejs');
const fg = require('fast-glob');

// 从public压缩复制文件到build
function minifyJs() {
  return pipeline(src('public/**/*.js'), uglifyjs(), dest('build'));
}

function minifyCss() {
  return src('public/**/*.css').pipe(cleanCSS()).pipe(dest('build'));
}

// 将打包好的js添加到html中
async function minifyIndexHtmlWithInjectJs() {
  const [[{ name: TKE_JS_NAME }], [{ name: PROJECT_JS_NAME }]] = await Promise.all([
    fg('build/static/js/index.tke.*.js', { objectMode: true }),
    fg('build/static/js/index.project.*.js', { objectMode: true })
  ]);

  console.log(TKE_JS_NAME, PROJECT_JS_NAME);

  return src('public/index.html')
    .pipe(ejs({ TKE_JS_NAME, PROJECT_JS_NAME }))
    .pipe(
      htmlmin({
        removeComments: true,
        collapseWhitespace: true,
        minifyJS: true,
        minifyCSS: true
      })
    )
    .pipe(dest('build'));
}

function minifyHtml() {
  return src(['public/**/*.html', '!public/index.html'])
    .pipe(
      htmlmin({
        removeComments: true,
        collapseWhitespace: true,
        minifyJS: true,
        minifyCSS: true
      })
    )
    .pipe(dest('build'));
}

// 复制其他不处理的文件
function copyAnother() {
  return src(['public/**/*', '!public/**/*.js', '!public/**/*.css', '!public/**/*.html']).pipe(dest('build'));
}

exports.default = parallel(minifyJs, minifyCss, minifyHtml, minifyIndexHtmlWithInjectJs, copyAnother);
