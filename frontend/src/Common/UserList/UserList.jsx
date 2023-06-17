
const UserList = (title,user) => {
    const sessionId = getCookie("sessionId"); 
    return (
        <div className="user-list">
            <div className="user-list-container">
                <div className="user-list-header">
                    <h2>Users</h2>
                </div>
                <div className="user-list-body">
                    <ul>
                        {props.users.map((user) => (
                            <li key={user.id}>
                                <div className="user-list-item">
                                    <div className="user-list-item-img">
                                        <img src={user.avatar} alt="user-img" />
                                    </div>
                                    <div className="user-list-item-info">
                                        <div className="user-list-item-name">
                                            <span>{user.firstName}</span>
                                            <span> </span>
                                            <span>{user.lastName}</span>
                                        </div>
                                    </div>
                                </div>
                            </li>
                        ))}
                    </ul>
                </div>
                <div className="user-list-footer">
                    <button onClick={props.onClose}>Close</button>
                </div>
            </div>
        </div>
    )
}

export default UserList;