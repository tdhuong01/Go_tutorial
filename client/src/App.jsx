import React, { useState } from 'react';
import './App.css';

function App() {
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({
    first_name: '',
    last_name: '',
    username: '',
    email: '',
    phone: '',
    password: '',
    confirm_password: ''
  });
  const [message, setMessage] = useState('');

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const url = isLogin ? '/api/login' : '/api/register';
    const data = isLogin ? { username: formData.username, password: formData.password } : formData;

    try {
      const response = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(data)
      });
      const result = await response.json();
      setMessage(result.message);
      if (result.success && isLogin) {
        window.location.href = '/';
      }
    } catch (error) {
      setMessage('Lỗi kết nối');
    }
  };

  return (
    <div className="container">
      <h2>{isLogin ? 'Sign In To Admin' : 'Register'}</h2>
      <p>{isLogin ? 'Enter your details to log in to your account:' : 'Enter your details to create an account:'}</p>
      <form id={isLogin ? "loginForm" : "registerForm"} onSubmit={handleSubmit}>
        {!isLogin && (
          <>
            <div className="input-group">
              <input type="text" name="first_name" placeholder="First Name" onChange={handleChange} required />
            </div>
            <div className="input-group">
              <input type="text" name="last_name" placeholder="Last Name" onChange={handleChange} required />
            </div>
            <div className="input-group">
              <input type="email" name="email" placeholder="Email" onChange={handleChange} required />
            </div>
            <div className="input-group">
              <input type="tel" name="phone" placeholder="Phone" onChange={handleChange} />
            </div>
          </>
        )}
        <div className="input-group">
          <input type="text" name="username" placeholder={isLogin ? "Email" : "Username"} onChange={handleChange} required />
        </div>
        <div className="input-group">
          <input type="password" name="password" placeholder="Password" onChange={handleChange} required />
        </div>
        {!isLogin && (
          <div className="input-group">
            <input type="password" name="confirm_password" placeholder="Confirm Password" onChange={handleChange} required />
          </div>
        )}
        <button className="btn" type="submit">{isLogin ? 'Sign In' : 'Register'}</button>
      </form>
      <div className="forgot">I've forgotten my password.</div>
      <div className="register-link">
        <a href="#" onClick={(e) => { e.preventDefault(); setIsLogin(!isLogin); }}>
          {isLogin ? "Don't have an account? Register here" : "Already have an account? Sign in"}
        </a>
      </div>
      <div id="result">{message}</div>
    </div>
  );
}

export default App;