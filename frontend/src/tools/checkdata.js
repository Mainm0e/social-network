// Desc: check data before send to server
/*
 */
export function checkPostData(data){
    console.log("incheck",data)
    if (data.title === "" || data.title === null || data.title === undefined) {
        return {
            title: "Error",
            message: ["Title is empty"],
            status: false,
        };
    } else if (data.content === "" || data.content === null || data.content === undefined) {
        return {
            title: "Error",
            message: ["Content is empty"],
            status: false,
        };
    } else if (data.privacy === "" || data.privacy === null || data.privacy === undefined) {
        return {
            title: "Error",
            message: ["Privacy is empty"],
            status: false,
        };
    } else {   
        return {
            title: "Success",
            message: ["Post successfully"],
            status: true,
        };
    }
    }
export function checkCommentData(data){
    /* 
     const commentData = {
    sessionId: sessionId,
    commentId: 0,
      postId: id,
      userId: parseInt(userId),
      creatorProfile: null,
      content: comment,
      image: image,
      date:"",
    }; */
    console.log("incheck",data)
    if (data.content === "" || data.content === null || data.content === undefined) {
        return {
            title: "Error",
            message: ["Content is empty"],
            status: false,
        };
    } else {
        return {
            title: "Success",
            message: ["Comment successfully"],
            status: true,
        };
    }
    }