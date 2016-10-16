[![Build Status](https://travis-ci.org/nim4/DBShield.svg?branch=master)](https://travis-ci.org/nim4/DBShield)
[![Go Report Card](https://goreportcard.com/badge/github.com/nim4/DBShield)](https://goreportcard.com/report/github.com/nim4/DBShield)
[![Dev chat](https://img.shields.io/badge/gitter-chat-20cc20.svg)](https://gitter.im/DBShield/Lobby)
[![GoDoc](https://godoc.org/github.com/nim4/DBShield?status.svg)](https://godoc.org/github.com/nim4/DBShield)

# DBShield

Protects your data by inspecting incoming queries from your application server and rejecting abnormal ones.



## How it works?

For example, this is how web server normally interacts with database server:

![Sample Web Server and DB](https://raw.githubusercontent.com/nim4/DBShield/master/misc/how_01.png)

By adding DBShield in front of database server we can protect it against abnormal queries. To detect abnormal queries we first run DBShield in learning mode. Learning mode lets any query pass but it records information about it (pattern, username, time and source) into internal database.

![Learning mode](https://raw.githubusercontent.com/nim4/DBShield/master/misc/how_02.png)


After collecting enough patterns we can run DBShield in protect mode. Protect mode can distinguish abnormal query patterns, user and source and take action based on configurations.

![Protect mode](https://raw.githubusercontent.com/nim4/DBShield/master/misc/how_03.png)

**Sample Output:**

```
$ go run main.go
2016/10/15 16:25:31 [INFO]  Config file: /etc/dbshield.yml
2016/10/15 16:25:31 [INFO]  Internal DB: /tmp/model/127.0.0.1_postgres.db
2016/10/15 16:25:31 [INFO]  Listening: 0.0.0.0:5000 (Threads: 5)
2016/10/15 16:25:31 [INFO]  Backend: postgres (127.0.0.1:5432)
2016/10/15 16:25:31 [INFO]  Protect: true
2016/10/15 16:25:33 [INFO]  Connected from: 127.0.0.1:35910
2016/10/15 16:25:33 [INFO]  Connected to: 127.0.0.1:5432
2016/10/15 16:25:33 [INFO]  SSL connection
2016/10/15 16:25:34 [DEBUG] Client handshake done
2016/10/15 16:25:34 [DEBUG] Server handshake done
2016/10/15 16:25:34 [INFO]  User: postgres
2016/10/15 16:25:34 [INFO]  Database: test
2016/10/15 16:25:34 [INFO]  Query: SELECT * FROM stocks where id=-1 or 1=1
2016/10/15 16:25:34 [WARN]  Pattern not found: [53 55 51 52 55 52 50 53 55 51 53 49 115 116 111 99 107 115 53 55 51 53 50 105 100 54 49 52 53 53 55 51 55 57 53 55 52 48 52 53 55 51 55 57 54 49 53 55 51 55 57] (SELECT * FROM stocks where id=-1 or 1=1)
2016/10/15 16:25:34 [WARN]  Dropping connection
```

## Installation

```
$ go get -u github.com/nim4/DBShield
```

then you can execute it like:
```
$ $GOPATH/bin/DBShield -c $GOPATH/src/github.com/nim4/DBShield/conf/dbshield.yml
```

## Supports:

| Database     | Protect | SSL |
|:------------:|:-------:|:---:|
| **MariaDB**  | ![Yes](https://raw.githubusercontent.com/nim4/DBShield/master/misc/yes.png) | ![Yes](https://raw.githubusercontent.com/nim4/DBShield/master/misc/yes.png) |
| **MySQL**    | ![Yes](https://raw.githubusercontent.com/nim4/DBShield/master/misc/yes.png) | ![Yes](https://raw.githubusercontent.com/nim4/DBShield/master/misc/yes.png) |
| **Oracle**   | ![Yes](https://raw.githubusercontent.com/nim4/DBShield/master/misc/yes.png) | ![No](https://raw.githubusercontent.com/nim4/DBShield/master/misc/no.png)  |
| **Postgres** | ![Yes](https://raw.githubusercontent.com/nim4/DBShield/master/misc/yes.png) | ![Yes](https://raw.githubusercontent.com/nim4/DBShield/master/misc/yes.png) |

## More
- [Sample configuration  file](https://github.com/nim4/DBShield/blob/master/conf/dbshield.yml)

## To Do

 - Add Microsoft SQL Server
 - Add more command-line arguments
 - Improve documentation
