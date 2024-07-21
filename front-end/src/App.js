import React from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import { GoogleOAuthProvider } from '@react-oauth/google';
import AppRoutes from './AppRoutes/AppRoutes';
import { GOOGLE_CLIENT_ID } from './Env';


const App = () => (
  <GoogleOAuthProvider clientId={GOOGLE_CLIENT_ID}>
    <Router>
      <AppRoutes />
    </Router>
  </GoogleOAuthProvider>
);

export default App;
