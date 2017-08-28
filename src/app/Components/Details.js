import React, { Component } from 'react';
import request from 'superagent';
import { store } from '../redux/reducers';
export default class Details extends Component {
  constructor() {
    super();
    this.state = { data: null, comments: null };
  }

  componentWillMount() {
    request.get('/api/postbyid/' + this.props.location.pathname.split('/').pop())
      .then(res => {
        this.setState({ data: res.body.post });
      })
      .catch(err => {
        console.log(err);
      })
      .then(() => {
        request.get('/api/commentsbyid/' + this.props.location.pathname.split('/').pop())
          .then(res => {
            this.setState({ comments: res.body.comments });
          })
          .catch(err => {
            console.log(err);
          });
      });
  }

  render() {
    return (
      <div className="details">
        {this.state.data &&
          <div>
            <h1>{this.state.data.title}</h1>
            <img src={this.state.data.src}/>
            <p>{this.state.data.description}</p>
            <textarea
              ref="comment"
              placeholder="Type here to post a comment"
              onKeyUp={(e) => {
                if (e.keyCode === 13) {
                  const text = this.refs.comment.value;
                  const postId = this.state.data.id;
                  const sender = store.getState().user.user;

                  const pac = { text, postId, sender };

                  request
                    .post('/comment')
                    .type('form')
                    .send(pac)
                    .set('Accept', 'application/json')
                    .then(res => {
                      console.log(res.body);
                    })
                    .catch(err => {
                      console.log(err);
                    });

                  this.refs.comment.value = '';

                  if (this.state.comments) {
                    const comments = this.state.comments;
                    const comment = Object.assign(
                      pac, { time: Date.now() / 1000 }
                    );
                    comments.push(comment);
                    this.setState({ comments });
                  } else {
                    const comments = [];
                    const comment = Object.assign(
                      pac, { time: Date.now() / 1000 }
                    );
                    comments.push(comment);
                    this.setState({ comments });
                  }

                }
              }}
              ></textarea>
          </div>
        }
        <div className="likes">
          <button
            onClick={() => {
              const postID = this.props.location.pathname.split('/').pop();
              const user = store.getState().user.user;
              const pac = { postID, user };

              request
                .post('/post_likes')
                .type('form')
                .send(pac)
                .set('Accept', 'application/json')
                .then(res => {
                  console.log(res.body);
                })
                .catch(err => {
                  console.log(err);
                });
            }}
            >{this.state.liked ? 'Liked' : 'Like'}
          </button>
        </div>
        {this.state.comments &&
          this.state.comments.map((c, i) => {
            const date = new Date(c.time * 1000).toDateString();
            const comment = (
              <div key={i} className="comment">
                <span>{c.text}</span>
                <span>
                  <span>{c.sender}</span>
                  <span>{date}</span>
                </span>
              </div>
            );
            return comment;
          })
        }
      </div>
    );
  }
}
