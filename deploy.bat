@echo off
REM 需求：ssh root@47.91.151.207 发布到这台机器的 /root/babyhabit 目录，内部包含前后端的代码。这台机器已经能够免密登录
REM 编译前端代码和后端代码都放到 build 目录下
REM 实现：
REM 1. 前端代码打包：在前端项目根目录执行 npm run build，将打包后的代码复制到后端项目的 static 目录下
REM 2. 后端代码打包：在后端项目根目录执行 go build -o babyhabit，将打包后的代码复制到 build 目录下
REM 3. 重启服务：在 /root/babyhabit 目录下执行 ./babyhabit restart 重启服务

setlocal enabledelayedexpansion

REM 清理并创建 build 目录
if exist build rmdir /s /q build
mkdir build
mkdir build\static

REM 1. 前端代码打包
echo Building frontend...
cd frontend
call npm run build
xcopy /e /i /y dist ..\build\static

REM 2. 后端代码打包
echo Building backend...
cd ..\backend
set GOOS=linux
set GOARCH=amd64
go build -o babyhabit
copy babyhabit ..\build\
copy .env ..\build\
REM xcopy /e /i /y files ..\build\files

REM 复制启动脚本
echo Copying startup script...
cd ..
copy build.sh build\

REM 复制files目录
echo Copying files directory...

REM 3. 上传到远程服务器
echo Deploying to remote server...
ssh -p 22 root@47.91.151.207 "supervisorctl stop babyhabit"
scp -P 22 -r build\* root@47.91.151.207:/root/babyhabit/

REM 4. 重启远程服务
echo Restarting remote service...
ssh -p 22 root@47.91.151.207 "supervisorctl start babyhabit"

echo Deployment completed successfully!
pause
