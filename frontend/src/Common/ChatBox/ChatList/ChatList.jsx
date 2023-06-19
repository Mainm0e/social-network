import {dummyUsers} from '../DummyData';
import './ChatList.css';
const ChatList = ({ onUserSelection }) => {
    const handleUserClick = (selectedUser) => {
      // Pass the selected user data to the parent component
      onUserSelection(selectedUser);
    };
  
    return (
      <div className='chat-list'>
        {dummyUsers.map((user) => (
          <div
            className='chat-list-item'
            key={user.id}
            onClick={() => handleUserClick(user)}
          >
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
    );
  };
  

export default ChatList;