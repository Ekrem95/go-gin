import React, { Component } from 'react';
import { store } from '../redux/reducers';
import request from 'superagent';
import { Link } from 'react-router-dom';
import $ from 'jquery';

export default class MyPosts extends Component {
  constructor() {
    super();
    this.state = { posts: null, pac: null };
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
                      const id = p.id;
                      const user = store.getState().user.user;

                      const pac = { id, user };
                      this.setState({ pac });

                      $('#dlgbox')
                      .css('display', 'flex')
                      .hide()
                      .fadeIn();
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

      <div id="dlgbox">
        <div id="dlg-body">Do you want to delete?</div>
        <div id="dlg-footer">
            <button
              onClick={() => {
                $('#dlgbox').fadeOut();
                request
                  .post('/delete/' + this.props.location.pathname.split('/').pop())
                  .type('form')
                  .send(this.state.pac)
                  .set('Accept', 'application/json')
                  .then(res => {
                    if (res.body.deleted === true) {
                      const posts = this.state.posts.filter(pos => pos.id !== this.state.pac.id);
                      this.setState({ posts });
                    } else {
                      // there was an error
                      return;
                    }
                  })
                  .catch(e => e);
              }}>
              OK
            </button>
            <button
              onClick={() => {
                this.setState({ pac: null });
                $('#dlgbox').fadeOut();
              }}>
              Cancel
            </button>
        </div>
      </div>
      </div>
    );
  }
}
