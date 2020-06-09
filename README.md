This service allows us to mock out the full ffserver + ffmpeg stream with a simple HTTP server and static files.

To use this create a data directory of the form:

```
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
docker run -ti -p 8090:8090 -v /path/to/data:/data:ro waggle/media-server
```

Data will be available at:

```sh
ffplay http://localhost:8090/bottom
ffplay http://localhost:8090/top
ffplay http://localhost:8090/live # alias for bottom
```
