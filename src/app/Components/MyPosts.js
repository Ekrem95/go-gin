import React, { Component } from 'react';
import { store } from '../redux/reducers';
import request from 'superagent';

export default class MyPosts extends Component {
  constructor() {
    super();
    this.state = { posts: null };
  }

  componentWillMount() {
    if (store.getState().user.user === 'anonymous') {
      store.subscribe(() => {
          if (store.getState().user.user !== 'anonymous') {
            this.getPostByUsername();
          }
        });
    } else {
      this.getPostByUsername();
    }
  }

  getPostByUsername() {
    request.get('/api/getpostbyusername/' + store.getState().user.user)
      .then(res => {
        let posts = [];

        Object.keys(res.body.p).map((e) => {
          posts.push({ id: e, title: res.body.p[e] });
        });

        this.setState({ posts });
      })
      .catch(e => e);
  }

  render() {
    return (
      <div>
        <h1>My Posts</h1>
        {this.state.posts &&
          this.state.posts.map(p => {
            const post = (
              <div key={p.id}>
                <h4>{p.title}</h4>
              </div>
            );
            return post;
          })
        }
      </div>
    );
  }
}
