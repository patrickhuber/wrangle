# Stop on error
$ErrorActionPreference = "Stop"

$platform = "windows"
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

$version = "0.10.0"
$archive = "wrangle-$platform-$architecture.zip"
$url = "https://github.com/patrickhuber/wrangle/releases/download/$version/$archive"

# download the archive
Invoke-WebRequest -uri $url -OutFile $archive

# expand the archive and remove the archive
Expand-Archive -Path $archive
Remove-Item -Path $archive

# run bootstrap command and cleanup downloaded executable
wrangle bootstrap
Remove-Item wrangle.exe