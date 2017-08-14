import React, { Component } from 'react';
import request from 'superagent';

export default class Home extends Component {
  componentWillMount() {
    request
      .get('/user')
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
