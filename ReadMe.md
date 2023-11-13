# Ring Captcha

_Ring Captcha_ is a library and a service for creating and checking ring captcha.

_Ring Captcha_ is a graphical type of captcha. A user must see an image and count 
rings drawn on the image. After the first guess attempt the captcha is 
destroyed to avoid getting an answer by brute force.

The product consists of an application receiving requests and providing response 
via a simple _JSON RPC 2.0_ interface running over the _HTTP_ protocol.

The service can either return the created captcha image in its _RPC_ response or 
it can save it to disk storage and provide access to saved images via a separate 
_HTTP_ server. The behaviour is configurable.

## List of functions

The service provides following functions and methods.

| Function           | Description               |
|--------------------|---------------------------|
| Ping               | Pings the server          |
| CreateCaptcha      | Creates a captcha         |
| CheckCaptcha       | Checks the captcha answer |
| ShowDiagnosticData | Shows diagnostic data     |

# Messages
Examples of requests and responses are provided below.

If the server uses default settings, then these messages should be sent as 
_HTTP POST_ requests to the following address:

> HTTP POST -> http://localhost:80

### Captcha creation request
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "CreateCaptcha",
  "params": {}
}

```

### Captcha creation response #1
```json
{
    "jsonrpc": "2.0",
    "result": {
        "taskId": "RCS-12345678-1234-1234-1234-123456789012",
        "imageFormat": "PNG",
        "isImageDataReturned": false,
        "timeSpent": 32
    },
    "id": 1
}
```

When the server is configured to save images to disk storage, the 
`isImageDataReturned` field is set to _False_ and the response contains no 
field named `imageDataB64`.

When the server is configured to respond with image data, the 
`isImageDataReturned` field is set to _True_ and the response contains an 
additional field named `imageDataB64` with binary data encoded with a standard `Base64` 
encoding. An example of such a response is shown below.

### Captcha creation response #2
```json
{
    "jsonrpc": "2.0",
    "result": {
        "taskId": "RCS-12345678-1234-1234-1234-123456789012",
        "imageFormat": "PNG",
        "isImageDataReturned": true,
        "imageDataB64": "...",
        "timeSpent": 35
    },
    "id": 1
}
```

To retrieve an image of captcha a user should make an _HTTP_ request using 
_GET_ method. The request should contain a parameter named `id`. 

Please note that _HTTP_ server that shows images is a separate server from the 
_JSON RPC_ server – they use different configuration. The second server is 
optional and can be easily disabled.

> HTTP GET -> http://localhost:81/?id=RCS-12345678-1234-1234-1234-123456789012

To check the correctness of user's answer, the following request is made.

### Captcha check request
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "CheckCaptcha",
  "params": {
    "taskId": "RCS-12345678-1234-1234-1234-123456789012",
    "value": 1
  }
}
```

### Captcha check response
```json
{
  "jsonrpc": "2.0",
  "result": {
    "taskId": "RCS-12345678-1234-1234-1234-123456789012",
    "isSuccess": false,
    "timeSpent": 0
  },
  "id": 2
}
```

Remember that after the first guess attempt the captcha is deleted and there is 
no sense in checking the same ID again.

# Configuration

An example of a configuration file is provided below.

```json
{
  "http": {
    "host": "localhost",
    "port": 80
  },
  "captcha": {
    "storeImages": true,
    "imagesFolder": "some_path\\img",
    "imageWidth": 256,
    "imageHeight": 256,
    "imageTTLSec": 60,
    "clearImagesFolderAtStart": true,
    "useHttpServerForImages": true,
    "httpServerHost": "localhost",
    "httpServerPort": 81,
    "httpServerName": "RCS"
  }
}
```

As can be seen from the file contents, image saving can be disabled as well as 
an _HTTP_ server for image sharing in case that you use other means of providing 
images to users. 

# Resource usage
Lifetime of images can be configured using the `imageTTLSec` parameter. The 
recommended value is one minute, i.e. 60 seconds. All the captcha that are 
older than the setting's value are automatically removed from disk storage 
and memory to free some space. 

When saving to disk storage is enabled, the performance of the service is 
limited by performance of your disk subsystem. Performance can be increased by 
either disabling saves to disk storage or by decreasing dimensions of captcha 
images. Some users may also find network storages and virtual disks useful for 
increasing performance. The service is not designed for super high loads by 
default.

# Liveness handler

Liveness handler is accessible with the same _JSON-RPC 2.0_ interface.

### Liveness request
```json
{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "Ping",
    "params": {}
}
```

### Liveness response
```json
{
    "jsonrpc": "2.0",
    "result": {
        "ok": true
    },
    "id": 3
}
```

## Diagnostics handler

Diagnostics handler is accessible with the same _JSON-RPC 2.0_ interface.

Note that a call to the diagnostics handler itself is also counted as a
request, and while the response of this handler is being prepared the request
is still in progress, i.e. not finished. This leads to a deviation of counters.

Request:
```json
{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "ShowDiagnosticData",
    "params": {}
}
```

Response:
```json
{
    "jsonrpc": "2.0",
    "result": {
        "timeSpent": 0,
        "totalRequestsCount": 100,
        "successfulRequestsCount": 99
    },
    "id": 4
}
```
