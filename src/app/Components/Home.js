import React, { Component } from 'react';
import { auth } from '../redux/reducers';
import request from 'superagent';
import $ from 'jquery';

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

    this.jquery();
  }

  jquery() {
    $(document).ready((e) => {
      let handler = (ev) => {
        var $target = $(ev.target);
        if ($target.parent().prop('className') === 'post') {
          const id = $target.parent().attr('id');
          const p = this.state.posts.filter(post => {
              const item = post.id == id;
              return item;
            });
          $('#popup').html(
            '<span>' + p[0].title + '</span>' +
            '<p>' + p[0].description + '</p>'
          );

          $(document).on('mousemove', function (event) {
            if (event.pageX + 290 > window.innerWidth) {
              $('#popup').css({
                top: event.pageY - 100, left: event.pageX - 300, position: 'absolute',
              });
            } else {
              $('#popup').css({
                top: event.pageY - 100, left: event.pageX + 30, position: 'absolute',
              });
            }
          });

          $('#popup').fadeIn();
        }
      };

      $('.post').mouseover(handler)
      .mouseout(() => {
        $('#popup').hide();
      });
    });
  }

  render() {
    return (
      <div className="home">
        {this.state.posts &&
          this.state.posts.map(p => {
            const post = (
              <div className="post" key={p.id} id={p.id}>
              <img
                onClick={() => {
                  this.props.history.push(`/p/${p.id}`);
                }}

                src={p.src}/>
              </div>
            );
            return post;
          })
        }
        <div id="popup"></div>
      </div>
    );
  }
}
