import React from 'react';
import { Route, Routes } from 'react-router-dom';
import Home from '../components/Home';
import HelloWorld from '../components/HelloWorld';
import Login from '../components/Login';
import Register from '../components/Register';
import Profile from '../components/Profile';
import Follows from '../components/Follows';
import UserQuestions from '../components/UserQuestions';

const AppRoutes = () => (
  <Routes>
    <Route exact path="/" element={<Home/>} />
    <Route exact path="/hello_world" element={<HelloWorld/>} />
    <Route exact path="/login" element={<Login/>} />
    <Route exact path="/register" element={<Register/>} />
    <Route exact path="/profile/:user" element={<Profile />} />
    <Route exact path="/profile/:user/followers" element={<Follows />} />
    <Route exact path="/profile/:user/following" element={<Follows />} />
    <Route exact path="/questions" element={<UserQuestions />} />
  </Routes>
);

export default AppRoutes;