@echo off
REM 获取 docker 容器ID
FOR /F "delims=" %%i IN ('docker inspect -f "{{.Id}}" imagehub-postgres') DO SET CONTAINER_ID=%%i

IF NOT DEFINED CONTAINER_ID (
    ECHO 错误：未能获取容器ID。请确认容器 "imagehub-postgres" 是否正在运行。
    EXIT /B 1
)

ECHO 容器ID: %CONTAINER_ID%

REM SQL 插入语句

ECHO 正在执行 INSERT 语句...

docker exec -i %CONTAINER_ID% ^
psql -U postgres -d image_hub -c ^
"INSERT INTO users (username, password_hash, email, is_admin, created_time, updated_time) VALUES ('张三', '$2a$10$t/OkkyKExIAjYRAklHQofesy9rmN6lSB74OyDCXZrYtfpXoIE1aIa', 'zhangsan@gmail.com', FALSE, NOW(), NOW()), ('李四', '$2a$10$t/OkkyKExIAjYRAklHQofesy9rmN6lSB74OyDCXZrYtfpXoIE1aIa', 'lisi@gmail.com', FALSE, NOW(), NOW()), ('王五', '$2a$10$t/OkkyKExIAjYRAklHQofesy9rmN6lSB74OyDCXZrYtfpXoIE1aIa', 'wangwu@gmail.com', TRUE, NOW(), NOW());"

IF ERRORLEVEL 1 (
    ECHO 错误：SQL 插入语句执行失败。
) ELSE (
    ECHO 成功：SQL 插入语句执行成功。
)

PAUSE