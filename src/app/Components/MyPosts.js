import React, { Component } from 'react';
import { store } from '../redux/reducers';
import request from 'superagent';
import { Link } from 'react-router-dom';

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
      <div className="myposts">
        <h1>My Posts</h1>
        {this.state.posts &&
          this.state.posts.map(p => {
            const post = (
              <div key={p.id} className="post">
                <Link
                  className="title"
                  to={`/p/${p.id}`}>
                  {p.title}
                </Link>
                <div className="buttons">
                  <button
                    onClick={() => {
                      this.props.history.push('/edit/' + p.id);
                    }}

                    type="button">
                      Edit
                  </button>
                  <button
                    onClick={() => {
                      console.log(p.id);
                    }}

                    type="button">
                      Delete
                  </button>
                </div>
              </div>
            );
            return post;
          })
        }
      </div>
    );
  }
}
