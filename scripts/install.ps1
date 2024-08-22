# Stop on error
$ErrorActionPreference = "Stop"
$version = "0.10.1"
$platform = "windows"
$destination = "wrangle-$version-$platform-$architecture"
$archive = "$destination.zip"

$computerInfo = Get-ComputerInfo
if ($computerInfo.OsType -ne "WINNT"){
    Write-Error -Message "Unkown os type ${computerInfo.OsType}. Expected 'WINNT'"
    return
}

$architecture = "amd64"
if ($computerInfo.OsArchitecture -ne "64-bit"){
    Write-Error -Message "Unkown os architecture ${computerInfo.OsArchitecture}. Expected '64-bit'"
    return
}

$url = "https://github.com/patrickhuber/wrangle/releases/download/$version/$archive"

"downloading $url to $archive"

# download the archive
Invoke-WebRequest -uri $url -OutFile $archive

# expand the archive and remove the archive file
Expand-Archive -Path $archive -DestinationPath $destination
Remove-Item -Path $archive -Force

# run bootstrap command and cleanup downloaded executable
Invoke-Expression "$destination/wrangle.exe bootstrap"
Remove-Item -Path $destination -Force -Recurse