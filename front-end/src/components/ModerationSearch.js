import React, { useState } from 'react';
import '../styles/style.css';
import { API_URL } from '../Env';
import InfiniteScroll from 'react-infinite-scroll-component';

function ModerationSearch() {
    const [inputValue, setInputValue] = useState('');
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [userInfo, setUserInfo] = useState(null);
    const [disciplineUser, setDisciplineUser] = useState('');
    const [userAnswerDisplay, setUserAnswerDisplay] = useState('');
    const cookieElement = document.cookie.split('; ').find(row => row.startsWith('token='));
    const token = cookieElement ? cookieElement.split('=')[1] : null;

    
    const [answers, setAnswers] = useState([]);
    const [hasMore, setHasMore] = useState(true);
    const [start, setStart] = useState(0);
    const count = 10; // Number of answers to load per request

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

    const fetchAnswers = async () => {
        try {
          const response = await fetch(`${API_URL}/get_user_profile/${userInfo.username}?start=${start}&count=${count}`);
          if (!response.ok) {
            setError('Failed to fetch answers');
          }
          const data = await response.json();
          const newAnswers = data.answers;
    
          if (newAnswers) {
            // Combine new answers with existing ones, avoiding duplicates
            const combinedAnswers = [...answers, ...newAnswers.filter(newAnswer => !answers.some(answer => answer.id === newAnswer.id))];
            setAnswers(combinedAnswers);
            setStart(prevStart => prevStart + count);
    
            if (newAnswers.length < count) {
              setHasMore(false);
            }
          } else {
            setHasMore(false);
          }
        } catch (error) {
            setError('Error fetching answers');
        }
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

    const handleUserAnswers = () => {
        if (!userAnswerDisplay) {
            setUserAnswerDisplay(true);
            fetchAnswers();
            return;
        } else {
            setUserAnswerDisplay(false);
            setAnswers([]);
            setHasMore(true);
            setStart(0);
            return;
        }
    };

    const resetHandlers = () => {
        setError('');
        setSuccess('');
        setUserInfo(null);
        setDisciplineUser(false);
        setUserAnswerDisplay(false);
        setAnswers([]);
        setHasMore(true);
        setStart(0);
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
                    <button onClick={handleUserAnswers}>List User answers</button>
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
            {userAnswerDisplay && (
                <InfiniteScroll
                    dataLength={answers.length}
                    next={fetchAnswers}
                    hasMore={hasMore}
                    endMessage={<p>No more answers</p>}
                >
                    {answers.map(answer => (
                        <div key={answer.id} className="answer">
                            <h3 className='question'>{answer.question_text}</h3>
                            {!answer.is_author_anonymous && answer.author.display_name && (
                                <p>Author display name: {answer.author.display_name} Author username: {answer.author.username}</p>
                            )}
                            <p>{answer.answer_text}</p>
                            <p>{answer.like_count} Likes</p>
                            <p className='date'>{new Date(answer.date_answered).toLocaleString()}</p>
                        </div>
                    ))}
                </InfiniteScroll>
            )}
        </div>
    );
}

export default ModerationSearch;