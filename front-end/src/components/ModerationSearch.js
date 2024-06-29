import React, { useState } from 'react';
import '../styles/style.css';
import { API_URL } from '../Env';

function ModerationSearch() {
    const [inputValue, setInputValue] = useState('');
    const [error, setError] = useState('');

    const handleInputChange = (e) => {
        setInputValue(e.target.value);
    };

    const CheckInputEmpty = (e) => {
        e.preventDefault();
        // Handle the form submission here
        if (inputValue.trim() === '') {
            setError('Please enter a value');
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
                    // handle the ok response
                } else if (response.status === 404) {
                    setError('User not found');
                } else if (response.headers.get('content-type').includes('application/json')) {
                    response.json().then(data => {
                        setError(data.message);
                    });
                } else {
                    setError('An error occurred');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                setError('Error reaching server');
            });
    }

    const handleSubmit = (e) => {
        if (!CheckInputEmpty(e)) {
            return;
        }
        setError('');
        console.log(LookUpUser(e));
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
        </div>
    );
}

export default ModerationSearch;