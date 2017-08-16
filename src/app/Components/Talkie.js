import React, { Component } from 'react';
import io from 'socket.io-client';

export default class Talkie extends Component {
  constructor() {
    super();
    this.sendMessage = this.sendMessage.bind(this);
    this.state = {
      messages: [
        'hello',
        'hi',
        'there',
        'hello',
        'hi',
        'there',
        'hello',
        'hi',
        'there',
        'hello',
        'hi',
        'there',
        'hello',
        'hi',
        'there',
        'hello',
        'hi',
        'there',
      ],
    };
  }

  componentWillMount() {
    let socket = io.connect('/');

    if (socket !== undefined) {
      socket.on('dist', msg => {
        console.log(msg);
      });
    }
  }

  sendMessage() {
    let socket = io.connect('/');

    if (socket !== undefined) {
      socket.emit('msg', this.state.val);
    }
  }

  render() {
    return (
      <div className="talkie">
        <div className="top">d</div>
        <div className="bottom">
          {this.state.messages.map((m, i) =>
            <div key={i}>{m}</div>
          )}
        </div>
        <input
          onChange={(e) => {
            this.setState({ val: e.target.value });
          }}

          onKeyDown={(e) => {
            if (e.keyCode === 13) {
              this.sendMessage();
            }
          }}

          />
      </div>
    );
  }
}
