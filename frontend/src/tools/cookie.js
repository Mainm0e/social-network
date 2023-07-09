export function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2){
        return parts[1]
    } else {
        return null
    }
    }

export function getUserId(name){
    const value = `${localStorage.getItem(name)}`;
    if (value){
        return parseInt(value)
    } else {
        return null
    }
}