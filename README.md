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
$ go get -u github.com/nim4/DBShield
```

then you can execute it like:
```
$ $GOPATH/bin/DBShield -f $GOPATH/src/github.com/nim4/DBShield/conf/dbshield.yml
```

## Supports:

| Database     | Protect | SSL |
|:------------:|:-------:|:---:|
| **MariaDB**  | ![Yes](https://cdn.rawgit.com/nim4/DBShield/master/misc/yes.png) | ![Yes](https://cdn.rawgit.com/nim4/DBShield/master/misc/yes.png) |
| **MySQL**    | ![Yes](https://cdn.rawgit.com/nim4/DBShield/master/misc/yes.png) | ![Yes](https://cdn.rawgit.com/nim4/DBShield/master/misc/yes.png) |
| **Oracle**   | ![Yes](https://cdn.rawgit.com/nim4/DBShield/master/misc/yes.png) | ![No](https://cdn.rawgit.com/nim4/DBShield/master/misc/no.png)  |
| **Postgres** | ![Yes](https://cdn.rawgit.com/nim4/DBShield/master/misc/yes.png) | ![Yes](https://cdn.rawgit.com/nim4/DBShield/master/misc/yes.png) |

## More
- [Sample configuration  file](https://github.com/nim4/DBShield/blob/master/conf/dbshield.yml)
- [GoDoc](https://godoc.org/github.com/nim4/DBShield/dbshield)

## To Do

 - Add Microsoft SQL Server
 - Add "exec" action
 - Add more command-line arguments
 - Improve documentation
