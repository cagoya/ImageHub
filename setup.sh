#!/bin/bash

# 获取docker容器ID
CONTAINER_ID=$(docker inspect -f '{{.Id}}' imagehub-postgres)

# 执行 INSERT 语句
docker exec -it $CONTAINER_ID \
psql -U postgres -d image_hub -c \
'INSERT INTO users (username, password_hash, email, is_admin, created_time, updated_time) VALUES 
(''张三'', ''$2a$10$t/OkkyKExIAjYRAklHQofesy9rmN6lSB74OyDCXZrYtfpXoIE1aIa'', ''zhangsan@gmail.com'', FALSE, NOW(), NOW()), 
(''李四'', ''$2a$10$t/OkkyKExIAjYRAklHQofesy9rmN6lSB74OyDCXZrYtfpXoIE1aIa'', ''lisi@gmail.com'', FALSE, NOW(), NOW()), 
(''王五'', ''$2a$10$t/OkkyKExIAjYRAklHQofesy9rmN6lSB74OyDCXZrYtfpXoIE1aIa'', ''wangwu@gmail.com'', TRUE, NOW(), NOW());'