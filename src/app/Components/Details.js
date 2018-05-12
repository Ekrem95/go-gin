import React, { Component } from 'react';
import { store } from '../redux/reducers';
export default class Details extends Component {
    constructor() {
        super();
        this.state = { data: null, comments: null, likes: null };
    }

    componentWillMount() {
        fetch('/api/postbyid/' + this.props.location.pathname.split('/').pop())
            .then(r => r.json())
            .then(r => this.setState({ data: r.post }))
            .catch(e => console.log(e))
            .then(() => {
                fetch(
                    '/api/commentsbyid/' +
                        this.props.location.pathname.split('/').pop()
                )
                    .then(r => r.json())
                    .then(r => this.setState({ comments: r.comments }))
                    .catch(e => console.log(e));
            })
            .then(() => {
                fetch(
                    '/get_likes/' +
                        this.props.location.pathname.split('/').pop()
                )
                    .then(r => r.json())
                    .then(r => {
                        if (r.users !== null) {
                            if (
                                r.users.indexOf(store.getState().user.user) > -1
                            ) {
                                this.setState({ liked: true });
                            }

                            this.setState({ likes: r.users.length });
                        } else {
                            this.setState({ likes: 0 });
                        }
                    });
            });
    }

    render() {
        return (
            <div className='details'>
                {this.state.data && (
                    <div>
                        <h1>{this.state.data.title}</h1>
                        <img src={this.state.data.src} />
                        <p>{this.state.data.description}</p>
                        <textarea
                            ref='comment'
                            placeholder='Type here to post a comment'
                            onKeyUp={e => {
                                if (e.keyCode === 13) {
                                    const text = this.refs.comment.value;
                                    const post_id = this.state.data.id.toString();
                                    const sender = store.getState().user.user;

                                    const pac = { text, post_id, sender };

                                    fetch('/comment', {
                                        method: 'post',
                                        body: JSON.stringify(pac)
                                    })
                                        .then(res => res.json())
                                        .then(res => console.log(res));

                                    this.refs.comment.value = '';

                                    if (this.state.comments) {
                                        const comments = this.state.comments;
                                        const comment = Object.assign(pac, {
                                            time: Date.now() / 1000
                                        });
                                        comments.push(comment);
                                        this.setState({ comments });
                                    } else {
                                        const comments = [];
                                        const comment = Object.assign(pac, {
                                            time: Date.now() / 1000
                                        });
                                        comments.push(comment);
                                        this.setState({ comments });
                                    }
                                }
                            }}
                        />
                    </div>
                )}
                <div className='likes'>
                    <button
                        onClick={() => {
                            const post_id = this.props.location.pathname
                                .split('/')
                                .pop();
                            const user = store.getState().user.user;
                            const pac = { post_id, user };

                            fetch('/post_likes', {
                                method: 'post',
                                body: JSON.stringify(pac)
                            })
                                .then(res => null)
                                .catch(e => console.log(e));

                            this.setState({ liked: !this.state.liked });

                            this.state.liked
                                ? this.setState({ likes: this.state.likes - 1 })
                                : this.setState({
                                    likes: this.state.likes + 1
                                });
                        }}
                    >
                        {this.state.liked ? 'Liked' : 'Like'}
                    </button>
                    <span>{this.state.likes} likes</span>
                </div>
                {this.state.comments &&
                    this.state.comments.map((c, i) => {
                        const date = new Date(c.time * 1000).toDateString();
                        const comment = (
                            <div key={i} className='comment'>
                                <span>{c.text}</span>
                                <span>
                                    <span>{c.sender}</span>
                                    <span>{date}</span>
                                </span>
                            </div>
                        );
                        return comment;
                    })}
            </div>
        );
    }
}
