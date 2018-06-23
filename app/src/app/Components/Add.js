import React, { Component } from 'react';
import { store } from '../redux/reducers';

export default class Add extends Component {
    constructor() {
        super();
        this.add = this.add.bind(this);
    }

    add() {
        const body = new FormData();
        body.append('title', this.refs.title.value);
        body.append('description', this.refs.desc.value);
        body.append('src', this.refs.src.value);
        body.append('posted_by', store.getState().user.user);

        fetch('/add', { method: 'post', body })
            .then(res => res.json())
            .then(res => {
                if (res.id) this.props.history.push('/');
            });
    }

    render() {
        return (
            <div className='add'>
                <h1>Add</h1>
                <form>
                    <input ref='title' type='text' placeholder='Title' />
                    <input ref='desc' type='text' placeholder='Description' />
                    <input ref='src' type='text' placeholder='Image Source' />
                    <button
                        onClick={() => this.add()}
                        type='button'
                    >
                        Add
                    </button>
                </form>
            </div>
        );
    }
}
