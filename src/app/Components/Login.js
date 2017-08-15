import React, { Component } from 'react';
import Form from '../common/Form';

export default class Login extends Component {
  render() {
    return (
      <Form
        header="Login"
        history={this.props.history}
        post={'/login'}
      />
    );
  }
}

// import React, { Component } from 'react';
// import request from 'superagent';
// import { store } from '../redux/reducers';
// import { NameInput, PasswordInput } from '../common';
//
// export default class Form extends Component {
//   constructor() {
//     super();
//     this.state = {
//       name: '', password: '',
//       nameErr: null, passwordErr: null, err: null,
//     };
//   }
//
//   render() {
//     return (
//       <div className="form">
//         <h2>Login</h2>
//         <p>{this.state.err}</p>
//         {NameInput(this)}
//         <p>{this.state.nameErr}</p>
//         {PasswordInput(this)}
//         <p>{this.state.passwordErr}</p>
//         <button
//           type="button"
//           onClick={() => {
//             if (
//               this.state.name.length > 2 &&
//               this.state.password.length > 5
//             ) {
//               const payload = {
//                 username: this.state.name, password: this.state.password,
//               };
//
//               request
//                 .post('/login')
//                 .type('form')
//                 .send(payload)
//                 .set('Accept', 'application/json')
//                 .then(res => {
//                   if (res.body.user) {
//                     store.dispatch({ type: 'AUTH' });
//                     this.props.history.push('/');
//                   } else if (res.body.err) {
//                     this.setState({
//                       err: 'Wrong username & password combination',
//                     });
//                   }
//                 });
//             }
//           }}
//           >Login</button>
//       </div>
//     );
//   }
// }
