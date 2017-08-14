import React, { Component } from 'react';
import { render } from 'react-dom';
import { BrowserRouter, Route, Switch } from 'react-router-dom';

// import Home from './Components/Home';
import { Home, Signup, Login, Nav } from './Components';

import style from './style.scss';

class App extends Component {
  render() {
    return (
      <BrowserRouter>
        <div>
        <Nav />
        <Switch>
          <Route exact path="/" component={Home} />
          <Route path="/login" component={Login} />
          <Route path="/signup" component={Signup} />
        </Switch>
        </div>
      </BrowserRouter>
    );
  }
}

render(<App />, document.getElementById('app'));
