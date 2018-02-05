@echo off
color 0a
title EasyDoc https://easydoc.089858.com

::EasyDoc https://easydoc.089858.com

::依赖：EasyDoc二进制文件要放在全局环境目录里(或与此脚本在同目录下)。
easydoc -build

@echo.
pause
