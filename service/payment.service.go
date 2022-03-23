package service

import (
	"github.com/joho/godotenv"
	midtrans "github.com/veritrans/go-midtrans"
	"go-fundraising/entity"
	"log"
	"os"
	"strconv"
)

type PaymentService interface {
	GetPaymentURL(user entity.User, transaction entity.Transaction) (string, error)
}

type paymentService struct {
}

func NewPaymentService() PaymentService {
	return &paymentService{}
}

func (s *paymentService) GetPaymentURL(user entity.User, transaction entity.Transaction) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err.Error())
	}
	MIDTRANS_SERVER_KEY := os.Getenv("MIDTRANS_SERVER_KEY")
	MIDTRANS_CLIENT_KEY := os.Getenv("MIDTRANS_CLIENT_KEY")

	midclient := midtrans.NewClient()
	midclient.ServerKey = MIDTRANS_SERVER_KEY
	midclient.ClientKey = MIDTRANS_CLIENT_KEY
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{Client: midclient}

	snapReq := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
		CustomerDetail: &midtrans.CustDetail{
			FName: user.Name,
			Email: user.Email,
		},
	}
	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}
	return snapTokenResp.RedirectURL, nil
}
