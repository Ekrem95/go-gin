import React, { Component } from 'react';
import request from 'superagent';
import { store } from '../redux/reducers';
export default class Details extends Component {
  constructor() {
    super();
    this.state = { data: null };
  }

  componentWillMount() {
    request.get('/api/postbyid/' + this.props.location.pathname.split('/').pop())
      .then(res => {
        this.setState({ data: res.body.post });
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
                  const from = store.getState().user.user;

                  const pac = { text, postId, from };

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
                }
              }}
              ></textarea>
          </div>
        }
      </div>
    );
  }
}
