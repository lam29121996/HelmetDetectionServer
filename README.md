## Config
```sh
{
    "port": 8080
}
```

## Request
```http
GET /helmetDetectionResult
``` 
## Expected output:
```sh
{
    "is_helmet_on": true,
    "photo_path": ""
}
```
```sh
{
    "is_helmet_on": false,
    "photo_path": "C:/xxx/xxx.jpg"
}
```