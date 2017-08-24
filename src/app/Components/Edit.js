import React, { Component } from 'react';
import request from 'superagent';

export default class Edit extends Component {
  constructor() {
    super();
    this.state = { post: null };
    this.edit = this.edit.bind(this);
  }

  componentWillMount() {
    request.get('/api/postbyid/' + this.props.location.pathname.split('/').pop())
      .then(res => {
        this.setState({ post: res.body.post });
      })
      .catch(e => e);
  }

  edit() {
    const title = this.refs.title.value;
    const description = this.refs.description.value;
    const src = this.refs.src.value;

    const pac = { title, description, src };

    request
      .post('/edit/' + this.props.location.pathname.split('/').pop())
      .type('form')
      .send(pac)
      .set('Accept', 'application/json')
      .then(res => {
        if (res.body.done === true) {
          this.props.history.push('/myposts');
        }
      })
      .catch(e => e);
  }

  render() {
    return (
      <div className="add">
        <h1>Edit</h1>
        {this.state.post &&
          <form>
            <input
              ref="title" type="text"
              placeholder="Title"
              defaultValue={this.state.post.title}
            />
            <input
              ref="description" type="text"
              placeholder="Description" defaultValue={this.state.post.description}
            />
            <input
              ref="src" type="text"
              placeholder="Image Source" defaultValue={this.state.post.src}
            />
            <button
              onClick={() => {
                this.edit();
              }}

              type="button"
              >Add
            </button>
          </form>
        }
      </div>
    );
  }
}
