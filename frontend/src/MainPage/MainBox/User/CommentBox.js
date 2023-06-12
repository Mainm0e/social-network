import React, {useState} from 'react';
import "./CommentBox.css"
const CommentBox = ({id}) => {
// return createcomment and commentlist button and default commentlist
const onSubmit = (commentData) => {
    console.log(commentData);
  };
    return(
        <div className="comment">
            <div className="comment-button">
           <button>Create Comment</button>
           <button>Comment List</button>
              </div>
            <CommentList id={id}/>
            <CreateComment id={id} onSubmit={onSubmit}/>
        </div>
    )
}

const CreateComment = ({ id, onSubmit }) => {
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
        <>
        </>
    )
} 




export default CommentBox;