import { z } from 'zod';

// 注册表单验证模式
export const registerSchema = z.object({
  username: z.string().min(1, '用户名不能为空'),
  password: z.string().min(6, '密码长度至少为6个字符'),
});

// 登录表单验证模式
export const loginSchema = z.object({
  username: z.string().min(1, '用户名不能为空'),
  password: z.string().min(1, '密码不能为空'),
});

// 商品表单验证模式
export const productSchema = z.object({
  name: z.string().min(1, '商品名称不能为空'),
  description: z.string().min(1, '商品描述不能为空'),
  price: z.number().positive('价格必须大于0'),
  stock: z.number().int().nonnegative('库存必须为非负整数'),
  imageUrl: z.string().url('请输入有效的图片链接').optional().or(z.literal('')),
});

// 订单表单验证模式
export const orderSchema = z.object({
  productId: z.union([z.string().min(1, '商品ID不能为空'), z.number().positive('商品ID不能为空')]),
  quantity: z.number().int().positive('购买数量必须为正整数'),
});

// 验证辅助函数
export function validateUsername(username: string): boolean {
  return username.trim().length > 0;
}

export function validatePassword(password: string): boolean {
  return password.length >= 6;
}

export function validateEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

export function validatePositiveNumber(value: number): boolean {
  return value > 0;
}

export function validateNonNegativeInteger(value: number): boolean {
  return Number.isInteger(value) && value >= 0;
}
