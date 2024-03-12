package upgradejob

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"tkestack.io/tke/pkg/util/log"

	"git.woa.com/kmetis/healthcheckpro/pb"
)

var components = []string{"apiserver", "etcd"}

func initHealthCheckerClient() pb.HealthCheckerClient {
	hcAddr := os.Getenv("HEALTH_CHECKER_ADDR")
	if hcAddr != "" {
		conn, err := grpc.Dial(hcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return pb.NewHealthCheckerClient(conn)
	}
	return nil
}

func checkHealth(client pb.HealthCheckerClient, region, cls, app string, nodes []string) error {
	if client == nil {
		// skip check
		return nil
	}

	if err := checkClusterHealth(client, region, cls); err != nil {
		return err
	}

	if err := checkClusterAppHealth(client, region, cls, app); err != nil {
		return err
	}

	/*  // TODO: 等联调完成后再放开
	if err := checkClusterNodes(client, region, cls, nodes); err != nil {
		return err
	}
	*/
	return nil
}

func checkClusterHealth(client pb.HealthCheckerClient, region, cls string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check basic components
	for _, component := range components {
		req := pb.ComponentHealthyRequest{Product: "tke", Region: region, ClusterId: cls, ComponentName: component}
		if rsp, err := client.IsComponentHealthy(ctx, &req); err != nil {
			return err
		} else {
			if err := rsp.GetErrMessage(); err != "" {
				return fmt.Errorf(err)
			}
		}
	}

	return nil
}

func checkClusterAppHealth(client pb.HealthCheckerClient, region, cls, app string) error {
	return nil
}

func checkClusterNodes(client pb.HealthCheckerClient, region, cls string, nodes []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for _, node := range nodes {
		// TODO: region一定要的吗？,node 可以用nodeName吗，它不一定是ip，可以同时查多个节点吗？
		req := pb.NodeHealthyRequest{Product: "tke", Region: region, ClusterId: cls, NodeIp: node}
		if rsp, err := client.IsNodeHealthy(ctx, &req); err != nil {
			return err
		} else {
			if err := rsp.GetErrMessage(); err != "" {
				return fmt.Errorf(err)
			}
		}
	}
	return nil
}
