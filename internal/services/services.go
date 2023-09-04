package services

import (
	"golang.org/x/crypto/bcrypt"
	"main/internal/configs/repositories"
	"main/pkg/models"
)

type Service struct {
	Repository *repositories.Repository
}

func NewService(repository *repositories.Repository) *Service {
	return &Service{Repository: repository}
}
func (s *Service) Create(note *models.Notes) error {
	return s.Repository.Create(note)
}
func (s *Service) Read(id *int) (models.Notes, error) {
	note, err := s.Repository.Read(id)
	return note, err
}
func (s *Service) Update(note *models.Notes) {
	s.Repository.Update(note)
}
func (s *Service) Delete(id *int) {
	s.Repository.Delete(id)
}
func (s *Service) ReadAll() ([]models.Notes, error) {
	notes, err := s.Repository.ReadAll()
	return notes, err
}
func (s *Service) UserRegistration(user *models.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = s.Repository.UserRegistration(user, hash)
	return err
}

func (s *Service) HashingThePassword(user *models.User) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	id, err := s.Repository.CheckUser(user, hash)
	return id, err
}
