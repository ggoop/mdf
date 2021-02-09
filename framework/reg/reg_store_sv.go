package reg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io"
	"os"
	"path"

	"github.com/ggoop/mdf/framework/configs"
	"github.com/ggoop/mdf/utils"
)

var reg_store *regStoreSv

type regStoreSv struct {
	data   map[string]*RegObject
	dbFile string
}

func Start() {
	reg_store = &regStoreSv{data: make(map[string]*RegObject)}
	reg_store.Init()
	reg_store.Register()
}

func (s *regStoreSv) Register() {
	address := configs.Default.App.Address
	if address == "" {
		address = fmt.Sprintf("http://127.0.0.1:%s", configs.Default.App.Port)
	}
	s.Add(RegObject{
		Code:          configs.Default.App.Code,
		Name:          configs.Default.App.Name,
		Address:       address,
		PublicAddress: configs.Default.App.PublicAddress,
		Configs:       configs.Default,
	})
}
func (s *regStoreSv) Add(item RegObject) *RegObject {
	item.Time = utils.NewTimePtr()
	s.data[item.Key()] = &item
	s.Store()

	setRegObjectCache(&item)

	return &item
}

func (s *regStoreSv) Get(item RegObject) *RegObject {
	if old, ok := s.data[item.Key()]; ok {
		return old
	} else {
		return nil
	}
}

func (s *regStoreSv) GetAll() []RegObject {
	items := make([]RegObject, 0)
	for _, item := range s.data {
		items = append(items, *item)
	}
	return items
}
func (s *regStoreSv) Init() {
	if s.dbFile == "" {
		s.dbFile = utils.JoinCurrentPath(path.Join(configs.Default.App.Storage, "uploads", "regs"))
	}
	if !utils.PathExists(s.dbFile) {
		return
	}
	fi, err := os.Open(s.dbFile)
	if err != nil {
		glog.Error(err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		item := RegObject{}
		err = json.Unmarshal(a, &item)
		if err != nil {
			glog.Error(err)
			return
		}
		s.Add(item)
	}
}
func (s *regStoreSv) Store() {
	items := s.GetAll()
	f, err := os.Create(s.dbFile)
	if err != nil {
		glog.Error(err)
		f.Close()
		return
	}
	for _, item := range items {
		b, err := json.Marshal(item)
		if err != nil {
			glog.Error(err)
			return
		}
		fmt.Fprintln(f, string(b))
		if err != nil {
			glog.Error(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		glog.Error(err)
		return
	}
}
