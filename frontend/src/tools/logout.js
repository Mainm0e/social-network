export function logout() {
  document.cookie = "sessionId=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
  // sent cookie to server to delete
  window.location.href = "/login";
}