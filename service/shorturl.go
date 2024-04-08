package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"byteurl/core/cache"
	"byteurl/dal/dao"
	"byteurl/dal/query"
	pb "byteurl/pb/api/leaf/v1"
	"byteurl/pkg/connect"
	"byteurl/pkg/encode"
	"byteurl/pkg/urltool"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
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
	defer conn.Close()

	cli = pb.NewLeafSegmentServiceClient(conn)
	// ctx := context.TODO()
	// resp, err := cli.GetSegmentId(ctx, &v1.GenIdsRequest{BizTag: "bbbb", Count: 1})
	// if err != nil {
	// 	fmt.Printf("c.GetSegmentId failed, err:%v\n", err)
	// 	return
	// }
	// fmt.Printf("resp:%v\n", resp.GetIds())
}

func GenShortURL(ctx context.Context, origin string) (shorturl string, err error) {
	// 空数据
	if len(origin) == 0 {
		return "", errors.New("invalid originURL")
	}
	// 不通
	if ok := connect.Get(origin); !ok {
		return "", errors.New("invalid originURL")
	}
	// 重复
	md5Value := encode.Sum([]byte(origin))
	_, err = query.Short.WithContext(ctx).Where(query.Short.Md5.Eq(md5Value)).First()
	if err != dao.ErrNotFound {
		if err == nil {
			return "", errors.New("duplicate originURL")
		}
		return "", err
	}
	// 循环
	basePath, err := urltool.GetBasePath(origin)
	if err != nil {
		return "", fmt.Errorf("urltool.GetBasePath failed, err: %v", err)
	}
	_, err = query.Short.WithContext(ctx).Where(query.Short.Surl.Eq(basePath)).First()
	if err != dao.ErrNotFound {
		if err == nil {
			return "", errors.New("duplicate shortURL")
		}
		return "", err
	}

	var short string
	for {
		ctx1, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// TODO: 取消硬编码
		// rpc取号
		resp, err := cli.GetSegmentId(ctx1, &pb.GenIdsRequest{BizTag: "bbbb", Count: 1})
		if err != nil {
			return "", err
		}
		short = encode.Int2String(uint64(resp.GetIds()[0]))
		//TODO: 黑名单

		break
	}

	fmt.Printf("short: %v\n", short)

	// 存储
	memcache.Set(short, origin)

	return "lsh7d1.com/" + short, nil
}
