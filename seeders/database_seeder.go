package seeders

import (
	"fmt"
	"motocare-dashboard/backend/models"
	"motocare-dashboard/backend/utils"
	"time"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	adminUser, err := seedUser(db, models.User{
		Username: "admin",
		Email:    "admin@motocare.test",
		Role:     "admin",
	})
	if err != nil {
		return err
	}

	regularUser, err := seedUser(db, models.User{
		Username: "user",
		Email:    "user@motocare.test",
		Role:     "user",
	})
	if err != nil {
		return err
	}

	categories, err := seedCategories(db)
	if err != nil {
		return err
	}

	services, err := seedServices(db, categories)
	if err != nil {
		return err
	}

	return seedBookings(db, regularUser.ID, adminUser.ID, services)
}

func seedUser(db *gorm.DB, user models.User) (*models.User, error) {
	var existingUser models.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return &existingUser, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func seedCategories(db *gorm.DB) ([]models.ServiceCategory, error) {
	seedData := []models.ServiceCategory{
		{Name: "Servis Berkala", Description: "Pemeriksaan rutin untuk menjaga performa motor."},
		{Name: "Ganti Oli", Description: "Layanan penggantian oli mesin dan oli transmisi."},
		{Name: "Tune Up", Description: "Penyetelan mesin, injeksi, dan komponen pembakaran."},
		{Name: "Rem", Description: "Perawatan kampas, cakram, dan sistem pengereman."},
		{Name: "Ban", Description: "Pemeriksaan, tambal, dan penggantian ban motor."},
		{Name: "Aki", Description: "Pengecekan dan penggantian aki motor."},
		{Name: "Kelistrikan", Description: "Perbaikan lampu, starter, kabel, dan sistem listrik."},
		{Name: "Rantai dan Gear", Description: "Pembersihan, penyetelan, dan penggantian rantai gear."},
		{Name: "Suspensi", Description: "Pemeriksaan shockbreaker dan kenyamanan berkendara."},
		{Name: "Detailing", Description: "Pembersihan mendalam dan perawatan tampilan motor."},
	}

	categories := make([]models.ServiceCategory, 0, len(seedData))
	for _, category := range seedData {
		var existingCategory models.ServiceCategory
		if err := db.Where("name = ?", category.Name).First(&existingCategory).Error; err == nil {
			categories = append(categories, existingCategory)
			continue
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		if err := db.Create(&category).Error; err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func seedServices(db *gorm.DB, categories []models.ServiceCategory) ([]models.Service, error) {
	seedData := []models.Service{
		{CategoryID: categories[0].ID, Name: "Servis Berkala 5.000 KM", Description: "Pemeriksaan dasar motor untuk pemakaian harian.", Price: 85000, DurationMinutes: 60, Status: "active"},
		{CategoryID: categories[1].ID, Name: "Ganti Oli Mesin", Description: "Penggantian oli mesin standar bengkel.", Price: 65000, DurationMinutes: 30, Status: "active"},
		{CategoryID: categories[2].ID, Name: "Tune Up Injeksi", Description: "Pembersihan throttle body dan pengecekan injeksi.", Price: 120000, DurationMinutes: 90, Status: "active"},
		{CategoryID: categories[3].ID, Name: "Ganti Kampas Rem", Description: "Penggantian kampas rem depan atau belakang.", Price: 95000, DurationMinutes: 45, Status: "active"},
		{CategoryID: categories[4].ID, Name: "Tambal Ban Tubeless", Description: "Perbaikan ban tubeless bocor ringan.", Price: 25000, DurationMinutes: 20, Status: "active"},
		{CategoryID: categories[5].ID, Name: "Ganti Aki Motor", Description: "Pengecekan kelistrikan dan penggantian aki.", Price: 250000, DurationMinutes: 30, Status: "active"},
		{CategoryID: categories[6].ID, Name: "Perbaikan Lampu", Description: "Pemeriksaan dan penggantian lampu motor.", Price: 55000, DurationMinutes: 35, Status: "active"},
		{CategoryID: categories[7].ID, Name: "Setel Rantai", Description: "Pembersihan dan penyetelan ketegangan rantai.", Price: 35000, DurationMinutes: 25, Status: "active"},
		{CategoryID: categories[8].ID, Name: "Servis Shockbreaker", Description: "Pemeriksaan dan servis shockbreaker motor.", Price: 180000, DurationMinutes: 120, Status: "active"},
		{CategoryID: categories[9].ID, Name: "Cuci Detailing Motor", Description: "Cuci detail bodi, mesin, dan velg motor.", Price: 75000, DurationMinutes: 75, Status: "inactive"},
	}

	services := make([]models.Service, 0, len(seedData))
	for _, service := range seedData {
		var existingService models.Service
		if err := db.Where("name = ?", service.Name).First(&existingService).Error; err == nil {
			services = append(services, existingService)
			continue
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		if err := db.Create(&service).Error; err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	return services, nil
}

func seedBookings(db *gorm.DB, userID uint, adminID uint, services []models.Service) error {
	statuses := []string{"pending", "confirmed", "in_progress", "completed", "cancelled", "pending", "confirmed", "in_progress", "completed", "pending"}
	customerNames := []string{"Budi Santoso", "Siti Rahma", "Andi Wijaya", "Dewi Lestari", "Rizky Pratama", "Nina Kartika", "Fajar Nugroho", "Maya Putri", "Agus Salim", "Tania Safitri"}
	phones := []string{"081234560001", "081234560002", "081234560003", "081234560004", "081234560005", "081234560006", "081234560007", "081234560008", "081234560009", "081234560010"}
	vehicleNames := []string{"Honda Beat", "Yamaha NMAX", "Honda Vario", "Suzuki Satria", "Yamaha Mio", "Honda Scoopy", "Kawasaki Ninja", "Yamaha Aerox", "Honda PCX", "Vespa Sprint"}

	for index := 0; index < 10; index++ {
		bookingUserID := userID
		if index == 9 {
			bookingUserID = adminID
		}

		booking := models.Booking{
			UserID:       bookingUserID,
			ServiceID:    services[index%len(services)].ID,
			CustomerName: customerNames[index],
			Phone:        phones[index],
			VehicleName:  vehicleNames[index],
			VehiclePlate: fmt.Sprintf("B %04d MTC", index+1),
			BookingDate:  time.Now().AddDate(0, 0, index+1),
			Status:       statuses[index],
			Notes:        "Data booking contoh untuk demo MotoCare Dashboard.",
		}

		var existingBooking models.Booking
		if err := db.Where("vehicle_plate = ?", booking.VehiclePlate).First(&existingBooking).Error; err == nil {
			continue
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		if err := db.Create(&booking).Error; err != nil {
			return err
		}
	}

	return nil
}
