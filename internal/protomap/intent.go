package protomap

import (
	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v1"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

func PBFIntent(i alpaca.PositionIntent) tradesvc.PositionIntent {
	switch i {
	case alpaca.BuyToClose:
		return tradesvc.PositionIntent_POSITION_INTENT_BUY_TO_CLOSE
	case alpaca.BuyToOpen:
		return tradesvc.PositionIntent_POSITION_INTENT_BUY_TO_OPEN
	case alpaca.SellToClose:
		return tradesvc.PositionIntent_POSITION_INTENT_SELL_TO_CLOSE
	case alpaca.SellToOpen:
		return tradesvc.PositionIntent_POSITION_INTENT_SELL_TO_OPEN
	default:
		return tradesvc.PositionIntent_POSITION_INTENT_UNSPECIFIED
	}
}

func Intent(i tradesvc.PositionIntent) alpaca.PositionIntent {
	switch i {
	case tradesvc.PositionIntent_POSITION_INTENT_BUY_TO_CLOSE:
		return alpaca.BuyToClose
	case tradesvc.PositionIntent_POSITION_INTENT_BUY_TO_OPEN:
		return alpaca.BuyToOpen
	case tradesvc.PositionIntent_POSITION_INTENT_SELL_TO_CLOSE:
		return alpaca.SellToClose
	case tradesvc.PositionIntent_POSITION_INTENT_SELL_TO_OPEN:
		return alpaca.SellToOpen
	default:
		return ""
	}
}
