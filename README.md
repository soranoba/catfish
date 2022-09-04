🎣 catfish
-----
Useful dummy server used for development.

# docker

```bash
docker build . -t catfish
docker run -p 8080:8080 -v ${YOUR_CONFIG}:/etc/catfish/config.yml catfish
```

## Specification
### Config

[A sample config is here.](/bin/config.yml)

#### Top level

|Field|Type|Required|Description|
|:---|:---|:---|:---|
|routes|`Array<Route>`|x|When Catfish receives a request, it checks for matching Routes in order from the top, and the first matching Route is used.|

#### Route

|Field|Type|Required|Example|Description|
|:---|:---|:---|:---|:---|
|method|`String`|o|`GET`|HTTP Request method.<br>Upper and lower cases are ignored.|
|path|`String`|o|`/users/:id`|HTTP path. It can include path parameters.|
|response|`Dictionary<String,Response>`|o| |The key is used as the response preset name.<br>When Catfish receives a request, it decides to whether to use the preset in order from the top. |

#### Response

|Field|Type|Required|Example|Description|
|:---|:---|:---|:---|:---|
|cond|`Float`|x|`0.8`|The probability that this preset will be used.(`[0.0, 1.0]`)|
|delay|`Float`|x|`0.1`|Delay time before response is returned. (sec)|
|status|`Integer`|o|`200`|HTTP Status code|
|header|`Dictionary<String,String>`|x| |HTTP response headers|
|body|`String`|x|`{"message":"OK"}`|HTTP response body|

### Path parameters

You can use two kinds of path parameters.

- `:` prefix: Always match one segment.
- `*` prefix: Match any segments. (includes 0)


|Path|Examples|
|:---|:---|
|`/users/:id`|✅`/users/1`<br>✅`/users/1/`<br>❌`/users/`<br>❌`/users/1/follow`|
|`/users/*path`|✅`/users`<br>✅`/users/1/follow`|


### Response headers

Catfish automatically add some headers in responses to easily debug with.

|Header name|Required|Description|
|:---|:---|:---|
|X-CATFISH-PATH|o|Indicates that Catfish returned the response according to the path setting.|
|X-CATFISH-RESPONSE-PRESET-NAME|o|Indicates that Catfish returned the response according to the preset setting.|
|X-CATFISH-ERROR|x|The descriptions of an error.|
