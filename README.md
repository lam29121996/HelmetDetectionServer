## Default config
```sh
{
    "port": 8080,
    "timeout(ms)": 5000,
    "images_file_path": "./images"
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
GET /images/
```
<pre>
<a href="five000000.png">five000000.png</a>
</pre>


![My image](/images/test1.jpg)