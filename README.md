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
<a href="test_1.jpg">test_1.jpg</a>
<a href="test_2.jpg">test_2.jpg</a>
<a href="test_3.jpg">test_3.jpg</a>
...
<a href="test_n.jpg">test_n.jpg</a>
</pre>

```http
GET /images/test_1.jpg
```
![My image](/images/test_1.jpg)