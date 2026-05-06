// 格式化工具函数

/**
 * 格式化日期，统一转为北京时间显示
 * 后端返回 "YYYY-MM-DD HH:mm:ss" 无时区标识，视为 UTC 解析后转北京时间
 */
export function formatDate(dateString: string): string {
  if (!dateString || dateString === '-') return '-';
  const trimmed = dateString.trim();

  let date: Date;
  if (/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/.test(trimmed)) {
    // 无时区标识，后端存的是 UTC，追加 Z 按 UTC 解析
    date = new Date(trimmed.replace(' ', 'T') + 'Z');
  } else {
    date = new Date(trimmed);
  }

  if (Number.isNaN(date.getTime())) return '-';

  // 手动偏移 +8 小时转为北京时间
  const bjDate = new Date(date.getTime() + 8 * 3600 * 1000);
  const y = bjDate.getUTCFullYear();
  const mo = String(bjDate.getUTCMonth() + 1).padStart(2, '0');
  const d = String(bjDate.getUTCDate()).padStart(2, '0');
  const h = String(bjDate.getUTCHours()).padStart(2, '0');
  const mi = String(bjDate.getUTCMinutes()).padStart(2, '0');
  const s = String(bjDate.getUTCSeconds()).padStart(2, '0');
  return `${y}-${mo}-${d} ${h}:${mi}:${s}`;
}

/**
 * 格式化日期（仅日期部分）
 * @param dateString ISO日期字符串
 * @returns 格式化后的日期字符串 (YYYY-MM-DD)
 */
export function formatDateOnly(dateString: string): string {
  const date = new Date(dateString);
  if (Number.isNaN(date.getTime())) {
    return '-';
  }
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');

  return `${year}-${month}-${day}`;
}

/**
 * 格式化价格
 * @param price 价格数值
 * @returns 格式化后的价格字符串 (¥123.45)
 */
export function formatPrice(price: number): string {
  return `¥${price.toFixed(2)}`;
}

/**
 * 格式化数字（千分位分隔）
 * @param num 数字
 * @returns 格式化后的数字字符串 (1,234,567)
 */
export function formatNumber(num: number): string {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
}

/**
 * 格式化相对时间
 * @param dateString ISO日期字符串
 * @returns 相对时间字符串 (刚刚、5分钟前、2小时前等)
 */
export function formatRelativeTime(dateString: string): string {
  const date = new Date(dateString);
  if (Number.isNaN(date.getTime())) {
    return '-';
  }
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffSeconds = Math.floor(diffMs / 1000);
  const diffMinutes = Math.floor(diffSeconds / 60);
  const diffHours = Math.floor(diffMinutes / 60);
  const diffDays = Math.floor(diffHours / 24);

  if (diffSeconds < 60) {
    return '刚刚';
  } else if (diffMinutes < 60) {
    return `${diffMinutes}分钟前`;
  } else if (diffHours < 24) {
    return `${diffHours}小时前`;
  } else if (diffDays < 7) {
    return `${diffDays}天前`;
  } else {
    return formatDateOnly(dateString);
  }
}
