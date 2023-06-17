// Desc: check data before send to server
/*
 */
export function checkPostData(data){
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
    } else if (data.privecy === "" || data.privecy === null || data.privecy === undefined) {
        return {
            title: "Error",
            message: ["Privecy is empty"],
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