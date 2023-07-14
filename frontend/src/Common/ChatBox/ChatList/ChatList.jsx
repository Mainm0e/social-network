import { useEffect } from 'react';
import {dummyUsers} from '../DummyData';
import './ChatList.css';
import { getCookie, getUserId } from '../../../tools/cookie';
import { fetchData } from '../../../tools/fetchData';
import { useState } from 'react';
const ChatList = ({ onUserSelection }) => {
  const handleUserClick = (selectedUser) => {
    // Pass the selected user data to the parent component
    onUserSelection(selectedUser);
  };
  const [contactList, setContactList] = useState([]);

  useEffect(() => {
    const fetchDataAsync = async () => {
      const users = await getContactList();
      setContactList(users);
    };

    fetchDataAsync();
  }, []);

  return (
    <div className='chat-list'>
      {contactList.map((user) => (
        <div
          className='chat-list-item'
          key={user.userId}
          onClick={() => handleUserClick(user)}
        >
          <div className='chat-list-item-avatar'>
            <img src={user.avatar} alt={user.firstName} />
          </div>
          <div className='chat-list-item-content'>
            <div className='chat-list-item-content-name'>{user.firstName}</div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default ChatList;

const getContactList = async () => {
  const method = 'POST';
  const type = 'profileList';
  const payload = {
    sessionId: getCookie('sessionId'),
    userId: getUserId('userId'),
    request: 'followers',
  };
  const payload2 = {
    sessionId: getCookie('sessionId'),
    userId: getUserId('userId'),
    request: 'followings',
  };

  // Fetch both sets of data concurrently
  const [followersData, followingsData] = await Promise.all([
    fetchData(method, type, payload),
    fetchData(method, type, payload2),
  ]);

  // Make contact list from both sets of data
  const contactList = [];
  if (followersData !== null) {
    followersData.forEach((user) => {
      contactList.push(user);
    });
  }
  if (followingsData !== null) {
    followingsData.forEach((user) => {
      const existingUser = contactList.find((contact) => contact.userId === user.userId);
      if (!existingUser) {
        contactList.push(user);
      }
    });
  }
  return contactList;
};
