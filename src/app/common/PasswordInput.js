import React from 'react';

export const PasswordInput = (component) =>
  <input
    type="password"
    placeholder="Password"
    onChange={(e) => {
      component.setState({ password: e.target.value });
    }}

    onBlur={(e) => {
      if (e.target.value.length < 6) {
        component.setState({
          passwordErr: 'Password must be longer than 5 chars',
        });
      } else {
        component.setState({
          passwordErr: null,
        });
      }
    }}

  />;
