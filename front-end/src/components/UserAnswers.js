import React, { useState, useEffect } from 'react';
import InfiniteScroll from 'react-infinite-scroll-component';
import { API_URL } from '../Env';


const UserAnswers = ({ user }) => {
  const [answers, setAnswers] = useState([]);
  const [hasMore, setHasMore] = useState(true);
  const [start, setStart] = useState(0);
  const count = 10; // Number of answers to load per request

  useEffect(() => {
    fetchAnswers();
  }, []);

  const fetchAnswers = async () => {
    try {
      const response = await fetch(`${API_URL}/get_user_profile/${user}?start=${start}&count=${count}`);
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

  const handleLike = (answerId) => {
    // Logic to handle liking an answer
  };

  const handleLikeCountClick = (answerId) => {
    // Logic to redirect to the page showing users who liked the answer
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
          <h3>Question: {answer.question_text}</h3>
          {!answer.is_author_anonymous && (
            <>
              {answer.author.display_name && <p>Author Display Name: {answer.author.display_name}</p>}
              {answer.author.username && <p>Author Username: {answer.author.username}</p>}
            </>
          )}
          <p>Answer: {answer.answer_text}</p>
          <p>Date Answered: {new Date(answer.date_answered).toLocaleString()}</p>
          <button onClick={() => handleLike(answer.id)}>Like</button>
          <span onClick={() => handleLikeCountClick(answer.id)}>
            {answer.like_count} Likes
          </span>
        </div>
      ))}
    </InfiniteScroll>
  );
};

export default UserAnswers;