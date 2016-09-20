<?php
$conn = oci_connect('test', 'xX123456', 'localhost:5000/XE');
if (!$conn) {
    $e = oci_error();
    trigger_error(htmlentities($e['message'], ENT_QUOTES), E_USER_ERROR);
}

$stid = oci_parse($conn, 'SELECT CUSTOMER_ID as myid, CUST_FIRST_NAME as name FROM DEMO_CUSTOMERS');


oci_execute($stid, OCI_DESCRIBE_ONLY); // Use OCI_DESCRIBE_ONLY if not fetching rows


$ncols = oci_num_fields($stid);

for ($i = 1; $i <= $ncols; $i++) {
    $column_name  = oci_field_name($stid, $i);
    $column_type  = oci_field_type($stid, $i);

    echo "$column_name\t$column_type\n";
}

oci_free_statement($stid);
oci_close($conn);
