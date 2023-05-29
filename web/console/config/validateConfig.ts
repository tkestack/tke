import { z } from 'zod';

export const nameRule = (key: string) => {
  return z
    .string()
    .min(1, { message: `${key}不能为空` })
    .max(63, { message: `${key}长度不能超过63个字符` })
    .regex(/^([A-Za-z0-9][-A-Za-z0-9_./]*)?[A-Za-z0-9]$/, { message: `${key}格式不正确` });
};
