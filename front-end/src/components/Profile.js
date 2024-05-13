import React from 'react';
import { useParams } from 'react-router-dom';
import UserInfo from './UserInfo';
import { useEffect } from 'react';
import { useState } from 'react';
import { API_URL } from '../Env';

const Profile = ({ }) => {
    const params = useParams();
    const user = params.user;

    const [userData, setUserData] = useState(null);

    useEffect(() => {
      const fetchUserData = async () => {
        try {
          const response = await fetch(`${API_URL}/get_user_profile/${user}`);
          const data = await response.json();
          setUserData(data);
        } catch (error) {
          console.error('Error fetching user data:', error);
        }
      };

      fetchUserData();
    }, [user]);

    if (!userData) {
      return <div>Loading...</div>;
    }

    return (
      <div>
        <UserInfo userInfo={userData} />
        {/* Render the UserAnswers component here */}
      </div>
    );
  };  

export default Profile;