# Git project root
$project_root = (git rev-parse --show-toplevel)
$front_pub = "$project_root\frontend\public"
$back_pub = "$project_root\backend\public"

# Build frontend
Write-Output "Building frontend..."
Start-Process -Wait -NoNewWindow -WorkingDirectory "$project_root\frontend" npm -ArgumentList "run build"

# Build backend
Write-Output "Building backend..."
if(Test-Path "$back_pub") { Remove-Item -Recurse -Force "$back_pub"}
Copy-Item -Recurse "$front_pub" "$back_pub"
Start-Process -Wait -NoNewWindow -WorkingDirectory "$project_root\backend" go -ArgumentList "build -o .\build"