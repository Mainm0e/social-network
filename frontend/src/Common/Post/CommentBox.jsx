import React, {useState} from 'react';
import { getCookie, getUserId } from '../../tools/cookie';
import "./CommentBox.css"
import { checkCommentData } from '../../tools/checkdata';
import { fetchData } from '../../tools/fetchData';
const CommentBox = ({id,comments,activePost}) => {
// return createcomment and commentlist button and default commentlist
    const [boxState, setBoxState] = useState(null);
    const [check, setCheck] = useState(null);
    const changeState = (e) => {
       if (e.target.innerText === "Create Comment" && check !== "create"){
            activePost(id)
            setCheck("create")
            setBoxState(<CreateComment id={id} showComment={showComment}/>)
       } else if (e.target.innerText === "Comment List"&& comments !== undefined && check !== "list"){
            activePost(id)
            setCheck("list")
            setBoxState(<CommentList comments={comments}/>)
       } else {
        setCheck(null)
        setBoxState(null)
        activePost(null)
       }
    }
    const showComment = () => {
        activePost(id)
        setBoxState(<CommentList comments={comments}/>)
    }

    return(
        <div className="comment">
        <div className="comment-button">
        <button
          onClick={changeState}
          className={check === 'create' ? 'active' : ''}
        >
          Create Comment
        </button>
        <button
          onClick={changeState}
          className={check === 'list' ? 'active' : ''}
        >
            Comment List
        </button>
        </div>
       {boxState}
        
        </div>
    )
}

const CreateComment = ({ id , showComment }) => {
    const userId = getUserId('userId');
    const sessionId = getCookie('sessionId');
    const [comment, setComment] = useState('');
    const [image, setImage] = useState(null);
    const [imagePreview, setImagePreview] = useState(null);

    const handleCommentChange = (e) => {
        setComment(e.target.value);
    };

    const handleImageChange = (e) => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.onload = () => {
            const base64 = reader.result;
        setImage(base64);
        setImagePreview(file)
        };
        reader.readAsDataURL(file);
    };
    
  const handleSubmit = (e) => {
    e.preventDefault();
    const commentData = {
    sessionId: sessionId,
    commentId: 0,
      postId: id,
      userId: parseInt(userId),
      creatorProfile: null,
      content: comment,
      image: image,
      date:"",
    };
    handleSubmitPost(commentData);
    setComment('');
    setImage(null);
    showComment();
  };

  const handleSubmitPost = (commentData) => {
    // Logic to handle the submission of the post data
  const check = checkCommentData(commentData);
  if
   (check.status === true ){
    const method = "POST";
    const type = "createComment"
    const payload = commentData
    fetchData(method,type,payload).then((data) => {

        /* !!todo!! if comment is sended it will go too comment list */
        console.log(data)
    })

  } else {
    alert(check.message)
  }
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
                {imagePreview && (
                    <div className="create_comment_image">
                        <img src={URL.createObjectURL(imagePreview)} alt="content" />
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
    console.log(comments)
    const loopComment = (comments) => {
        return comments.map((comment) => (
        <div className="comment_list_item" key={comment.commentId}>
                <div className="comment_list_item_header">
                    <div className="comment_list_item_body">
                    <p className="content">{comment.content}</p>
                </div>
                    <div className="comment_list_item_header_right">
                        <div className="comment_list_item_header_user">
                            <img src={comment.creatorProfile.avatar} alt="avatar" />
                            <p>{comment.creatorProfile.firstName} {comment.creatorProfile.lastName}</p>
                        </div>
                    </div>
                </div>
                <div className="comment_list_item_header_left">
                        <div className="comment_list_item_header_info">
                            <p>{comment.Date}</p>
                        </div>
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