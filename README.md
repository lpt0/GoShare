# GoShare

ShareX-compatible server written in Go. 

URL shortening will additionally grab a WARC of the provided URL (requires a copy of `wget` in the same directory or path).

Originally written while I was learning the language, for practice.

## Installation and usage

### Server

1. Build the server using `go build -o goshare`.
    - If you are on Windows and do not wish to build it, you can find a copy in [releases](https://github.com/lpt0/GoShare/releases).
2. Copy `config.json.example` to `config.json`.
3. Modify the parameters as required, ensuring that the directory specified in `FilePath` exists.
4. Run `goshare`.

### ShareX

1. Modify `GoShare.sxcu`, replacing `localhost:8080` with the host and port specified in `config.json`.
    - Additionally, modify the `Authorization` header if required.
2. Open ShareX.
3. Click `Destinations > Custom uploader settings > Import > From file...`
4. Select the modified `GoShare.sxcu` file.
5. (optional) Change the default upload destination by clicking `Destinations > Image Uploader > Custom image uploader` (or any other type of uploader).
