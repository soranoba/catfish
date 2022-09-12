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

| Field  | Type           | Required | Description                                                                                                                 |
|:-------|:---------------|:---------|:----------------------------------------------------------------------------------------------------------------------------|
| routes | `Array<Route>` | x        | When Catfish receives a request, it checks for matching Routes in order from the top, and the first matching Route is used. |

#### Route

| Field    | Type                          | Required | Example      | Description                                                                                                                                                                                 |
|:---------|:------------------------------|:---------|:-------------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| method   | `String`                      | o        | `GET`        | HTTP method.<br>Allowed values are `GET`, `POST`, `PUT`, `DELETE` or `*`.<br>`*` means any HTTP method.                                                                                     |
| path     | `String`                      | o        | `/users/:id` | HTTP path. You can include path parameters.                                                                                                                                                 |
| parser   | `String`                      | x        | `json`       | HTTP request body parser.<br>Allowed value is `json`.                                                                                                                                       |
| response | `Dictionary<String,Response>` | o        |              | The key is used as the response preset name. (See also `X-CATFISH-RESPONSE-PRESET-NAME`)<br>When Catfish receives a request, it decides to whether to use the preset in order from the top. |

#### Response

| Field  | Type                        | Required | Example                 | Description                                                                     |
|:-------|:----------------------------|:---------|:------------------------|:--------------------------------------------------------------------------------|
| cond   | `String`                    | x        | `totalRequestCount < 2` | A conditional expression indicating the probability of returning this response. |
| delay  | `Float`                     | x        | `0.1`                   | Delay time before response is returned. (sec)                                   |
| status | `Integer`                   | o        | `200`                   | HTTP Status code                                                                |
| header | `Dictionary<String,String>` | x        |                         | HTTP response headers                                                           |
| body   | `String`                    | x        | `{"message":"OK"}`      | HTTP response body                                                              |

### Condition expression

You can use formulas with variables.<br>
The formula should return a probability value which range of [0.0, 1.0].
If the formula returns Boolean value, it means 1.0 or 0.0.

#### Available variables

| Variable name     | Description                                                                        |
|:------------------|:-----------------------------------------------------------------------------------|
| routeRequestCount | Number of requests for the route. On the first request for route, this value is 1. |
| totalRequestCount | Total number of requests. On the first request, this value is 1.                   |

### Path parameters

You can use two kinds of path parameters.

- `:` prefix: Always match one segment.
- `*` prefix: Match any segments. (includes 0)

| Config path    | Request path                                                    |
|:---------------|:----------------------------------------------------------------|
| `/users/:id`   | ✅`/users/1`<br>✅`/users/1/`<br>❌`/users/`<br>❌`/users/1/follow` |
| `/users/*path` | ✅`/users`<br>✅`/users/1/follow`                                 |

### Response headers

Catfish automatically add some headers in responses to easily debug with.

| Header name                    | Required | Description                                                                   |
|:-------------------------------|:---------|:------------------------------------------------------------------------------|
| X-CATFISH-PATH                 | o        | Indicates that Catfish returned the response according to the path setting.   |
| X-CATFISH-RESPONSE-PRESET-NAME | o        | Indicates that Catfish returned the response according to the preset setting. |
| X-CATFISH-ERROR                | x        | The descriptions of an error.                                                 |

### Response body

You can use variables, etc with [text/template](https://pkg.go.dev/text/template) format.<br>
The data passed to the template engine is [main.Context]().
