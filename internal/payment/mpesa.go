package payment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Provider interface {
	InitiateCheckout(ctx context.Context, request ProviderRequest) (*ProviderResponse, error)
}

type SandboxMPESAClient struct{}

func NewSandboxMPESAClient() *SandboxMPESAClient {
	return &SandboxMPESAClient{}
}

func (c *SandboxMPESAClient) InitiateCheckout(_ context.Context, request ProviderRequest) (*ProviderResponse, error) {
	reference := fmt.Sprintf("MPESA-%s-%s", strings.ReplaceAll(request.AccountRef, " ", ""), uuid.NewString()[:8])

	return &ProviderResponse{
		Provider:          "mpesa",
		ProviderReference: reference,
		Status:            StatusInitiated,
	}, nil
}

type Clock interface {
	Now() time.Time
}
