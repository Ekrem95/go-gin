import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';
import request from 'superagent';
import { store } from '../redux/reducers';

export default class Nav extends Component {
  constructor() {
    super();
    this.state = { loggedIn: Boolean };
  }

  componentWillMount() {
    console.log('store');
    store.subscribe(() => {
      console.log(store.getState());
      const state = store.getState();
      switch (state) {
        case 1:
          this.setState({ loggedIn: true });
          break;
        default:
          this.setState({ loggedIn: false });
      }
    });
    store.dispatch({ type: 'AUTH' });
  }

  render() {
    return (
      <div className="nav">
        {this.state.loggedIn ?
          <div>
          <NavLink to="/">Home</NavLink>
          <span
            onClick={() => {
              request
                .post('/logout')
                .then(res => {
                  console.log(res.body);
                });
            }}
            >Logout
          </span>
          </div>
          :
          <div>
          <NavLink to="/signup" activeClassName="activeRoute">Sign up</NavLink>
          <NavLink to="/login" activeClassName="activeRoute">Login</NavLink>
          </div>
        }
      </div>
    );
  }
}
