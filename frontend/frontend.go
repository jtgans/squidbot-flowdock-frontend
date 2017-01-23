package frontend

import (
	"log"
	"expvar"
	"time"

	"golang.org/x/net/context"
	"github.com/wm/go-flowdock/flowdock"
	"google.golang.org/grpc"

	pb "github.com/jtgans/squidbot-grpc"
)

var (
	flowdockRequests = expvar.NewInt("flowdock_requests")
	flowdockConnectionErrors = expvar.NewInt("flowdock_connection_errors")
	flowdockErrors = expvar.NewMap("flowdock_errors")
	flowdockMessagesHandled = expvar.NewMap("flowdock_messages_handled")
	flowdockMessagesReceived = expvar.NewMap("flowdock_messages_received")

	brainReconnects = expvar.NewInt("brain_reconnects")
	brainConnectionErrors = expvar.NewInt("brain_connection_errors")
	brainErrors = expvar.NewMap("brain_errors")
)

type Frontend struct {
	brainHostPort string
	brainClient pb.BrainClient
	
	personalToken string
	fdClient *flowdock.Client
	flows []flowdock.Flow
}

const version = "0.1.0"

func NewFrontend(brainHostPort string, personalToken string) *Frontend {
	return &Frontend {
		brainHostPort: brainHostPort,
		personalToken: personalToken,
	}
}

func (frontend *Frontend) IsOk() string {
	return "ok"
}

func (frontend *Frontend) UpdateCachedFlows() {
	flowdockRequests.Add(1)
	flows, response, err := frontend.fdClient.Flows.List(true, &flowdock.FlowsListOptions{ User: false })
	if err != nil {
		log.Printf("Unable to update flow list with flowdock: %v", err)
		log.Printf("Additionally, the following response was received: %v", response)
		flowdockConnectionErrors.Add(1)
		flowdockErrors.Add(err.Error(), 1)
	}

	frontend.flows = flows
}

func (frontend *Frontend) StartFlowdockConnector() {
	frontend.fdClient = flowdock.NewClientWithToken(nil, frontend.personalToken)
	frontend.UpdateCachedFlows()
}

func (frontend *Frontend) StartBrainConnector() {
	var resp *pb.FrontendResponse
	var attempts int64 = 0

	for connected := false; !connected; attempts++ {
		brainReconnects.Add(1)
		time.Sleep(time.Duration(attempts) * time.Second)

		conn, err := grpc.Dial(frontend.brainHostPort, grpc.WithInsecure())
		if err != nil {
			log.Printf("Unable to connect to brain: %v", err)
			brainConnectionErrors.Add(1)
			continue
		}

		frontend.brainClient = pb.NewBrainClient(conn)

		req := &pb.FrontendRequest {
			FrontendVersion: version,
			FrontendName: "squidbot-flowdock-frontend",
		}
		resp, err = frontend.brainClient.FrontendStarted(context.Background(), req)

		if err != nil {
			log.Printf("Unable to register with brain: %v", err)
			brainErrors.Add(err.Error(), 1)
			continue
		}
	}

	log.Printf("Registered with Squidbot brain version %v", resp.BrainVersion)
}

func (frontend *Frontend) Start() {
	log.Printf("Squidbot Flowdock Frontend v%v starting", version)
	go frontend.StartFlowdockConnector();
	go frontend.StartBrainConnector();
}
