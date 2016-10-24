[![Linux & OS X](https://travis-ci.org/nim4/DBShield.svg?branch=master "Linux & OS X")](https://travis-ci.org/nim4/DBShield)
[![Windows](https://ci.appveyor.com/api/projects/status/github/nim4/DBShield?branch=master&svg=true "Windows")](https://ci.appveyor.com/project/nim4/DBShield/branch/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/nim4/DBShield)](https://goreportcard.com/report/github.com/nim4/DBShield)
[![codecov](https://codecov.io/gh/nim4/DBShield/branch/master/graph/badge.svg)](https://codecov.io/gh/nim4/DBShield)
[![Dev chat](https://img.shields.io/badge/gitter-chat-20cc20.svg "Dev chat")](https://gitter.im/DBShield/Lobby)
[![GoDoc](https://godoc.org/github.com/nim4/DBShield?status.svg)](https://godoc.org/github.com/nim4/DBShield)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/nim4/DBShield/master/LICENSE)
# DBShield

Protects your data by inspecting incoming queries from your application server and rejecting abnormal ones.


---
## How it works?

For example, this is how web server normally interacts with database server:

![Sample Web Server and DB](https://raw.githubusercontent.com/nim4/DBShield/master/misc/how_01.png)

By adding DBShield in front of database server we can protect it against abnormal queries. To detect abnormal queries we first run DBShield in learning mode. Learning mode lets any query pass but it records information about it (pattern, username, time and source) into the internal database.

![Learning mode](https://raw.githubusercontent.com/nim4/DBShield/master/misc/how_02.png)


After collecting enough patterns we can run DBShield in protect mode. Protect mode can distinguish abnormal query pattern, user and source and take action based on configurations.

![Protect mode](https://raw.githubusercontent.com/nim4/DBShield/master/misc/how_03.png)

---

## Sample Outputs

**CLI**

```
$ go run main.go
2016/10/15 16:25:31 [INFO]  Config file: /etc/dbshield.yml
2016/10/15 16:25:31 [INFO]  Internal DB: /tmp/model/10.0.0.21_postgres.db
2016/10/15 16:25:31 [INFO]  Listening: 0.0.0.0:5000
2016/10/15 16:25:31 [INFO]  Backend: postgres (10.0.0.21:5432)
2016/10/15 16:25:31 [INFO]  Protect: true
2016/10/15 16:25:31 [INFO]  Web interface on https://127.0.0.1:8070/
2016/10/15 16:25:33 [INFO]  Connected from: 10.0.0.20:35910
2016/10/15 16:25:33 [INFO]  Connected to: 10.0.0.21:5432
2016/10/15 16:25:33 [INFO]  SSL connection
2016/10/15 16:25:34 [DEBUG] Client handshake done
2016/10/15 16:25:34 [DEBUG] Server handshake done
2016/10/15 16:25:34 [INFO]  User: postgres
2016/10/15 16:25:34 [INFO]  Database: test
2016/10/15 16:25:34 [INFO]  Query: SELECT * FROM stocks where id=-1 or 1=1
2016/10/15 16:25:34 [WARN]  Pattern not found: [53 55 51 52 55 52 50 53 55 51 53 49 115 116 111 99 107 115 53 55 51 53 50 105 100 54 49 52 53 53 55 51 55 57 53 55 52 48 52 53 55 51 55 57 54 49 53 55 51 55 57] (SELECT * FROM stocks where id=-1 or 1=1)
2016/10/15 16:25:34 [WARN]  Dropping connection
```


**Web Interface**

![Web UI](https://raw.githubusercontent.com/nim4/DBShield/master/misc/graph.png)

---
## Demo
For demo, we are using [sqlmap](https://github.com/sqlmapproject/sqlmap)(automatic SQL injection and database takeover tool) to exploit the SQL injection vulnerability at `user.php`

In the first scenario, the sqlmap successfully exploits the SQL injection when web application connected directly to the database(MySQL), In the second scenario, we modify the `user.php` so DBShield gets between the web application and database which will drop the injection attempt and make sqlmap fail.

![Demo](misc/demo.gif)
---
## Installation

Get it
```
$ go get -u github.com/nim4/DBShield
```

Then you can see help using "-h" argument:
```
$ $GOPATH/bin/DBShield -h
DBShield 1.0.0-beta3
Usage of DBShield:
  -c string
    	Config file (default "/etc/dbshield.yml")
  -l	Get list of captured patterns
  -h	Show help
  -k	Show parsed config and exit
  -version
    	Show version

```

and run it with your configuration, like:
```
$ $GOPATH/bin/DBShield -c config.yml
```
see [sample configuration  file](https://github.com/nim4/DBShield/blob/master/conf/dbshield.yml)


>:warning: **WARNING:**
> Do NOT use default certificates in production environments!


---
## Supports:

| Database     | Protect | SSL |
|:------------:|:-------:|:---:|
| **DB2**   | ![Yes][YesImg] | ![No][NoImg]  |
| **MariaDB**  | ![Yes][YesImg] | ![Yes][YesImg] |
| **MySQL**    | ![Yes][YesImg] | ![Yes][YesImg] |
| **Oracle**   | ![Yes][YesImg] | ![No][NoImg]  |
| **Postgres** | ![Yes][YesImg] | ![Yes][YesImg] |

---
## To Do

(Sorted by priority)

 - Improve documentation
 - Add Microsoft SQL Server
 - Add more command-line arguments
 - Get 90% test coverage
 - Support Oracle SSL

 [YesImg]: https://raw.githubusercontent.com/nim4/DBShield/master/misc/yes.png
 [NoImg]: https://raw.githubusercontent.com/nim4/DBShield/master/misc/no.png
