package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/apabramov/anti-bruteforce/internal/server/pb"
)

var (
	host     string
	port     string
	typelist string
)

var errlistNotExist = errors.New("list not exists")

const (
	black = "black"
	white = "white"
)

func main() {
	rootCmd := &cobra.Command{
		Use:  "cli",
		Long: "CLI interface for anti-bruteforce service",
	}
	rootCmd.PersistentFlags().StringVar(&host, "host", "localhost", "gRPC server host ")
	rootCmd.PersistentFlags().StringVar(&port, "port", "12000", "gRPC server port")

	rootCmd.AddCommand(AddCmd(), DeleteCmd(), ResetCmd())

	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}

func AddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "add 192.168.1.1",
		Long: "Add to black list or white lists",
		Args: cobra.ExactArgs(1),
		RunE: Add,
	}
	cmd.Flags().StringVar(&typelist, "type_list", black, "select type list")

	return cmd
}

func DeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "del 192.168.1.1",
		Long: "Delete from black list or white lists",
		Args: cobra.ExactArgs(1),
		RunE: Del,
	}
	cmd.Flags().StringVar(&typelist, "type_list", black, "select type list")

	return cmd
}

func ResetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "reset login password IP",
		Long: "reset login password ip from bucket",
		Args: cobra.ExactArgs(3),
		RunE: Reset,
	}

	return cmd
}

func Add(cmd *cobra.Command, args []string) error {
	if typelist != black && typelist != white {
		return errlistNotExist
	}
	subnet := args[0]

	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewEventServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.SubnetRequest{Subnet: subnet}
	res := &pb.ResultResponse{}

	switch typelist {
	case black:
		res, err = client.AddBlackList(ctx, req)
	case white:
		res, err = client.AddWhiteList(ctx, req)
	}

	if err != nil {
		return err
	}

	if res.GetError() != "" {
		fmt.Printf("error add subnet %s err: %s", subnet, res.GetError())
		return errors.New(res.GetError())
	}
	fmt.Printf("subnet %s success add", subnet)

	return nil
}

func Del(cmd *cobra.Command, args []string) error {
	if typelist != black && typelist != white {
		return errlistNotExist
	}
	subnet := args[0]

	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewEventServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.SubnetRequest{Subnet: subnet}
	res := &pb.ResultResponse{}

	switch typelist {
	case black:
		res, err = client.DeleteBlackList(ctx, req)
	case white:
		res, err = client.DeleteWhiteList(ctx, req)
	}

	if err != nil {
		return err
	}

	if res.GetError() != "" {
		fmt.Printf("error delete subnet %s err:  %s", subnet, res.GetError())
		return errors.New(res.GetError())
	}
	fmt.Printf("subnet %s success delete", subnet)

	return nil
}

func Reset(cmd *cobra.Command, args []string) error {
	login := args[0]
	pass := args[1]
	ip := args[2]

	conn, err := getGRPCClient()
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewEventServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.Reset(ctx, &pb.AuthRequest{Login: login, Password: pass, Ip: ip})

	if err != nil {
		return err
	}

	if res.GetError() != "" {
		fmt.Printf("error reset bucket login: %s, password: %s, IP: %s. Error: %s", login, pass, ip, res.GetError())
		return errors.New(res.GetError())
	}
	fmt.Printf("Reset bucket %s %s %s", login, pass, ip)

	return nil
}

func getGRPCClient() (*grpc.ClientConn, error) {
	clientOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial(net.JoinHostPort(host, port), clientOptions...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
