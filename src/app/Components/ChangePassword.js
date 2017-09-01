import React, { Component } from 'react';
import request from 'superagent';

export default class ChangePassword extends Component {
  constructor() {
    super();
    this.state = { errors: [] };
  }

  render() {
    return (
      <div className="add">
        <h1>Change Password</h1>
        <form>
          <input
            ref="current"
            type="password"
            placeholder="Current Password"
          />
          <input
            ref="password"
            type="password"
            placeholder="New Password"
          />
          <input
            ref="password2"
            type="password"
            placeholder="Please Repeat New Password"
          />
          <button
            onClick={() => {
              const current = this.refs.current.value;
              const newPassword = this.refs.password.value;
              const assert = this.refs.password2.value;

              if (
                current.length > 5 &&
                newPassword.length > 5 &&
                assert.length > 5
              ) {
                if (newPassword === assert) {
                  if (newPassword !== current) {
                    this.setState({ errors: [] });

                    const http = new XMLHttpRequest();
                    const url = '/changepassword';
                    const params = `current=${current}&newPassword=${newPassword}`;
                    http.open('POST', url, true);

                    http.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');

                    http.onreadystatechange = function () {
                        if (http.readyState == 4 && http.status == 200) {
                          console.log(JSON.parse(http.responseText));
                        }
                      };

                    http.send(params);
                  } else {
                    const errors = [];
                    const err = 'New and old passwords should be different.';
                    errors.push(err);
                    this.setState({ errors });
                  }
                } else {
                  const errors = [];
                  const err = 'Passwords do not match';
                  errors.push(err);
                  this.setState({ errors });
                }
              } else {
                const errors = this.state.errors;
                const err = 'Fields must have at least 6 chars';
                if (errors.indexOf(err) < 0) {
                  errors.push(err);
                  this.setState({ errors });
                }
              }
            }}

            type="button">
            Confirm
          </button>
        </form>
        {this.state.errors.length > 0 &&
          this.state.errors.map((e, i) => <p key={i}>{e}</p>)
        }
      </div>
    );
  }
}
