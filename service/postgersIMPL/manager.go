package postgersIMPL

import (
	"errors"
	"fmt"
	_ "github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"log"
	"math/rand"
	"slices"
)

type id_seg struct {
	PersonID int            `gorm:"primaryKey,unique"`
	Segments pq.StringArray `gorm:"type:text[]"`
}

type seg_id struct {
	Segment string        `gorm:"primaryKey,unique"`
	Ids     pq.Int64Array `gorm:"type:integer[]"`
}

func NewManager(dsn string) *Manager {
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&id_seg{})
	db.AutoMigrate(&seg_id{})

	return &Manager{
		db,
	}
}

type Manager struct {
	db *gorm.DB
}

func (m *Manager) PostPaS(personID int, segments []string) error {
	var cnt int
	err := m.db.Model(&seg_id{}).Where("segment = ANY(?)", pq.StringArray(segments)).Count(&cnt).Error
	if err != nil {
		return err
	}
	if cnt != len(segments) {
		return fmt.Errorf("there are unregistered segments in the request")
	}

	var item id_seg
	req := m.db.Where("person_id = ?", personID).First(&item)
	if errors.Is(req.Error, gorm.ErrRecordNotFound) {
		err = m.db.Transaction(func(tx *gorm.DB) error {
			if err = tx.Create(&id_seg{personID, segments}).Error; err != nil {
				return err
			}

			for _, seg := range segments {
				var si seg_id
				if err = m.db.Where("segment = ?", seg).First(&si).Error; err != nil {
					return err
				}
				var oldId = []int64(si.Ids)
				oldId = append(oldId, int64(personID))
				if err = tx.Model(&seg_id{}).Where("segment = ?", seg).Update("ids", oldId).Error; err != nil {
					return err
				}
			}
			return nil
		})
	} else if err == nil {
		var oldSeg = []string(item.Segments)
		var newSeg []string
		for _, seg := range segments {
			if !slices.Contains(oldSeg, seg) {
				newSeg = append(newSeg, seg)
			}
		}
		oldSeg = append(oldSeg, newSeg...)

		err = m.db.Transaction(func(tx *gorm.DB) error {
			if err = tx.Model(&id_seg{}).Where("person_id = ?", personID).Update("segments", oldSeg).Error; err != nil {
				return err
			}

			for _, seg := range newSeg {
				var si seg_id
				if err = m.db.Where("segment = ?", seg).First(&si).Error; err != nil {
					return err
				}
				var oldId = []int64(si.Ids)
				oldId = append(oldId, int64(personID))
				if err = tx.Model(&seg_id{}).Where("segment = ?", seg).Update("ids", oldId).Error; err != nil {
					return err
				}
			}
			return nil
		})
	}
	return err
}

func (m *Manager) GetSegments(personID int) ([]string, error) {
	var item id_seg
	err := m.db.Where("person_id = ?", personID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return item.Segments, nil
}

func (m *Manager) GetIDs(segment string) ([]int64, error) {
	var item seg_id
	err := m.db.Where("segment = ?", segment).First(&item).Error
	if err != nil {
		return nil, err
	}
	return item.Ids, nil
}

func (m *Manager) PostS(segments []string) error {
	var item seg_id
	err := m.db.Transaction(func(tx *gorm.DB) error {
		for _, segment := range segments {
			err := tx.Where("segment = ?", segment).First(&item).Error
			if err == nil {
				return fmt.Errorf("record with name \"%s\" exists", segment)
			} else {
				res := tx.Create(&seg_id{segment, pq.Int64Array{}})
				if res.Error != nil {
					return res.Error
				}
			}
		}
		return nil
	})
	return err
}

func removeElemStr(segments []string, removeSeg string) []string {
	var newItems []string

	for _, i := range segments {
		if i != removeSeg {
			newItems = append(newItems, i)
		}
	}
	return newItems
}

func (m *Manager) DeleteSegment(segment string) error {
	var item seg_id
	err := m.db.Where("segment = ?", segment).First(&item).Error
	if err != nil {
		return err
	}

	err = m.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Where("segment = ?", segment).Delete(&seg_id{}).Error; err != nil {
			return err
		}
		ids := []int64(item.Ids)
		for _, id := range ids {
			var t id_seg
			if err = tx.Where("person_id = ?", id).First(&t).Error; err != nil {
				return err
			}

			newSlice := removeElemStr(t.Segments, segment)
			log.Println("NEWSLICE: ", newSlice)
			if err = tx.Model(&id_seg{}).Where("person_id = ?", id).Update("segments", newSlice).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func removeElements(segments []string, deleteSeg []string) []string {
	lenArr := len(segments)
	var k int = 0
	for i := 0; i < lenArr; {
		if !slices.Contains(deleteSeg, segments[i]) {
			segments[k] = segments[i]
			k++
		}
		i++
	}
	return segments[0:k]
}

func removeElemInt(segments []int64, removeSeg int64) []int64 {
	var newItems []int64

	for _, i := range segments {
		if i != removeSeg {
			newItems = append(newItems, i)
		}
	}
	return newItems
}

func (m *Manager) DeleteSegments(personID int, segments []string) error {
	var item id_seg
	err := m.db.Where("person_id = ?", personID).First(&item).Error
	if err != nil {
		return err
	}

	newSeg := removeElements(item.Segments, segments)
	err = m.db.Transaction(func(tx *gorm.DB) error {
		if err = m.db.Model(&id_seg{}).Where("person_id = ?", personID).Update("segments", newSeg).Error; err != nil {
			return err
		}

		for _, segment := range segments {
			var t seg_id
			if err = tx.Where("segment = ?", segment).First(&t).Error; err != nil {
				return err
			}

			newSlice := removeElemInt(t.Ids, int64(personID))
			if err = tx.Model(&seg_id{}).Where("segment = ?", segment).Update("ids", newSlice).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (m *Manager) PostWithPer(segment string, ids []int64, per int) ([]int64, error) {
	var item seg_id
	err := m.db.Where("segment = ?", segment).First(&item).Error
	if err != nil {
		return nil, err
	}

	var newIds []int64
	for _, val := range ids {
		random := rand.Intn(100)
		if random <= per && !slices.Contains(item.Ids, val) && !slices.Contains(newIds, val) {
			newIds = append(newIds, val)
		}
	}
	item.Ids = append(item.Ids, newIds...)
	err = m.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Model(&seg_id{}).Where("segment = ?", segment).Update("ids", item.Ids).Error; err != nil {
			return err
		}

		for _, id := range item.Ids {
			var check id_seg
			err = m.db.Where("person_id = ?", id).First(&check).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err = tx.Create(&id_seg{int(id), []string{segment}}).Error; err != nil {
					return err
				}
			} else if err == nil {
				check.Segments = append(check.Segments, segment)
				if err = tx.Model(&id_seg{}).Where("person_id = ?", id).Update("segments", check.Segments).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return nil
	})

	return newIds, nil
}
