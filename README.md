## Default config
```sh
{
    "port": 8080,
    "timeout(ms)": 5000,
    "images_file_path": "/Users/kelvinlam/Downloads/images"
}
```

## Request & Expected output
```http
GET /helmetDetectionResult
```
```sh
{
    "is_helmet_on": true,
    "photo_path": ""
}
```
```sh
{
    "is_helmet_on": false,
    "photo_path": "images/xxx.jpg"
}
```


```http
GET /images
```