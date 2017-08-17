import React, { Component } from 'react';
import io from 'socket.io-client';
import { store } from '../redux/reducers';
import $ from 'jquery';
import request from 'superagent';

export default class Talkie extends Component {
  constructor() {
    super();
    this.sendMessage = this.sendMessage.bind(this);
    this.state = {
      messages: null,
      val: '',
    };
  }

  componentWillMount() {
    let socket = io.connect('/');

    if (socket !== undefined) {
      socket.on('dist', msg => {
        const messages = this.state.messages || [];
        messages.push(msg);
        this.setState({ messages });
      });
    }

    this.jquery();

    request.get('/messages')
      .then(res => {
        let messages = [];
        res.body.messages.map(m => {
          m = JSON.parse(m);
          messages.push(m);
        });
        this.setState({ messages });
      });
  }

  jquery() {
    $(document).ready(function () {
      $('#hide-chat').on('click', () => {
          $('.talkie').fadeToggle();
          $('#show-chat').fadeToggle();
        });
      $('#show-chat').on('click', () => {
          $('.talkie').fadeToggle();
          $('#show-chat').fadeToggle();
        });
    });
  }

  sendMessage() {
    let socket = io.connect('/');

    if (socket !== undefined) {
      const username = store.getState().user.user;
      socket.emit('msg', {
        text: this.state.val,
        time: Date.now().toString(),
        sender: username,
      });
    }
  }

  render() {
    return (
      <div>
      <div className="talkie">
        <div className="top">
          <div>Messages</div>
          <span id="hide-chat">Hide</span>
        </div>
        <div className="bottom">
          {
            this.state.messages &&
            this.state.messages.map((m, i) => {
              const first = new Date(Number(m.Time)).toString().slice(4, 10);
              const second = new Date(Number(m.Time)).toString().slice(16, 21);
              const message = (
                <div className="message" key={i}>
                  <div className="text">{m.Text}</div>
                  <div className="details">
                    <div>{`${first}, ${second}`}</div>
                    <div>{m.Sender}</div>
                  </div>
                </div>
              );
              return message;
            })
          }
        </div>
        <input
          id="chat"
          onChange={(e) => {
            this.setState({ val: e.target.value });
          }}

          onKeyDown={(e) => {
            if (e.keyCode === 13 &&
              this.state.val.length > 0
            ) {
              this.sendMessage();
              $('#chat').val('');
              this.setState({ val: '' });
            }
          }}

          />
      </div>
      <div id="show-chat">Chat</div>
      </div>
    );
  }
}
