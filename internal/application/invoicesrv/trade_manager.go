package invoicesrv

import (
	"context"
	"invoice-service/internal/domain"
)

func (s *invoiceService) ApproveOrRejectTrade(ctx context.Context, tradeID string, isApproved bool) error {
	// Get trade to operate
	trade, err := s.databaseRepository.GetTradeByID(ctx, tradeID)
	if err != nil {
		return err
	}

	// Update the trade status in the Trades table
	err = s.updateTradeStatus(ctx, trade, isApproved)
	if err != nil {
		return err
	}

	// Retrieve the invoice and investor IDs associated with the trade
	invoice, issuer, err := s.getInvoiceAndIssuerID(ctx, trade)
	if err != nil {
		return err
	}

	// Update the balances based on the trade approval or rejection
	if isApproved {
		// Update balances for issuer
		err = s.updateBalancesOnTradeApproval(ctx, issuer, invoice.AskingPrice())
		if err != nil {
			return err
		}
	} else {
		// Update balances for investor
		err = s.updateBalanceOnTradeRejection(ctx, trade)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *invoiceService) updateTradeStatus(ctx context.Context, trade *domain.Trade, isApproved bool) error {
	var newTrade = trade.UpgradeTradeStatus(isApproved)

	_, err := s.databaseRepository.SaveTrade(ctx, newTrade)
	if err != nil {
		return err
	}

	return nil
}

func (s *invoiceService) getInvoiceAndIssuerID(ctx context.Context, trade *domain.Trade) (*domain.Invoice, *domain.Issuer, error) {
	invoice, err := s.databaseRepository.GetInvoiceByID(ctx, trade.InvoiceID())
	if err != nil {
		return nil, nil, err
	}

	issuer, err := s.databaseRepository.GetIssuerByID(ctx, invoice.IssuerID())
	if err != nil {
		return nil, nil, err
	}

	return invoice, issuer, nil
}

func (s *invoiceService) updateBalancesOnTradeApproval(ctx context.Context, issuer *domain.Issuer, invoiceAmount float64) error {
	issuer.UpdateBalance(*issuer.AvailableFunds() + invoiceAmount)
	_, err := s.databaseRepository.SaveIssuer(ctx, issuer)
	if err != nil {
		return err
	}

	return nil
}

func (s *invoiceService) updateBalanceOnTradeRejection(ctx context.Context, trade *domain.Trade) error {
	for _, investorID := range trade.InvestorsIDs() {
		investor, err := s.databaseRepository.GetInvestorByID(ctx, investorID)
		if err != nil {
			return err
		}

		bids, err := s.databaseRepository.GetBidsFromInvoiceAndInvestor(ctx, trade.InvoiceID(), &investorID)
		if err != nil {
			return err
		}

		for _, b := range bids {
			investor.UpdateBalance(*investor.AvailableFunds() + b.BidAmount())
		}

		_, err = s.databaseRepository.SaveInvestor(ctx, investor)
		if err != nil {
			return err
		}
	}

	return nil
}
