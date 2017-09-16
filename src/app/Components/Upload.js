import React, { Component } from 'react';

export default class Upload extends Component {
  constructor() {
    super();
    this.onDrop = this.onDrop.bind(this);
  }

  onDrop(files) {
    const formData = new FormData();
    formData.append('photo', files[0], files[0].name);

    const xhr = new XMLHttpRequest();
    xhr.open('POST', '/upload', true);

    xhr.upload.addEventListener('progress', function (evt) {
          if (evt.lengthComputable) {
            const progress = document.getElementById('progress');
            progress.setAttribute('style', `width: ${evt.loaded / evt.total * 100 * 2}px`);
          }
        }, false);

    xhr.onloadstart = function (e) {
      const progress = document.getElementById('progress');
      progress.setAttribute('style', `width: 0px`);
    };

    xhr.onloadend = function (e) {
          // console.log('end');
        };

    xhr.onload = function () {
      if (xhr.status !== 200) {
        console.log('An error occurred!');
      }
    };

    xhr.send(formData);
  }

  render() {
    return (
      <div className="upload">
        <h1>Upload</h1>
        <p
          onClick={() => {
            document.getElementById('file').click();
          }}

          onDragLeave={(e) => {
            e.preventDefault();
            e.target.style.background = '#000';
          }}

          onDragOver={(e) => {
            e.preventDefault();
          }}

          onDragEnter={(e) => {
            e.preventDefault();
            e.target.style.background = '#333';
          }}

          onDrop={(e) => {
            e.preventDefault();
            e.target.style.background = '#000';
            const files = e.dataTransfer.files;
            this.onDrop(files);
          }}

          id="dropzone">
          Drag & Drop here to upload
        </p>
        <div id="bar">
          <div id="progress"></div>
        </div>
        <input
          type="file"
          id="file"
          onChange={() => {
            var input = document.getElementById('file');
            var curFiles = input.files;
            this.onDrop(curFiles);
          }}

         />
      </div>
    );
  }
}
