package rpc

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/33cn/chain33/blockchain"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/common/limits"
	"github.com/33cn/chain33/common/log"
	"github.com/33cn/chain33/executor"
	"github.com/33cn/chain33/mempool"
	"github.com/33cn/chain33/p2p"
	"github.com/33cn/chain33/queue"
	"github.com/33cn/chain33/rpc"
	"github.com/33cn/chain33/store"
	"github.com/33cn/chain33/system/consensus/solo"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/util"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	_ "github.com/33cn/chain33/system"
	_ "github.com/33cn/plugin/plugin/store/init"

	jsonrpc "github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	cty "github.com/33cn/chain33/system/dapp/coins/types"
	gty "github.com/33cn/plugin/plugin/dapp/guess/types"
	pty "github.com/33cn/plugin/plugin/dapp/norm/types"
)

var (
	secp crypto.Crypto

	config = `# Title为local，表示此配置文件为本地单节点的配置。此时本地节点所在的链上只有这一个节点，共识模块一般采用solo模式。
Title="local"
TestNet=true
FixTime=false
[crypto]
[log]
# 日志级别，支持debug(dbug)/info/warn/error(eror)/crit
loglevel = "info"
logConsoleLevel = "info"
# 日志文件名，可带目录，所有生成的日志文件都放到此目录下
logFile = "logs/chain33.log"
# 单个日志文件的最大值（单位：兆）
maxFileSize = 300
# 最多保存的历史日志文件个数
maxBackups = 100
# 最多保存的历史日志消息（单位：天）
maxAge = 28
# 日志文件名是否使用本地时间（否则使用UTC时间）
localTime = true
# 历史日志文件是否压缩（压缩格式为gz）
compress = true
# 是否打印调用源文件和行号
callerFile = false
# 是否打印调用方法
callerFunction = false
[blockchain]
# 缓存区块的个数
defCacheSize=128
# 同步区块时一次最多申请获取的区块个数
maxFetchBlockNum=128
# 向对端节点请求同步区块的时间间隔
timeoutSeconds=5
# 使用的数据库类型
driver="leveldb"
# 数据库文件目录
dbPath="datadir"
# 数据库缓存大小
dbCache=64
# 是否为单节点
singleMode=true
# 同步区块批量写数据库时，是否需要立即写磁盘，非固态硬盘的电脑可以设置为false，以提高性能
batchsync=false
# 是否记录添加或者删除区块的序列，若节点作为主链节点，为平行链节点提供服务，需要设置为true
isRecordBlockSequence=true
# 是否为平行链节点
isParaChain=false
# 是否开启交易快速查询索引
enableTxQuickIndex=false
[p2p]
types=["dht"]
enable=true
# 使用的数据库类型
driver="leveldb"
# 数据库文件目录
dbPath="datadir/addrbook"
# 数据库缓存大小
dbCache=4
# GRPC请求日志文件
grpcLogFile="grpc33.log"

[rpc]
# jrpc绑定地址
jrpcBindAddr="localhost:9801"
# grpc绑定地址
grpcBindAddr="localhost:9802"
# 白名单列表，允许访问的IP地址，默认是“*”，允许所有IP访问
whitelist=["127.0.0.1"]
# jrpc方法请求白名单，默认是“*”，允许访问所有RPC方法
jrpcFuncWhitelist=["*"]
# jrpc方法请求黑名单，禁止调用黑名单里配置的rpc方法，一般和白名单配合使用，默认是空
# jrpcFuncBlacklist=["xxxx"]
# grpc方法请求白名单，默认是“*”，允许访问所有RPC方法
grpcFuncWhitelist=["*"]
# grpc方法请求黑名单，禁止调用黑名单里配置的rpc方法，一般和白名单配合使用，默认是空
# grpcFuncBlacklist=["xxx"]
# 是否开启https
enableTLS=false
# 证书文件，证书和私钥文件可以用cli工具生成
certFile="cert.pem"
# 私钥文件
keyFile="key.pem"
[mempool]
# mempool队列名称，可配，timeline，score，price
name="timeline"
# mempool缓存容量大小，默认10240
poolCacheSize=10240
# 最小得交易手续费用，这个没有默认值，必填，一般是100000
minTxFeeRate=100000
# 每个账户在mempool中得最大交易数量，默认100
maxTxNumPerAccount=10000
# timeline 是默认的先来先进的按时间排序
[mempool.sub.timeline]
# mempool缓存容量大小，默认10240
poolCacheSize=10240
# 最小得交易手续费用，这个没有默认值，必填，一般是100000
minTxFeeRate=100000
# 每个账户在mempool中得最大交易数量，默认100
maxTxNumPerAccount=10000
# score是分数队列模式(分数=常量a*手续费/交易字节数-常量b*时间*定量c，按分数排队，高的优先，常量a，b和定量c可配置)，按分数来排序
[mempool.sub.score]
# mempool缓存容量大小，默认10240
poolCacheSize=10240
# 最小得交易手续费用，这个没有默认值，必填，一般是100000
minTxFeeRate=100000
# 每个账户在mempool中得最大交易数量，默认100
maxTxNumPerAccount=10000
# 时间占价格比例
timeParam=1
# 手续费相对于时间的一个合适的常量，取当前unix时间戳前四位数，排队时手续费高1e-5的分数~=快1s的分数
priceConstant=1544
# 常量比例
pricePower=1
# price是价格队列模式(价格=手续费/交易字节数，价格高者优先，同价则时间早优先)
[mempool.sub.price]
# mempool缓存容量大小，默认10240
poolCacheSize=10240
# 最小得交易手续费用，这个没有默认值，必填，一般是100000
minTxFeeRate=100000
# 每个账户在mempool中得最大交易数量，默认100
maxTxNumPerAccount=10000
[consensus]
#共识名,可选项有solo,ticket,raft,tendermint,para
name="solo"
#是否开启挖矿,开启挖矿才能创建区块
minerstart=true
#创世区块时间(UTC时间)
genesisBlockTime=1514533394
#创世交易地址
genesis="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
[mver.consensus]
#基金账户地址
fundKeyAddr = "1BQXS6TxaYYG5mADaWij4AxhZZUTpw95a5"
#用户回报
coinReward = 18
#发展基金回报
coinDevFund = 12
#ticket价格
ticketPrice = 10000
#挖矿难度
powLimitBits = "0x1f00ffff"
#每次调整难度的最大的范围，如果设置成 4 的话，范围是 (1/4 - 4)，一次不能增加 4倍以上的难度，或者难度减少为 原来的 1/4 ，这个参数，是为了难度不会突然爆增加或者减少
retargetAdjustmentFactor = 4
#表示如果区块时间大于当前时间 16s ，那么区块就会判定为无效区块。
futureBlockTime = 16
#ticket冻结时长
ticketFrozenTime = 5    #5s only for test
ticketWithdrawTime = 10 #10s only for test
ticketMinerWaitTime = 2 #2s only for test
#区块包含最多交易数
maxTxNumber = 1600      #160
#调整挖矿难度的间隔，(ps:难度不是每个区块都调整的，而是每隔 targetTimespan / targetTimePerBlock 块调整一次)
targetTimespan = 2304
#每个区块打包的目标时间
targetTimePerBlock = 16
# 仅保留这一项，其他consensus相关的配置全部删除
[consensus.sub.solo]
#创世交易地址
genesis="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
#创世区块时间(UTC时间)
genesisBlockTime=1514533394
#获取交易间隔时长,单位纳秒
waitTxMs=10
[store]
# 数据存储格式名称，目前支持mavl,kvdb,kvmvcc,mpt
name="mavl"
# 数据存储驱动类别，目前支持leveldb,goleveldb,memdb,gobadgerdb,ssdb,pegasus
driver="leveldb"
# 数据文件存储路径
dbPath="datadir/mavltree"
# Cache大小
dbCache=128
# 数据库版本
localdbVersion="1.0.0"
[store.sub.mavl]
# 是否使能mavl加前缀
enableMavlPrefix=false
# 是否使能MVCC,如果mavl中enableMVCC为true此处必须为true
enableMVCC=false
# 是否使能mavl数据裁剪
enableMavlPrune=false
# 裁剪高度间隔
pruneHeight=10000
[wallet]
# 交易发送最低手续费，单位0.00000001BTY(1e-8),默认100000，即0.001BTY
minFee=100000
# walletdb驱动名，支持leveldb/memdb/gobadgerdb/ssdb/pegasus
driver="leveldb"
# walletdb路径
dbPath="wallet"
# walletdb缓存大小
dbCache=16
# 钱包发送交易签名方式
signType="secp256k1"
[wallet.sub.ticket]
# 是否关闭ticket自动挖矿，默认false
minerdisable=false
# 允许购买ticket挖矿的白名单地址，默认配置“*”，允许所有地址购买
minerwhitelist=["*"]
[exec]
#是否开启stat插件
enableStat=false
#是否开启MVCC插件
enableMVCC=false
alias=["token1:token","token2:token","token3:token"]
[exec.sub.token]
#是否保存token交易信息
saveTokenTxList=true
#token审批人地址
tokenApprs = [
    "1Bsg9j6gW83sShoee1fZAt9TkUjcrCgA9S",
    "1Q8hGLfoGe63efeWa8fJ4Pnukhkngt6poK",
    "1LY8GFia5EiyoTodMLfkB5PHNNpXRqxhyB",
    "1GCzJDS6HbgTQ2emade7mEJGGWFfA15pS9",
    "1JYB8sxi4He5pZWHCd3Zi2nypQ4JMB6AxN",
    "12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv",
]
[exec.sub.cert]
# 是否启用证书验证和签名
enable=false
# 加密文件路径
cryptoPath="authdir/crypto"
# 带证书签名类型，支持"auth_ecdsa", "auth_sm2"
signType="auth_ecdsa"
[exec.sub.relay]
#relay执行器保存BTC头执行权限地址
genesis="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
[exec.sub.manage]
#manage执行器超级管理员地址
superManager=[
    "1Bsg9j6gW83sShoee1fZAt9TkUjcrCgA9S", 
    "12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv", 
    "1Q8hGLfoGe63efeWa8fJ4Pnukhkngt6poK"
]
`
)

var (
	random *rand.Rand

	loopCount = 1
	conn      *grpc.ClientConn
	c         types.Chain33Client
	adminPriv = "CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"
	adminAddr = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"

	//userAPubkey = "03EF0E1D3112CF571743A3318125EDE2E52A4EB904BCBAA4B1F75020C2846A7EB4"
	userAAddr = "15LsTP6tkYGZcN7tc1Xo2iYifQfowxot3b"
	userAPriv = "5A6A14DA6F5A42835E529D75D87CC8904544F59EEE5387A37D87EEAD194D7EB2"

	//userBPubkey = "027848E7FA630B759DB406940B5506B666A344B1060794BBF314EB459D40881BB3"
	userBAddr = "14w5JaGDrXTZwF5Wv51UAtuGgAupenLAok"
	userBPriv = "754F53FCEA0CB1F528918726A49B3551B7F1284D802A1D6AAF4522E8A8DA1B5B"
)

const fee = 1e6

func init() {
	err := limits.SetLimits()
	if err != nil {
		panic(err)
	}
	log.SetLogLevel("info")
	random = rand.New(rand.NewSource(types.Now().UnixNano()))

	cr2, err := crypto.New(types.GetSignName("", types.SECP256K1))
	if err != nil {
		fmt.Println("crypto.New failed for types.ED25519")
		return
	}
	secp = cr2
}

func Init() {
	fmt.Println("=======Init Data1!=======")
	os.RemoveAll("datadir")
	os.RemoveAll("wallet")
	os.Remove("chain33.test.toml")

	ioutil.WriteFile("chain33.test.toml", []byte(config), 0664)
}

func clearTestData() {
	fmt.Println("=======start clear test data!=======")

	os.Remove("chain33.test.toml")
	os.RemoveAll("wallet")
	err := os.RemoveAll("datadir")
	if err != nil {
		fmt.Println("delete datadir have a err:", err.Error())
	}

	fmt.Println("test data clear successfully!")
}

func TestGuess(t *testing.T) {
	Init()
	testGuessImp(t)
	fmt.Println("=======start clear test data!=======")
	clearTestData()
}

func testGuessImp(t *testing.T) {
	fmt.Println("=======start guess test!=======")
	q, chain, s, mem, exec, cs, p2p := initEnvGuess()
	cfg := q.GetConfig()
	defer chain.Close()
	defer mem.Close()
	defer exec.Close()
	defer s.Close()
	defer q.Close()
	defer cs.Close()
	defer p2p.Close()
	err := createConn()
	for err != nil {
		err = createConn()
	}
	time.Sleep(2 * time.Second)
	fmt.Println("=======start NormPut!=======")

	for i := 0; i < loopCount; i++ {
		NormPut(cfg)
		time.Sleep(time.Second)
	}

	fmt.Println("=======start sendTransferTx!=======")
	//从创世地址向测试地址A转入代币
	sendTransferTx(cfg, adminPriv, userAAddr, 2000000000000)
	sendTransferTx(cfg, adminPriv, userBAddr, 2000000000000)

	time.Sleep(2 * time.Second)
	in := &types.ReqBalance{}
	in.Addresses = append(in.Addresses, userAAddr)
	acct, err1 := c.GetBalance(context.Background(), in)
	if err1 != nil || len(acct.Acc) == 0 {
		fmt.Println("no balance for ", userAAddr)
	} else {
		fmt.Println(userAAddr, " balance:", acct.Acc[0].Balance, "frozen:", acct.Acc[0].Frozen)
	}
	assert.Equal(t, true, acct.Acc[0].Balance == 2000000000000)

	in2 := &types.ReqBalance{}
	in2.Addresses = append(in.Addresses, userBAddr)
	acct2, err2 := c.GetBalance(context.Background(), in2)
	if err2 != nil || len(acct2.Acc) == 0 {
		fmt.Println("no balance for ", userBAddr)
	} else {
		fmt.Println(userBAddr, " balance:", acct2.Acc[0].Balance, "frozen:", acct2.Acc[0].Frozen)
	}
	assert.Equal(t, true, acct2.Acc[0].Balance == 2000000000000)

	fmt.Println("=======start sendTransferToExecTx!=======")
	//从测试地址向dos合约转入代币
	sendTransferToExecTx(cfg, userAPriv, "guess", 1000000000000)
	sendTransferToExecTx(cfg, userBPriv, "guess", 1000000000000)
	time.Sleep(2 * time.Second)

	fmt.Println("=======start GetBalance!=======")

	in3 := &types.ReqBalance{}
	in3.Addresses = append(in3.Addresses, userAAddr)
	acct3, err3 := c.GetBalance(context.Background(), in3)
	if err3 != nil || len(acct3.Acc) == 0 {
		fmt.Println("no balance for ", userAAddr)
	} else {
		fmt.Println(userAAddr, " balance:", acct3.Acc[0].Balance, "frozen:", acct3.Acc[0].Frozen)
	}
	assert.Equal(t, true, acct3.Acc[0].Balance == 1000000000000)

	in4 := &types.ReqBalance{}
	in4.Addresses = append(in4.Addresses, userBAddr)
	acct4, err4 := c.GetBalance(context.Background(), in4)
	if err4 != nil || len(acct4.Acc) == 0 {
		fmt.Println("no balance for ", userBAddr)
	} else {
		fmt.Println(userBAddr, " balance:", acct4.Acc[0].Balance, "frozen:", acct4.Acc[0].Frozen)
	}
	assert.Equal(t, true, acct4.Acc[0].Balance == 1000000000000)

	fmt.Println("=======start sendGuessStartTx!=======")
	ok, gameid := sendGuessStartTx(cfg, "WorldCup Final", "A:France;B:Claodia", "football", adminPriv)
	if !ok {
		panic("Guess start failed.")
	} else {
		fmt.Println("txid: ", hex.EncodeToString(gameid))
	}
	time.Sleep(2 * time.Second)

	strGameID1 := "0x" + hex.EncodeToString(gameid)

	reply := queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusStart)

	fmt.Println("=======start sendGuessBetTx!=======")
	ok, txid := sendGuessBetTx(cfg, strGameID1, "A", 5e8, userAPriv)
	if !ok {
		panic("Guess A bet failed.")
	} else {
		fmt.Println("Guess A bet txid: ", hex.EncodeToString(txid))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusBet && reply.Games[0].BetStat.TotalBetTimes == 1)

	ok, txid = sendGuessBetTx(cfg, strGameID1, "B", 5e8, userBPriv)
	if !ok {
		panic("Guess B bet failed.")
	} else {
		fmt.Println("Guess B bet txid: ", hex.EncodeToString(txid))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusBet && reply.Games[0].BetStat.TotalBetTimes == 2)

	fmt.Println("=======start sendGuessStopTx failed!=======")
	ok, txid = sendGuessStopTx(cfg, strGameID1, userBPriv)
	if !ok {
		panic("Guess stop failed,only admin can stop.")
	} else {
		fmt.Println("Guess stop txid: ", hex.EncodeToString(txid))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusBet && reply.Games[0].BetStat.TotalBetTimes == 2)

	fmt.Println("=======start sendGuessStopTx!=======")
	ok, txid = sendGuessStopTx(cfg, strGameID1, adminPriv)
	if !ok {
		panic("Guess stop failed.")
	} else {
		fmt.Println("Guess stop txid: ", hex.EncodeToString(txid))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusStopBet && reply.Games[0].BetStat.TotalBetTimes == 2)

	fmt.Println("=======start sendGuessBetTx failed!=======")
	ok, txid = sendGuessBetTx(cfg, strGameID1, "A", 5e8, userAPriv)
	if !ok {
		fmt.Println("Guess stopped, bet failed.")
	} else {
		fmt.Println(fmt.Sprintf("Guess A bet txid: %s", hex.EncodeToString(txid)))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusStopBet && reply.Games[0].BetStat.TotalBetTimes == 2)

	fmt.Println("=======start sendGuessPublishTx failed!=======")
	ok, txid = sendGuessPublishTx(cfg, strGameID1, "A", userAPriv)
	if !ok {
		fmt.Println("sendGuessPublishTx failed,only admin can publish.")
	} else {
		fmt.Println(fmt.Sprintf("publish ok, but it's not correct, txid: %s", hex.EncodeToString(txid)))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusStopBet && reply.Games[0].BetStat.TotalBetTimes == 2)

	fmt.Println("=======start sendGuessPublishTx!=======")
	ok, txid = sendGuessPublishTx(cfg, strGameID1, "A", adminPriv)
	if !ok {
		fmt.Println("sendGuessPublishTx failed.")
	} else {
		fmt.Println(fmt.Sprintf("publish ok, txid: %s", hex.EncodeToString(txid)))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusPublish && reply.Games[0].BetStat.TotalBetTimes == 2)

	fmt.Println("=======start sendGuessAbortTx!=======")
	ok, txid = sendGuessAbortTx(cfg, strGameID1, adminPriv)
	if !ok {
		fmt.Println("Guess abort failed, already published.")
	} else {
		fmt.Println(fmt.Sprintf("Guess abort txid: %s", hex.EncodeToString(txid)))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID1)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusPublish && reply.Games[0].BetStat.TotalBetTimes == 2)

	//再来一次，测试异常流程:start->abort->stop
	fmt.Println("=======start sendGuessStartTx!=======")
	ok, gameid = sendGuessStartTx(cfg, "WorldCup Final", "A:France;B:Claodia", "football", adminPriv)
	if !ok {
		fmt.Println("Guess start failed.")
	} else {
		fmt.Println("txid: ", hex.EncodeToString(gameid))
	}
	time.Sleep(2 * time.Second)

	strGameID2 := "0x" + hex.EncodeToString(gameid)

	reply = queryGuessByIds(strGameID2)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusStart)

	fmt.Println("=======start sendGuessAbortTx!=======")
	ok, txid = sendGuessAbortTx(cfg, strGameID2, adminPriv)
	if !ok {
		fmt.Println("Guess abort failed.")
	} else {
		fmt.Println(fmt.Sprintf("Guess abort txid: %s", hex.EncodeToString(txid)))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID2)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusAbort)

	fmt.Println("=======start sendGuessStopTx failed!=======")
	ok, txid = sendGuessStopTx(cfg, strGameID2, adminPriv)
	if !ok {
		fmt.Println("Guess stop failed,it's already aborted.")
	} else {
		fmt.Println(fmt.Sprintf("Guess stop txid: %s", hex.EncodeToString(txid)))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID2)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusAbort)

	//再来一次，测试流程:start->stop->abort
	fmt.Println("=======start sendGuessStartTx!=======")
	ok, gameid = sendGuessStartTx(cfg, "WorldCup Final", "A:France;B:Claodia", "football", adminPriv)
	if !ok {
		fmt.Println("Guess start failed.")
	} else {
		fmt.Println("txid: ", hex.EncodeToString(gameid))
	}
	time.Sleep(2 * time.Second)

	strGameID3 := "0x" + hex.EncodeToString(gameid)
	reply = queryGuessByIds(strGameID3)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusStart)

	fmt.Println("=======start sendGuessStopTx!=======")
	ok, txid = sendGuessStopTx(cfg, strGameID3, adminPriv)
	if !ok {
		fmt.Println("Guess stop failed")
	} else {
		fmt.Println("Guess stop txid: ", hex.EncodeToString(txid))
	}
	time.Sleep(2 * time.Second)
	reply = queryGuessByIds(strGameID3)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusStopBet)

	fmt.Println("=======start sendGuessAbortTx!=======")
	ok, txid = sendGuessAbortTx(cfg, strGameID3, adminPriv)
	if !ok {
		fmt.Println("Guess abort failed.")
	} else {
		fmt.Println(fmt.Sprintf("Guess abort txid: %s", hex.EncodeToString(txid)))
	}
	time.Sleep(2 * time.Second)

	//以下测试查询接口
	fmt.Println("=======start queryGuessByIds!=======")
	reply = queryGuessByIds(strGameID3)
	assert.Equal(t, true, reply.Games[0].Status == gty.GuessGameStatusAbort)

	fmt.Println("=======start queryGuessByID!=======")
	reply2 := queryGuessByID(strGameID1)
	assert.Equal(t, true, reply2.Game.Status == gty.GuessGameStatusPublish && reply2.Game.BetStat.TotalBetTimes == 2)

	fmt.Println("=======start queryGuessByAddr!=======")
	record := queryGuessByAddr(userAAddr)
	assert.Equal(t, true, record.Records[0].GameID == strGameID1)

	fmt.Println("=======start queryGuessByStatus!=======")

	record = queryGuessByStatus(gty.GuessGameStatusPublish)
	assert.Equal(t, true, record.Records[0].GameID == strGameID1)

	record = queryGuessByStatus(gty.GuessGameStatusAbort)
	assert.Equal(t, true, len(record.Records) == 2)

	fmt.Println("=======start queryGuessByAdminAddr!=======")
	record = queryGuessByAdminAddr(adminAddr)
	assert.Equal(t, true, len(record.Records) == 3)

	fmt.Println("=======start queryGuessByAddrStatus!=======")
	record = queryGuessByAddrStatus(userBAddr, gty.GuessGameStatusPublish)
	assert.Equal(t, true, len(record.Records) == 1 && record.Records[0].GameID == strGameID1)
	record = queryGuessByAddrStatus(userBAddr, 10)
	assert.Equal(t, true, len(record.Records) == 0)

	fmt.Println("=======start queryGuessByAdminAddrStatus!=======")
	record = queryGuessByAdminAddrStatus(adminAddr, gty.GuessGameStatusAbort)
	assert.Equal(t, true, len(record.Records) == 2)

	record = queryGuessByAdminAddrStatus(adminAddr, gty.GuessGameStatusPublish)
	assert.Equal(t, true, len(record.Records) == 1)

	fmt.Println("=======start queryGuessByCategoryStatus!=======")
	record = queryGuessByCategoryStatus("football", gty.GuessGameStatusPublish)
	assert.Equal(t, true, len(record.Records) == 1)

	record = queryGuessByCategoryStatus("football", gty.GuessGameStatusAbort)
	assert.Equal(t, true, len(record.Records) == 2)

	time.Sleep(2 * time.Second)
}

func initEnvGuess() (queue.Queue, *blockchain.BlockChain, queue.Module, queue.Module, *executor.Executor, queue.Module, queue.Module) {
	flag.Parse()
	chain33Cfg := types.NewChain33Config(types.ReadFile("chain33.test.toml"))
	var q = queue.New("channel")
	q.SetConfig(chain33Cfg)
	cfg := chain33Cfg.GetModuleConfig()
	cfg.Log.LogFile = ""
	sub := chain33Cfg.GetSubConfig()
	chain := blockchain.New(chain33Cfg)
	chain.SetQueueClient(q.Client())

	exec := executor.New(chain33Cfg)
	exec.SetQueueClient(q.Client())
	chain33Cfg.SetMinFee(0)
	s := store.New(chain33Cfg)
	s.SetQueueClient(q.Client())

	cs := solo.New(cfg.Consensus, sub.Consensus["solo"])
	cs.SetQueueClient(q.Client())

	mem := mempool.New(chain33Cfg)
	mem.SetQueueClient(q.Client())
	network := p2p.NewP2PMgr(chain33Cfg)

	network.SetQueueClient(q.Client())

	rpc.InitCfg(cfg.RPC)
	gapi := rpc.NewGRpcServer(q.Client(), nil)
	go gapi.Listen()

	japi := rpc.NewJSONRPCServer(q.Client(), nil)
	go japi.Listen()

	return q, chain, s, mem, exec, cs, network
}

func createConn() error {
	var err error
	url := "127.0.0.1:9802"
	fmt.Println("grpc url:", url)
	conn, err = grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	c = types.NewChain33Client(conn)
	return nil
}

func generateKey(i, valI int) string {
	key := make([]byte, valI)
	binary.PutUvarint(key[:10], uint64(valI))
	binary.PutUvarint(key[12:24], uint64(i))
	if _, err := rand.Read(key[24:]); err != nil {
		os.Exit(1)
	}
	return string(key)
}

func generateValue(i, valI int) string {
	value := make([]byte, valI)
	binary.PutUvarint(value[:16], uint64(i))
	binary.PutUvarint(value[32:128], uint64(i))
	if _, err := rand.Read(value[128:]); err != nil {
		os.Exit(1)
	}
	return string(value)
}

func getprivkey(key string) crypto.PrivKey {
	bkey, err := hex.DecodeString(key)
	if err != nil {
		panic(err)
	}
	priv, err := secp.PrivKeyFromBytes(bkey)
	if err != nil {
		panic(err)
	}
	return priv
}

func prepareTxList(cfg *types.Chain33Config) *types.Transaction {
	var key string
	var value string
	var i int

	key = generateKey(i, 32)
	value = generateValue(i, 180)

	nput := &pty.NormAction_Nput{Nput: &pty.NormPut{Key: []byte(key), Value: []byte(value)}}
	action := &pty.NormAction{Value: nput, Ty: pty.NormActionPut}
	tx := &types.Transaction{Execer: []byte("norm"), Payload: types.Encode(action), Fee: fee}
	tx.To = address.ExecAddress("norm")
	tx.Nonce = random.Int63()
	tx.ChainID = cfg.GetChainID()
	tx.Sign(types.SECP256K1, getprivkey("CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"))
	return tx
}

func NormPut(cfg *types.Chain33Config) {
	tx := prepareTxList(cfg)

	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		return
	}
}

func sendTransferTx(cfg *types.Chain33Config, fromKey, to string, amount int64) bool {
	signer := util.HexToPrivkey(fromKey)
	var tx *types.Transaction
	transfer := &cty.CoinsAction{}
	v := &cty.CoinsAction_Transfer{Transfer: &types.AssetsTransfer{Amount: amount, Note: []byte(""), To: to}}
	transfer.Value = v
	transfer.Ty = cty.CoinsActionTransfer
	execer := []byte(cfg.GetCoinExec())
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(transfer), To: to, Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("in sendTransferTx formatTx failed")
		return false
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendTransferTx SendTransaction failed")

		return false
	}
	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendTransferTx SendTransaction failed,reply not ok.")

		return false
	}
	fmt.Println("sendTransferTx ok")

	return true
}

func sendTransferToExecTx(cfg *types.Chain33Config, fromKey, execName string, amount int64) bool {
	signer := util.HexToPrivkey(fromKey)
	var tx *types.Transaction
	transfer := &cty.CoinsAction{}
	execAddr := address.ExecAddress(execName)
	v := &cty.CoinsAction_TransferToExec{TransferToExec: &types.AssetsTransferToExec{Amount: amount, Note: []byte(""), ExecName: execName, To: execAddr}}
	transfer.Value = v
	transfer.Ty = cty.CoinsActionTransferToExec
	execer := []byte(cfg.GetCoinExec())
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(transfer), To: address.ExecAddress("guess"), Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("sendTransferToExecTx formatTx failed.")

		return false
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendTransferToExecTx SendTransaction failed")

		return false
	}
	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendTransferToExecTx SendTransaction failed,reply not ok.")

		return false
	}

	fmt.Println("sendTransferToExecTx ok")

	return true
}

func sendGuessStartTx(cfg *types.Chain33Config, topic, option, category, privKey string) (bool, []byte) {
	signer := util.HexToPrivkey(privKey)
	var tx *types.Transaction
	action := &gty.GuessGameAction{}

	v := &gty.GuessGameAction_Start{
		Start: &gty.GuessGameStart{
			Topic:          topic,
			Options:        option,
			Category:       category,
			MaxBetsOneTime: 100e8,
			MaxBetsNumber:  1000e8,
			DevFeeFactor:   5,
			DevFeeAddr:     "1D6RFZNp2rh6QdbcZ1d7RWuBUz61We6SD7",
			PlatFeeFactor:  5,
			PlatFeeAddr:    "1PHtChNt3UcfssR7v7trKSk3WJtAWjKjjX",
		},
	}

	action.Value = v
	action.Ty = gty.GuessGameActionStart
	execer := []byte("guess")
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(action), To: address.ExecAddress(string(execer)), Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("sendGuessStartTx formatTx failed.")

		return false, nil
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendGuessStartTx SendTransaction failed")

		return false, nil
	}

	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendGuessStartTx SendTransaction failed,reply not ok.")

		return false, nil
	}

	fmt.Println("sendGuessStartTx ok")

	return true, reply.Msg
}

func sendGuessBetTx(cfg *types.Chain33Config, gameID, option string, betsNum int64, privKey string) (bool, []byte) {
	signer := util.HexToPrivkey(privKey)
	var tx *types.Transaction
	action := &gty.GuessGameAction{}

	v := &gty.GuessGameAction_Bet{
		Bet: &gty.GuessGameBet{
			GameID:  gameID,
			Option:  option,
			BetsNum: betsNum,
		},
	}

	action.Value = v
	action.Ty = gty.GuessGameActionBet
	execer := []byte("guess")
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(action), To: address.ExecAddress(string(execer)), Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("sendGuessBetTx formatTx failed.")

		return false, nil
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendGuessBetTx SendTransaction failed")

		return false, nil
	}

	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendGuessBetTx SendTransaction failed,reply not ok.")

		return false, nil
	}

	fmt.Println("sendGuessBetTx ok")

	return true, reply.Msg
}

func sendGuessStopTx(cfg *types.Chain33Config, gameID, privKey string) (bool, []byte) {
	signer := util.HexToPrivkey(privKey)
	var tx *types.Transaction
	action := &gty.GuessGameAction{}

	v := &gty.GuessGameAction_StopBet{
		StopBet: &gty.GuessGameStopBet{
			GameID: gameID,
		},
	}

	action.Value = v
	action.Ty = gty.GuessGameActionStopBet
	execer := []byte("guess")
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(action), To: address.ExecAddress(string(execer)), Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("sendGuessStopTx formatTx failed.")

		return false, nil
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendGuessStopTx SendTransaction failed")

		return false, nil
	}

	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendGuessStopTx SendTransaction failed,reply not ok.")

		return false, nil
	}

	fmt.Println("sendGuessStopTx ok")

	return true, reply.Msg
}

func sendGuessAbortTx(cfg *types.Chain33Config, gameID, privKey string) (bool, []byte) {
	signer := util.HexToPrivkey(privKey)
	var tx *types.Transaction
	action := &gty.GuessGameAction{}

	v := &gty.GuessGameAction_Abort{
		Abort: &gty.GuessGameAbort{
			GameID: gameID,
		},
	}

	action.Value = v
	action.Ty = gty.GuessGameActionAbort
	execer := []byte("guess")
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(action), To: address.ExecAddress(string(execer)), Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("sendGuessAbortTx formatTx failed.")

		return false, nil
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendGuessAbortTx SendTransaction failed")

		return false, nil
	}

	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendGuessAbortTx SendTransaction failed,reply not ok.")

		return false, nil
	}

	fmt.Println("sendGuessAbortTx ok")

	return true, reply.Msg
}

func sendGuessPublishTx(cfg *types.Chain33Config, gameID, result, privKey string) (bool, []byte) {
	signer := util.HexToPrivkey(privKey)
	var tx *types.Transaction
	action := &gty.GuessGameAction{}

	v := &gty.GuessGameAction_Publish{
		Publish: &gty.GuessGamePublish{
			GameID: gameID,
			Result: result,
		},
	}

	action.Value = v
	action.Ty = gty.GuessGameActionPublish
	execer := []byte("guess")
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(action), To: address.ExecAddress(string(execer)), Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("sendGuessPublishTx formatTx failed.")

		return false, nil
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendGuessPublishTx SendTransaction failed")

		return false, nil
	}

	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendGuessPublishTx SendTransaction failed,reply not ok.")

		return false, nil
	}

	fmt.Println("sendGuessPublishTx ok")

	return true, reply.Msg
}

func queryGuessByIds(gameIDs string) *gty.ReplyGuessGameInfos {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	gameIds := strings.Split(gameIDs, ";")
	req := &gty.QueryGuessGameInfos{
		GameIDs: gameIds,
	}
	params.FuncName = gty.FuncNameQueryGamesByIDs
	params.Payload = types.MustPBToJSON(req)
	var res gty.ReplyGuessGameInfos
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "Chain33.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "Chain33.Query", params, &res)
	ctx.Run()
	return &res
}

func queryGuessByID(gameID string) *gty.ReplyGuessGameInfo {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	req := &gty.QueryGuessGameInfo{
		GameID: gameID,
	}
	params.FuncName = gty.FuncNameQueryGameByID
	params.Payload = types.MustPBToJSON(req)
	var res gty.ReplyGuessGameInfo
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "Chain33.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "Chain33.Query", params, &res)
	ctx.Run()
	return &res
}

func queryGuessByAddr(addr string) *gty.GuessGameRecords {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	req := &gty.QueryGuessGameInfo{
		Addr: addr,
	}
	params.FuncName = gty.FuncNameQueryGameByAddr
	params.Payload = types.MustPBToJSON(req)
	var res gty.GuessGameRecords
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "Chain33.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "Chain33.Query", params, &res)
	ctx.Run()
	return &res
}

func queryGuessByStatus(status int32) *gty.GuessGameRecords {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	req := &gty.QueryGuessGameInfo{
		Status: status,
	}
	params.FuncName = gty.FuncNameQueryGameByStatus
	params.Payload = types.MustPBToJSON(req)
	var res gty.GuessGameRecords
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "Chain33.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "Chain33.Query", params, &res)
	ctx.Run()
	return &res
}

func queryGuessByAdminAddr(addr string) *gty.GuessGameRecords {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	req := &gty.QueryGuessGameInfo{
		AdminAddr: addr,
	}
	params.FuncName = gty.FuncNameQueryGameByAdminAddr
	params.Payload = types.MustPBToJSON(req)
	var res gty.GuessGameRecords
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "Chain33.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "Chain33.Query", params, &res)
	ctx.Run()
	return &res
}

func queryGuessByAddrStatus(addr string, status int32) *gty.GuessGameRecords {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	req := &gty.QueryGuessGameInfo{
		Addr:   addr,
		Status: status,
	}
	params.FuncName = gty.FuncNameQueryGameByAddrStatus
	params.Payload = types.MustPBToJSON(req)
	var res gty.GuessGameRecords
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "Chain33.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "Chain33.Query", params, &res)
	ctx.Run()
	return &res
}

func queryGuessByAdminAddrStatus(addr string, status int32) *gty.GuessGameRecords {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	req := &gty.QueryGuessGameInfo{
		AdminAddr: addr,
		Status:    status,
	}
	params.FuncName = gty.FuncNameQueryGameByAdminStatus
	params.Payload = types.MustPBToJSON(req)
	var res gty.GuessGameRecords
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "Chain33.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "Chain33.Query", params, &res)
	ctx.Run()
	return &res
}

func queryGuessByCategoryStatus(category string, status int32) *gty.GuessGameRecords {
	var params rpctypes.Query4Jrpc
	params.Execer = gty.GuessX

	req := &gty.QueryGuessGameInfo{
		Category: category,
		Status:   status,
	}
	params.FuncName = gty.FuncNameQueryGameByCategoryStatus
	params.Payload = types.MustPBToJSON(req)
	var res gty.GuessGameRecords
	//ctx := jsonrpc.NewRPCCtx("http://"+types.Conf("config.rpc").GStr("jrpcBindAddr"), "Chain33.Query", params, &res)
	ctx := jsonrpc.NewRPCCtx("http://127.0.0.1:9801", "Chain33.Query", params, &res)
	ctx.Run()
	return &res
}
