# Container debugging image (CDI)

Have you ever been in a situation where you need to debug/test/demo a container and you someshow got stuck trying to figure out how the requests/responses are getting lost in the networking space? Or redirected into oblivion by an ingress controller? Or wanted to make sure each of the requests is nicely loadbalanced between multiple hosts? If you answered "yes" to any of those questions - this image may be of interest to you.

## How to use it?

Look at `compose.yml` - all the important settings are there.

## Example response?

```none
HTTP URL                    : /
HTTP Host                   : localhost:8080
HTTP Listen port            : 8080
HTTP Referer                :
HTTP User agent             : curl/7.81.0

Server hostname             : afbd81cba828
Server's address            : 172.22.0.2                              afbd81cba828

Client's IP                 : 172.22.0.1 devbox
```

Each section can be turned on/off depending on your needs/requirements (again - see `compose.yml`).

## License

MIT. See `LICENSE` file for details.
