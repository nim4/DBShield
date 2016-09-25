# DBShield

Protects your data by inspecting incoming queries from your application server and rejecting abnormal ones. it also supports data masking to avoid data leaks.

# How it works?

For example this is how your web server normally interacts with database:

![Sample Web Server and DB](https://cdn.rawgit.com/nim4/DBShield/master/misc/how_01.png)


![Learning mode](https://cdn.rawgit.com/nim4/DBShield/master/misc/how_02.png)
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
