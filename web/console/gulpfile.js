/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
const { src, dest, parallel } = require('gulp');
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
  const rsp = await fg(['static/js/index.(tke|project).main.*.js', 'static/js/common-vendor.*.js'], { cwd: 'build' });

  const TKE_JS_NAME = rsp.find(p => p.includes('tke'));
  const PROJECT_JS_NAME = rsp.find(p => p.includes('project'));
  const COMMON_VENDOR_JS_NAME = rsp.find(p => p.includes('common-vendor'));

  return src('public/index.tmpl.html')
    .pipe(ejs({ TKE_JS_NAME, PROJECT_JS_NAME, COMMON_VENDOR_JS_NAME }))
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
  return src(['public/**/*.html', '!public/index.tmpl.html'])
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
