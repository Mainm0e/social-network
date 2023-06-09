import "./user.css"
import CreatePost from "./CreatePost"
import PostList from "./PostList"
const Body = ({id}) => {
    const handleSubmitPost = (postData) => {
        // Logic to handle the submission of the post data
        console.log('Submitted post:', postData);
        // Make API requests or perform other operations here
      };
    
    
    return (
        <div className="main_body">
            <CreatePost id={id} onSubmit={handleSubmitPost}/>
            <PostList/>
            
        </div>
    )
}

export default Body