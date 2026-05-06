// LocalStorage封装工具

const AUTH_TOKEN_KEY = 'auth_token';
const USER_INFO_KEY = 'user_info';

/**
 * 存储认证Token
 */
export function setAuthToken(token: string): void {
  localStorage.setItem(AUTH_TOKEN_KEY, token);
}

/**
 * 获取认证Token
 */
export function getAuthToken(): string | null {
  const token = localStorage.getItem(AUTH_TOKEN_KEY);

  // 如果 token 为 "undefined" 或 "null" 字符串，返回 null
  if (!token || token === 'undefined' || token === 'null') {
    return null;
  }

  return token;
}

/**
 * 移除认证Token
 */
export function removeAuthToken(): void {
  localStorage.removeItem(AUTH_TOKEN_KEY);
}

/**
 * 存储用户信息
 */
export function setUserInfo(user: any): void {
  localStorage.setItem(USER_INFO_KEY, JSON.stringify(user));
}

/**
 * 获取用户信息
 */
export function getUserInfo<T = any>(): T | null {
  const userStr = localStorage.getItem(USER_INFO_KEY);

  // 如果没有数据或数据为 "undefined" 字符串，返回 null
  if (!userStr || userStr === 'undefined' || userStr === 'null') {
    return null;
  }

  try {
    return JSON.parse(userStr) as T;
  } catch (error) {
    console.error('Failed to parse user info:', error);
    // 清除无效数据
    removeUserInfo();
    return null;
  }
}

/**
 * 移除用户信息
 */
export function removeUserInfo(): void {
  localStorage.removeItem(USER_INFO_KEY);
}

/**
 * 清除所有认证相关数据
 */
export function clearAuthData(): void {
  removeAuthToken();
  removeUserInfo();
}

/**
 * 通用存储方法
 */
export function setItem(key: string, value: any): void {
  try {
    const serialized = typeof value === 'string' ? value : JSON.stringify(value);
    localStorage.setItem(key, serialized);
  } catch (error) {
    console.error(`Failed to set item ${key}:`, error);
  }
}

/**
 * 通用获取方法
 */
export function getItem<T = any>(key: string): T | null {
  try {
    const item = localStorage.getItem(key);
    if (!item) return null;

    // 尝试解析JSON，如果失败则返回原始字符串
    try {
      return JSON.parse(item) as T;
    } catch {
      return item as T;
    }
  } catch (error) {
    console.error(`Failed to get item ${key}:`, error);
    return null;
  }
}

/**
 * 通用移除方法
 */
export function removeItem(key: string): void {
  localStorage.removeItem(key);
}

/**
 * 清除所有存储
 */
export function clearAll(): void {
  localStorage.clear();
}
