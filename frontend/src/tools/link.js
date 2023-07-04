import { getUserId } from "./cookie";


export const home = () => {
    const userId = String(getUserId("userId"));
    window.location.href = "/user?id=" + userId+"#postlist";
  };
  
export const link_following = () => {
    window.location.href = "#followings";
  }
export const link_followers = () => {
    window.location.href = "#followers";
  }