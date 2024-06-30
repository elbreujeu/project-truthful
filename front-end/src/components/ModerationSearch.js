import React, { useState } from 'react';
import '../styles/style.css';
import { API_URL } from '../Env';

function ModerationSearch() {
    const [inputValue, setInputValue] = useState('');
    const [error, setError] = useState('');
    const [userInfo, setUserInfo] = useState(null);

    const handleInputChange = (e) => {
        setInputValue(e.target.value);
    };

    const CheckInputEmpty = (e) => {
        e.preventDefault();
        // Handle the form submission here
        if (inputValue.trim() === '') {
            setError('Please enter a value');
            setUserInfo(null);
            return false;
        }
        return true;
    };

    const LookUpUser = (e) => {
        e.preventDefault();
        
        // Do a query to API_URL/get_user_profile/:username to check if the user exists
        fetch(`${API_URL}/get_user_profile/${inputValue}`)
            .then(response => {
                if (response.ok) {
                    response.json().then(data => {
                        setUserInfo(data);
                    });
                } else if (response.status === 404) {
                    setError('User not found');
                    setUserInfo(null);
                } else if (response.headers.get('content-type').includes('application/json')) {
                    response.json().then(data => {
                        setError(data.message);
                        setUserInfo(null);
                    });
                } else {
                    setError('An error occurred');
                    setUserInfo(null);
                }
            })
            .catch(error => {
                console.error('Error:', error);
                setError('Error reaching server');
                setUserInfo(null);
            });
    };

    const handleSubmit = (e) => {
        if (!CheckInputEmpty(e)) {
            return;
        }
        setError('');
        LookUpUser(e);
    };

    return (
        <div>
            {error && <div className="error_box" style={{marginBottom: '1rem'}}>{error}</div>}
            <h1>Moderation</h1>
            <form onSubmit={handleSubmit}>
                <label htmlFor="input">Enter Username:</label>
                <input
                    type="text"
                    id="input"
                    value={inputValue}
                    onChange={handleInputChange}
                />
                <button type="submit">Submit</button>
            </form>
            {userInfo && (
                <div>
                    <h2>User Information</h2>
                    <p>Username: {userInfo.username}</p>
                    <p>Display name: {userInfo.display_name}</p>
                    <p>Number of answered questions: {userInfo.answer_count}</p>
                    <p>Number of followers: {userInfo.follower_count}</p>
                    <p>Number of followings: {userInfo.following_count}</p>
                </div>
            )}
        </div>
    );
}

export default ModerationSearch;