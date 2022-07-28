package replicationqueue

import (
	"fmt"

	"github.com/inpher/sb/internal/models"
	"github.com/inpher/sb/internal/replicationqueue/googlepubsub"
	"github.com/inpher/sb/internal/types"
	"github.com/pkg/errors"
)

type ReplicationQueue interface {
	PushToQueue(*models.Replication) error
	ConsumeQueue(func(*models.Replication) error) error
}

func GetReplicationQueue(config *types.ReplicationQueueConfig, hostname string) (rq ReplicationQueue, err error) {

	if !config.Enabled {
		err = fmt.Errorf("TTYRecs offloading is disabled")
		return
	}

	switch config.QueueType {
	case "googlepubsub":
		rq, err = googlepubsub.NewReplicationQueuePubSub(config.QueueOptions, hostname)
	default:
		err = fmt.Errorf("storage type %s is not implemented", config.QueueType)
	}

	if err != nil {
		err = errors.Wrap(err, "error while initializing ttyrecsoffloading")
	}

	return
}
