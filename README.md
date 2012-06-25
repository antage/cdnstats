# CDNStats

CDNStat is a daemon collecting various statistics from nginx requests: request
count, sent bytes, referer, URI path.

## Installation

```
mkdir app
GOPATH=app go get -u github.com/antage/cdnstats
app/bin/cdnstats -h x.x.x.x -p pppp
```

Where x.x.x.x is host name, pppp is port number.
Default values: 127.0.0.1:9090

Now you can open browser at http://x.x.x.x:pppp/ for web-page displaying
statistics.

## How to setup Nginx?

Configuration example:

```
server {
    location / {
        root /var/www;
        post_action @stats; # after each request send information to cdnstats
    }

    location @stats {
        proxy_pass http://x.x.x.x:pppp/collect?bucket=[bucket name]&s=[hostname]&uri=$uri;
        # don't wait too long
        proxy_send_timeout 5s;
        proxy_read_timeout 5s;

        # optional header if you use domain name instead of ip-address x.x.x.x 
        # proxy_set_header Host cdnstat.example.org;

        # this headers are used by cdnstats
        proxy_set_header X-Bytes-Sent $body_bytes_sent;
        # Referer header is sent implicitly

        # delete unused headers
        proxy_set_header Accept "";
        proxy_set_header Accept-Encoding "";
        proxy_set_header Accept-Language "";
        proxy_set_header Accept-Charset "";
        proxy_set_header User-Agent "";
        proxy_set_header Cookie "";

        # don't send POST-request body
        proxy_pass_request_body off;
    }
}
```
