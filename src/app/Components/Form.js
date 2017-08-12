import React, { Component } from 'react';
import request from 'superagent';

export default class Form extends Component {
  constructor() {
    super();
    this.state = { val: null };
  }

  render() {
    return (
      <div>
      <input
        autoFocus
        type="text"
        onChange={(e) => {
          this.setState({ val: e.target.value });
        }}

      />
      <button
        type="button"
        onClick={() => {
          const payload = { name: this.state.val, age: 22 };

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
