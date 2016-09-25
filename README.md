# DBShield

Protects your data by inspecting incoming queries from your application server and rejecting abnormal ones. it also supports data masking to avoid data leaks.

# How it works?

For example this is how web server normally interacts with database server:

![Sample Web Server and DB](https://cdn.rawgit.com/nim4/DBShield/master/misc/how_01.png)

by adding DBShield in front of database server we can protect it against abnormal queries and apply masking to database responses. to detect abnormal queries we first run DBShield in learning mode. learning mode lets any query pass but it records information about it(pattern, username, time and source) into internal database.

![Learning mode](https://cdn.rawgit.com/nim4/DBShield/master/misc/how_02.png)


after collecting enough patterns we can run DBShield in protect mode. protect mode can distinguish abnormal query patterns, request time, user and source and take action based on configurations. it can also replace results of normal queries by applying data masking(ex. replace each digit of CC number by X except the last 3 digits)

![Protect mode](https://cdn.rawgit.com/nim4/DBShield/master/misc/how_03.png)


# Supported Databases:

| Database      | Protect       | Data Masking  |
| ------------- |:-------------:| -------------:|
| MariaDB       | Yes           |      Yes      |
| MySQL         | Yes           |      Yes      |
| Oracle        |     Yes       |      No      |

## To Do
 - Add Masking support to Oracle Database
 - Add Microsoft SQL Server support
