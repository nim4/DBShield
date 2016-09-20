<?php
$servername = "localhost:5000";
$username = "root";
$password = "xX123456";

// Create connection
$conn = new mysqli($servername, $username, $password, "test");

// Check connection
if ($conn->connect_error) {
    die("Connection failed: " . $conn->connect_error);
}

$name = uniqid();
if ($result = $conn->query("INSERT INTO Persons(Name,City) VALUES(\"$name\",'".uniqid()."')")) {

}

if ($result = $conn->query("SELECT * FROM Persons WHERE name=\"$name\" and id>1 and 2.2>1.5 # or name=N1m4")) {
  foreach($result as $k => $v){
    print_r($v) ;
  }
  /* free result set */
  mysqli_free_result($result);
}

if ($result = $conn->query("SELECT * FROM Persons;")) {
  foreach($result as $k => $v){
    print_r($v) ;
  }
  /* free result set */
  mysqli_free_result($result);
}
$conn->close();
?>
