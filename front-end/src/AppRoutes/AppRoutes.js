import React from 'react';
import { Route, Routes } from 'react-router-dom';
import Home from '../components/Home';
import HelloWorld from '../components/HelloWorld';


const AppRoutes = () => (
  <Routes>
    <Route exact path="/" element={<Home/>} />
    <Route exact path="/hello_world" element={<HelloWorld/>} />
  </Routes>
);

export default AppRoutes;