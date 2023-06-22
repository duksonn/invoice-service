package invoicesrv

import (
	"context"
	"errors"
	"invoice-service/internal/domain"
)

func (s *invoiceService) PlaceBid(ctx context.Context, invoiceID string, investorID int, bidAmount float64) error {
	// Check if the investor has available funds
	investor, err := s.databaseRepository.GetInvestorByID(ctx, investorID)
	if err != nil {
		return err
	}

	// Check if the bid amount exceeds the available funds
	if bidAmount > *investor.AvailableFunds() {
		return errors.New("insufficient funds to place the bid")
	}

	// Reduce the available funds of the investor
	err = s.updateAvailableFunds(ctx, investor, bidAmount)
	if err != nil {
		return err
	}

	// Insert the bid into the Bids table
	err = s.insertBid(ctx, invoiceID, investorID, bidAmount)
	if err != nil {
		return err
	}

	// Check if the placed bids fill 100% of the invoice price
	totalBidsAmount, err := s.databaseRepository.GetTotalBidsAmount(ctx, invoiceID)
	if err != nil {
		return err
	}

	// If the bids reach 100% of the invoice price, purchase the invoice
	invoice, err := s.databaseRepository.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return err
	}
	if *totalBidsAmount >= invoice.AskingPrice() {
		err = s.purchaseInvoice(ctx, invoice, investor)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *invoiceService) updateAvailableFunds(ctx context.Context, investor *domain.Investor, bidAmount float64) error {
	var availableFounds = *investor.AvailableFunds() - bidAmount
	var newInvestor = domain.NewInvestor(investor.ID(), investor.Name(), &availableFounds)

	_, err := s.databaseRepository.SaveInvestor(ctx, newInvestor)
	if err != nil {
		return err
	}
	return nil
}

func (s *invoiceService) insertBid(ctx context.Context, invoiceID string, investorID int, bidAmount float64) error {
	var newBid = domain.CreateBid(investorID, invoiceID, bidAmount)

	_, err := s.databaseRepository.SaveBid(ctx, newBid)
	if err != nil {
		return err
	}
	return nil
}

func (s *invoiceService) purchaseInvoice(ctx context.Context, invoice *domain.Invoice, investor *domain.Investor) error {
	res, err := invoice.UpdateInvoice(domain.LOCKED.String())
	if err != nil {
		return err
	}

	_, err = s.databaseRepository.SaveInvoice(ctx, res)
	if err != nil {
		return err
	}

	bids, err := s.databaseRepository.GetBidsFromInvoiceAndInvestor(ctx, invoice.ID(), nil)
	if err != nil {
		return err
	}

	var investorsIDs []int
	for _, b := range bids {
		investorsIDs = append(investorsIDs, b.InvestorID())
	}

	// Register trade with WAITING_APPROVAL status
	var trade = domain.CreateTrade(invoice.ID(), investorsIDs)
	_, err = s.databaseRepository.SaveTrade(ctx, trade)
	if err != nil {
		return err
	}

	return nil
}
