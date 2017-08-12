import React, { Component } from 'react';
import { render } from 'react-dom';

// import Home from './Components/Home';
import { Home } from './Components';

class App extends Component {
  render() {
    return (
      <Home />
    );
  }
}

render(<App />, document.getElementById('app'));
