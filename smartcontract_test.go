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

func TestCreateProduct(t *testing.T) {
	contract := &SupplyChainContract{}
	cc, err := contractapi.NewChaincode(contract)
	assert.NoError(t, err)

	stub := shimtest.NewMockStub("supplychain", cc)
	ctx := &mockContext{stub: stub}

	exists, err := contract.ProductExists(ctx, "p3")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Start a mock transaction to set TxTimestamp
	stub.MockTransactionStart("tx1")
	err = contract.CreateProduct(ctx, "p3", "Tablet", "CompanyC", "High-performance tablet", "Electronics")
	stub.MockTransactionEnd("tx1")
	assert.NoError(t, err)

	exists, err = contract.ProductExists(ctx, "p3")
	assert.NoError(t, err)
	assert.True(t, exists)

	product, err := contract.QueryProduct(ctx, "p3")
	assert.NoError(t, err)
	assertProduct(t, product, "Tablet", "Manufactured", "CompanyC", "High-performance tablet", "Electronics")
}

func TestUpdateProduct(t *testing.T) {

	contract := &SupplyChainContract{}
	cc, err := contractapi.NewChaincode(contract)
	assert.NoError(t, err)

	stub := shimtest.NewMockStub("supplychain", cc)
	ctx := &mockContext{stub: stub}

	stub.MockTransactionStart("tx1")
	err = contract.CreateProduct(ctx, "p5", "Wireless Headphones", "CompanyE", "Noise-cancelling headphones", "Electronics")
	stub.MockTransactionEnd("tx1")
	assert.NoError(t, err)

	stub.MockTransactionStart("tx2")
	err = contract.UpdateProduct(ctx, "p5", "Sold", "CustomerA", "High-quality wireless headphones", "Advance Electronics")
	stub.MockTransactionEnd("tx2")
	assert.NoError(t, err)

	product, err := contract.QueryProduct(ctx, "p5")
	assert.NoError(t, err)
	assertProduct(t, product, "Wireless Headphones", "Sold", "CustomerA", "High-quality wireless headphones", "Advance Electronics")
}

func TestQueryProduct(t *testing.T) {
	contract := &SupplyChainContract{}
	cc, err := contractapi.NewChaincode(contract)
	assert.NoError(t, err)

	stub := shimtest.NewMockStub("supplychain", cc)
	ctx := &mockContext{stub: stub}

	stub.MockTransactionStart("tx1")
	err = contract.CreateProduct(ctx, "p4", "Smartwatch", "CompanyD", "Feature-rich smartwatch", "Electronics")
	stub.MockTransactionEnd("tx1")
	assert.NoError(t, err)

	product, err := contract.QueryProduct(ctx, "p4")
	assert.NoError(t, err)
	assertProduct(t, product, "Smartwatch", "Manufactured", "CompanyD", "Feature-rich smartwatch", "Electronics")

	_, err = contract.QueryProduct(ctx, "nonexistent")
	assert.Error(t, err, "the product does not exist")
	assert.Contains(t, err.Error(), "the product does not exist")
}
