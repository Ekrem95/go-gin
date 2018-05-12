import React from 'react';

export const NameInput = component => (
    <input
        autoFocus
        type='text'
        placeholder='Username'
        onChange={e => {
            component.setState({ name: e.target.value });
        }}
        onBlur={e => {
            if (e.target.value.length < 3) {
                component.setState({
                    nameErr: 'Username must be longer than 2 chars'
                });
            } else if (e.target.value.length > 50) {
                component.setState({
                    nameErr: 'Username must be shorter than 50 chars'
                });
            } else {
                component.setState({
                    nameErr: null
                });
            }
        }}
    />
);
