package s3

// RedisSuite is a base test suite type which launches a local docker container
// with Redis.
type RedisSuite struct {
	Suite
	Address string
}

func (s *RedisSuite) SetupSuite() {
	s.Start = StartRedisContainer
	s.Suite.SetupSuite()
	s.Address = s.Container.Addr()
}
