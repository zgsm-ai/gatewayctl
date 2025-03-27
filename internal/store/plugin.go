package store

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Plugin struct {
	ID        string         `gorm:"type:text;primaryKey"`
	URL       string         `gorm:"type:text;unique;not null"`
	Extra     datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'::jsonb"`
	CreatedAt time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	Deleted   bool           `gorm:"type:boolean;"`
}

func (Plugin) TableName() string {
	return "plugin"
}

// PluginInf defines functions used to handle user request.
type PluginInf interface {
	Create(*Plugin) error
	Delete(string) error
	Update() error
	Get(string) (Plugin, error)
	List() ([]string, error)
}

type pluginModel struct {
	DB *gorm.DB
}

var PluginModel PluginInf

func NewPluginModel(db *gorm.DB) {
	PluginModel = &pluginModel{DB: db}
	err := db.AutoMigrate(&Plugin{})
	if err != nil {
		panic(err)
	}
}

func (u *pluginModel) Create(plugin *Plugin) error {
	err := u.DB.Save(plugin).Error
	return err
}

func (u *pluginModel) Delete(id string) error {
	err := u.DB.Model(&Plugin{}).Where("id = ?", id).Where("deleted = false").Update("Deleted", true).Error
	return err
}

func (u *pluginModel) Update() error {
	return nil
}

// List returns user list in the storage. This function has a good performance.
func (u *pluginModel) List() ([]string, error) {
	var ids []string
	err := u.DB.Model(&Plugin{}).Select("id").Where("deleted = false").Find(&ids).Error
	return ids, err
}

func (u *pluginModel) Get(id string) (Plugin, error) {
	var plugin Plugin
	err := u.DB.Model(&Plugin{}).Omit("created_at").First(&plugin, "id = ? and deleted = false", id).Error
	return plugin, err
}
