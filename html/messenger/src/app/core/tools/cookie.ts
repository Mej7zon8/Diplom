export function setApiCookie(name: string, value: string) {
  document.cookie = `${name}=${value}; path=/api/v3/`
}
export function deleteApiCookie(name: string) {
  document.cookie = `${name}=; path=/api/v3/; expires=Thu, 01 Jan 1970 00:00:00 GMT`
}
export function deleteCookie(name: string) {
  document.cookie = `${name}=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT`
}
