<?php
$servername = "localhost:5005";
$username = "root";
$password = "xX123456";

// Create connection
$conn = new mysqli($servername, $username, $password);

// Check connection
if ($conn->connect_error) {
    die("Connection failed: " . $conn->connect_error);
}


if (!$conn->select_db("test")) {
    die("Selecting DB failed.");
}


if ($result = $conn->query("show tables;")) {
  foreach($result as $k => $v){
    print_r($v) ;
  }
  /* free result set */
  mysqli_free_result($result);
}
$conn->close();
?>
