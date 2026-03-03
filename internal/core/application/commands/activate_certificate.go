package commands

import (
	"context"
	"github.com/mijgona/salon-crm/internal/core/ports"

	"github.com/google/uuid"
)

// ActivateCertificateCommand holds data for activating a certificate.
type ActivateCertificateCommand struct {
	CertificateID uuid.UUID
	ClientID      uuid.UUID
}

// ActivateCertificateHandler handles certificate activation.
type ActivateCertificateHandler struct {
	certificateRepo ports.CertificateRepository
	txManager       ports.TxManager
}

// NewActivateCertificateHandler creates a new handler.
func NewActivateCertificateHandler(certificateRepo ports.CertificateRepository, txManager ports.TxManager) *ActivateCertificateHandler {
	return &ActivateCertificateHandler{certificateRepo: certificateRepo, txManager: txManager}
}

// Handle executes the activate certificate command.
func (h *ActivateCertificateHandler) Handle(ctx context.Context, cmd ActivateCertificateCommand) error {
	return h.txManager.Execute(ctx, func(tx interface{}) error {
		cert, err := h.certificateRepo.Get(ctx, tx, cmd.CertificateID)
		if err != nil {
			return err
		}

		if err := cert.Activate(cmd.ClientID); err != nil {
			return err
		}

		return h.certificateRepo.Update(ctx, tx, cert)
	})
}
