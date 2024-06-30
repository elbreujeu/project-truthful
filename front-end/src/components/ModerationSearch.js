import React, { useState } from 'react';
import '../styles/style.css';
import { API_URL } from '../Env';

function ModerationSearch() {
    const [inputValue, setInputValue] = useState('');
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [userInfo, setUserInfo] = useState(null);
    const [disciplineUser, setDisciplineUser] = useState('');
    const cookieElement = document.cookie.split('; ').find(row => row.startsWith('token='));
    const token = cookieElement ? cookieElement.split('=')[1] : null;

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
                        resetHandlers();
                        setUserInfo(data);
                    });
                } else if (response.status === 404) {
                    resetHandlers();
                    setError('User not found');
                } else if (response.headers.get('content-type').includes('application/json')) {
                    response.json().then(data => {
                        resetHandlers();
                        setError(data.error);
                    });
                } else {
                    resetHandlers();
                    setError('An error occurred');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                setError('Error reaching server');
                setUserInfo(null);
            });
    };

    const handleUserSubmit = (e) => {
        if (!CheckInputEmpty(e)) {
            return;
        }
        resetHandlers();
        LookUpUser(e);
    };

    const handleBanSubmit = (e) => {
        e.preventDefault();

        const reason = document.getElementById('reason').value;
        const isPermanent = document.getElementById('permanent').checked;
        const duration = isPermanent ? 0 : document.getElementById('duration').value;

        const durationInt = parseInt(duration);
        if (isNaN(durationInt) || durationInt < 0) {
            setError('Invalid duration');
            return;
        }

        console.log('Banning user', userInfo.id, 'for', durationInt, 'days with reason:', reason)

        fetch(`${API_URL}/moderation/ban_user`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                user_id: userInfo.id,
                reason: reason,
                duration: durationInt
            })
        }).then(response => {
            if (response.ok) {
                setSuccess('User banned successfully');
            } else if (response.headers.get('content-type').includes('application/json')) {
                response.json().then(data => {
                    setError(data.error);
                });
            } else {
                setError('An error occurred');
            }
        } ).catch(error => {
            console.error('Error:', error);
            setError('Error reaching server');
        });
    }

    const handleDisciplineUser = () => {
        if (!disciplineUser) {
            setDisciplineUser(true);
            return;
        } else {
            setDisciplineUser(false);
            return;
        }
    };

    const resetHandlers = () => {
        setError('');
        setSuccess('');
        setUserInfo(null);
        setDisciplineUser(false);
    }

    return (
        <div>
            {error && <div className="error_box" style={{marginBottom: '1rem'}}>{error}</div>}
            {success && <div className="success_box" style={{marginBottom: '1rem'}}>{success}</div>}
            <h1>Moderation</h1>
            <form onSubmit={handleUserSubmit}>
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
                    <button onClick={handleDisciplineUser}>Discipline User</button>
                </div>
            )}
            {disciplineUser && (
                <div>
                    <h2>Discipline User</h2>
                    <form onSubmit={handleBanSubmit}>
                        <label htmlFor="reason">Reason:</label>
                        <input type="text" id="reason" />
                        <label htmlFor="duration">Duration:</label>
                        <input type="text" id="duration" />
                        <label htmlFor="permanent">Permanent:</label>
                        <input type="checkbox" id="permanent" />
                        <button type="submit">Submit</button>
                    </form>
                </div>
            )}
        </div>
    );
}

export default ModerationSearch;