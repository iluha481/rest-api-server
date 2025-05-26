package tests

import (
	"net"
	"server/tests/suite"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	s := suite.New(t)
	url := net.JoinHostPort(s.Cfg.Host, s.Cfg.Port)
	url = "http://" + url
	resp, err := s.HttpClient.Get(url + "/ping")
	require.NoError(t, err)
	require.NotEmpty(t, resp)
}
