package model

import (
	"github.com/google/uuid"
	liberr "github.com/konveyor/controller/pkg/error"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	path "path"
)

type BucketOwner struct {
	Bucket string `gorm:"<-:create"`
}

func (m *BucketOwner) BeforeCreate(db *gorm.DB) (err error) {
	err = m.Create()
	return
}

func (m *BucketOwner) BeforeDelete(db *gorm.DB) (err error) {
	err = m.DeleteBucket()
	return
}

//
// Create associated storage.
func (m *BucketOwner) Create() (err error) {
	uid := uuid.New()
	m.Bucket = path.Join(
		Settings.Hub.Bucket.Path,
		uid.String())
	err = os.MkdirAll(m.Bucket, 0777)
	if err != nil {
		err = liberr.Wrap(
			err,
			"path",
			m.Bucket)
		return
	}
	return
}

//
// EmptyBucket delete bucket content.
func (m *BucketOwner) EmptyBucket() (err error) {
	content, _ := ioutil.ReadDir(m.Bucket)
	for _, n := range content {
		p := path.Join(m.Bucket, n.Name())
		_ = os.RemoveAll(p)
	}
	return
}

//
// DeleteBucket associated storage.
func (m *BucketOwner) DeleteBucket() (err error) {
	err = os.RemoveAll(m.Bucket)
	return
}
