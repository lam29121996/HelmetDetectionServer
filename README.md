## Build
```
$ GOOS=windows GOARCH=amd64 go build -o helmet_detection_server.exe main.go
```

## Default config
```sh
{
    "port": 8080,
    "timeout(ms)": 3500,
    "images_hiu_ming_folder_path": "./hiuMingImages",
    "images_hiu_kwong_folder_path": "./hiuKwongImages",
    "record_from": "07:30",
    "record_to": "18:30"
}
```

## Request & Expected output
```http
GET /helmetDetectionResult?id=1
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
