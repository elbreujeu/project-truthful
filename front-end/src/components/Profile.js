import React from 'react';
import { useParams } from 'react-router-dom';
import UserInfo from './UserInfo';
import { useEffect } from 'react';
import { useState } from 'react';
import { API_URL } from '../Env';

const Profile = ({ }) => {
    const params = useParams();
    const user = params.user;
  
    // Step 1: Create a state variable to store the user data
    // Step 2: Create a useEffect hook to fetch the user data
    // Step 3 : Render the user info component with the user data (beware, a new component "UserAnswers" will be added in the future)

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