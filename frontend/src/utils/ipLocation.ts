/**
 * 从IP地址获取地理位置信息
 * @param ip IP地址
 * @returns 地理位置，如果无法获取则返回'Unknown'
 */
export async function getLocationFromIp(ip: string): Promise<string> {
  try {
    // 使用免费的IP地理位置API
    const response = await fetch(`https://ipapi.co/${ip}/json/`)
    const data = await response.json()
    return data.city ? `${data.city}, ${data.country_name}` : 'Unknown'
  } catch (error) {
    console.error('Failed to get location from IP:', error)
    return 'Unknown'
  }
}