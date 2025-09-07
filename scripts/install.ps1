# Stop on error
$ErrorActionPreference = "Stop"
$latest = Invoke-WebRequest 'https://api.github.com/repos/patrickhuber/wrangle/releases/latest'
$json = $latest.Content | ConvertFrom-Json
$version = $json.tag_name.TrimStart("v")
"version: $version"

$computerInfo = Get-ComputerInfo

$platform = ""
switch ($computerInfo.OsType){
    "WINNT" { $platform = "windows" }
    "LINUX" { $platform = "linux" }
    "MACROS" { $platform = "darwin" }
    default  { Write-Error -Message "Unkown os type ${computerInfo.OsType}. Expected 'windows', 'linux' or 'darwin'" }
}
"platform: $platform"

$architecture = "amd64"
if ($computerInfo.OsArchitecture -ne "64-bit"){
    Write-Error -Message "Unkown os architecture ${computerInfo.OsArchitecture}. Expected '64-bit'"
    return
}
"architecture: $architecture"

foreach ($asset in $json.assets){
    if (-not $asset.name.Contains($version)){
        continue
    }
    if (-not $asset.name.Contains($platform)){
        continue
    }
    if (-not $asset.name.Contains($architecture)){
        continue
    }
    
    $url = $asset.browser_download_url
    $archive = $asset.name
    $destination = [System.IO.Path]::GetFileNameWithoutExtension($archive)
    
    "downloading $url to $archive"

    # download the archive
    Invoke-WebRequest -uri $url -OutFile $archive

    # expand the archive and remove the archive file
    Expand-Archive -Path $archive -DestinationPath $destination
    Remove-Item -Path $archive -Force

    # run bootstrap command and cleanup downloaded executable
    $env:WRANGLE_LOG_LEVEL="debug"
    Invoke-Expression "$destination/wrangle.exe bootstrap"
    Remove-Item -Path $destination -Force -Recurse
    break
}