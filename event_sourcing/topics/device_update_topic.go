package topics

import (
	"encoding/json"
	"github.com/IntelXLabs-LLC/go-common/datatypes"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        *uuid.UUID      `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	UpdatedAt *time.Time      `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"-"`
	CreatedAt *time.Time      `json:"created_at"`
}

// DeviceInfoCache device info cache
type DeviceInfoCache struct {
	BaseModel
	DeviceID *uuid.UUID     `json:"device_id" gorm:"not null;index:idx_device_id"`
	State    datatypes.JSON `json:"state" gorm:"type:jsonb"`
}

type DeviceUpdateTopic struct {
	BaseModel
	ExternalDeviceID string     `json:"external_device_id" gorm:"uniqueIndex:idx_external_device_id"`
	HomeID           string     `json:"home_id"`
	UserID           string     `json:"user_id"`
	RoomID           string     `json:"room_id"`
	BrandID          *uuid.UUID `json:"brand_id"`

	Metadata        datatypes.JSON   `json:"metadata" gorm:"type:jsonb"`
	DeviceInfoCache *DeviceInfoCache `json:"device_info_cache"`
}

func (d *DeviceUpdateTopic) GetTopicName() string {
	return "device_update_topic"
}

func (d *DeviceUpdateTopic) GetBody() []byte {
	dataBytes, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}

	return dataBytes
}
