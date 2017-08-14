import React, { Component } from 'react';
import { Link } from 'react-router-dom';

export default class Nav extends Component {
  render() {
    return (
      <div className="nav">
        <Link to="/">Home</Link>
        <Link to="/signup">Sign up</Link>
        <Link to="/login">Login</Link>
      </div>
    );
  }
}
