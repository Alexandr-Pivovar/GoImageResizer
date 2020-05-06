# Go Image Resizer

### Resize image:

>Method: Post
>
>Path:**/api/v1/resize**
>
>Body:

```json5
{
  "data": "iVBORw0KGgoAAAANSUhEUgAAAZAAAAGQCAIAAAAP3a", // base64 format
  "width": 800,
  "height": 600,
}
```
> Response body:

```json5
 {
    "origin_url": "http://example.com/origin.png",
    "resized_url": "http://example.com/resized.png",
    "width": 800,
    "height": 600
 }
```

### Get all history resized images:

>Method: Get
>
>Path:**/api/v1/history**
>
> Response body:

```json5
{
  "076eb660-aa44-4216-a4b8-6353dc1623eb": {
    "origin_url": "http://example.com/origin.png",
    "resized_url": "http://example.com/resized.png",
    "width": 1000,
    "height": 1000
  },
  "12367d4e-33bc-4148-9c70-d538e15ccb10": {
    "origin_url": "http://example.com/origin.png",
    "resized_url": "http://example.com/resized.png",
    "width": 1600,
    "height": 1600
  }
}
```

### Get resized image by id:

>Method: Get
>
>Path:**/api/v1/history/{id}**
>
> Response body:

```json5
 {
    "origin_url": "http://example.com/origin.png",
    "resized_url": "http://example.com/resized.png",
    "width": 800,
    "height": 600
 }
```

### Update image by id:

>Method: Post
>
>Path:**/api/v1/update**
>
> body:

```json5
 {
    "id": "076eb660-aa44-4216-a4b8-6353dc1623eb",
    "width": 800,
    "height": 600
 }
```

> Response body:

```json5
 {
    "origin_url": "http://example.com/origin.png",
    "resized_url": "http://example.com/resized.png",
    "width": 800,
    "height": 600
 }
```

## Run tests

./test.sh
