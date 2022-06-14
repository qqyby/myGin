echo "building..."
echo ""

gitBranch=$(git symbolic-ref --short -q HEAD)
gitTagCommit=$(git describe --tags --always)
gitCommit="$gitBranch-$gitTagCommit"
buildDate=$(TZ=Asia/Shanghai date +%FT%T%z)
buildVersion="1.0.0"

build()
{
    go build -o bravo_zhurong -ldflags "-X main.buildTime=${buildDate} -X main.buildVersion=${buildVersion} -X main.gitCommitID=${gitCommit}"
    ret=$?; if [ 0 -ne $ret ]; then echo "build failed. ret=$ret"; exit $ret; fi
}

build

go vet ./

echo "build success"
echo ""
echo "* 导入数据库:"
echo "    mysql -uroot -p < docs/db.sql"
echo "* 启动服务:"
echo "    ./bravo_zhurong [-c configs/config.yaml]"