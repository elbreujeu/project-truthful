import React from 'react';
import '../styles/style.css';
import '../styles/UserInfo.css';

const UserInfo = (userData) => {
    const userInfo = userData.userInfo;
    const followUser = () => {
        console.log('Follow user')
        // Backend POST request to API_URL/follow_user
        // TODO : Write backend call here
    };
    return (
        <div className="user-info">
            <div className="profile-picture"> {/* Wrap the image in a .profile-picture element */}
                <img src="https://i.imgur.com/Jvh1OQm.jpeg" alt="Profile" />
            </div>
            <p>@{userInfo.username}</p> {/* Username */}
            <p>{userInfo.display_name}</p> {/* Display name */}
            <div className="profile-stats">
                <a href={`/profile/${userInfo.username}/followers`}>{userInfo.follower_count} followers</a> {/* Follower count */}
                <p>{userInfo.answer_count} answers</p> {/* Answer count */}
                <a href={`/profile/${userInfo.username}/following`}>{userInfo.following_count} following</a> {/* Following count */}
            </div>
            <button className="button" onClick={followUser}>Follow</button> {/* Follow button TODO : add a route in backend to check if the user is following the user. If so, turn the button into an "unfollow button */}
        </div>
    );
};

export default UserInfo;