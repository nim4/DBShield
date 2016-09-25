# DBShield

Protects your data by inspecting incoming queries from your application server and rejecting abnormal ones. It also supports data masking to avoid data leaks.

## How it works?

For example, this is how web server normally interacts with database server:

![Sample Web Server and DB](https://cdn.rawgit.com/nim4/DBShield/master/misc/how_01.png)

By adding DBShield in front of database server we can protect it against abnormal queries and apply masking to database responses. To detect abnormal queries we first run DBShield in learning mode. Learning mode lets any query pass but it records information about it (pattern, username, time and source) into internal database.

![Learning mode](https://cdn.rawgit.com/nim4/DBShield/master/misc/how_02.png)


After collecting enough patterns we can run DBShield in protect mode. Protect mode can distinguish abnormal query patterns, request time, user and source and take action based on configurations. It can also replace results of normal queries by applying data masking (ex. replace each digit of CC number by "X" except the last three digits)

![Protect mode](https://cdn.rawgit.com/nim4/DBShield/master/misc/how_03.png)

## Installation

```bash
$ go get github.com/nim4/DBShield
$ go install github.com/nim4/DBShield
# Copy the sample config file
$ sudo cp $GOPATH/src/github.com/nim4/DBShield/conf/dbshield.yml /etc/dbshield.yml
$ $GOPATH/bin/DBShield
```



## Supported Databases:

| Database      | Protect       | Data Masking  |
| ------------- |:-------------:| -------------:|
| MariaDB       | Yes           |      Yes      |
| MySQL         | Yes           |      Yes      |
| Oracle        | Yes           |      No       |

## To Do
 - Add Masking support for Oracle Database
 - Add Microsoft SQL Server
 - Add Postgres
