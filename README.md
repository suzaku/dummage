# dummage

Dummage is a tiny service that generates dummy images. 

Once dummage is started, we can ask it for images of sizes and background colors.

It can be used in development environments where you need access to image URLs. For example, you may need to insert some fake products into your test DB, and set the product images to `dummage` URLs:

```python
product.image = "http://localhost:8000/500x500.jpg"
product.thumbnail = "http://localhost:8000/100x100.jpg"
```

## Installation

### Install with Go

If you have Go installed, you can simply get dummage with the following command:

```bash
go get github.com/suzaku/dummage
```

### Install with Docker

If you prefer to run dummage with docker, you can do so with the `satoru/dummage` image:

```bash
docker run -p 8000:8000 --name dummage --rm satoru/dummage
```

## Usage

Once installed, you can use the `dummage` command to start a server.

```bash
$ dummage
2016/03/09 09:18:05 Starting server on port 8000
```

If you want to use a port other than 8000, you can specify it with the `-port` option:

```bash
$ dummage -port 9090
2016/03/09 09:19:32 Starting server on port 9090
```

Now `dummage` is ready to generate images for you, try the following URLS:

* [http://localhost:8000/200x200.jpg](http://localhost:8000/200x200.jpg)
* [http://localhost:8000/300x500-d5d5d5.jpg](http://localhost:8000/300x500-d5d5d5.jpg)
* [http://localhost:8000/300x200.png](http://localhost:8000/300x200.png)

The dimension and background color of the requested image is parsed from the resource name,

`{width}x{height}-{background}.{jpg or png}`

If the `-{background}` part is left out, a random background color will be used.
