import React, { Component } from 'react';
import { auth } from '../redux/reducers';
import request from 'superagent';

export default class Home extends Component {
  constructor() {
    super();
    this.state = { posts: null };
  }

  componentWillMount() {
    auth()
    .then(res => {
      if (res.auth.auth === 0) {
        this.props.history.push('/login');
      }
    });

    request.get('/api/posts')
      .then(res => {
        if (res.body !== null) {
          const posts = res.body.posts.reverse();
          this.setState({ posts });
        }
      });
  }

  render() {
    return (
      <div className="home">
        {this.state.posts &&
          this.state.posts.map(p => {
            const post = (
              <div className="post" key={p.id}>
              <h3>{p.title}</h3>
              <p>{p.description}</p>
              <img src={p.src}/>
              </div>
            );
            return post;
          })
        }
      </div>
    );
  }
}
