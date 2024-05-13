import React from 'react';
import '../styles/style.css';

const UserInfo = (userData) => {
    const userInfo = userData.userInfo;
    return (
        <div>
            <h2>Username: {userInfo.username}</h2>
            <h2>Email Address: {userInfo.display_name}</h2>
        </div>
    );
};

export default UserInfo;