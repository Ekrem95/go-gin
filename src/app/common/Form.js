import React, { Component } from 'react';
import { store } from '../redux/reducers';
import { NameInput, PasswordInput } from './index';

export default class Form extends Component {
  constructor() {
    super();
    this.state = {
      name: '', password: '',
      nameErr: null, passwordErr: null, err: null,
    };
  }

  render() {
    return (
      <div className="form">
        <h2>{this.props.header}</h2>
        <p>{this.state.err}</p>
        {NameInput(this)}
        <p>{this.state.nameErr}</p>
        {PasswordInput(this)}
        <p>{this.state.passwordErr}</p>
        <button
          type="button"
          onClick={() => {
            if (
              this.state.name.length > 2 &&
              this.state.password.length > 5
            ) {
              const payload = {
                username: this.state.name, password: this.state.password,
              };

              var data = new FormData();
              data.append('username', this.state.name);
              data.append('password', this.state.password);

              fetch(this.props.post, {
                method: 'POST',
                credentials: 'same-origin',
                body: data,
              })
              .then(r => r.json())
              .then(res => {
                if (res.user) {
                  store.dispatch({ type: 'AUTH' });
                  store.dispatch({ type: 'USER', payload: res.user });
                  this.props.history.push('/');
                } else if (res.err) {
                  this.setState({
                    err: 'Wrong username & password combination',
                  });
                } else if (res.error) {
                  this.setState({
                    err: res.error,
                  });
                }
              });
            }
          }}
          >{this.props.header}</button>
      </div>
    );
  }
}
