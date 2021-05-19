const { src, dest, series, parallel } = require('gulp');
const uglifyjs = require('gulp-uglify');
const { pipeline } = require('readable-stream');
const cleanCSS = require('gulp-clean-css');
const htmlmin = require('gulp-htmlmin');
const ejs = require('gulp-ejs');
const fg = require('fast-glob');

// 从public压缩复制文件到build
function minifyJs() {
  return pipeline(src('public/static/js/*.js'), uglifyjs(), dest('build/static/js'));
}

function minifyCss() {
  return src('public/static/css/**/*.css').pipe(cleanCSS()).pipe(dest('build/static/css'));
}

// 将打包好的js添加到html中
async function minifyHtmlWithInjectJs() {
  const [{ name: TKE_JS_NAME }, { name: PROJECT_JS_NAME }] = await fg(
    ['build/index.tke.*.js', 'build/index.project.*.js'],
    {
      objectMode: true
    }
  );

  return src('public/index.html').pipe(ejs({ TKE_JS_NAME, PROJECT_JS_NAME })).pipe(dest('build'));
}

// 复制其他不处理的文件
function copyAnother() {
  return src(['public/static/**/*', '!public/static/**/*.js', '!public/static/**/*.css']).pipe(dest('build/static'));
}

exports.default = parallel(minifyHtmlWithInjectJs);
