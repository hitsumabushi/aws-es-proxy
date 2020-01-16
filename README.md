# aws-es-proxy-go

aws-es-proxy-go is the proxy for signing AWS signature v4.

## setup & configuration

1. build aws-es-proxy-go
  - run make build, then aws-es-proxy-go created under build dir
2. Write config file. see example directory
3. Configure AWS credentials: instance profile/env vars/credential files
4. run `./build/run -config config.json`

## How does it work

If config is followings:
```json
{
    "list_path": "/_list",
    "server_map": {
        "/id001": {
            "region": "us-west-2",
            "host": "id001-foo.us-west-2.es.amazonaws.com"
        },
        "/id002": {
            "region": "us-west-2",
            "host":  "id002-bar.us-west-2.es.amazonaws.com"
        }
    }
}
```

* The request `http://{this_proxy}/id001/foo/bar` transfer to `https://id001-foo.us-west-2.es.amazonaws.com/foo/bar`, and the request signed by sigv4.
* The request `http://{this_proxy}/_list` return endpoint list.
