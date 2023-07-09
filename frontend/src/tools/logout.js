export function logout() {
  document.cookie = "sessionId=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
  localStorage.removeItem("userId");
  // sent cookie to server to delete
  window.location.href = "/login";
}