package main

import (
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
)

type mockContext struct {
	contractapi.TransactionContextInterface
	stub *shimtest.MockStub
}

func (m *mockContext) GetStub() shim.ChaincodeStubInterface {
	return m.stub
}

func TestInitLedgerAndGetAllProducts(t *testing.T) {
	contract := &SupplyChainContract{}
	cc, err := contractapi.NewChaincode(contract)
	assert.NoError(t, err)

	stub := shimtest.NewMockStub("supplychain", cc)
	ctx := &mockContext{stub: stub}

	// Start a mock transaction to set TxTimestamp
	stub.MockTransactionStart("tx1")
	err = contract.InitLedger(ctx)
	stub.MockTransactionEnd("tx1")
	assert.NoError(t, err)

	products, err := contract.GetAllProducts(ctx)
	assert.NoError(t, err)
	assert.Len(t, products, 2)

	for _, p := range products {
		if p.ID == "p1" {
			assertProduct(t, p, "Laptop", "Manufactured", "CompanyA", "High-end gaming laptop", "Electronics")
		} else if p.ID == "p2" {
			assertProduct(t, p, "Smartphone", "Manufactured", "CompanyB", "Latest model smartphone", "Electronics")
		}
	}
}

func assertProduct(t *testing.T, product *Product, name string, status string, owner string, description string, category string) {
	assert.Equal(t, name, product.Name)
	assert.Equal(t, status, product.Status)
	assert.Equal(t, owner, product.Owner)
	assert.Equal(t, description, product.Description)
	assert.Equal(t, category, product.Category)
}
