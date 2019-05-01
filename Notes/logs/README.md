# logs

## `access_log.gz`

1.5 million lines of Apache access logs from uc3-mrtstore1-prd, covering 15
July 2018 – 28 Apr 2019. Use `gzcat` or `gunzip -c` to read without
uncompressing. (Uncompressed size: 330 MB) 

## `requests.gz`

The extract request from each line, except for the following invalid request from `ingest02`:

```
"-" - - [25/Feb/2019:10:34:33 -0800] "\xff\xf4\xff\xfd\x06" 400 319 "-" "-" 0 7 0 "172.30.8.76"
```

## Requests

### GET requests

#### Content requests

GET requests are in one of the following forms:

| Count   | Form                                                                       | Purpose                          |
| ---     | ---                                                                        | ---                              |
| 928,812 | `GET /content/<node>/<encoded-ark>/<version>/<encoded-pathname>`           | file download                    |
| 389,130 | `GET /content/<node>/<encoded-ark>/<version>/<encoded-pathname>?fixity=no` | file download for audit          |
| 79,107  | `GET /manifest/<node>/<encoded-ark>`                                       | manifest request                 |
| 1236    | `GET /producer/<node>/<encoded-ark>/<version>?t=zip`                       | version download (user-friendly) |
| 1226    | `GET /content/<node>/<encoded-ark>?t=zip`                                  | object download                  |
| 549     | `GET /ping?t=xml&SLEEP=<time>`                                             | ping                             |
| 326     | `GET /?t=xml&SLEEP=`                                                       | *                                |
| 174     | `GET /content/<node>/<encoded-ark>/<version>?t=zip`                        | version download                 |
| 11      | `GET /state?t=xml`                                                         | state                            |
| 3       | `GET /state/<node>/<encoded-ark>`                                          | state (object)                   |
| 2       | `GET /state/<node>?t=anvl`                                                 | state (node)                     |
| 1       | `GET /state/<node>/<encoded-ark>/<version>?t=xml`                          | state (version)                  |
| 1       | `GET /producer/<node>/<encoded-ark>?t=zip`                                 | object download (user-friendly)  |
| 1       | `GET /cloudcontainer/`                                                     | *                                |

<sup>* Response code always 404 Not Found</sup>

### POST requests

POST requests are in one of the following forms:

| Count  | Form                                             | Purpose                               |
| ---    | ---                                              | ---                                   |
| 69,451 | `POST /add/<node>/<encoded-ark>`                 |                                       |
| 7,210  | `POST /update/<node>/<encoded-ark>`               |                                       |
| 88     | `POST /producerasync/<node>/<encoded-ark>`       | async object download (user-friendly) |
| 1      | `POST /copy/<from-node>/<to-node>/<encoded-ark>` | *                                     |
| 0      | `POST /async/<node>/<encoded-ark>`               | async object download†                |

<sup>* Response code 404 Not Found</sup>

<sup>† Not present in dataset, but hypothetically possible</sup>

### DELETE requests

DELETE requests are in the following form:

| Count | Form | Purpose |
| --- | --- | --- |
| 15,764 | `DELETE /content/<node>/<encoded-ark>` | delete object |

Note that 4,598 of 15,764 DELETE requests resulted in 404 Not Found.

### OPTIONS requests

OPTIONS requests are in the following form:

| Count | Form | Purpose |
| --- | --- | --- |
| 6904 | `OPTIONS *` | |

### HEAD requests

| Count | Form                                  | Purpose |
| ---   | ---                                   | ---     |
| 1     | `HEAD /state/<node>/<encoded-ark>`    |         |
| 1     | `HEAD /manifest/<node>/<encoded-ark>` | *        |

<sup>* Response code 500 Internal Server Error</sup>

<!--
egrep 'GET /content/[0-9]+/[^/]+/[0-9]+/[^?/]+$' get-requests.txt | wc -l
  928812

egrep 'GET /content/[0-9]+/[^/]+/[0-9]+/[^/]+\?[^?]+$' get-requests.txt | wc -l
  389130

egrep 'GET /content/[0-9]+/[^?/]+\?t=zip$' get-requests.txt | wc -l
  1226

egrep 'GET /content/[0-9]+/[^?/]+/[0-9]+\?t=zip$' get-requests.txt | wc -l
  174

egrep 'GET /producer/[0-9]+/[^?/]+\?t=zip$' get-requests.txt | wc -l
  1

egrep 'GET /producer/[0-9]+/[^?/]+/[0-9]+\?t=zip$' get-requests.txt | wc -l
  1236

egrep 'GET /manifest/[0-9]+/[^?/]+$' get-requests.txt | wc -l
  79107

egrep 'GET /state/[0-9]+/[^?/]+$' get-requests.txt | wc -l
  3
-->

