package healthchecker

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"go.temporal.io/api/operatorservice/v1"
	"go.temporal.io/api/workflowservice/v1"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"sync"
	"time"
)

const (
	fullWorkflowServiceName = "temporal.api.workflowservice.v1.WorkflowService"
	fullHistoryServiceName  = "temporal.api.workflowservice.v1.HistoryService"
	fullMatchingServiceName = "temporal.api.workflowservice.v1.MatchingService"
)

// HealthChecker is used to check health of Temporal services
type HealthChecker struct {
	cfg HealthCheckConfig
	mu  sync.Mutex // protects following
	grpcConn
	wfs workflowservice.WorkflowServiceClient
	ops operatorservice.OperatorServiceClient
}

// grpcConn contains grpc connections to Temporal services
// TODO: add support for TLS
type grpcConn struct {
	frontend *grpc.ClientConn
	history  *grpc.ClientConn
	matching *grpc.ClientConn
}

// HealthCheckConfig contains configuration for HealthChecker
type HealthCheckConfig struct {
	FrontendService HealthCheckServiceConfig
	HistoryService  HealthCheckServiceConfig
	MatchingService HealthCheckServiceConfig
}

// HealthCheckServiceConfig contains configuration for a single Temporal service
// TODO: add support for TLS
type HealthCheckServiceConfig struct {
	IsEnabled bool
	Address   string
	TimeOut   int
}

// FormatError formats error message for health check
func FormatError(err error) string {
	if multiErr, ok := err.(*multierror.Error); ok {
		for i, err := range multiErr.Errors {
			multiErr.Errors[i] = fmt.Errorf("Health check failed: %v", err)
		}
	}
	return err.Error()
}

// grpcIsServing checks if a grpc service is serving
func (hc *HealthChecker) grpcIsServing(svc string, conn *grpc.ClientConn) error {
	hcli := healthpb.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(hc.cfg.FrontendService.TimeOut)*time.Second)
	defer cancel()

	req := &healthpb.HealthCheckRequest{
		Service: svc,
	}

	check, err := hcli.Check(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to check health of %s: %v", svc, err)
	}

	if check.Status != healthpb.HealthCheckResponse_SERVING {
		return fmt.Errorf("health check of %s failed with status %s", svc, check.Status)
	}

	return nil
}

// NewHealthChecker creates a new HealthChecker
// cfg is the configuration for HealthChecker
// multiple errors can be returned if multiple services are enabled but failed to connect
func NewHealthChecker(cfg HealthCheckConfig) (*HealthChecker, error) {
	hc := &HealthChecker{
		cfg: cfg,
	}

	var err error
	var errs *multierror.Error

	if cfg.FrontendService.IsEnabled {
		if cfg.FrontendService.Address == "" {
			errs = multierror.Append(errs, fmt.Errorf("frontend service address is empty"))
		} else {
			hc.grpcConn.frontend, err = grpc.Dial(cfg.FrontendService.Address, grpc.WithInsecure())
			if err != nil {
				errs = multierror.Append(errs, err)
			}
		}
	}

	if cfg.HistoryService.IsEnabled {
		if cfg.HistoryService.Address == "" {
			errs = multierror.Append(errs, fmt.Errorf("history service address is empty"))
		} else {
			hc.grpcConn.history, err = grpc.Dial(cfg.HistoryService.Address, grpc.WithInsecure())
			if err != nil {
				errs = multierror.Append(errs, err)
			}
		}
	}

	if cfg.MatchingService.IsEnabled {
		if cfg.MatchingService.Address == "" {
			errs = multierror.Append(errs, fmt.Errorf("matching service address is empty"))
		} else {
			hc.grpcConn.matching, err = grpc.Dial(cfg.MatchingService.Address, grpc.WithInsecure())
			if err != nil {
				errs = multierror.Append(errs, err)
			}

		}
	}

	if errs != nil {
		return nil, fmt.Errorf("failed to create health checker: %v", errs.ErrorOrNil())
	}

	return hc, nil
}

// BasicCheck checks if the Temporal services are serving
// multiple errors can be returned if multiple services are enabled but failed to connect
func (hc *HealthChecker) BasicCheck() error {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	var errs *multierror.Error
	if hc.cfg.FrontendService.IsEnabled {
		errs = multierror.Append(errs, hc.grpcIsServing(fullWorkflowServiceName, hc.grpcConn.frontend))
	}
	if hc.cfg.HistoryService.IsEnabled {
		errs = multierror.Append(errs, hc.grpcIsServing(fullHistoryServiceName, hc.grpcConn.history))
	}
	if hc.cfg.MatchingService.IsEnabled {
		errs = multierror.Append(errs, hc.grpcIsServing(fullMatchingServiceName, hc.grpcConn.matching))
	}

	return errs.ErrorOrNil()
}

// FullCheck checks if the Temporal services are serving and if the cluster is healthy
// This will check if the cluster is healthy by checking if the cluster information is available,
// system namespace is available, and if the system workflows are running
func (hc *HealthChecker) FullCheck() error {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	var errs *multierror.Error
	if hc.cfg.FrontendService.IsEnabled {

		err := hc.initialFullCheckClients()
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to initialize full check clients: %v", err))
		}

		err = hc.checkClusterInfo()
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to check cluster info: %v", err))
		}
		err = hc.checkSystemInfo()
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to check system info: %v", err))
		}
		err = hc.checkNamespaces()
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to check namespaces: %v", err))
		}
		err = hc.checkListClusters()
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to check clusters: %v", err))
		}
	}

	return errs.ErrorOrNil()
}

// initialFullCheckClients initializes the clients needed for full check
func (hc *HealthChecker) initialFullCheckClients() error {
	if hc.cfg.FrontendService.IsEnabled != true {
		return fmt.Errorf("frontend service is not enabled")
	}

	var errs *multierror.Error
	if hc.cfg.FrontendService.IsEnabled {
		hc.wfs = workflowservice.NewWorkflowServiceClient(hc.grpcConn.frontend)
		if hc.wfs == nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to initialize WorkflowServiceClient"))
		}
		hc.ops = operatorservice.NewOperatorServiceClient(hc.grpcConn.frontend)
		if hc.ops == nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to initialize OperatorServiceClient"))
		}
	}

	return errs.ErrorOrNil()
}

// checkSystemInfo checks if the system namespace is available and if the system workflows are running
func (hc *HealthChecker) checkClusterInfo() error {
	info, err := hc.wfs.GetClusterInfo(context.Background(), &workflowservice.GetClusterInfoRequest{})
	if err != nil {
		return fmt.Errorf("failed to get cluster info: %v", err)
	}

	var errs *multierror.Error
	if info.GetClusterId() == "" {
		errs = multierror.Append(errs, fmt.Errorf("Cluster Id is empty"))
	}

	if info.GetClusterName() == "" {
		errs = multierror.Append(errs, fmt.Errorf("Cluster Name is empty"))
	}

	log.Printf("Cluster Id: %s", info.GetClusterId())
	log.Printf("Version Info: %s", info.GetVersionInfo())
	log.Printf("Cluster Name: %s", info.GetClusterName())
	log.Printf("History Shard Count: %d", info.GetHistoryShardCount())

	return errs.ErrorOrNil()
}

// checkSystemInfo checks if the system namespace is available and if the system workflows are running
func (hc *HealthChecker) checkSystemInfo() error {
	systemInfo, err := hc.wfs.GetSystemInfo(context.Background(), &workflowservice.GetSystemInfoRequest{})
	if err != nil {
		return fmt.Errorf("failed to get system info: %v", err)
	}

	cp := systemInfo.GetCapabilities()
	log.Printf("Capabilities ActivityFailureIncludeHeartbeat: %v", cp.ActivityFailureIncludeHeartbeat)
	log.Printf("Capabilities SdkMetadata: %v", cp.SdkMetadata)
	log.Printf("Capabilities BuildIdBasedVersioning: %v", cp.BuildIdBasedVersioning)
	log.Printf("Capabilities UpsertMemo: %v", cp.UpsertMemo)
	log.Printf("ServerVersion: %s", systemInfo.GetServerVersion())

	return nil
}

// checkNamespaces checks if the namespaces are available
func (hc *HealthChecker) checkNamespaces() error {
	// TODO: add pagination
	namespaces, err := hc.wfs.ListNamespaces(context.Background(), &workflowservice.ListNamespacesRequest{
		PageSize: 10,
	})
	if err != nil {
		return fmt.Errorf("failed to list namespaces: %v", err)
	}

	for _, ns := range namespaces.Namespaces {
		nsInfo := ns.GetNamespaceInfo()
		log.Printf("Namespace: %s, State: %s, Description: %s", nsInfo.GetName(), nsInfo.GetState(), nsInfo.GetDescription())
	}

	return nil
}

func (hc *HealthChecker) checkListClusters() error {
	clusters, err := hc.ops.ListClusters(context.Background(), &operatorservice.ListClustersRequest{})
	if err != nil {
		return fmt.Errorf("failed to list clusters: %v", err)
	}

	for i, c := range clusters.Clusters {
		log.Printf("Cluster %d: %s, %s", i, c.GetClusterName(), c.GetClusterId())
		if c.GetIsConnectionEnabled() == true {
			log.Printf("Cluster %d is connected", i)
		} else {
			log.Printf("Cluster %d is not connected", i)
			return fmt.Errorf("cluster %d is not connected", i)
		}
	}

	return nil
}

func (hc *HealthChecker) Close() {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	if hc.grpcConn.frontend != nil {
		hc.grpcConn.frontend.Close()
	}
	if hc.grpcConn.history != nil {
		hc.grpcConn.history.Close()
	}
	if hc.grpcConn.matching != nil {
		hc.grpcConn.matching.Close()
	}
}
