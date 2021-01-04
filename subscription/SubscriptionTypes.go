package exactonline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	types "github.com/leapforce-libraries/go_types"
	utilities "github.com/leapforce-libraries/go_utilities"
)

// SubscriptionType stores SubscriptionType from exactonline
//
type SubscriptionType struct {
	ID               types.GUID  `json:"ID"`
	Code             string      `json:"Code"`
	Created          *types.Date `json:"Created"`
	Creator          types.GUID  `json:"Creator"`
	CreatorFullName  string      `json:"CreatorFullName"`
	Description      string      `json:"Description"`
	Division         int32       `json:"Division"`
	Modified         *types.Date `json:"Modified"`
	Modifier         types.GUID  `json:"Modifier"`
	ModifierFullName string      `json:"ModifierFullName"`
}

// SubscriptionTypeUpdate stores SubscriptionType values for insert/update
//
type SubscriptionTypeUpdate struct {
	Code        *string `json:"Code,omitempty"`
	Description *string `json:"Description,omitempty"`
}

type GetSubscriptionTypesCall struct {
	urlNext string
	client  *Client
}

type GetSubscriptionTypesCallParams struct {
	ModifiedAfter *time.Time
}

func (c *Client) NewGetSubscriptionTypesCall(params *GetSubscriptionTypesCallParams) *GetSubscriptionTypesCall {
	call := GetSubscriptionTypesCall{}
	call.client = c

	selectFields := utilities.GetTaggedTagNames("json", SubscriptionType{})
	call.urlNext = fmt.Sprintf("%s/SubscriptionTypes?$select=%s", c.BaseURL(), selectFields)

	filter := []string{}

	if params != nil {
		if params.ModifiedAfter != nil {
			filter = append(filter, c.DateFilter("Modified", "gt", params.ModifiedAfter, true, "&"))
		}
	}

	if len(filter) > 0 {
		call.urlNext = fmt.Sprintf("%s&$filter=%s", call.urlNext, strings.Join(filter, " AND "))
	}

	return &call
}

func (call *GetSubscriptionTypesCall) Do() (*[]SubscriptionType, *errortools.Error) {
	if call.urlNext == "" {
		return nil, nil
	}

	subscriptionTypes := []SubscriptionType{}

	next, err := call.client.Get(call.urlNext, &subscriptionTypes)
	if err != nil {
		return nil, err
	}

	call.urlNext = next

	return &subscriptionTypes, nil
}

func (call *GetSubscriptionTypesCall) DoAll() (*[]SubscriptionType, *errortools.Error) {
	subscriptionTypes := []SubscriptionType{}

	for true {
		_subscriptionTypes, e := call.Do()
		if e != nil {
			return nil, e
		}

		if _subscriptionTypes == nil {
			break
		}

		if len(*_subscriptionTypes) == 0 {
			break
		}

		subscriptionTypes = append(subscriptionTypes, *_subscriptionTypes...)
	}

	return &subscriptionTypes, nil
}

func (c *Client) CreateSubscriptionType(subscriptionType *SubscriptionTypeUpdate) (*SubscriptionType, *errortools.Error) {
	url := fmt.Sprintf("%s/SubscriptionTypes", c.BaseURL())

	b, err := json.Marshal(subscriptionType)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	subscriptionTypeNew := SubscriptionType{}

	e := c.Post(url, bytes.NewBuffer(b), &subscriptionTypeNew)
	if e != nil {
		return nil, e
	}
	return &subscriptionTypeNew, nil
}

func (c *Client) UpdateSubscriptionType(id types.GUID, subscriptionType *SubscriptionTypeUpdate) *errortools.Error {
	url := fmt.Sprintf("%s/SubscriptionTypes(guid'%s')", c.BaseURL(), id.String())

	b, err := json.Marshal(subscriptionType)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	e := c.Put(url, bytes.NewBuffer(b))
	if e != nil {
		return e
	}
	return nil
}

func (c *Client) DeleteSubscriptionType(id types.GUID) *errortools.Error {
	url := fmt.Sprintf("%s/SubscriptionTypes(guid'%s')", c.BaseURL(), id.String())

	err := c.Delete(url)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetSubscriptionTypesCount(createdBefore *time.Time) (int64, *errortools.Error) {
	return c.GetCount("SubscriptionTypes", createdBefore)
}