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
export const TYPES = {
  string: {
    type: 'string',
    value: undefined
  },
  stringWithInit: {
    type: 'string',
    value: '*'
  },
  number: {
    type: 'number',
    isTypeNumber: true,
    value: undefined
  },

  duration: {
    type: 'duration',
    isTypeNumber: true,
    unit: 'ms',
    value: undefined
  },
  stringArray: {
    type: 'stringArray',
    value: undefined
  },
  labels: {
    type: 'labels',
    value: undefined,
    map: {} as any
  },
  boolean: {
    type: 'boolean',
    value: undefined
  }
};
let gid = 0;
export function getState(schema, component: React.Component, rootObj?, obj?) {
  if (typeof schema.properties === 'object') {
    let properties = {};

    for (let key in schema.properties) {
      if (key !== 'generateName') {
        properties[key] = getState(schema.properties[key], component, rootObj, obj !== undefined ? obj[key] : obj);
      }
    }

    if (schema.type === 'pickOne') {
      let json = schemaObjToJSON({ properties });
      if (json) {
        schema.pick = Object.keys(json).find(key => !schema.properties[key].required);
      }
    }

    if (schema.isChecked !== undefined) {
      return { ...schema, properties, component, isChecked: Boolean(obj) };
    }

    return { ...schema, properties, component };
  }

  if (schema.valueSchema) {
    return {
      ...schema,
      properties: obj
        ? Object.keys(obj).reduce(
            (sum, key) => ({ ...sum, [key]: getState(schema.valueSchema, component, rootObj, obj[key]) }),
            {}
          )
        : {},
      component
    };
  }

  if (schema.type === 'array') {
    return {
      ...schema,
      value:
        obj && obj.map
          ? obj.map(item => getState(schema.item, component, rootObj, item))
          : schema.isInitFirstItem && !rootObj
          ? [schema.getInitFirstItem ? schema.getInitFirstItem() : getState(schema.item, component, rootObj, undefined)]
          : [],
      component
    };
  }
  if (schema.type === 'labels') {
    let value = obj
      ? Object.keys(obj).map(key => ({ __key: ++gid, key, value: obj[key] }))
      : schema.isInitFirstItem
      ? [
          schema.getInitFirstItem
            ? schema.getInitFirstItem()
            : {
                __key: ++gid,
                key: '',
                value: ''
              }
        ]
      : [];
    return {
      ...schema,
      value,
      map: obj || {},
      activeTimestamp: schema.activeTimestamp === 0 ? Boolean(value.length) : schema.activeTimestamp,
      component
    };
  }
  if (schema.type === 'stringArray') {
    let value = obj || [];
    return {
      ...schema,
      value,
      activeTimestamp: schema.activeTimestamp === 0 ? Boolean(value) : schema.activeTimestamp,
      component
    };
  }

  if (schema.type === 'stringMatch') {
    let value = Object.values(obj || {})[0];
    return {
      ...schema,
      key: Object.keys(obj || {})[0] || 'exact',
      value,
      activeTimestamp: schema.activeTimestamp === 0 ? Boolean(value) : schema.activeTimestamp,
      component
    };
  }

  if (schema.type === 'map') {
    let key = Object.keys(obj || {})[0] || '';
    return {
      ...schema,
      key,
      value: getState(schema.value, component, rootObj, Object.values(obj || {})[0]),
      activeTimestamp: schema.activeTimestamp === 0 ? Boolean(key) : schema.activeTimestamp,
      component
    };
  }

  if (schema.type === 'duration') {
    return {
      ...schema,
      value: obj !== undefined ? /\d+/.exec(obj)[0] : schema.value,
      unit: obj !== undefined ? /[a-zA-Z]+/.exec(obj)[0] : schema.unit,
      activeTimestamp: schema.activeTimestamp === 0 ? Boolean(obj) : schema.activeTimestamp,
      component
    };
  }

  if (schema.isChecked !== undefined) {
    return {
      ...schema,
      value: obj !== undefined ? obj : schema.value,
      isChecked: Boolean(obj !== undefined ? obj : schema.value),
      component
    };
  }

  return {
    ...schema,
    value: obj !== undefined ? obj : schema.value,
    component,
    activeTimestamp: schema.activeTimestamp === 0 ? Boolean(obj) : schema.activeTimestamp
  };
}

export function onChange(jsonObj, index?, objKey?) {
  return value => {
    switch (jsonObj.type) {
      case 'array':
      case 'labels':
        if (objKey !== undefined) {
          jsonObj.value[index][objKey] = value;
          if (jsonObj.type === 'labels') {
            jsonObj.map = jsonObj.map || {};
            jsonObj.map[objKey] = value;
          }
        } else {
          jsonObj.value[index] = value;
        }
        break;
      case 'stringMatch':
      case 'map':
        jsonObj[index] = value;
        break;
      case 'number':
        jsonObj.value = +value;
        if (isNaN(jsonObj.value) || (value.endsWith && value.endsWith('.'))) {
          jsonObj.value = value;
        }
        break;
      case 'duration':
        jsonObj.value = +value;
        break;
      case 'pickOne':
        jsonObj.pick = value;
        break;
      default:
        jsonObj[index || 'value'] = value;
    }
    updateComponent(jsonObj);
  };
}

export function onAdd(jsonObj) {
  return () => {
    jsonObj.value = jsonObj.value || [];
    switch (jsonObj.type) {
      case 'array':
        jsonObj.value.push({ __key: ++gid, ...getState(jsonObj.item, jsonObj.component, undefined, undefined) });
        break;
      case 'labels':
        jsonObj.value.push({ __key: ++gid, key: undefined, value: undefined });
    }
    updateComponent(jsonObj);
  };
}

export function onDelete(jsonObj, index) {
  return () => {
    jsonObj.value.splice(index, 1);
    updateComponent(jsonObj);
  };
}

function updateComponent(jsonObj) {
  let { component } = jsonObj;
  if (!component) {
    /* eslint-disable */
    debugger;
  }
  component.setState({}, () => {
    component.props.onChange && component.props.onChange(component.state);
  });
}

export function schemaObjToJSON(obj, skipPrivateValue = true) {
  if (!obj) {
    /* eslint-disable no-debugger*/
    debugger;
  }
  if (obj.toJSON) {
    return obj.toJSON(obj, schemaObjToJSON);
  }
  if (obj.isChecked === false) {
    return undefined;
  }
  if (obj.properties || obj.valueSchema) {
    let json = {};
    for (let key in obj.properties) {
      if (obj.pick && obj.pick !== key && !obj.properties[key].required) {
        continue;
      }
      if (skipPrivateValue && key[0] === '_') {
        continue;
      }
      let value = schemaObjToJSON(obj.properties[key], skipPrivateValue);
      if (value !== undefined) {
        if (typeof value === 'object') {
          if (Array.isArray(value) && value.length !== 0) {
            json[key] = value;
          } else if (Object.keys(value).length) {
            json[key] = value;
          }
        } else {
          json[key] = value;
        }
      }
    }
    if (!Object.keys(json).length) {
      return undefined;
    }
    return json;
  }

  if (obj.type === 'stringArray') {
    let value = typeof obj.value === 'string' ? obj.value.split('\n') : obj.value;
    value = value.map(l => l.trim()).filter(l => l);
    return value;
  }

  if (obj.type === 'array') {
    let value = obj.value.map(d => schemaObjToJSON(d, skipPrivateValue)).filter(v => v !== undefined);
    if (!value.length) {
      return undefined;
    }
    return value;
  }

  if (obj.type === 'stringMatch') {
    if (obj.value === undefined || !obj.key) {
      return undefined;
    }
    return { [obj.key]: obj.value };
  }

  if (obj.type === 'map') {
    if (!obj.key || !(obj.value && obj.value.key && obj.value.value !== undefined)) {
      return undefined;
    }
    return {
      [obj.key]: {
        [obj.value.key]: obj.value.value
      }
    };
  }

  if (obj.type === 'labels') {
    if (!obj.value.length) {
      return undefined;
    }
    return obj.value.reduce((sum, item) => ({ ...sum, [item.key]: item.value }), {});
  }

  if (typeof obj.value === 'object') {
    if (!Object.keys(obj.value).length) {
      return undefined;
    }
  }

  if (obj.type === 'duration' && obj.value !== undefined) {
    return obj.value + obj.unit;
  }

  if (obj.isChecked !== undefined) {
    if (!obj.isChecked) {
      return undefined;
    }
  }

  return obj.value;
}

export function cloneSchemaObj(schema, obj, component: React.Component) {
  let json = schemaObjToJSON(obj);
  return getState(schema, component, json, json);
}
