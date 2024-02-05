package pipectl

import api "github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"

type APIClient interface {
	api.APIServiceClient
	Close() error
}
