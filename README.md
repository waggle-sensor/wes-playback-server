# Playback Server

This service allows us to mock out the full ffserver + ffmpeg stream with a simple HTTP server and static files.

## Usage

To use this create a data directory of the form:

```txt
data/
  bottom/
    live.mp4
    images/
      image1.jpg
      image2.jpg
      ...
  top/
    live.mp4
    images/
      image1.jpg
      image2.jpg
      ...
```

Now we can run the server with:

```sh
docker run -ti -p 8090:8090 -v /path/to/data:/data:ro waggle/playback-server
```

## Resources

Data will be available at:

```sh
# alias for ffserver bottom mp4 stream
ffplay http://localhost:8090/live

# serves mp4 video stream
ffplay http://localhost:8090/bottom/live.mp4
ffplay http://localhost:8090/top/live.mp4

# serves sequence of images from data directory
wget http://localhost:8090/bottom/image.jpg
wget http://localhost:8090/top/image.jpg
```
