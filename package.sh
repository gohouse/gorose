#!/bin/bash

# 工具函数
go get github.com/gohouse/utils

# 以下为数据库驱动, 按需安装

# mysql 驱动
go get github.com/go-sql-driver/mysql

# sqlite 驱动
go get github.com/mattn/go-sqlite3

# postgresql 驱动
go get github.com/lib/pq

# Oracle 驱动
go get github.com/mattn/go-oci8