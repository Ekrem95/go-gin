import React, { Component } from 'react';

export default class Edit extends Component {
    constructor() {
        super();
        this.state = { post: null };
        this.edit = this.edit.bind(this);
    }

    componentWillMount() {
        fetch('/api/postbyid/' + this.props.location.pathname.split('/').pop())
            .then(res => res.json())
            .then(res => this.setState({ post: res.post }))
            .catch(e => e);
    }

    edit() {
        const title = this.refs.title.value;
        const description = this.refs.description.value;
        const src = this.refs.src.value;

        const pac = { title, description, src };

        fetch('/edit/' + this.props.location.pathname.split('/').pop(), {
            method: 'post',
            body: JSON.stringify(pac)
        })
            .then(res => res.json())
            .then(res => {
                if (res.done === true) {
                    this.props.history.push('/myposts');
                }
            })
            .catch(e => console.log(e));
    }

    render() {
        return (
            <div className='add'>
                <h1>Edit</h1>
                {this.state.post && (
                    <form>
                        <input
                            ref='title'
                            type='text'
                            placeholder='Title'
                            defaultValue={this.state.post.title}
                        />
                        <input
                            ref='description'
                            type='text'
                            placeholder='Description'
                            defaultValue={this.state.post.description}
                        />
                        <input
                            ref='src'
                            type='text'
                            placeholder='Image Source'
                            defaultValue={this.state.post.src}
                        />
                        <button
                            onClick={() => {
                                this.edit();
                            }}
                            type='button'
                        >
                            Edit
                        </button>
                    </form>
                )}
            </div>
        );
    }
}
