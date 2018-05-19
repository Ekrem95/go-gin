import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';
import { store } from '../redux/reducers';

export default class Nav extends Component {
    constructor() {
        super();
        this.state = { loggedIn: Boolean };
    }

    componentWillMount() {
        fetch('/user', {
            credentials: 'same-origin'
        })
            .then(r => r.json())
            .then(res => {
                if (res.user === null) {
                    store.dispatch({ type: 'UNAUTH' });
                } else {
                    store.dispatch({ type: 'AUTH' });
                    store.dispatch({ type: 'USER', payload: res.user });
                }
            });
        store.subscribe(() => {
            const state = store.getState();
            switch (state.auth.auth) {
                case 1:
                    this.setState({ loggedIn: true });
                    break;
                default:
                    this.setState({ loggedIn: false });
            }
        });
    }

    render() {
        return (
            <div className='nav'>
                {this.state.loggedIn ? (
                    <div>
                        <NavLink to='/'>Home</NavLink>
                        <NavLink to='/add'>Add</NavLink>
                        <NavLink to='/upload'>Upload</NavLink>
                        <NavLink to='/myposts'>My Posts</NavLink>
                        <span
                            onClick={() => {
                                fetch('/logout', {
                                    method: 'POST',
                                    credentials: 'same-origin'
                                }).then(() => {
                                    store.dispatch({ type: 'UNAUTH' });
                                    window.location.href = '/login';
                                });
                            }}
                        >
                            Logout
                        </span>
                    </div>
                ) : (
                    <div>
                        <NavLink to='/signup' activeClassName='activeRoute'>
                            Sign up
                        </NavLink>
                        <NavLink to='/login' activeClassName='activeRoute'>
                            Login
                        </NavLink>
                    </div>
                )}
            </div>
        );
    }
}
