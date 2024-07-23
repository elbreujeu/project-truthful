import React, { useState, useEffect } from 'react';
import InfiniteScroll from 'react-infinite-scroll-component';
import { API_URL } from '../Env';
import '../styles/UserAnswers.css';


const UserAnswers = ({ user }) => {
  const [answers, setAnswers] = useState([]);
  const [hasMore, setHasMore] = useState(true);
  const [start, setStart] = useState(0);
  const count = 10; // Number of answers to load per request

  const cookieElement = document.cookie.split('; ').find(row => row.startsWith('token='));
  const token = cookieElement ? cookieElement.split('=')[1] : null;

  useEffect(() => {
    fetchAnswers();
  }, []);

  const fetchAnswers = async () => {
    try {
      const response = await fetch(`${API_URL}/get_user_profile/${user}?start=${start}&count=${count}`, {
        headers: token !== null ? {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        } : {}
      });
      if (!response.ok) {
        throw new Error('Network response was not ok');
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
      console.error('Error fetching answers', error);
    }
  };

  const handleLike = (answerId, isLiking) => {
    if (!token) {
        console.error('No token found');
        window.location.href = '/login';
    }

    fetch(`${API_URL}/like_answer`, {
      method: 'POST',
      body: JSON.stringify({ answer_id: answerId, like: isLiking }),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      }
    })

    const likeButton = document.getElementById(`like-button-${answerId}`);
    likeButton.innerText = isLiking ? 'Unlike' : 'Like';

    const likeCount = document.getElementById(`like-count-${answerId}`);
    likeCount.innerText = `${isLiking ? parseInt(likeCount.innerText) + 1 : parseInt(likeCount.innerText) - 1} Likes`;
  };

  const handleLikeCountClick = (answerId) => {
    window.location.href = `/answers/${answerId}/likes`;
  };

  const handleAuthorClick = (authorUsername) => {
    window.location.href = `/profile/${authorUsername}`;
  };

  return (
    <InfiniteScroll
      dataLength={answers.length}
      next={fetchAnswers}
      hasMore={hasMore}
      loader={<h4>Loading...</h4>}
      endMessage={<p>No more answers</p>}
    >
      {answers.map(answer => (
        <div key={answer.id} className="answer">
          <h3 className='question'>{answer.question_text}</h3>
          {!answer.is_author_anonymous && answer.author.display_name && (
            <>
              <span onClick={() => handleAuthorClick(answer.author.username)} className='author'>{answer.author.display_name}</span>
            </>
          )}
          <p>{answer.answer_text}</p>
          <button id={`like-button-${answer.id}`} onClick={() => {
            handleLike(answer.id, !answer.liked_by_requester);
            answer.liked_by_requester = !answer.liked_by_requester;
          }}>
            {answer.liked_by_requester ? 'Unlike' : 'Like'}
          </button>
          <span id={`like-count-${answer.id}`} onClick={() => handleLikeCountClick(answer.id)}>
            {answer.like_count} Likes
          </span>
          <p className='date'>{new Date(answer.date_answered).toLocaleString()}</p>
        </div>
      ))}
    </InfiniteScroll>
  );
};

export default UserAnswers;