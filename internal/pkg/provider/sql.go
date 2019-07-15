package provider

import (
	"context"
	"encoding/base64"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/errors"
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
			return nil, errors.WrapWithExtra(ErrNotFound, "post not found", map[string]interface{}{
				"token": token,
			})
		}
		return nil, errors.WrapWithExtra(err, "could not read post from db", map[string]interface{}{
			"token": token,
		})
	}

	result, err := p.modelToProto(postInstance)
	if err != nil {
		return nil, errors.WrapWithExtra(err, "could not convert model to proto", map[string]interface{}{
			"token": token,
		})
	}

	return result, nil
}

func (p sqlProvider) AddPost(ctx context.Context, protoPost *postview.Post) error {
	modelInstance, err := p.protoToModel(protoPost)
	if err != nil {
		return errors.WrapWithExtra(err, "could not convert proto to model", map[string]interface{}{
			"post": protoPost,
		})
	}

	err = p.db.Create(modelInstance).Error
	if err != nil {
		return errors.WrapWithExtra(err, "could not add model to db", map[string]interface{}{
			"post": protoPost,
		})
	}

	return nil
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
		return nil, errors.Wrap(err, "could not marshal proto")
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
		return nil, errors.Wrap(err, "could not decode base64")
	}

	var result postview.Post
	err = proto.Unmarshal(data, &result)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal proto")
	}

	return &result, nil
}
