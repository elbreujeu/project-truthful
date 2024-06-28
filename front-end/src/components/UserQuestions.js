import React, { useState, useEffect } from 'react';
import InfiniteScroll from 'react-infinite-scroll-component';
import { API_URL } from '../Env';
import { useParams } from 'react-router-dom';
import '../styles/UserQuestions.css'


const QuestionList = ({ user }) => {
    // gets the "token" cookie
    const cookieElement = document.cookie.split('; ').find(row => row.startsWith('token='));
    const token = cookieElement ? cookieElement.split('=')[1] : null;
    if (!token) {
        console.error('No token found');
        window.location.href = '/login';
    }
  const [questions, setQuestions] = useState([]);
  const [hasMore, setHasMore] = useState(true);
  const [start, setStart] = useState(0);
  const count = 10; // Number of questions to load per request

  useEffect(() => {
    fetchQuestions();
  }, []);

  const fetchQuestions = async () => {
    try {
    const response = await fetch(`${API_URL}/get_questions?start=${start}&count=${count}`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      }
    });
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      const newQuestions = await response.json();

      if (newQuestions) {
        // Combine new questions with existing ones, avoiding duplicates
        const combinedQuestions = [...questions, ...newQuestions.filter(newQuestion => !questions.some(question => question.id === newQuestion.id))];
        setQuestions(combinedQuestions);
        setStart(prevStart => prevStart + count);

        if (newQuestions.length < count) {
          setHasMore(false);
        }
      } else {
        setHasMore(false);
      }
    } catch (error) {
      console.error('Error fetching questions', error);
    }
  };
  
  const sendAnswer = async (questionId) => {
    // selects the question-box from the questionId
    const questionBox = document.querySelector(`.question-box[question-id="${questionId}"]`);
    const answer = questionBox.querySelector('input').value;
    
    // makes a POST request to the API to send the answer
    const response = await fetch(`${API_URL}/answer_question`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        question_id: questionId,
        text: answer
      })
    });
    if (!response.ok) {
      console.error('Error sending answer');
      return;
    }
    // if the answer was sent successfully, delete the question from the DOM
    if (response.ok){
    questionBox.remove();
  }
  };

  return (
    <InfiniteScroll
      dataLength={questions.length}
      next={fetchQuestions}
      hasMore={hasMore}
    >
      {questions.length === 0 ? (
        <p>No questions available</p>
      ) : (
        questions.map(question => (
          <div className='question-box' key={question.id} question-id={question.id}>
            <p className='question-text'>{question.text}</p>
            {!question.is_author_anonymous && (
              <p className='author-name' onClick={() => window.location.href = `/profile/${question.author.username}`}>
                {question.author.display_name}
              </p>
            )}
            <p className='date'>{new Date(question.created_at).toLocaleString()}</p>
            <div className='answer-box'>
              <input type='text' placeholder='Type your answer...' />
              <button onClick={() => sendAnswer(question.id)}>Send</button>
            </div>
          </div>
        ))
      )}
    </InfiniteScroll>
  );
};

const UserQuestions = ({ }) => {
  const params = useParams();
  const user = params.user;

  return (
    <QuestionList user={user} />
  );
}

export default UserQuestions;