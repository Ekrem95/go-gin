import React, { Component } from 'react';
import { auth } from '../redux/reducers';

export default class Home extends Component {
  componentWillMount() {
    auth()
    .then(res => {
      if (res === 0) {
        window.location.href = '/login';
      }
    });
  }

  render() {
    return (
      <div>
        <h1>Home</h1>
      </div>
    );
  }
}
