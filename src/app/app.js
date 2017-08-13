import React, { Component } from 'react';
import { render } from 'react-dom';

// import Home from './Components/Home';
import { Signup } from './Components';

import style from './style.scss';

class App extends Component {
  render() {
    return (
      <Signup />
    );
  }
}

render(<App />, document.getElementById('app'));
