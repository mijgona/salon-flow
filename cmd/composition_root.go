package cmd

import (
	"sync"

	"github.com/mijgona/salon-crm/internal/adapters/out/postgres"
	"github.com/mijgona/salon-crm/internal/adapters/out/postgres/appointmentrepo"
	"github.com/mijgona/salon-crm/internal/adapters/out/postgres/certificaterepo"
	"github.com/mijgona/salon-crm/internal/adapters/out/postgres/clientrepo"
	"github.com/mijgona/salon-crm/internal/adapters/out/postgres/loyaltyrepo"
	"github.com/mijgona/salon-crm/internal/adapters/out/postgres/outboxrepo"
	"github.com/mijgona/salon-crm/internal/adapters/out/postgres/schedulerepo"
	"github.com/mijgona/salon-crm/internal/core/application/commands"
	"github.com/mijgona/salon-crm/internal/core/application/eventhandlers"
	"github.com/mijgona/salon-crm/internal/core/application/queries"
	"github.com/mijgona/salon-crm/internal/core/domain/model/client"
	"github.com/mijgona/salon-crm/internal/core/domain/model/scheduling"
	"github.com/mijgona/salon-crm/internal/core/ports"
	"github.com/mijgona/salon-crm/internal/pkg/ddd"

	"github.com/jackc/pgx/v5/pgxpool"

	httpHandlers "github.com/mijgona/salon-crm/internal/adapters/in/http"
)

// CompositionRoot wires all dependencies together using lazy initialization.
type CompositionRoot struct {
	pool *pgxpool.Pool

	// Infrastructure
	txManager ports.TxManager
	mediatr   ddd.Mediatr

	// Repositories
	clientRepo      ports.ClientRepository
	appointmentRepo ports.AppointmentRepository
	scheduleRepo    ports.MasterScheduleRepository
	loyaltyRepo     ports.LoyaltyRepository
	certificateRepo ports.CertificateRepository
	outboxRepo      ports.OutboxRepository

	// Command Handlers
	registerClientHandler      *commands.RegisterClientHandler
	updateClientProfileHandler *commands.UpdateClientProfileHandler
	bookAppointmentHandler     *commands.BookAppointmentHandler
	cancelAppointmentHandler   *commands.CancelAppointmentHandler
	completeAppointmentHandler *commands.CompleteAppointmentHandler
	earnPointsHandler          *commands.EarnPointsHandler
	activateCertHandler        *commands.ActivateCertificateHandler

	// Query Handlers
	getClientHandler         *queries.GetClientHandler
	getClientHistoryHandler  *queries.GetClientHistoryHandler
	getAvailableSlotsHandler *queries.GetAvailableSlotsHandler
	getLoyaltyAccountHandler *queries.GetLoyaltyAccountHandler
	getCalendarHandler       *queries.GetCalendarHandler

	// HTTP Handlers
	clientHandler      *httpHandlers.ClientHandler
	appointmentHandler *httpHandlers.AppointmentHandler
	loyaltyHandler     *httpHandlers.LoyaltyHandler
	calendarHandler    *httpHandlers.CalendarHandler

	// Sync once instances
	onceTxManager           sync.Once
	onceMediatr             sync.Once
	onceClientRepo          sync.Once
	onceAppointmentRepo     sync.Once
	onceScheduleRepo        sync.Once
	onceLoyaltyRepo         sync.Once
	onceCertificateRepo     sync.Once
	onceOutboxRepo          sync.Once
	onceRegisterClient      sync.Once
	onceUpdateClientProfile sync.Once
	onceBookAppointment     sync.Once
	onceCancelAppointment   sync.Once
	onceCompleteAppointment sync.Once
	onceEarnPoints          sync.Once
	onceActivateCert        sync.Once
	onceGetClient           sync.Once
	onceGetClientHistory    sync.Once
	onceGetAvailableSlots   sync.Once
	onceGetLoyaltyAccount   sync.Once
	onceGetCalendar         sync.Once
	onceClientHandler       sync.Once
	onceAppointmentHandler  sync.Once
	onceLoyaltyHandler      sync.Once
	onceCalendarHandler     sync.Once
}

// NewCompositionRoot creates a new CompositionRoot backed by a pgx pool.
func NewCompositionRoot(pool *pgxpool.Pool) *CompositionRoot {
	return &CompositionRoot{pool: pool}
}

func (cr *CompositionRoot) TxManager() ports.TxManager {
	cr.onceTxManager.Do(func() {
		cr.txManager = postgres.NewTxManager(cr.pool)
	})
	return cr.txManager
}

func (cr *CompositionRoot) Mediatr() ddd.Mediatr {
	cr.onceMediatr.Do(func() {
		m := ddd.NewInProcessMediatr()

		// Subscribe event handlers
		accrueHandler := eventhandlers.NewAccruePointsOnCompletedHandler(cr.LoyaltyRepo(), cr.TxManager())
		m.Subscribe(accrueHandler, scheduling.AppointmentCompleted{})

		visitHandler := eventhandlers.NewAddVisitRecordOnCompletedHandler(cr.ClientRepo(), cr.TxManager())
		m.Subscribe(visitHandler, scheduling.AppointmentCompleted{})

		loyaltyHandler := eventhandlers.NewCreateLoyaltyOnRegisteredHandler(cr.LoyaltyRepo(), cr.TxManager())
		m.Subscribe(loyaltyHandler, client.ClientRegistered{})

		cr.mediatr = m
	})
	return cr.mediatr
}

func (cr *CompositionRoot) ClientRepo() ports.ClientRepository {
	cr.onceClientRepo.Do(func() {
		cr.clientRepo = clientrepo.NewPostgresClientRepository(cr.pool)
	})
	return cr.clientRepo
}

func (cr *CompositionRoot) AppointmentRepo() ports.AppointmentRepository {
	cr.onceAppointmentRepo.Do(func() {
		cr.appointmentRepo = appointmentrepo.NewPostgresAppointmentRepository(cr.pool)
	})
	return cr.appointmentRepo
}

func (cr *CompositionRoot) ScheduleRepo() ports.MasterScheduleRepository {
	cr.onceScheduleRepo.Do(func() {
		cr.scheduleRepo = schedulerepo.NewPostgresScheduleRepository(cr.pool)
	})
	return cr.scheduleRepo
}

func (cr *CompositionRoot) LoyaltyRepo() ports.LoyaltyRepository {
	cr.onceLoyaltyRepo.Do(func() {
		cr.loyaltyRepo = loyaltyrepo.NewPostgresLoyaltyRepository(cr.pool)
	})
	return cr.loyaltyRepo
}

func (cr *CompositionRoot) CertificateRepo() ports.CertificateRepository {
	cr.onceCertificateRepo.Do(func() {
		cr.certificateRepo = certificaterepo.NewPostgresCertificateRepository(cr.pool)
	})
	return cr.certificateRepo
}

func (cr *CompositionRoot) OutboxRepo() ports.OutboxRepository {
	cr.onceOutboxRepo.Do(func() {
		cr.outboxRepo = outboxrepo.NewPostgresOutboxRepository(cr.pool)
	})
	return cr.outboxRepo
}

func (cr *CompositionRoot) RegisterClientHandler() *commands.RegisterClientHandler {
	cr.onceRegisterClient.Do(func() {
		cr.registerClientHandler = commands.NewRegisterClientHandler(cr.ClientRepo(), cr.TxManager())
	})
	return cr.registerClientHandler
}

func (cr *CompositionRoot) UpdateClientProfileHandler() *commands.UpdateClientProfileHandler {
	cr.onceUpdateClientProfile.Do(func() {
		cr.updateClientProfileHandler = commands.NewUpdateClientProfileHandler(cr.ClientRepo(), cr.TxManager())
	})
	return cr.updateClientProfileHandler
}

func (cr *CompositionRoot) BookAppointmentHandler() *commands.BookAppointmentHandler {
	cr.onceBookAppointment.Do(func() {
		// BookAppointment needs a ServiceCatalogClient — nil for now (to be wired later)
		cr.bookAppointmentHandler = commands.NewBookAppointmentHandler(
			cr.AppointmentRepo(), cr.ScheduleRepo(), nil, cr.TxManager(),
		)
	})
	return cr.bookAppointmentHandler
}

func (cr *CompositionRoot) CancelAppointmentHandler() *commands.CancelAppointmentHandler {
	cr.onceCancelAppointment.Do(func() {
		cr.cancelAppointmentHandler = commands.NewCancelAppointmentHandler(
			cr.AppointmentRepo(), cr.ScheduleRepo(), cr.TxManager(),
		)
	})
	return cr.cancelAppointmentHandler
}

func (cr *CompositionRoot) CompleteAppointmentHandler() *commands.CompleteAppointmentHandler {
	cr.onceCompleteAppointment.Do(func() {
		cr.completeAppointmentHandler = commands.NewCompleteAppointmentHandler(cr.AppointmentRepo(), cr.TxManager())
	})
	return cr.completeAppointmentHandler
}

func (cr *CompositionRoot) EarnPointsHandler() *commands.EarnPointsHandler {
	cr.onceEarnPoints.Do(func() {
		cr.earnPointsHandler = commands.NewEarnPointsHandler(cr.LoyaltyRepo(), cr.TxManager())
	})
	return cr.earnPointsHandler
}

func (cr *CompositionRoot) ActivateCertificateHandler() *commands.ActivateCertificateHandler {
	cr.onceActivateCert.Do(func() {
		cr.activateCertHandler = commands.NewActivateCertificateHandler(cr.CertificateRepo(), cr.TxManager())
	})
	return cr.activateCertHandler
}

func (cr *CompositionRoot) GetClientHandler() *queries.GetClientHandler {
	cr.onceGetClient.Do(func() {
		cr.getClientHandler = queries.NewGetClientHandler(cr.ClientRepo())
	})
	return cr.getClientHandler
}

func (cr *CompositionRoot) GetClientHistoryHandler() *queries.GetClientHistoryHandler {
	cr.onceGetClientHistory.Do(func() {
		cr.getClientHistoryHandler = queries.NewGetClientHistoryHandler(cr.AppointmentRepo())
	})
	return cr.getClientHistoryHandler
}

func (cr *CompositionRoot) GetAvailableSlotsHandler() *queries.GetAvailableSlotsHandler {
	cr.onceGetAvailableSlots.Do(func() {
		cr.getAvailableSlotsHandler = queries.NewGetAvailableSlotsHandler(cr.ScheduleRepo())
	})
	return cr.getAvailableSlotsHandler
}

func (cr *CompositionRoot) GetLoyaltyAccountHandler() *queries.GetLoyaltyAccountHandler {
	cr.onceGetLoyaltyAccount.Do(func() {
		cr.getLoyaltyAccountHandler = queries.NewGetLoyaltyAccountHandler(cr.LoyaltyRepo())
	})
	return cr.getLoyaltyAccountHandler
}

func (cr *CompositionRoot) GetCalendarHandler() *queries.GetCalendarHandler {
	cr.onceGetCalendar.Do(func() {
		cr.getCalendarHandler = queries.NewGetCalendarHandler(cr.AppointmentRepo())
	})
	return cr.getCalendarHandler
}

func (cr *CompositionRoot) ClientHTTPHandler() *httpHandlers.ClientHandler {
	cr.onceClientHandler.Do(func() {
		cr.clientHandler = httpHandlers.NewClientHandler(
			cr.RegisterClientHandler(),
			cr.UpdateClientProfileHandler(),
			cr.GetClientHandler(),
			cr.GetClientHistoryHandler(),
		)
	})
	return cr.clientHandler
}

func (cr *CompositionRoot) AppointmentHTTPHandler() *httpHandlers.AppointmentHandler {
	cr.onceAppointmentHandler.Do(func() {
		cr.appointmentHandler = httpHandlers.NewAppointmentHandler(
			cr.BookAppointmentHandler(),
			cr.CancelAppointmentHandler(),
			cr.CompleteAppointmentHandler(),
			cr.GetAvailableSlotsHandler(),
		)
	})
	return cr.appointmentHandler
}

func (cr *CompositionRoot) LoyaltyHTTPHandler() *httpHandlers.LoyaltyHandler {
	cr.onceLoyaltyHandler.Do(func() {
		cr.loyaltyHandler = httpHandlers.NewLoyaltyHandler(
			cr.EarnPointsHandler(),
			cr.GetLoyaltyAccountHandler(),
		)
	})
	return cr.loyaltyHandler
}

func (cr *CompositionRoot) CalendarHTTPHandler() *httpHandlers.CalendarHandler {
	cr.onceCalendarHandler.Do(func() {
		cr.calendarHandler = httpHandlers.NewCalendarHandler(cr.GetCalendarHandler())
	})
	return cr.calendarHandler
}
