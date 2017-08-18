import React, { Component } from 'react';
import { auth } from '../redux/reducers';
import request from 'superagent';

export default class Home extends Component {
  componentWillMount() {
    auth()
    .then(res => {
      if (res.auth.auth === 0) {
        this.props.history.push('/login');
      }
    });

    request.get('/api/posts')
      .then(res => {
        console.log(res.body);
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
