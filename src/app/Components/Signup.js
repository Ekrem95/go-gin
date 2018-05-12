import React, { Component } from 'react';
import Form from '../common/Form';

export default class Signup extends Component {
    render() {
        return (
            <Form
                header='Sign up'
                history={this.props.history}
                post={'/signup'}
            />
        );
    }
}
