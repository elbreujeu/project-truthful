import React, { useState, useEffect } from 'react';
import { API_URL } from '../Env';

const HelloWorld = () => {
  const [message, setMessage] = useState('');

  useEffect(() => {
    fetch(API_URL + '/hello_world')
      .then(res => res.json())
      .then(data => {
        setMessage(data.message);
      })
      .catch(error => console.error(error));
  }, []);

  return (
    <div>
      <h2>Hello World</h2>
      <p>{message}</p>
    </div>
  );
};

export default HelloWorld;