package handlers

// RegisterDoc godoc
// @Summary Register user
// @Description Create a new user account. Role is optional and defaults to user.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequestDoc true "Register request body"
// @Success 201 {object} RegisterResponseDoc "register berhasil"
// @Failure 400 {object} ErrorResponseDoc "validasi gagal"
// @Failure 409 {object} ErrorResponseDoc "username atau email sudah digunakan"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /register [post]
func RegisterDoc() {}

// LoginDoc godoc
// @Summary Login user
// @Description Authenticate user and return JWT token. Use the returned token in Swagger Authorize as Bearer <token>.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequestDoc true "Login request body"
// @Success 200 {object} LoginResponseDoc "login berhasil"
// @Failure 400 {object} ErrorResponseDoc "validasi gagal"
// @Failure 401 {object} ErrorResponseDoc "email atau password salah"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /login [post]
func LoginDoc() {}

// ChangePasswordDoc godoc
// @Summary Change password
// @Description Change own password. Admin may change another user's password by sending user_id.
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequestDoc true "Change password request body"
// @Success 200 {object} MessageResponseDoc "password berhasil diubah"
// @Failure 400 {object} ErrorResponseDoc "validasi gagal"
// @Failure 401 {object} ErrorResponseDoc "token tidak valid atau current_password salah"
// @Failure 403 {object} ErrorResponseDoc "akses ditolak"
// @Failure 404 {object} ErrorResponseDoc "user tidak ditemukan"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /change-password [put]
func ChangePasswordDoc() {}

// MeDoc godoc
// @Summary Get current user
// @Description Return authenticated user profile from JWT token.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponseDoc "data user berhasil diambil"
// @Failure 401 {object} ErrorResponseDoc "token tidak valid"
// @Failure 404 {object} ErrorResponseDoc "user tidak ditemukan"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /me [get]
func MeDoc() {}

// ListCategoriesDoc godoc
// @Summary List categories
// @Description Public endpoint to list service categories.
// @Tags Categories
// @Produce json
// @Success 200 {object} CategoryListResponseDoc "success"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/categories [get]
func ListCategoriesDoc() {}

// DetailCategoryDoc godoc
// @Summary Get category detail
// @Description Public endpoint to get one service category by ID.
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} CategoryResponseDoc "kategori berhasil diambil"
// @Failure 400 {object} ErrorResponseDoc "id tidak valid"
// @Failure 404 {object} ErrorResponseDoc "kategori tidak ditemukan"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/categories/{id} [get]
func DetailCategoryDoc() {}

// CreateCategoryDoc godoc
// @Summary Create category
// @Description Admin only. Create a service category.
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CategoryRequestDoc true "Category request body"
// @Success 201 {object} CategoryResponseDoc "kategori berhasil dibuat"
// @Failure 400 {object} ErrorResponseDoc "validasi gagal"
// @Failure 401 {object} ErrorResponseDoc "token tidak valid"
// @Failure 403 {object} ErrorResponseDoc "akses ditolak"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/categories [post]
func CreateCategoryDoc() {}

// UpdateCategoryDoc godoc
// @Summary Update category
// @Description Admin only. Update a service category.
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param request body CategoryRequestDoc true "Category request body"
// @Success 200 {object} CategoryResponseDoc "kategori berhasil diubah"
// @Failure 400 {object} ErrorResponseDoc "validasi gagal"
// @Failure 401 {object} ErrorResponseDoc "token tidak valid"
// @Failure 403 {object} ErrorResponseDoc "akses ditolak"
// @Failure 404 {object} ErrorResponseDoc "kategori tidak ditemukan"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/categories/{id} [put]
func UpdateCategoryDoc() {}

// DeleteCategoryDoc godoc
// @Summary Delete category
// @Description Admin only. Delete a service category.
// @Tags Categories
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} MessageResponseDoc "kategori berhasil dihapus"
// @Failure 400 {object} ErrorResponseDoc "id tidak valid"
// @Failure 401 {object} ErrorResponseDoc "token tidak valid"
// @Failure 403 {object} ErrorResponseDoc "akses ditolak"
// @Failure 404 {object} ErrorResponseDoc "kategori tidak ditemukan"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/categories/{id} [delete]
func DeleteCategoryDoc() {}

// ListServicesDoc godoc
// @Summary List services
// @Description Public endpoint to list services with pagination, search, filters, and sorting.
// @Tags Services
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by name or description"
// @Param category_id query int false "Filter by category ID"
// @Param status query string false "Filter by status" Enums(active,inactive)
// @Param sort_by query string false "Sort column" Enums(id,name,price,duration_minutes,status,created_at)
// @Param sort_order query string false "Sort order" Enums(asc,desc)
// @Success 200 {object} ServiceListResponseDoc "success"
// @Failure 400 {object} ErrorResponseDoc "query tidak valid"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/services [get]
func ListServicesDoc() {}

// DetailServiceDoc godoc
// @Summary Get service detail
// @Description Public endpoint to get one service by ID.
// @Tags Services
// @Produce json
// @Param id path int true "Service ID"
// @Success 200 {object} ServiceResponseDoc "layanan berhasil diambil"
// @Failure 400 {object} ErrorResponseDoc "id tidak valid"
// @Failure 404 {object} ErrorResponseDoc "layanan tidak ditemukan"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/services/{id} [get]
func DetailServiceDoc() {}

// CreateServiceDoc godoc
// @Summary Create service
// @Description Admin only. Create a motorcycle service.
// @Tags Services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ServiceRequestDoc true "Service request body"
// @Success 201 {object} ServiceResponseDoc "layanan berhasil dibuat"
// @Failure 400 {object} ErrorResponseDoc "validasi gagal"
// @Failure 401 {object} ErrorResponseDoc "token tidak valid"
// @Failure 403 {object} ErrorResponseDoc "akses ditolak"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/services [post]
func CreateServiceDoc() {}

// UpdateServiceDoc godoc
// @Summary Update service
// @Description Admin only. Update a motorcycle service.
// @Tags Services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Service ID"
// @Param request body ServiceRequestDoc true "Service request body"
// @Success 200 {object} ServiceResponseDoc "layanan berhasil diubah"
// @Failure 400 {object} ErrorResponseDoc "validasi gagal"
// @Failure 401 {object} ErrorResponseDoc "token tidak valid"
// @Failure 403 {object} ErrorResponseDoc "akses ditolak"
// @Failure 404 {object} ErrorResponseDoc "layanan tidak ditemukan"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/services/{id} [put]
func UpdateServiceDoc() {}

// DeleteServiceDoc godoc
// @Summary Delete service
// @Description Admin only. Delete a motorcycle service.
// @Tags Services
// @Produce json
// @Security BearerAuth
// @Param id path int true "Service ID"
// @Success 200 {object} MessageResponseDoc "layanan berhasil dihapus"
// @Failure 400 {object} ErrorResponseDoc "id tidak valid"
// @Failure 401 {object} ErrorResponseDoc "token tidak valid"
// @Failure 403 {object} ErrorResponseDoc "akses ditolak"
// @Failure 404 {object} ErrorResponseDoc "layanan tidak ditemukan"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/services/{id} [delete]
func DeleteServiceDoc() {}

// ListBookingsDoc godoc
// @Summary List bookings
// @Description JWT required. Admin sees all bookings, user sees only their own bookings.
// @Tags Bookings
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search customer, phone, vehicle name, or plate"
// @Param status query string false "Filter by status" Enums(pending,confirmed,in_progress,completed,cancelled)
// @Param sort_by query string false "Sort column" Enums(id,customer_name,booking_date,status,created_at)
// @Param sort_order query string false "Sort order" Enums(asc,desc)
// @Success 200 {object} BookingListResponseDoc "success"
// @Failure 400 {object} ErrorResponseDoc "query tidak valid"
// @Failure 401 {object} ErrorResponseDoc "token tidak valid"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/bookings [get]
func ListBookingsDoc() {}

// DetailBookingDoc godoc
// @Summary Get booking detail
// @Description JWT required. Admin can view any booking, user can view only their own booking.
// @Tags Bookings
// @Produce json
// @Security BearerAuth
// @Param id path int true "Booking ID"
// @Success 200 {object} BookingResponseDoc "booking berhasil diambil"
// @Failure 400 {object} ErrorResponseDoc "id tidak valid"
// @Failure 401 {object} ErrorResponseDoc "token tidak valid"
// @Failure 403 {object} ErrorResponseDoc "akses ditolak"
// @Failure 404 {object} ErrorResponseDoc "booking tidak ditemukan"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/bookings/{id} [get]
func DetailBookingDoc() {}

// CreateBookingDoc godoc
// @Summary Create booking
// @Description JWT required. User creates a booking for their own account. User-created bookings must be pending.
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BookingRequestDoc true "Booking request body"
// @Success 201 {object} BookingResponseDoc "booking berhasil dibuat"
// @Failure 400 {object} ErrorResponseDoc "validasi gagal"
// @Failure 401 {object} ErrorResponseDoc "token tidak valid"
// @Failure 403 {object} ErrorResponseDoc "akses ditolak"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/bookings [post]
func CreateBookingDoc() {}

// UpdateBookingDoc godoc
// @Summary Update booking status
// @Description Admin only. Update booking status.
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Booking ID"
// @Param request body BookingStatusRequestDoc true "Booking status request body"
// @Success 200 {object} BookingResponseDoc "status booking berhasil diubah"
// @Failure 400 {object} ErrorResponseDoc "validasi gagal"
// @Failure 401 {object} ErrorResponseDoc "token tidak valid"
// @Failure 403 {object} ErrorResponseDoc "akses ditolak"
// @Failure 404 {object} ErrorResponseDoc "booking tidak ditemukan"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/bookings/{id} [put]
func UpdateBookingDoc() {}

// DeleteBookingDoc godoc
// @Summary Delete booking
// @Description Admin only. Delete booking.
// @Tags Bookings
// @Produce json
// @Security BearerAuth
// @Param id path int true "Booking ID"
// @Success 200 {object} MessageResponseDoc "booking berhasil dihapus"
// @Failure 400 {object} ErrorResponseDoc "id tidak valid"
// @Failure 401 {object} ErrorResponseDoc "token tidak valid"
// @Failure 403 {object} ErrorResponseDoc "akses ditolak"
// @Failure 404 {object} ErrorResponseDoc "booking tidak ditemukan"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/bookings/{id} [delete]
func DeleteBookingDoc() {}

// DashboardStatsDocEndpoint godoc
// @Summary Dashboard statistics
// @Description JWT required. Admin sees all data, user sees only their own booking statistics.
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} DashboardStatsDoc "dashboard statistics"
// @Failure 401 {object} ErrorResponseDoc "token tidak valid"
// @Failure 500 {object} ErrorResponseDoc "server error"
// @Router /api/dashboard/stats [get]
func DashboardStatsDocEndpoint() {}
