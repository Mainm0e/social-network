import { getUserId } from "./cookie";


export const home = () => {
    const userId = String(getUserId("userId"));
    window.location.href = "/user?id=" + userId;
  };
  