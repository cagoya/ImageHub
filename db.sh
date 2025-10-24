#!/bin/bash

# MySQL 连接信息 - 根据实际情况修改
DB_USER="root" # 修改为实际用户
DB_PASS="2023678ZZh#" # 修改为实际密码
DB_NAME="image_hub"

DB_EXISTS=$(mysql -u"${DB_USER}" -p"${DB_PASS}" -s -N -e "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME='${DB_NAME}'")

if [ -z "$DB_EXISTS" ]; then
    # $DB_EXISTS 为空，表示数据库不存在
    echo "数据库 '${DB_NAME}' 不存在，正在创建..."

    # 创建数据库
    # 使用 CREATE DATABASE IF NOT EXISTS 可以避免如果数据库在并发操作中被创建时产生的错误
    CREATE_STATUS=$(mysql -u"${DB_USER}" -p"${DB_PASS}" -e "CREATE DATABASE IF NOT EXISTS ${DB_NAME};" 2>&1)

    if [ $? -eq 0 ]; then
        echo "数据库 '${DB_NAME}' 创建成功。"
    else
        echo "数据库 '${DB_NAME}' 创建失败！"
        echo "错误信息: ${CREATE_STATUS}"
        exit 1
    fi
else
    # $DB_EXISTS 非空，表示数据库已存在
    echo "数据库 '${DB_NAME}' 已存在，无需创建。"
fi

exit 0