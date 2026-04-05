runner build
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

cd bin
.\runner.exe install
cd ..