import React, { Component } from 'react';
import Form from '../common/Form';

export default class Login extends Component {
  render() {
    return (
      <Form
        header="Login"
        history={this.props.history}
        post={'/login'}
      />
    );
  }
}
