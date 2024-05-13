import React from 'react';
import '../styles/style.css';
import '../styles/UserInfo.css';
import { API_URL } from '../Env';
import { useState } from 'react';

const UserInfo = (userData) => {
    const userInfo = userData.userInfo;
    const [errorBox, setErrorBox]  = useState('');
    const [successBox, setSuccessBox]  = useState('');
    const followUser = () => {
        console.log('Follow user')
        // gets the "token" cookie
        const cookieElement = document.cookie.split('; ').find(row => row.startsWith('token='));
        const token = cookieElement ? cookieElement.split('=')[1] : null;
        if (!token) {
            console.error('No token found');
            window.location.href = '/login';
        }
        // Backend POST request to API_URL/follow_user
        // if request is unsuccessful, console.error the error and put red textbox above the page
        fetch(`${API_URL}/follow_user`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ user_id: userInfo.id, follow: true})
        })
            .then(response => {
                // gets the response json
                if (!response.ok) {
                    response.json().then(data => {
                        console.error(data.error);
                        // sets the error box to the error message
                        setErrorBox(data.error);
                    });
                } else {
                    response.json().then(data => {
                        console.log(data);
                        // sets the success box to the success message
                        setSuccessBox(data.message);
                    });
                }
                // Handle successful response here
            })
            .catch(error => {
                console.error(error);
                // Handle error here
            });
    };
    const unfollowUser = () => {
        console.log('Unfollow user')
        // gets the "token" cookie
        const cookieElement = document.cookie.split('; ').find(row => row.startsWith('token='));
        const token = cookieElement ? cookieElement.split('=')[1] : null;
        if (!token) {
            console.error('No token found');
            window.location.href = '/login';
        }
        // Backend POST request to API_URL/follow_user
        // if request is unsuccessful, console.error the error and put red textbox above the page
        fetch(`${API_URL}/follow_user`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ user_id: userInfo.id, follow: false})
        })
            .then(response => {
                // gets the response json
                if (!response.ok) {
                    response.json().then(data => {
                        console.error(data.error);
                        // sets the error box to the error message
                        setErrorBox(data.error);
                    });
                } else {
                    response.json().then(data => {
                        console.log(data);
                        // sets the success box to the success message
                        setSuccessBox(data.message);
                    });
                }
                // Handle successful response here
            })
            .catch(error => {
                console.error(error);
                // Handle error here
            });
    };
    return (
        <div className="user-info">
            {errorBox && <div className="error_box" style={{marginBottom: '1rem'}}>{errorBox}</div>}
            {successBox && <div className="success_box" style={{marginBottom: '1rem'}}>{successBox}</div>}
            <div className="profile-picture"> {/* Wrap the image in a .profile-picture element */}
                <img src="https://i.imgur.com/Jvh1OQm.jpeg" alt="Profile" />
            </div>

            <p><strong>{userInfo.display_name}</strong></p> {/* Display name */}
            <p>@{userInfo.username}</p> {/* Username */}

            <div className="profile-stats">
                <a href={`/profile/${userInfo.username}/followers`}>{userInfo.follower_count} followers</a> {/* Follower count */}
                {userInfo.answer_count} answers {/* Answer count */}
                <a href={`/profile/${userInfo.username}/following`}>{userInfo.following_count} following</a> {/* Following count */}
            </div>
            {/* Follow button TODO : add a route in backend to check if the user is following the user. If so, turn the button into an "unfollow button */}
            <button className="button" onClick={followUser}>Follow</button>
            <p></p>
            <button className="button" onClick={unfollowUser}>Unfollow</button>
        </div>
    );
};

export default UserInfo;