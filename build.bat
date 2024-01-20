@echo off
title Build Production File
IF EXIST ../_cloud_production/chs_cloud_backend.zip DEL /F ../_cloud_production/chs_cloud_backend.zip
@echo Compress build output
powershell -Command "& {Compress-Archive -LiteralPath 'cmd','config','database','docker','migrations','internal','pkg', 'go.mod','go.sum', 'docker-compose.yml'  -DestinationPath ../_cloud_production/chs_cloud_backend.zip -Verbose  -Force}"
@echo Compress done
@echo Build has complete