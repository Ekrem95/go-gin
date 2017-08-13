import React, { Component } from 'react';
import request from 'superagent';

export default class Form extends Component {
  constructor() {
    super();
    this.state = { name: null, password: null };
  }

  render() {
    return (
      <div className="form">
        <h2>Sign up</h2>
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
          const payload = { name: this.state.name, password: this.state.password };

          request
            .post('/form_post')
            .type('form')
            .send(payload)
            .set('Accept', 'application/json')
            .then(res => {
              console.log(res.body);
            });
        }}
        >Send</button>
      </div>
    );
  }
}
