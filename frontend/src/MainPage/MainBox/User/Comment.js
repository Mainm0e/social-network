
const Comment = ({id}) => {

    // return createcomment and commentlist button and default commentlist

    return(
        <div className="comment">
           <button className="comment-button">Create Comment</button>
           <button className="comment-button">Comment List</button>
            <CommentList id={id}/>
            <CreateComment id={id}/>
        </div>
    )
}
export default Comment;

const CreateComment = ({id}) => {
   // return form for creating comment
   //form should have input for comment and button for submitting and add img
       return(
        <div className="create-comment">
            <form>
                <input type="text" placeholder="Comment"/>
                <button type="submit">Submit</button>
                
                
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