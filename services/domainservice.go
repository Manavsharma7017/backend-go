package services

import (
	"backend/database" // Assuming you have a `database` package for DB connection
	"backend/models"
	"errors"
)

// CreateDomain creates a new domain in the database
func CreateDomain(domain *models.Domain) (*models.Domain, error) {
	// Assuming `database.DB` is your DB instance
	if err := database.DB.Create(domain).Error; err != nil {
		return nil, err
	}
	return domain, nil
}

// GetDomainByID retrieves a domain by its ID
func GetDomainByID(id string) (*models.Domain, error) {
	var domain models.Domain
	if err := database.DB.First(&domain, "id = ?", id).Error; err != nil {
		return nil, errors.New("domain not found")
	}

	return &domain, nil
}

// GetAllDomains retrieves all domains from the database
func GetAllDomains() ([]models.Domain, error) {
	var domains []models.Domain
	if err := database.DB.Find(&domains).Error; err != nil {
		return nil, err
	}
	return domains, nil
}

// UpdateDomain updates an existing domain
func UpdateDomain(id string, domain *models.Domain) (*models.Domain, error) {
	var existingDomain models.Domain
	if err := database.DB.First(&existingDomain, "id = ?", id).Error; err != nil {
		return nil, errors.New("domain not found")
	}

	// Update the domain fields (assuming only name is updatable, add more fields as needed)
	existingDomain.Name = domain.Name
	// existingDomain.Description = domain.Description
	if err := database.DB.Save(&existingDomain).Error; err != nil {
		return nil, err
	}
	return &existingDomain, nil
}

// DeleteDomain deletes a domain by its ID
func DeleteDomain(id string) error {
	var domain models.Domain
	if err := database.DB.First(&domain, "id = ?", id).Error; err != nil {
		return errors.New("domain not found")
	}

	if err := database.DB.Delete(&domain).Error; err != nil {
		return err
	}
	return nil
}
