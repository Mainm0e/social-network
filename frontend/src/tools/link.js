import { getUserId } from "./cookie";


export const home = () => {
    const userId = String(getUserId("userId"));
    window.location.href = "/user?id=" + userId+"#postlist";
  };
export const profile = (id) => {
    window.location.href = "/user?id=" + id+"#postlist";
};
export const exploreGroup = () => {
    window.location.href = "/group";
  };
export const link_following = () => {
    window.location.href = "#followings";
  }
export const link_followers = () => {
    window.location.href = "#followers";
  }
export const link_notifications = () => {
    window.location.href = "#notifications";
  }
  export const link_eventlist = () => {
    window.location.href = "#eventlist";
  }