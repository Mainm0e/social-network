import React, {useState} from 'react';
import "./CommentBox.css"
const CommentBox = ({id,comments}) => {
// return createcomment and commentlist button and default commentlist
    const [boxState, setBoxState] = useState(null);
    const changeState = (e) => {
       if (e.target.innerText === "Create Comment"){
           setBoxState(<CreateComment id={id} showComment={showComment}/>)
       } else if (e.target.innerText === "Comment List"&& comments !== undefined){
            setBoxState(<CommentList comments={comments}/>)
       }
    }
    const showComment = () => {
        setBoxState(<CommentList comments={comments}/>)
    }


   
    return(
        <div className="comment">
        <div className="comment-button">
            <button onClick={changeState}>Create Comment</button>
            <button onClick={changeState}>Comment List</button>
        </div>
       {boxState}
        
        </div>
    )
}

const CreateComment = ({ id , showComment }) => {
    const [comment, setComment] = useState('');
    const [image, setImage] = useState(null);

    const handleCommentChange = (e) => {
        setComment(e.target.value);
    };

    const handleImageChange = (e) => {
        const file = e.target.files[0];
        setImage(file);
    };
    
  const handleSubmit = (e) => {
    e.preventDefault();
    const commentData = {
      comment: comment,
      image: image,
    };
    onSubmit(commentData);
    setComment('');
    setImage(null);
    showComment();
  };

const onSubmit = (commentData) => {
    console.log(commentData);
  };

       return(
        <div className="create_comment">
            <form onSubmit={handleSubmit}>
                <div className="create_comment_top">
                    <textarea
                        placeholder="Comment"
                        className="create_comment_content"
                        value={comment}
                        onChange={handleCommentChange}
                    />
                    <input
                        type="file"
                        accept="image/*"
                        onChange={handleImageChange}
                    />
                </div>
                {image && (
                    <div className="create_comment_image">
                        <img src={URL.createObjectURL(image)} alt="content" />
                    </div>
                )}
                <div className="create_comment_button">
                    <button className="create_comment_submit">Submit</button>
                </div>
            </form>
        </div>
    )
}

const CommentList = ({comments}) => {
    const loopComment = (comments) => {

        return comments.map((comment) => (
        <div className="comment_list_item" key={comment.id}>
                <div className="comment_list_item_header">
                    <div className="comment_list_item_header_left">
                        <div className="comment_list_item_header_info">
                            <p>{comment.time}</p>
                        </div>
                    </div>
                    <div className="comment_list_item_header_right">
                        <div className="comment_list_item_header_user">
                          {/*   <img src={comment.user} alt="avatar" /> */}
                            <p>{comment.user.email}</p>
                        </div>
                    </div>
                </div>
                <div className="comment_list_item_body">
                    <p className="content">{comment.content}</p>
                </div>
                {checkImage(comment.image)}
            </div>
        ));
        };

    const checkImage = (image) => {
        if (image === ''|| image === null|| image === undefined) {
            return null;
        } else {
           return <div className="comment_list_item_image">
                <img src={image} alt="content" />
            </div>;
        }
    };
    return(
        <div className="comment_list">
            {loopComment(comments)}
        </div>
    )
} 




export default CommentBox;