import React, { Component } from 'react';
import request from 'superagent';

export default class Form extends Component {
  constructor() {
    super();
    this.state = { name: '', password: '' };
  }

  render() {
    return (
      <div className="form">
        <h2>Login</h2>
        <input
          autoFocus
          type="text"
          onChange={(e) => {
            this.setState({ name: e.target.value });
          }}

        />
        <input
          type="password"
          onChange={(e) => {
            this.setState({ password: e.target.value });
          }}

        />
        <button
          type="button"
          onClick={() => {
            if (
              this.state.name.length > 2 &&
              this.state.password.length > 5
            ) {
              const payload = {
                username: this.state.name, password: this.state.password,
              };

              request
                .post('/login')
                .type('form')
                .send(payload)
                .set('Accept', 'application/json')
                .then(res => {
                  console.log(res.body);
                });
            }
          }}
          >Login</button>
      </div>
    );
  }
}
