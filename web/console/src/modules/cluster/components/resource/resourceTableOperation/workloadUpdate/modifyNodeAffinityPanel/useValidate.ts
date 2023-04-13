const state = [
  {
    weight: {
      value: 1,
      validate(value) {
        if (!value) return { status: 'error', message: '权重不能为空' };
      }
    },
    subRules: [
      {
        key: {
          value: '',
          validate(value) {
            if (!value) return { status: 'error', message: 'key不能为空' };
          }
        },
        operator: '',
        value: {
          value: '',
          validate(value) {
            if (!value) return { status: 'error', message: 'value不能为空' };
          }
        }
      }
    ]
  }
];
