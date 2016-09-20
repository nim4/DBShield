<?php

$link = pg_Connect("host=127.0.0.1 port=5432 dbname=postgres user=postgres password=xX123456") or die("failed");

$result = pg_exec($link, "select * from playground;");
echo pg_numrows($result);

