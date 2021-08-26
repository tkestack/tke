/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

/* eslint-disable */
export const cloneDeep = function(item) {
  if (!item) {
    return item;
  } // null, undefined values check

  var types = [Number, String, Boolean],
    result;

  // normalizing primitives if someone did new String('aaa'), or new Number('444');
  types.forEach(function(type) {
    if (item instanceof type) {
      result = type(item);
    }
  });

  if (typeof result == 'undefined') {
    if (Object.prototype.toString.call(item) === '[object Array]') {
      result = [];
      item.forEach(function(child, index, array) {
        result[index] = cloneDeep(child);
      });
    } else if (typeof item == 'object') {
      // testing that this is DOM
      if (item.nodeType && typeof item.cloneNode == 'function') {
        var result = item.cloneNode(true);
      } else if (!item.prototype) {
        // check that this is a literal
        if (item instanceof Date) {
          result = new Date(item);
        } else {
          // it is an object literal
          result = {};
          for (var i in item) {
            result[i] = cloneDeep(item[i]);
          }
        }
      } else {
        result = item;
      }
    } else {
      result = item;
    }
  }

  return result;
};
/* eslint-enable */
