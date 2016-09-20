<?php
$servername = "localhost:5005";
$username = "root";
$password = "xX123456";

// Create connection
$conn = new mysqli($servername, $username, $password, "test");
$conn.query("show tables");
