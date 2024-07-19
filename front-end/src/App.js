import React from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import { GoogleOAuthProvider } from '@react-oauth/google';
import AppRoutes from './AppRoutes/AppRoutes';
import { clientId } from './Env';


const App = () => (
  <GoogleOAuthProvider clientId={clientId}>
    <Router>
      <AppRoutes />
    </Router>
  </GoogleOAuthProvider>
);

export default App;
