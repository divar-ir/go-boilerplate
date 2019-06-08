package provider

import (
	"context"
	"encoding/base64"

	"github.com/golang/protobuf/proto"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"
	gormigrate "gopkg.in/gormigrate.v1"

	"github.com/jinzhu/gorm"
)

type sqlProvider struct {
	db *gorm.DB
}

type post struct {
	gorm.Model

	Token string `gorm:"type:varchar(8);primary_key;"`
	Data  string `gorm:"type:text;"`
}

func NewSQL(db *gorm.DB) PostProvider {
	return sqlProvider{
		db: db,
	}
}

func (p sqlProvider) Close() error {
	return p.db.Close()
}

func (p sqlProvider) GetPost(ctx context.Context, token string) (*postview.Post, error) {
	postInstance := &post{}
	err := p.db.Where("token = ?", token).First(postInstance).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	result, err := p.modelToProto(postInstance)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p sqlProvider) AddPost(ctx context.Context, protoPost *postview.Post) error {
	modelInstance, err := p.protoToModel(protoPost)
	if err != nil {
		return err
	}

	return p.db.Create(modelInstance).Error
}

func (p sqlProvider) Migrate() error {
	m := gormigrate.New(p.db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// create model table
		{
			ID: "201906082327",
			Migrate: func(tx *gorm.DB) error {
				type post struct {
					gorm.Model

					Token string `gorm:"type:varchar(8);primary_key;"`
					Data  string `gorm:"type:text;"`
				}
				return tx.AutoMigrate(&post{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("posts").Error
			},
		},
	})

	return m.Migrate()
}

func (p sqlProvider) protoToModel(protoPost *postview.Post) (*post, error) {
	binaryData, err := proto.Marshal(protoPost)
	if err != nil {
		return nil, err
	}

	data := base64.StdEncoding.EncodeToString(binaryData)

	return &post{
		Token: protoPost.Token,
		Data:  data,
	}, nil
}

func (p sqlProvider) modelToProto(m *post) (*postview.Post, error) {
	data, err := base64.StdEncoding.DecodeString(m.Data)
	if err != nil {
		return nil, err
	}

	var result postview.Post
	err = proto.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
