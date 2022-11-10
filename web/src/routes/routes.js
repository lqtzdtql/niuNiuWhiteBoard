import ForgetPassword from '@Src/view/forget/ForgetPassword';
import Home from '@Src/view/home/Home';
import Login from '@Src/view/login/Login';
import Register from '@Src/view/register/Register';
import Room from '@Src/view/room/Room';
import { Navigate } from 'react-router-dom';

export const routes = [
  {
    path: '/',
    element: <Login />,
  },
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/register',
    element: <Register />,
  },
  {
    path: '/home',
    element: <Home />,
  },
  {
    path: '/forget',
    element: <ForgetPassword />,
  },
  {
    path: '/room',
    element: <Room />,
  },
  {
    path: '*',
    element: <Navigate to="/" />,
  },
];
