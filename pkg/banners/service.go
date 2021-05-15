package banners

import (
	"io/ioutil"
	"fmt"
	"mime/multipart"
	"context"
	"errors"
	"sync"
)

var startId int64

//Service представляет собой сервис по управления баннерами
type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

//NewService создаёт сервис
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

//Banner предсталяет собой баннер
type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
	Image   string
}

//All ...
func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items, nil
}

//ByID ...
func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}

	return nil, errors.New("item not found")
}

//Save ...
func (s *Service) Save(ctx context.Context, item *Banner, file multipart.File) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if item.ID == 0 {
		startId++
		item.ID=startId

		if item.Image != "" {
			item.Image = fmt.Sprint(item.ID) + "." + item.Image
			err := uploadFile(file, "./web/banners/"+item.Image)
			if err != nil {
				return nil, err
			}
		}

		s.items=append(s.items, item)
		return item, nil
	}

	for key, banner := range s.items {
		if banner.ID == item.ID {
			if item.Image != "" {
				item.Image = fmt.Sprint(item.ID) + "." + item.Image
				err := uploadFile(file, "./web/banners/"+item.Image)
				if err != nil {
					return nil, err
				}
			} else {
				item.Image = s.items[key].Image
			}
			s.items[key]=item
			return item, nil
		}
	}

	return nil, errors.New("item not found")
}

//RemoveByID ... Метод для удаления
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, banner := range s.items {
		if banner.ID == id {
			s.items = append(s.items[:k], s.items[k+1:]...)
			return banner, nil
		}
	}

	return nil, errors.New("item not found")
}

func uploadFile(file multipart.File, path string) error {
	var data, err = ioutil.ReadAll(file)

	if err != nil {
		return errors.New("not readble data")
	}

	err = ioutil.WriteFile(path, data, 0666)

	if err != nil {
		return errors.New("not saved from folder ")
	}

	return nil
}
