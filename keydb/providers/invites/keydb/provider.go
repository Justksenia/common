package keydb

import (
	"fmt"

	"gitlab.com/adstail/ts-common/keydb/redis"
)

const (
	expirationTime = redis.PersistentTTL
	instanceName   = "invite_links_storage"
	keyPrefixList  = "invite-link-channel-list"
	keyPrefixLink  = "invite_link-hash-link"
)

type InviteLinksKeyDBProvider struct {
	client *redis.Instance
}

func New(client *redis.KeyDBFactory) *InviteLinksKeyDBProvider {
	instance := client.NewInstance(instanceName, expirationTime)
	return &InviteLinksKeyDBProvider{
		client: instance,
	}
}

func channelIDKey(channelID int64) string {
	return fmt.Sprintf("%s-%d", keyPrefixList, channelID)
}

func hashLinkKey(hash string) string {
	return fmt.Sprintf("%s-%s", keyPrefixLink, hash)
}
