package v3

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

type AttributeGetter interface {
	Attribute() AttributeInterface
}
type AttributeInterface interface {
	CreateObjectAttribute(data types.MapStr) (int, error)
	DeleteObjectAttribute(cond common.Condition) error
	UpdateObjectAttribute(data types.MapStr, cond common.Condition) error
	SearchObjectAttributes(cond common.Condition) ([]types.MapStr, error)
}

type Attribute struct {
	cli *Client
}

func newAttribute(cli *Client) *Attribute {
	return &Attribute{
		cli: cli,
	}
}

// CreateObjectAttribute create a new model object attribute
func (m *Attribute) CreateObjectAttribute(data types.MapStr) (int, error) {

	targetURL := fmt.Sprintf("%s/api/v3/object/attr", m.cli.GetAddress())

	rst, err := m.cli.httpCli.POST(targetURL, nil, data.ToJSON())
	if nil != err {
		return 0, err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		return 0, errors.New(gs.Get("bk_error_msg").String())
	}

	// parse id
	id := gs.Get("data.id").Int()

	return int(id), nil
}

// DeleteObjectAttribute delete a object attribute by condition
func (m *Attribute) DeleteObjectAttribute(cond common.Condition) error {

	data := cond.ToMapStr()
	id, err := data.Int("id")
	if nil != err {
		return err
	}

	targetURL := fmt.Sprintf("%s/api/v3/object/attr/%d", m.cli.GetAddress(), id)

	rst, err := m.cli.httpCli.DELETE(targetURL, nil, nil)
	if nil != err {
		return err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		return errors.New(gs.Get("bk_error_msg").String())
	}

	return nil
}

// UpdateObjectAttribute update a object attribute by condition
func (m *Attribute) UpdateObjectAttribute(data types.MapStr, cond common.Condition) error {

	dataCond := cond.ToMapStr()
	id, err := dataCond.Int("id")
	if nil != err {
		return err
	}

	targetURL := fmt.Sprintf("%s/api/v3/object/attr/%d", m.cli.GetAddress(), id)

	rst, err := m.cli.httpCli.PUT(targetURL, nil, data.ToJSON())
	if nil != err {
		return err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		return errors.New(gs.Get("bk_error_msg").String())
	}
	return nil
}

// SearchObjectAttributes search some object attributes by condition
func (m *Attribute) SearchObjectAttributes(cond common.Condition) ([]types.MapStr, error) {

	data := cond.ToMapStr()

	targetURL := fmt.Sprintf("%s/api/v3/object/attr/search", m.cli.GetAddress())

	rst, err := m.cli.httpCli.POST(targetURL, nil, data.ToJSON())
	if nil != err {
		return nil, err
	}

	gs := gjson.ParseBytes(rst)

	// check result
	if !gs.Get("result").Bool() {
		return nil, errors.New(gs.Get("bk_error_msg").String())
	}

	dataStr := gs.Get("data").String()
	if 0 == len(dataStr) {
		return nil, errors.New("data is empty")
	}

	resultMap := make([]types.MapStr, 0)
	err = json.Unmarshal([]byte(dataStr), &resultMap)
	return resultMap, err
}
