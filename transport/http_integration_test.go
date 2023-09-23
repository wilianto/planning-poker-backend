package transport_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/wilianto/planning-poker-backend/common"
	"github.com/wilianto/planning-poker-backend/model/schema/ent"
	"github.com/wilianto/planning-poker-backend/room"
	"github.com/wilianto/planning-poker-backend/transport"

	_ "github.com/lib/pq"
)

type IntegrationTestSuite struct {
	suite.Suite
	app         *fiber.App
	client      *ent.Client
	roomService *room.Service
}

func TestHttpTransportSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	// TODO: init from main app
	err := godotenv.Load("../.env")
	require.NoError(s.T(), err)

	client, err := common.InitTestDB()
	require.NoError(s.T(), err)
	s.client = client

	s.roomService = room.NewService(client)

	s.app = fiber.New()
	transport.HttpRouting(s.app, s.roomService)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	// clean up DB
	s.client.Room.Delete().ExecX(context.Background())
	s.client.Game.Delete().ExecX(context.Background())
	s.client.Player.Delete().ExecX(context.Background())
	s.client.Card.Delete().ExecX(context.Background())

	s.client.Close()
}

func (s *IntegrationTestSuite) TestCreate() {
	reqBody := `{"name": "room 1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rooms", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	require.NoError(s.T(), err)
	require.Equal(s.T(), fiber.StatusCreated, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(s.T(), err)

	var respBody room.CreateResponse
	err = json.Unmarshal(body, &respBody)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), respBody.ID)
	require.Equal(s.T(), "room 1", respBody.Name)
	require.NotEmpty(s.T(), respBody.CreatedAt)
	require.NotEmpty(s.T(), respBody.UpdatedAt)
}
