export async function fetchData(method,type, payload) {
const response = await fetch("http://localhost:8080/api", {
  method: method,
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({ 
    type: type,
    payload: payload,
  }),
});
const responseData = await response.json();
console.log(responseData)
if (responseData.statusCode === 200) {
  return responseData.event.payload
}
}