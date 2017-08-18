import React, { Component } from 'react';
import request from 'superagent';

export default class Add extends Component {
  constructor() {
    super();
    this.add = this.add.bind(this);
  }

  add() {
    const title = this.refs.title.value;
    const desc = this.refs.desc.value;
    const src = this.refs.src.value;

    const pac = { title, desc, src };

    request
      .post('/add')
      .type('form')
      .send(pac)
      .set('Accept', 'application/json')
      .end(function (err, res) {
      // Calling the end function will send the request
    });

  }

  render() {
    return (
      <div className="add">
        <h1>Add</h1>
        <form>
          <input ref="title" type="text"/>
          <input ref="desc" type="text"/>
          <input ref="src" type="text"/>
          <button
            onClick={() => {
              this.add();
            }}

            type="button"
            >Add
          </button>
        </form>
      </div>
    );
  }
}
