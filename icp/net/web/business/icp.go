package business

import (
	"bytes"
	"encoding/json"
	"fmt"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nGPU/icp/common"
	"github.com/nGPU/icp/header"
	log4plus "github.com/nGPU/include/log4go"
	"net/http"
	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/principal"
	"github.com/aviate-labs/agent-go/candid/idl"
)

const (
	NodeBase              = header.Base + 350
	NotFoundBalance = NodeBase + 1
)

type ICP20 struct {
}

var gICP20 *ICP20

type Account struct {
	Owner      string `candid:"owner"`
	Subaccount []byte `candid:"subaccount"`
}

func GenerateAccount(principal string) Account {
	subaccount := make([]byte, 32)
	return Account{
		Owner:      principal,
		Subaccount: subaccount,
	}
}

func getBalance(url string, account Account) (string, error) {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"args": []interface{}{account},
	})
	resp, err := http.Post(
		fmt.Sprintf("%s", url),
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct {
		Balance string `json:"balance"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Balance, nil
}

func getBalance2(CanisterID string, userName string) (string, error) {
	// 1. Initialize agent（main network）
	client, err := agent.New(agent.DefaultConfig)
	if err != nil {
		errString:=fmt.Sprintf("agent.New failed err=[%s]", err.Error())
		return "", errors.New(errString)
	}
	// 2. Construct target Canister
	canisterID := principal.MustDecode(CanisterID)
	// 3. Construct user account
	userPrincipal, _ := principal.Decode(userName)
	account := struct {
		Owner      principal.Principal `candid:"owner"`
		Subaccount [32]byte            `candid:"subaccount"`
	}{
		Owner:      userPrincipal,
		Subaccount: [32]byte{},
	}
	// 4. call icrc1_balance_of
	var r0 idl.Nat
	if err := client.Query(canisterID, "icrc1_balance_of", []interface{}{account}, []any{&r0},); err != nil {
		errString:=fmt.Sprintf("client.Query failed err=[%s]", err.Error())
		return "", errors.New(errString)
	}
	return r0.String(), nil
}

func (w *ICP20) getICP20(c *gin.Context) {
	funName := "getICP20"
	// user Principal
	principal := c.DefaultQuery("principal", "")
	// Canister ID（for example: "mxzaz-hqaaa-aaaar-qaada-cai"）
	canisterID := "pz3cl-aaaaa-aaaaj-qngsq-cai"
	log4plus.Info("%s principal=[%s] canisterID=[%s]", funName, principal, canisterID)
	// Balance
	balance, err := getBalance2(canisterID, principal)
	if err != nil {
		errString := fmt.Sprintf("%s getBalance not found principal=[%s] canisterID=[%s]", funName, principal, canisterID)
		log4plus.Error(errString)
		common.SendError(c, NotFoundBalance, errString)
		return
	}
	response := struct {
		CodeId     int64  `json:"codeID"`
		Msg        string `json:"msg"`
		Principal  string `json:"principal"`
		CanisterID string `json:"canisterID"`
		Balance    string `json:"balance"`
	}{
		CodeId:     http.StatusOK,
		Msg:        "success",
		Principal:  principal,
		CanisterID: canisterID,
		Balance:    balance,
	}
	c.JSON(http.StatusOK, response)
}

func (w *ICP20) Start(nodeGroup *gin.RouterGroup) {
	nodeGroup.GET("/getICP20", w.getICP20)
}

func SingletonIcpERC20() *ICP20 {
	if gICP20 == nil {
		gICP20 = &ICP20{}
	}
	return gICP20
}
