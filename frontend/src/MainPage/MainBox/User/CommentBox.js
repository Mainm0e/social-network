import React, {useState} from 'react';
import "./CommentBox.css"
const CommentBox = ({id}) => {
// return createcomment and commentlist button and default commentlist
    const [boxState, setBoxState] = useState(null);
    const changeState = (e) => {
       if (e.target.innerText === "Create Comment"){
           setBoxState(<CreateComment id={id} showComment={showComment}/>)
       } else if (e.target.innerText === "Comment List"){
            setBoxState(<CommentList id={id}/>)
       }
    }
    const showComment = () => {
        setBoxState(<CommentList id={id}/>)
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

const CommentList = ({id}) => {
    return(
        <div className="comment_list">
            <div className="comment_list_top">
                <div className="comment_list_top_left">
                    <img src="https://i.imgur.com/1qkK1Q6.jpg" alt="profile" />
                </div>
                <div className="comment_list_top_right">
                    <div className="comment_list_top_right_name">
                        <span>name</span>
                    </div>
                    <div className="comment_list_top_right_date">
                        <span>date</span>
                    </div>
                </div>
            </div>
            <div className="comment_list_content">

            </div>
        </div>
    )
} 




export default CommentBox;