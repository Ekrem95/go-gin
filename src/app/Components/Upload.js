import React, { Component } from 'react';
import Dropzone from 'react-dropzone';
import request from 'superagent';

export default class Upload extends Component {
  constructor() {
    super();
    this.state = { files: [], msg: null };
    this.onDrop = this.onDrop.bind(this);
    this.onOpenClick = this.onOpenClick.bind(this);
  }

  onDrop(files) {
    var photo = new FormData();
    photo.append('photo', files[0]);

    request.post('/upload')
      .send(photo)
      .end((err, resp) => {
        if (err) {
          console.error(err);
          this.setState({ msg: 'Error Occured While Uploading Image.' });
        } else {
          this.setState({ msg: 'Image Uploaded Succesfully.' });
        }
      });
  }

  onOpenClick(files) {
    this.refs.dropzone.open();
  }

  render() {
    return (
      <div>
        <Dropzone ref="dropzone" multiple={false} accept={'image/*'} onDrop={this.onDrop}>
           <div>Try dropping some files here, or click to select files to upload.</div>
         </Dropzone>
         {this.state.msg &&
           <p>{this.state.msg}</p>
         }
         {/* <button type="button" onClick={this.onOpenClick}>
             Open Dropzone
         </button> */}
         {/* {this.state.files ? <div>
         <h2>Uploading {this.state.files.length} files...</h2>
         <div>{this.state.files.map(file => <img src={file.preview} />)}</div>
         </div> : null} */}
      </div>
    );
  }
}
