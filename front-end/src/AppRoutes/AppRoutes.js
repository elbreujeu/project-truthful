import React from 'react';
import { Route, Routes } from 'react-router-dom';
import Home from '../components/Home';
import HelloWorld from '../components/HelloWorld';
import Login from '../components/Login';
import Register from '../components/Register';

const AppRoutes = () => (
  <Routes>
    <Route exact path="/" element={<Home/>} />
    <Route exact path="/hello_world" element={<HelloWorld/>} />
    <Route exact path="/login" element={<Login/>} />
    <Route exact path="/register" element={<Register/>} />
  </Routes>
);

export default AppRoutes;