param(
    [string]$param1
)

git add .
git commit -m $param1
git push --set-upstream origin main -f
