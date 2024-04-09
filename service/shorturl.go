package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"byteurl/core/cache"
	"byteurl/dal/dao"
	"byteurl/dal/model"
	"byteurl/mq"
	pb "byteurl/pb/api/leaf/v1"
	"byteurl/pkg/connect"
	"byteurl/pkg/encode"
	"byteurl/pkg/urltool"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrExistData = errors.New("existing data")

	cli      pb.LeafSegmentServiceClient
	_        = cli
	memcache = cache.NewCache("byteurl", time.Hour, cache.WithAroundCapLimit(1e6))
)

func init() {
	conn, err := grpc.NewClient(
		// consul服务
		"consul://localhost:8500/segment?healthy=true",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("grpc.NewClient failed, err: %v", err)
	}

	cli = pb.NewLeafSegmentServiceClient(conn)
}

func GenShortURL(ctx context.Context, longurl string) (shorturl string, err error) {
	// 空数据 TODO: 路由层做
	if len(longurl) == 0 {
		return "", errors.New("invalid longurl")
	}
	// 长链不通
	if ok := connect.Get(longurl); !ok {
		return "", errors.New("invalid longurl")
	}

	// 重复创建
	md5Value := encode.Sum([]byte(longurl))
	// 1.查内存缓存
	if _, ok := memcache.Get(md5Value); ok {
		return "", ErrExistData
	}
	// 2.查缓存
	err = dao.Get(ctx, md5Value)
	if err != dao.RedisErrNotFound {
		return "", ErrExistData
	}
	// 3.查数据库
	if err = dao.PeekMd5NotExist(ctx, md5Value); err != nil {
		return "", err
	}

	// 循环创建
	basePath, err := urltool.GetBasePath(longurl)
	if err != nil {
		return "", fmt.Errorf("urltool.GetBasePath failed, err: %v", err)
	}
	if err = dao.PeekShortURLNotExist(ctx, basePath); err != nil {
		return "", err
	}

	var short string
	for {
		// TODO: 取消硬编码
		// rpc取号
		resp, err := cli.GetSegmentId(ctx, &pb.GenIdsRequest{BizTag: "bbbb", Count: 1})
		if err != nil {
			return "", err
		}
		short = encode.Int2String(uint64(resp.GetIds()[0]))
		//TODO: 黑名单

		break
	}

	fmt.Printf("short: %v\n", short)

	// 存储
	s := &model.Short{
		Lurl:     longurl,
		Md5:      md5Value,
		Surl:     short,
		CreateBy: "lsh7d1",
	}
	// 1.写缓存
	if err := dao.Set(ctx, []string{short}, []string{longurl}); err != nil {
		return "", err
	}
	// 2.写mq
	shortJson, _ := json.Marshal(s)
	mq.WriteMessage(ctx, []mq.MqEntity{{Key: []byte(short), Value: []byte(shortJson)}})
	// 3.写内存缓存
	memcache.Set(short, longurl)

	return "lsh7d1.com/" + short, nil
}
