package transformer

import (
	"github.com/DioGolang/home-broker/internal/market/dto"
	"github.com/DioGolang/home-broker/internal/market/entity"
)

func TrnasformImput(input dto.TradeInput) *entity.Order {
	asset := entity.NewAsset(input.AssetID, input.AssetID, 100)
	investor := entity.NewInvestor(input.InvestorID)
	order := entity.NewOrder(input.OrderId, investor, asset, input.Shares, input.Price, input.OrderType)
	if input.CurrentShares > 0 {
		assetPosition := entity.NewInvestorAssetPosition(input.AssetID, input.CurrentShares)
		investor.AddAssetPosition(assetPosition)
	}
	return order
}
