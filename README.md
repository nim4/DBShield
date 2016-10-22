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

By adding DBShield in front of database server we can protect it against abnormal queries. To detect abnormal queries we first run DBShield in learning mode. Learning mode lets any query pass but it records information about it (pattern, username, time and source) into internal database.

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

For testing we have a vulnerable page at `http://192.168.22.1/user.php`

`user.php` contents:
```php
<?php
include('config.php');
// Create connection
$conn = new mysqli($servername, $username, $password, "test");

// Check connection
if ($conn->connect_error) {
    die("Connection failed: " . $conn->connect_error);
}

if (!empty($_GET['id'])){
  if ($result = $conn->query("SELECT * FROM Persons WHERE id=".$_GET['id'])) {
    foreach($result as $k => $v){
      echo "Name: <b>{$v['Name']}</b><br />City: <b>{$v['City']}</b>" ;
    }
    mysqli_free_result($result);
  }
 }
$conn->close();
```

We are using [sqlmap](https://github.com/sqlmapproject/sqlmap) for exploiting the vulnerability, result are as below:

```
$ sqlmap -u http://192.168.22.1/user.php?id=1
```
```
[12:14:31] [INFO] GET parameter 'id' is 'Generic UNION query (NULL) - 1 to 20 columns' injectable
GET parameter 'id' is vulnerable. Do you want to keep testing the others (if any)? [y/N]
sqlmap identified the following injection point(s) with a total of 53 HTTP(s) requests:
---
Parameter: id (GET)
    Type: boolean-based blind
    Title: AND boolean-based blind - WHERE or HAVING clause
    Payload: id=1 AND 8909=8909

    Type: AND/OR time-based blind
    Title: MySQL >= 5.0.12 AND time-based blind (SELECT)
    Payload: id=1 AND (SELECT * FROM (SELECT(SLEEP(5)))eIyW)

    Type: UNION query
    Title: Generic UNION query (NULL) - 3 columns
    Payload: id=1 UNION ALL SELECT NULL,NULL,CONCAT(0x71786b7071,0x64666b56715965797a6e654141634c765a6575674b79686461476c5556766671584f74486c5a5a58,0x717a717a71)-- -
---
[12:14:33] [INFO] the back-end DBMS is MySQL
web server operating system: Linux Ubuntu
web application technology: PHP 7.0.8
back-end DBMS: MySQL 5.0.12
```

Then we try to exploiting it again using the same tool, this time DBShield is protecting the database:

```
$ sqlmap -u http://192.168.22.1/user.php?id=1
```

```
[12:20:36] [INFO] testing 'Oracle AND time-based blind'
[12:20:37] [INFO] testing 'Generic UNION query (NULL) - 1 to 10 columns'
[12:20:37] [WARNING] using unescaped version of the test because of zero knowledge of the back-end DBMS. You can try to explicitly set it using option '--dbms'
[12:20:43] [INFO] testing 'MySQL UNION query (NULL) - 1 to 10 columns'
[12:20:47] [WARNING] GET parameter 'id' is not injectable
```
---
## Installation

Get it
```
$ go get -u github.com/nim4/DBShield
```

Then you can get help using "-h" argument:
```
$ $GOPATH/bin/DBShield -h
DBShield 1.0.0-beta2
Usage of DBShield:
  -c string
    	Config file (default "/etc/dbshield.yml")
  -d	Get list of captured patterns
  -h	Show help
  -k	Show parsed config and exit
  -version
    	Show version

```

and run it with your configuration like:
```
$ $GOPATH/bin/DBShield -c config.yml
```
see [sample configuration  file](https://github.com/nim4/DBShield/blob/master/conf/dbshield.yml)

---
## Supports:

| Database     | Protect | SSL |
|:------------:|:-------:|:---:|
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
 - Add IBM DB2
 - Get 90% test coverage
 - Support Oracle SSL

 [YesImg]: https://raw.githubusercontent.com/nim4/DBShield/master/misc/yes.png
 [NoImg]: https://raw.githubusercontent.com/nim4/DBShield/master/misc/no.png
