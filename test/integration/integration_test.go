// nolint
package integration_test

import (
	"context"
	faker "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"

	"github.com/apabramov/anti-bruteforce/internal/server/pb"
)

type TestSuite struct {
	suite.Suite
	client pb.EventServiceClient
}

func (s *TestSuite) SetupSuite() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial("127.0.0.1:12000", opts...)
	s.Require().NoError(err)
	s.client = pb.NewEventServiceClient(conn)
}

func (s *TestSuite) TestAuthIPInvalid() {
	_, err := s.client.Auth(context.Background(), &pb.AuthRequest{
		Login:    faker.Username(),
		Password: faker.Password(),
		Ip:       faker.Name(),
	})
	s.Require().Error(err)
}

func (s *TestSuite) TestIPBlackList() {
	login := faker.Username()
	pass := faker.Password()
	subnet := "192.168.1.0/24"
	s.AddBlackList(subnet)
	ok := s.Auth(login, pass, "192.168.1.1")
	s.Require().False(ok)

	s.DeleteBlackList(subnet)

	ok = s.Auth(login, pass, "192.168.1.1")
	s.Require().True(ok)
}

func (s *TestSuite) TestAddBlackListInvalid() {
	_, err := s.client.AddBlackList(context.Background(), &pb.SubnetRequest{
		Subnet: faker.Word(),
	})
	s.Require().Error(err)
}

func (s *TestSuite) TestDeleteBlackListInvalid() {
	_, err := s.client.DeleteBlackList(context.Background(), &pb.SubnetRequest{
		Subnet: faker.Word(),
	})
	s.Require().Error(err)
}

func (s *TestSuite) TestIPWhiteList() {
	var ok bool
	login := faker.Username()
	pass := faker.Password()
	ip := "10.10.1.1"
	for i := 0; i < 10; i++ {
		ok = s.Auth(login, pass, ip)
		s.Require().True(ok)
	}
	ok = s.Auth(login, pass, ip)
	s.Require().False(ok)

	subnet := "10.10.1.0/24"
	s.AddWhiteList(subnet)

	ok = s.Auth(login, pass, "10.10.1.1")
	s.Require().True(ok)

	s.DeleteWhiteList(subnet)
}

func (s *TestSuite) TestAddWhiteListInvalid() {
	_, err := s.client.AddWhiteList(context.Background(), &pb.SubnetRequest{
		Subnet: faker.Word(),
	})
	s.Require().Error(err)
}

func (s *TestSuite) TestDeleteWhiteListInvalid() {
	_, err := s.client.DeleteBlackList(context.Background(), &pb.SubnetRequest{
		Subnet: faker.Word(),
	})
	s.Require().Error(err)
}

func (s *TestSuite) TestLoginLimit() {
	var ok bool
	login := faker.Username()
	pass := faker.Password()
	ip := faker.IPv4()
	for i := 0; i < 10; i++ {
		ok = s.Auth(login, pass, ip)
		s.Require().True(ok)
	}
	ok = s.Auth(login, pass, ip)
	s.Require().False(ok)
}

func (s *TestSuite) TestPasswordLimit() {
	var ok bool
	pass := faker.Password()
	ip := faker.IPv4()
	for i := 0; i < 100; i++ {
		ok = s.Auth(faker.Username(), pass, ip)
		s.Require().True(ok)
	}
	ok = s.Auth(faker.Username(), pass, ip)
	s.Require().False(ok)
}

func (s *TestSuite) TestIPLimit() {
	var ok bool
	ip := faker.IPv4()
	for i := 0; i < 1000; i++ {
		ok = s.Auth(faker.Username(), faker.Password(), ip)
		s.Require().True(ok)
	}
	ok = s.Auth(faker.Username(), faker.Password(), ip)
	s.Require().False(ok)
}

func (s *TestSuite) TestReset() {
	var ok bool
	login := faker.Username()
	pass := faker.Password()
	ip := faker.IPv4()
	for i := 0; i < 10; i++ {
		ok = s.Auth(login, pass, ip)
		s.Require().True(ok)
	}
	ok = s.Auth(login, pass, ip)
	s.Require().False(ok)

	s.Reset(login, pass, ip)

	ok = s.Auth(login, pass, ip)
	s.Require().True(ok)
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) AddWhiteList(subnet string) {
	s.T().Helper()
	res, err := s.client.AddWhiteList(context.Background(), &pb.SubnetRequest{Subnet: subnet})
	s.Require().NoError(err)
	s.Require().NotNil(res)
}

func (s *TestSuite) AddBlackList(subnet string) {
	s.T().Helper()
	res, err := s.client.AddBlackList(context.Background(), &pb.SubnetRequest{Subnet: subnet})
	s.Require().NoError(err)
	s.Require().NotNil(res)
}

func (s *TestSuite) DeleteWhiteList(subnet string) {
	s.T().Helper()
	res, err := s.client.DeleteWhiteList(context.Background(), &pb.SubnetRequest{Subnet: subnet})
	s.Require().NoError(err)
	s.Require().NotNil(res)
}

func (s *TestSuite) DeleteBlackList(subnet string) {
	s.T().Helper()
	res, err := s.client.DeleteBlackList(context.Background(), &pb.SubnetRequest{Subnet: subnet})
	s.Require().NoError(err)
	s.Require().NotNil(res)
}

func (s *TestSuite) Auth(login, password, ip string) bool {
	s.T().Helper()
	res, err := s.client.Auth(context.Background(), &pb.AuthRequest{Login: login, Password: password, Ip: ip})
	s.Require().NoError(err, ip)
	s.Require().NotNil(res)
	return res.GetResult()
}

func (s *TestSuite) Reset(login, password, ip string) {
	s.T().Helper()
	res, err := s.client.Reset(context.Background(), &pb.AuthRequest{Login: login, Password: password, Ip: ip})
	s.Require().NoError(err)
	s.Require().NotNil(res)
}
