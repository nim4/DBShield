# DBShield

Protects your data by inspecting incoming queries from your application server and rejecting abnormal ones.

## How it works?

For example, this is how web server normally interacts with database server:

![Sample Web Server and DB](https://cdn.rawgit.com/nim4/DBShield/master/misc/how_01.png)

By adding DBShield in front of database server we can protect it against abnormal queries. To detect abnormal queries we first run DBShield in learning mode. Learning mode lets any query pass but it records information about it (pattern, username, time and source) into internal database.

![Learning mode](https://cdn.rawgit.com/nim4/DBShield/master/misc/how_02.png)


After collecting enough patterns we can run DBShield in protect mode. Protect mode can distinguish abnormal query patterns, request time, user and source and take action based on configurations.

![Protect mode](https://cdn.rawgit.com/nim4/DBShield/master/misc/how_03.png)

## Installation

```
$ go get github.com/nim4/DBShield
```

then you can execute it like:
```
$ $GOPATH/bin/DBShield -f $GOPATH/src/github.com/nim4/DBShield/conf/dbshield.yml
```



## Supported Databases:

| Database    | Protect | TLS |
| ----------- |:-------:| ---:|
| **MariaDB** | Yes     | Yes |
| **MySQL**   | Yes     | Yes |
| **Oracle**  | Yes     | Yes |


## To Do
 - Add Postgres
 - Add Microsoft SQL Server
