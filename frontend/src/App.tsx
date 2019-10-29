import React, { FC, useState, useCallback } from 'react';
import axios from 'axios';
import { useDropzone } from 'react-dropzone';
import './App.css';

const App: React.FC = () => {
  const [image, setImage] = useState('');
  const maxSize = 1048576; // 1MB

  const onDrop = useCallback((acceptedFiles: any) => {
    setImage(acceptedFiles[0]);
  }, []);

  const {
    isDragActive,
    getRootProps,
    getInputProps,
    isDragReject,
    acceptedFiles,
    rejectedFiles,
  } = useDropzone({
    onDrop,
    accept: 'image/png',
    minSize: 0,
    maxSize,
  });

  const isFileTooLarge =
    rejectedFiles.length > 0 && rejectedFiles[0].size > maxSize;

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>): void => {
    e.preventDefault();

    const submitData = new FormData();
    submitData.append('file', image);

    // TODO: api host
    axios
      .post('http://localhost:9292/', submitData)
      .then(response => {
        alert(response.data.image_url);
      })
      .catch(error => console.log(error));
  };

  return (
    <>
      <div className="App">
        <h2>画像アップロード</h2>
        <form onSubmit={handleSubmit}>
          {acceptedFiles.length < 1 && (
            <div {...getRootProps()} className="drag-area">
              <input {...getInputProps()} />
              {!isDragActive &&
                <div className="drag-text">
                  ここをクリックするか、画像をドラッグしてください
                </div>
              }
              {isDragActive && !isDragReject && 'ドロップしてね'}
              {isDragReject &&
                'ごめんなさい。画像は png のみアップロード可能です'}
              {isFileTooLarge && (
                <div className="text-danger mt-2">
                  ファイルサイズが大きすぎます
                </div>
              )}
            </div>
          )}
          <ul className="preview">
            {acceptedFiles.length > 0 &&
              acceptedFiles.map(acceptedFile => (
                <li key={acceptedFile.name}>
                  <img
                    src={URL.createObjectURL(acceptedFile)}
                    alt=""
                    height="100"
                  />
                </li>
              ))}
          </ul>
          <button
            type="submit"
          >
            アップロード
          </button>
        </form>
      </div>
    </>
  );
};

export default App;
