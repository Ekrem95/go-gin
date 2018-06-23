import React, { Component } from 'react';
import { store } from '../redux/reducers';
import { Link } from 'react-router-dom';
import $ from 'jquery';

export default class MyPosts extends Component {
    constructor() {
        super();
        this.state = { posts: null, pac: null };
    }

    componentWillMount() {
        if (store.getState().user.user === 'anonymous') {
            store.subscribe(() => {
                if (store.getState().user.user !== 'anonymous') {
                    this.getPostByUsername();
                }
            });
        } else {
            this.getPostByUsername();
        }
    }

    getPostByUsername() {
        fetch('/api/getpostbyusername/' + store.getState().user.user)
            .then(res => res.json())
            .then(res => {
                let posts = [];

                Object.keys(res.posts).map(i => posts.push({ id: i, title: res.posts[i] }));

                this.setState({ posts });
            })
            .catch(e => e);
    }

    render() {
        return (
            <div className='myposts'>
                <h1>My Posts</h1>
                {this.state.posts &&
                    this.state.posts.map(p => {
                        const post = (
                            <div key={p.id} className='post'>
                                <Link className='title' to={`/p/${p.id}`}>
                                    {p.title}
                                </Link>
                                <div className='buttons'>
                                    <button
                                        onClick={() => {
                                            this.props.history.push('/edit/' + p.id);
                                        }}
                                        type='button'
                                    >
                                        Edit
                                    </button>
                                    <button
                                        onClick={() => {
                                            const id = p.id;
                                            const user = store.getState().user.user;

                                            const pac = { id, user };
                                            this.setState({ pac });

                                            $('#dlgbox')
                                                .css('display', 'flex')
                                                .hide()
                                                .fadeIn();
                                        }}
                                        type='button'
                                    >
                                        Delete
                                    </button>
                                </div>
                            </div>
                        );
                        return post;
                    })}

                <div id='dlgbox'>
                    <div id='dlg-body'>Do you want to delete?</div>
                    <div id='dlg-footer'>
                        <button
                            onClick={() => {
                                $('#dlgbox').fadeOut();

                                var data = new FormData();
                                data.append('user', this.state.pac.user);
                                data.append('id', this.state.pac.id);

                                fetch(
                                    '/delete/' +
                                    this.props.location.pathname
                                        .split('/')
                                        .pop(),
                                    {
                                        method: 'post',
                                        credentials: 'same-origin',
                                        body: data
                                    }
                                )
                                    .then(res => res.json())
                                    .then(res => {
                                        if (res.deleted === true) {
                                            const posts = this.state.posts.filter(
                                                pos => pos.id !== this.state.pac.id
                                            );
                                            this.setState({ posts });
                                        } else {
                                            // there was an error
                                            return;
                                        }
                                    })
                                    .catch(e => e);
                            }}
                        >
                            OK
                        </button>
                        <button
                            onClick={() => {
                                this.setState({ pac: null });
                                $('#dlgbox').fadeOut();
                            }}
                        >
                            Cancel
                        </button>
                    </div>
                </div>
            </div>
        );
    }
}
