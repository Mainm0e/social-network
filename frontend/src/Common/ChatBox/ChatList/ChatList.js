import {dummyUsers} from '../DummyData';
import './ChatList.css';
const ChatList = () => {
    return (
        <div className='chat-list'>
            {dummyUsers.map((user) => (
                <div className='chat-list-item' key={user.id}>
                    <div className='chat-list-item-avatar'>
                        <img src={user.avatar} alt={user.name} />
                    </div>
                    <div className='chat-list-item-content'>
                        <div className='chat-list-item-content-name'>
                            {user.name}
                        </div>
                    </div>
                </div>
            ))}
        </div>
    )
}

export default ChatList;