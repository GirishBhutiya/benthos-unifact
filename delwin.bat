@echo OFF
IF EXIST tmp (
    Echo deleting tmp folder
    CD /D "tmp"
    for /F "delims=" %%G in ('dir /b') do (
        REM check if it is a directory or file
        IF EXIST "%%G\" (
            rmdir "%%G" /s /q
        ) else (
            del "%%G" /q
        )
    )
    CD ..
    rmdir tmp
)
IF EXIST target (
    Echo deleting target folder
    CD /D "target"
    for /F "delims=" %%G in ('dir /b') do (
        REM check if it is a directory or file
        IF EXIST "%%G\" (
            rmdir "%%G" /s /q
        ) else (
            del "%%G" /q
        )
    )
    CD ..
    rmdir target
)


