package main

import (
	sap_api_caller "sap-api-integrations-production-order-conf-reads/SAP_API_Caller"
	sap_api_input_reader "sap-api-integrations-production-order-conf-reads/SAP_API_Input_Reader"
	"sap-api-integrations-production-order-conf-reads/config"

	"github.com/latonaio/golang-logging-library-for-sap/logger"
	sap_api_request_client_header_setup "github.com/latonaio/sap-api-request-client-header-setup"
	sap_api_time_value_converter "github.com/latonaio/sap-api-time-value-converter"
)

func main() {
	l := logger.NewLogger()
	conf := config.NewConf()
	fr := sap_api_input_reader.NewFileReader()
	gc := sap_api_request_client_header_setup.NewSAPRequestClientWithOption(conf.SAP)
	caller := sap_api_caller.NewSAPAPICaller(
		conf.SAP.BaseURL(),
		"100",
		gc,
		l,
	)
	inputSDC := fr.ReadSDC("./Inputs/SDC_Production_Order_Confirmation_Conf_By_OrderID_Seq_Op_sample.json")
	sap_api_time_value_converter.ChangeTimeFormatToSAPFormatStruct(&inputSDC)
	accepter := inputSDC.Accepter
	if len(accepter) == 0 || accepter[0] == "All" {
		accepter = []string{
			"ConfByOrderID", "MaterialMovements", "BatchCharacteristic",
			"ConfByOrderIDConfGroup", "ConfByOrderIDSeqOp",
		}
	}

	caller.AsyncGetProductionOrderConfirmation(
		inputSDC.ProductionOrderConfirmation.OrderID,
		inputSDC.ProductionOrderConfirmation.MaterialMovements.Batch,
		inputSDC.ProductionOrderConfirmation.ConfirmationGroup,
		inputSDC.ProductionOrderConfirmation.Sequence,
		inputSDC.ProductionOrderConfirmation.OrderOperation,
		accepter,
	)
}
