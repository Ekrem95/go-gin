import React, { Component } from 'react';
// import request from 'superagent';
import { store } from '../redux/reducers';

export default class Home extends Component {
  componentWillMount() {
    console.log(store.getState());
    if (store.getState() === 0) {
      window.location.href = '/login';
    }
  }

  render() {
    return (
      <div>
        <h1>Home</h1>
      </div>
    );
  }
}
