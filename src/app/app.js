import React, { Component } from 'react';
import { render } from 'react-dom';
import { Provider } from 'react-redux';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import { store } from './redux/reducers';

// import Home from './Components/Home';
import {
   Home, Signup, Login, Nav, Talkie, Add, Details, Upload, MyPosts
 } from './Components';

import style from './style.scss';

class App extends Component {
  render() {
    return (
      <Provider store={store}>
        <BrowserRouter>
          <div>
          <Nav />
          <Switch>
            <Route exact path="/" component={Home} />
            <Route path="/login" component={Login} />
            <Route path="/signup" component={Signup} />
            <Route path="/add" component={Add} />
            <Route path="/p/:id" component={Details} />
            <Route path="/upload" component={Upload} />
            <Route path="/myposts" component={MyPosts} />
          </Switch>
          <Talkie />
          </div>
        </BrowserRouter>
      </Provider>
    );
  }
}

render(<App />, document.getElementById('app'));
