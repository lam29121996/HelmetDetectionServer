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
<a href="test1.jpg">test1.jpg</a>
<a href="test2.jpg">test2.jpg</a>
<a href="test3.jpg">test3.jpg</a>
...
<a href="testn.jpg">testn.jpg</a>
</pre>


![My image](/images/test1.jpg)