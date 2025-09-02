package invites

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInviteLink_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		link        InviteLink
		expectedErr error
	}{
		{
			name: "valid full link",
			link: "https://t.me/+5V23yMex8GY5ZWFi",
		},
		{
			name: "valid hash of link",
			link: "+5V23yMex8GY5ZWFi",
		},
		{
			name:        "not enough len",
			link:        "https://t.me/+5V23yMex8GY5ZWF",
			expectedErr: ErrNotEnoughLenInviteLink,
		},
		{
			name:        "invalid schema",
			link:        "ws://t.me/+5V23yMex8GY5ZWFi",
			expectedErr: ErrInvalidInviteLink,
		},
		{
			name:        "invalid host",
			link:        "https://vk.com/+5V23yMex8GY5ZWFi",
			expectedErr: ErrInvalidInviteLink,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.ErrorIs(t, tc.link.Validate(), tc.expectedErr)
		})
	}
}

func TestInviteLinkMeta_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		meta        InviteLinkMeta
		expectedErr error
	}{
		{
			name: "valid meta",
			meta: InviteLinkMeta{
				ChannelID:       1,
				Name:            "name",
				ApproveRequired: false,
				CreatedAt:       time.Now(),
			},
		},
		{
			name: "zero channel id",
			meta: InviteLinkMeta{
				ChannelID:       0,
				Name:            "name",
				ApproveRequired: false,
				CreatedAt:       time.Now(),
			},
			expectedErr: ErrChannelIDIsEmpty,
		},
		{
			name: "zero created at",
			meta: InviteLinkMeta{
				ChannelID:       1,
				Name:            "name",
				ApproveRequired: false,
			},
			expectedErr: ErrInvalidCreatedAt,
		},
		{
			name: "invalid valid_to",
			meta: InviteLinkMeta{
				ChannelID: 1,
				CreatedAt: time.Now(),
				ValidTo:   time.Now().Add(-time.Hour * 24),
			},
			expectedErr: ErrValidToEarlierCreatedAt,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.ErrorIs(t, tc.meta.Validate(), tc.expectedErr)
		})
	}
}

func TestConvertToHash(t *testing.T) {
	testCases := []struct {
		name string
		link InviteLink
		hash Hash
	}{
		{
			name: "link with plus",
			link: "https://t.me/+5V23yMex8GY5ZWFi",
			hash: "5V23yMex8GY5ZWFi",
		},
		{
			name: "link without plus",
			link: "https://t.me/joinchat/5V23yMex8GY5ZWFi",
			hash: "5V23yMex8GY5ZWFi",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.link.Hash(), tc.hash)
		})
	}
}
