//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"tinder-geo/api/tinderpbv1"
	"tinder-geo/internal/app"
	"tinder-geo/internal/domain/model"
	"tinder-geo/internal/infrastructure/messaging/dto"

	. "github.com/ahmetb/go-linq/v3"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

// go test -v --tags=e2e -count=1 ./... -coverprofile=cover.out

const SECRET_KEY string = "fjg847sdjvnjxcFHdsag38d_d8sj3aqQwfdsph3456v0bjz45ty54gpo3vhjs7234f09Odp"
const TOPIC string = "profile.updates"

/*
Test steps:
1. Runs infrastructure
2. Creates profiles and set locations using kafka producer and grpc client.
3. Gets profiles using grpc call and check them.
*/
func TestGetProfiles(t *testing.T) {
	ctx, producer, client := runInfrastructure(t)

	var users []userInfo = []userInfo{
		{ID: 1, Gender: model.M, Location: locations[0][0]},
		{ID: 2, Gender: model.M, Location: locations[1][0]},
		{ID: 3, Gender: model.M, Location: locations[2][0]},
		{ID: 4, Gender: model.M, Location: locations[2][1]},
		{ID: 5, Gender: model.F, Location: locations[0][0]},
		{ID: 6, Gender: model.F, Location: locations[0][1]},
		{ID: 7, Gender: model.F, Location: locations[2][0]},
		{ID: 8, Gender: model.F, Location: locations[2][1]},
	}

	var err error
	for i, u := range users {
		users[i].Token, err = GenerateToken(u.ID, u.Gender)
		require.NoError(t, err)

		// Create profile using kafka
		err = CreateOrUpdateProfile(u.ID, u.Gender, producer, TOPIC)
		require.NoError(t, err)
	}

	// Wait kafka processing
	// todo: warning: flaky test! Check topic offset instead of time.Sleep
	time.Sleep(15 * time.Second)

	// Set users locations using GRPC call
	for _, v := range users {
		fmt.Println(v.Token)
		err = UpdateLocation(ctx, v.Token, v.Location.Latitude, v.Location.Longitude, client)
		require.NoError(t, err)
	}

	expected := []struct {
		ID          int64
		ProfilesIds []int64
	}{
		{
			ID:          1,
			ProfilesIds: []int64{5, 6},
		}, {
			ID:          2,
			ProfilesIds: []int64{},
		}, {
			ID:          3,
			ProfilesIds: []int64{7, 8},
		}, {
			ID:          4,
			ProfilesIds: []int64{7, 8},
		}, {
			ID:          5,
			ProfilesIds: []int64{1},
		}, {
			ID:          6,
			ProfilesIds: []int64{1},
		}, {
			ID:          7,
			ProfilesIds: []int64{3, 4},
		}, {
			ID:          8,
			ProfilesIds: []int64{3, 4},
		},
	}

	for i, v := range users {
		resp, err := GetProfilesByLocation(ctx, v.Token, v.Location.Latitude, v.Location.Longitude, client)
		require.NoError(t, err)

		exp := expected[i]

		assert.Equal(t, exp.ID, v.ID)

		var actualIds []int64
		From(resp.Profiles).Select(func(p interface{}) interface{} {
			return p.(*tinderpbv1.LocationProfileDto).ProfileId
		}).ToSlice(&actualIds)

		assert.ElementsMatch(t, exp.ProfilesIds, actualIds)
	}
}

// go test -run ^TestConcurrentUpdates$ -v -count=1 .
func TestConcurrentUpdates(t *testing.T) {
	ctx, producer, client := runInfrastructure(t)

	var users []userInfo = []userInfo{
		{ID: 1, Gender: model.M, Location: locations[0][0]},
		{ID: 2, Gender: model.F, Location: locations[1][0]},
	}

	var err error
	for i, u := range users {
		users[i].Token, err = GenerateToken(u.ID, u.Gender)
		require.NoError(t, err)

		// Create profile using kafka
		err = CreateOrUpdateProfile(u.ID, u.Gender, producer, TOPIC)
		require.NoError(t, err)
	}

	// Wait kafka processing
	// todo: warning: flaky test! Check topic offset instead of time.Sleep
	time.Sleep(15 * time.Second)

	calls := []func(){func() {
		err = CreateOrUpdateProfile(users[0].ID, users[0].Gender, producer, TOPIC)
		require.NoError(t, err)
	}, func() {
		err = UpdateLocation(ctx, users[0].Token, users[0].Location.Latitude, users[0].Location.Longitude, client)
		require.NoError(t, err)
	}, func() {
		err = UpdateLocation(ctx, users[0].Token, users[0].Location.Latitude, users[0].Location.Longitude, client)
		require.NoError(t, err)
	}, func() {
		err = CreateOrUpdateProfile(users[0].ID, users[0].Gender, producer, TOPIC)
		require.NoError(t, err)
	}, func() {
		err = CreateOrUpdateProfile(users[0].ID, users[0].Gender, producer, TOPIC)
		require.NoError(t, err)
	}, func() {
		err = CreateOrUpdateProfile(users[1].ID, users[1].Gender, producer, TOPIC)
		require.NoError(t, err)
	}, func() {
		err = CreateOrUpdateProfile(users[1].ID, users[1].Gender, producer, TOPIC)
		require.NoError(t, err)
	}, func() {
		err = UpdateLocation(ctx, users[1].Token, users[1].Location.Latitude, users[1].Location.Longitude, client)
		require.NoError(t, err)
	}, func() {
		err = UpdateLocation(ctx, users[1].Token, users[1].Location.Latitude, users[1].Location.Longitude, client)
		require.NoError(t, err)
	}, func() {
		err = UpdateLocation(ctx, users[0].Token, users[0].Location.Latitude, users[0].Location.Longitude, client)
		require.NoError(t, err)
	}, func() {
		err = UpdateLocation(ctx, users[0].Token, users[0].Location.Latitude, users[0].Location.Longitude, client)
		require.NoError(t, err)
	}, func() {
		err = UpdateLocation(ctx, users[0].Token, users[0].Location.Latitude, users[0].Location.Longitude, client)
		require.NoError(t, err)
	}, func() {
		err = UpdateLocation(ctx, users[1].Token, users[1].Location.Latitude, users[1].Location.Longitude, client)
		require.NoError(t, err)
	}, func() {
		err = UpdateLocation(ctx, users[0].Token, users[0].Location.Latitude, users[0].Location.Longitude, client)
		require.NoError(t, err)
	}, func() {
		err = CreateOrUpdateProfile(users[0].ID, users[0].Gender, producer, TOPIC)
		require.NoError(t, err)
	}, func() {
		err = UpdateLocation(ctx, users[0].Token, users[0].Location.Latitude, users[0].Location.Longitude, client)
		require.NoError(t, err)
	}, func() {
		err = CreateOrUpdateProfile(users[0].ID, users[0].Gender, producer, TOPIC)
		require.NoError(t, err)
	}}

	for _, fn := range calls {
		go fn()
	}

	// Wait kafka processing
	// todo: warning: flaky test! Check topic offset instead of time.Sleep
	time.Sleep(15 * time.Second)

	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:7379",
		Password: "",
		DB:       0,
	})
	defer redis.Close()

	user1Version, err := redis.Get(ctx, "profile:version:1").Int()
	require.NoError(t, err)
	user2Version, err := redis.Get(ctx, "profile:version:2").Int()
	require.NoError(t, err)

	assert.Equal(t, user1Version, 13)
	assert.Equal(t, user2Version, 6)
}

/*
1. Runs redis and kafka using docker-compose.yml.
2. Runs app including grpc server using app.Run() call.
3. Checks grpc server with healthcheck
*/
func runInfrastructure(t *testing.T) (context.Context, *kafka.Producer, tinderpbv1.GeoServiceClient) {
	compose, err := tc.NewDockerCompose("testdata/docker-compose.yml")
	require.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		require.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	require.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")

	bootstrapServers := "localhost:39092"

	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	require.NoError(t, err)

	_, err = admin.CreateTopics(ctx, []kafka.TopicSpecification{{
		Topic:             TOPIC,
		NumPartitions:     1,
		ReplicationFactor: 1}})
	require.NoError(t, err)

	t.Setenv("CONFIG_PATH", "testdata/config.yaml")

	closer := app.Run()
	t.Cleanup(closer)

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	require.NoError(t, err)

	conn, err := grpc.NewClient("localhost:3342", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	// Check grpc server
	err = retry(func() error {
		grpcHealthClient := grpc_health_v1.NewHealthClient(conn)
		resp, err := grpcHealthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
		if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
			return fmt.Errorf("failed: %w", err)
		}
		return nil
	}, 5, 1*time.Second)
	require.NoError(t, err)

	client := tinderpbv1.NewGeoServiceClient(conn)

	return ctx, producer, client
}

func retry(f func() error, attempts int, timeout time.Duration) (err error) {
	for i := 0; i < attempts; i++ {
		err = f()
		if err == nil {
			return
		}
		if i != attempts-1 {
			time.Sleep(timeout)
		}
	}
	return
}

type location struct {
	Latitude  float64
	Longitude float64
}

var locations [][]location = [][]location{{
	// gcpvj
	{51.509, -0.118},
	{51.506, -0.119},
}, {
	// gcpuv
	{51.500, -0.118},
	{51.501, -0.120},
}, {
	// gcpvm
	{51.561, -0.120},
	{51.562, -0.121},
}}

type userInfo struct {
	ID       int64
	Gender   model.Gender
	Token    string
	Location location
}

func GenerateToken(id int64, gender model.Gender) (string, error) {
	key := []byte(SECRET_KEY)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"ProfileId": strconv.FormatInt(id, 10),
			"Gender":    gender.String(),
		})
	return token.SignedString(key)
}

func CreateOrUpdateProfile(id int64, gender model.Gender, p *kafka.Producer, topic string) error {
	profile := dto.ProfileDto{}
	profile.ID = id
	profile.Gender = gender
	profile.Age = 30
	profile.Name = "Some name"
	profile.Description = "Some description"
	profile.Photos = []string{"http://photo-1", "http://photo-2"}

	b, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          b},
		nil,
	)
	return err
}

func UpdateLocation(ctx context.Context, token string, latitude float64, longitude float64, client tinderpbv1.GeoServiceClient) error {
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %s", token)})
	ctxGrpc := metadata.NewOutgoingContext(ctx, md)
	_, err := client.ChangeLocation(ctxGrpc, &tinderpbv1.ChangeLocationRequest{Latitude: latitude, Longitude: longitude})
	return err
}

func GetProfilesByLocation(ctx context.Context, token string, latitude float64, longitude float64, client tinderpbv1.GeoServiceClient) (*tinderpbv1.GetProfilesByLocationResponse, error) {
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %s", token)})
	ctxGrpc := metadata.NewOutgoingContext(ctx, md)
	return client.GetProfilesByLocation(ctxGrpc, &tinderpbv1.GetProfilesByLocationRequest{Latitude: latitude, Longitude: longitude})
}
