import React, { Component } from 'react';
import { store } from '../redux/reducers';

export default class Add extends Component {
  constructor() {
    super();
    this.add = this.add.bind(this);
  }

  add() {
    const title = this.refs.title.value;
    const description = this.refs.desc.value;
    const src = this.refs.src.value;
    const postedBy = store.getState().user.user;

    fetch('/add', {
      method: 'post',
      body: JSON.stringify({ title, description, src, postedBy }),
    })
    .then(res => res.json())
    .then(res => {
      if (res.done === true) {
        this.props.history.push('/');
      }
    });
  }

  render() {
    return (
      <div className="add">
        <h1>Add</h1>
        <form>
          <input ref="title" type="text" placeholder="Title"/>
          <input ref="desc" type="text" placeholder="Description"/>
          <input ref="src" type="text" placeholder="Image Source"/>
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
